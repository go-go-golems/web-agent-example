package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type wsOptions struct {
	baseOptions
	FilterType   string
	Raw          bool
	PingInterval time.Duration
	Timeout      time.Duration
}

func runWS(args []string) error {
	opts := &wsOptions{}
	fs := flag.NewFlagSet("ws", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&opts.Backend, "backend", "http://localhost:8080", "Backend origin (http://host:port)")
	fs.StringVar(&opts.ConvID, "conv-id", "", "Conversation ID (required)")
	fs.StringVar(&opts.Profile, "profile", "", "Profile slug (optional)")
	fs.StringVar(&opts.FilterType, "filter-type", "", "Only print SEM events with this prefix")
	fs.BoolVar(&opts.Pretty, "pretty", true, "Pretty-print SEM event data")
	fs.BoolVar(&opts.Raw, "raw", false, "Print raw JSON frames only")
	fs.DurationVar(&opts.PingInterval, "ping-interval", 5*time.Second, "Send ping every interval (0 to disable)")
	fs.DurationVar(&opts.Timeout, "timeout", 0, "Stop after this duration (0 = no timeout)")
	fs.Func("cookie", "Cookie in name=value form (repeatable)", func(v string) error {
		if strings.TrimSpace(v) == "" {
			return errors.New("cookie cannot be empty")
		}
		opts.Cookies = append(opts.Cookies, v)
		return nil
	})
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(opts.ConvID) == "" {
		return errors.New("missing --conv-id")
	}

	wsURL, err := buildWSURL(opts.Backend, opts.ConvID, opts.Profile)
	if err != nil {
		return err
	}

	headers := map[string][]string{}
	if len(opts.Cookies) > 0 {
		cookieHeader, err := cookieHeaderValue(opts.Cookies)
		if err != nil {
			return err
		}
		headers["Cookie"] = []string{cookieHeader}
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	fmt.Fprintf(os.Stdout, "ws connected: %s\n", wsURL)

	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"), time.Now().Add(time.Second))
		_ = conn.Close()
	}()

	if opts.PingInterval > 0 {
		go pingLoop(ctx, conn, opts.PingInterval)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, data, err := conn.ReadMessage()
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return err
			}
			if err := renderWSFrame(data, opts); err != nil {
				return err
			}
		}
	}
}

func buildWSURL(backend, convID, profile string) (string, error) {
	base := normalizeBackend(backend)
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported backend scheme: %s", u.Scheme)
	}
	q := u.Query()
	q.Set("conv_id", convID)
	if strings.TrimSpace(profile) != "" {
		q.Set("profile", strings.TrimSpace(profile))
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/ws"
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func cookieHeaderValue(cookies []string) (string, error) {
	parts := make([]string, 0, len(cookies))
	for _, raw := range cookies {
		kv := strings.SplitN(raw, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid cookie %q (expected key=value)", raw)
		}
		parts = append(parts, strings.TrimSpace(kv[0])+"="+strings.TrimSpace(kv[1]))
	}
	return strings.Join(parts, "; "), nil
}

func pingLoop(ctx context.Context, conn *websocket.Conn, every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		}
	}
}

func renderWSFrame(data []byte, opts *wsOptions) error {
	if opts.Raw {
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}
	if sem, ok := payload["sem"].(bool); ok && sem {
		ev, ok := payload["event"].(map[string]any)
		if !ok {
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}
		etype, _ := ev["type"].(string)
		if opts.FilterType != "" && !strings.HasPrefix(etype, opts.FilterType) {
			return nil
		}
		eid, _ := ev["id"].(string)
		fmt.Fprintf(os.Stdout, "sem %s %s\n", etype, eid)
		if dataField, ok := ev["data"]; ok {
			if opts.Pretty {
				b, err := json.MarshalIndent(dataField, "", "  ")
				if err == nil {
					fmt.Fprintln(os.Stdout, string(b))
					return nil
				}
			}
			b, err := json.Marshal(dataField)
			if err == nil {
				fmt.Fprintln(os.Stdout, string(b))
			}
		}
		return nil
	}
	if opts.Pretty {
		if err := prettyPrintJSON(data, os.Stdout); err == nil {
			return nil
		}
	}
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}
