package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type chatOptions struct {
	baseOptions
	Prompt       string
	ThinkingMode string
	JSON         bool
}

func runChat(args []string) error {
	opts := &chatOptions{}
	fs := flag.NewFlagSet("chat", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&opts.Backend, "backend", "http://localhost:8080", "Backend origin (http://host:port)")
	fs.StringVar(&opts.ConvID, "conv-id", "", "Conversation ID (generated if empty)")
	fs.StringVar(&opts.Profile, "profile", "default", "Profile slug (default if empty)")
	fs.StringVar(&opts.Prompt, "prompt", "", "Prompt to send")
	fs.StringVar(&opts.ThinkingMode, "thinking-mode", "", "Thinking mode override (fast/slow/etc)")
	fs.BoolVar(&opts.JSON, "json", false, "Print raw JSON response")
	fs.BoolVar(&opts.Pretty, "pretty", true, "Pretty-print JSON responses")
	fs.DurationVar(&opts.Timeout, "timeout", 30*time.Second, "HTTP timeout")
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

	_, status, data, err := submitChat(opts)
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("/chat failed: status=%d body=%s", status, strings.TrimSpace(string(data)))
	}

	if opts.JSON {
		fmt.Fprintln(os.Stdout, string(data))
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

func submitChat(opts *chatOptions) (string, int, []byte, error) {
	if opts == nil {
		return "", 0, nil, errors.New("missing chat options")
	}
	convID := ensureConvID(opts.ConvID)
	backend := normalizeBackend(opts.Backend)
	path := chatPath(opts.Profile)

	payload := map[string]any{
		"conv_id": convID,
		"prompt":  opts.Prompt,
	}
	if overrides := buildOverrides(opts.ThinkingMode); overrides != nil {
		payload["overrides"] = overrides
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return convID, 0, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, backend+path, bytes.NewReader(body))
	if err != nil {
		return convID, 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if opts.Profile != "" && opts.Profile != "default" && !hasCookie(opts.Cookies, "chat_profile") {
		req.AddCookie(&http.Cookie{Name: "chat_profile", Value: opts.Profile})
	}
	if err := addCookies(req, opts.Cookies); err != nil {
		return convID, 0, nil, err
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return convID, 0, nil, err
	}
	data, err := readResponseBody(resp)
	if err != nil {
		return convID, 0, nil, err
	}
	return convID, resp.StatusCode, data, nil
}
