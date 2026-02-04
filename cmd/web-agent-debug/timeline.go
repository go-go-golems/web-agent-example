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
	ConvID       string          `json:"convId"`
	Version      json.RawMessage `json:"version"`
	ServerTimeMs json.RawMessage `json:"serverTimeMs"`
	Entities     []struct {
		EntityID string `json:"entityId"`
		Kind     string `json:"kind"`
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

	snap, data, err := fetchTimeline(opts)
	if err != nil {
		return err
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
	if snap == nil {
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}
	printTimelineSummary(os.Stdout, snap)
	return nil
}

func fetchTimeline(opts *timelineOptions) (*timelineSnapshot, []byte, error) {
	if opts == nil {
		return nil, nil, errors.New("missing timeline options")
	}
	backend := normalizeBackend(opts.Backend)
	u, err := url.Parse(backend + "/timeline")
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}
	if err := addCookies(req, opts.Cookies); err != nil {
		return nil, nil, err
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	data, err := readResponseBody(resp)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, data, fmt.Errorf("/timeline failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var snap timelineSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, data, nil
	}
	return &snap, data, nil
}

func printTimelineSummary(out *os.File, snap *timelineSnapshot) {
	if snap == nil || out == nil {
		return
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

	fmt.Fprintf(out, "conv_id: %s\n", snap.ConvID)
	fmt.Fprintf(out, "version: %s\n", formatJSONScalar(snap.Version))
	fmt.Fprintf(out, "server_time_ms: %s\n", formatJSONScalar(snap.ServerTimeMs))
	fmt.Fprintf(out, "entities: %d\n", len(snap.Entities))
	for _, k := range kinds {
		fmt.Fprintf(out, "  %s: %d\n", k, counts[k])
	}
}

func formatJSONScalar(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	s := strings.TrimSpace(string(raw))
	return strings.Trim(s, "\"")
}
