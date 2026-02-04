package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type baseOptions struct {
	Backend string
	ConvID  string
	Profile string
	Cookies []string
	Pretty  bool
	Timeout time.Duration
}

func normalizeBackend(backend string) string {
	b := strings.TrimSpace(backend)
	if b == "" {
		return "http://localhost:8080"
	}
	return strings.TrimRight(b, "/")
}

func ensureConvID(id string) string {
	id = strings.TrimSpace(id)
	if id != "" {
		return id
	}
	return uuid.NewString()
}

func chatPath(profile string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" || profile == "default" {
		return "/chat"
	}
	return "/chat/" + profile
}

func buildOverrides(thinkingMode string) map[string]any {
	thinkingMode = strings.TrimSpace(thinkingMode)
	if thinkingMode == "" {
		return nil
	}
	return map[string]any{
		"middlewares": []map[string]any{{
			"name": "webagent-thinking-mode",
			"config": map[string]any{
				"mode": thinkingMode,
			},
		}},
	}
}

func addCookies(req *http.Request, cookies []string) error {
	for _, raw := range cookies {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid cookie %q (expected key=value)", raw)
		}
		req.AddCookie(&http.Cookie{Name: parts[0], Value: parts[1]})
	}
	return nil
}

func hasCookie(cookies []string, name string) bool {
	for _, raw := range cookies {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) == 2 && parts[0] == name {
			return true
		}
	}
	return false
}

func prettyPrintJSON(raw []byte, out *os.File) error {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = out.Write(append(b, '\n'))
	return err
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, errors.New("missing response body")
	}
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}
