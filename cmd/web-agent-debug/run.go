//glazedclilint:file-ignore legacy debug Cobra command uses raw flags; migrate to Glazed fields in a follow-up
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type runOptions struct {
	baseOptions
	Prompt        string
	ThinkingMode  string
	FilterType    string
	Raw           bool
	PingInterval  time.Duration
	Timeout       time.Duration
	TimelineDelay time.Duration
	SinceVersion  uint64
	Limit         int
}

func runRun(args []string) error {
	opts := &runOptions{}
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&opts.Backend, "backend", "http://localhost:8080", "Backend origin (http://host:port)")
	fs.StringVar(&opts.ConvID, "conv-id", "", "Conversation ID (generated if empty)")
	fs.StringVar(&opts.Profile, "profile", "default", "Profile slug")
	fs.StringVar(&opts.Prompt, "prompt", "", "Prompt to send")
	fs.StringVar(&opts.ThinkingMode, "thinking-mode", "", "Thinking mode override (fast/slow/etc)")
	fs.StringVar(&opts.FilterType, "filter-type", "", "Only print SEM events with this prefix")
	fs.BoolVar(&opts.Pretty, "pretty", true, "Pretty-print SEM event data")
	fs.BoolVar(&opts.Raw, "raw", false, "Print raw JSON frames only")
	fs.DurationVar(&opts.PingInterval, "ping-interval", 5*time.Second, "Send ping every interval (0 to disable)")
	fs.DurationVar(&opts.Timeout, "timeout", 20*time.Second, "Overall run timeout")
	fs.DurationVar(&opts.TimelineDelay, "timeline-delay", 2*time.Second, "Wait before fetching /timeline")
	fs.Uint64Var(&opts.SinceVersion, "since-version", 0, "Fetch entities since version")
	fs.IntVar(&opts.Limit, "limit", 0, "Limit number of entities")
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
	if opts.Prompt == "" {
		if fs.NArg() > 0 {
			opts.Prompt = strings.Join(fs.Args(), " ")
		} else {
			return errors.New("missing --prompt")
		}
	}

	convWasEmpty := strings.TrimSpace(opts.ConvID) == ""
	convID := ensureConvID(opts.ConvID)
	opts.ConvID = convID
	if convWasEmpty {
		fmt.Fprintf(os.Stdout, "conv_id: %s\n", convID)
	}

	wsOpts := &wsOptions{
		baseOptions: baseOptions{
			Backend: opts.Backend,
			ConvID:  convID,
			Profile: opts.Profile,
			Cookies: opts.Cookies,
			Pretty:  opts.Pretty,
		},
		FilterType:   opts.FilterType,
		Raw:          opts.Raw,
		PingInterval: opts.PingInterval,
		Timeout:      opts.Timeout,
	}

	wsURL, err := buildWSURL(wsOpts.Backend, wsOpts.ConvID, wsOpts.Profile)
	if err != nil {
		return err
	}
	headers := map[string][]string{}
	if len(wsOpts.Cookies) > 0 {
		cookieHeader, err := cookieHeaderValue(wsOpts.Cookies)
		if err != nil {
			return err
		}
		headers["Cookie"] = []string{cookieHeader}
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
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

	if wsOpts.PingInterval > 0 {
		go pingLoop(ctx, conn, wsOpts.PingInterval)
	}

	readErr := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				readErr <- nil
				return
			default:
				_, data, err := conn.ReadMessage()
				if err != nil {
					readErr <- err
					return
				}
				if err := renderWSFrame(data, wsOpts); err != nil {
					readErr <- err
					return
				}
			}
		}
	}()

	chatOpts := &chatOptions{
		baseOptions: baseOptions{
			Backend: opts.Backend,
			ConvID:  convID,
			Profile: opts.Profile,
			Cookies: opts.Cookies,
			Pretty:  opts.Pretty,
			Timeout: 30 * time.Second,
		},
		Prompt:       opts.Prompt,
		ThinkingMode: opts.ThinkingMode,
		JSON:         false,
	}

	_, status, data, err := submitChat(chatOpts)
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("/chat failed: status=%d body=%s", status, strings.TrimSpace(string(data)))
	}

	if opts.TimelineDelay > 0 {
		time.Sleep(opts.TimelineDelay)
	}

	tlOpts := &timelineOptions{
		baseOptions: baseOptions{
			Backend: opts.Backend,
			ConvID:  convID,
			Cookies: opts.Cookies,
			Pretty:  opts.Pretty,
			Timeout: 15 * time.Second,
		},
		SinceVersion: opts.SinceVersion,
		Limit:        opts.Limit,
		JSON:         false,
	}

	snap, _, err := fetchTimeline(tlOpts)
	if err != nil {
		return err
	}
	if snap != nil {
		printTimelineSummary(os.Stdout, snap)
	}

	select {
	case err := <-readErr:
		return err
	default:
		return nil
	}
}
