package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type timelineOptions struct {
	baseOptions
	SinceVersion uint64
	Limit        int
	JSON         bool
}

type timelineSnapshot struct {
	ConvID       string `json:"conv_id"`
	Version      uint64 `json:"version"`
	ServerTimeMs int64  `json:"server_time_ms"`
	Entities     []struct {
		EntityID    string `json:"entity_id"`
		Kind        string `json:"kind"`
		CreatedAtMs int64  `json:"created_at_ms"`
		UpdatedAtMs int64  `json:"updated_at_ms"`
	} `json:"entities"`
}

func runTimeline(args []string) error {
	opts := &timelineOptions{}
	fs := flag.NewFlagSet("timeline", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&opts.Backend, "backend", "http://localhost:8080", "Backend origin (http://host:port)")
	fs.StringVar(&opts.ConvID, "conv-id", "", "Conversation ID (required)")
	fs.Uint64Var(&opts.SinceVersion, "since-version", 0, "Fetch entities since version")
	fs.IntVar(&opts.Limit, "limit", 0, "Limit number of entities")
	fs.BoolVar(&opts.JSON, "json", false, "Print raw JSON response")
	fs.BoolVar(&opts.Pretty, "pretty", true, "Pretty-print JSON responses")
	fs.DurationVar(&opts.Timeout, "timeout", 15*time.Second, "HTTP timeout")
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

	backend := normalizeBackend(opts.Backend)
	u, err := url.Parse(backend + "/timeline")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("conv_id", opts.ConvID)
	if opts.SinceVersion > 0 {
		q.Set("since_version", fmt.Sprintf("%d", opts.SinceVersion))
	}
	if opts.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	if err := addCookies(req, opts.Cookies); err != nil {
		return err
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	data, err := readResponseBody(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("/timeline failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(data)))
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

	var snap timelineSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}
	counts := map[string]int{}
	for _, e := range snap.Entities {
		counts[e.Kind]++
	}
	kinds := make([]string, 0, len(counts))
	for k := range counts {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	fmt.Fprintf(os.Stdout, "conv_id: %s\n", snap.ConvID)
	fmt.Fprintf(os.Stdout, "version: %d\n", snap.Version)
	fmt.Fprintf(os.Stdout, "server_time_ms: %d\n", snap.ServerTimeMs)
	fmt.Fprintf(os.Stdout, "entities: %d\n", len(snap.Entities))
	for _, k := range kinds {
		fmt.Fprintf(os.Stdout, "  %s: %d\n", k, counts[k])
	}
	return nil
}
