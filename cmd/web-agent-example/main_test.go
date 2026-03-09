package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRewritePathPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		from    string
		to      string
		path    string
		rawPath string
		query   string
		want    string
	}{
		{
			name:  "timeline root",
			from:  "/timeline",
			to:    "/api/timeline",
			path:  "/timeline",
			query: "conv_id=conv-1",
			want:  "/api/timeline?conv_id=conv-1",
		},
		{
			name:  "timeline subtree",
			from:  "/timeline/",
			to:    "/api/timeline/",
			path:  "/timeline/history",
			query: "limit=2",
			want:  "/api/timeline/history?limit=2",
		},
		{
			name:  "turns legacy alias",
			from:  "/turns",
			to:    "/api/debug/turns",
			path:  "/turns",
			query: "conv_id=conv-1",
			want:  "/api/debug/turns?conv_id=conv-1",
		},
		{
			name:    "raw path preserved",
			from:    "/debug/",
			to:      "/api/debug/",
			path:    "/debug/conversations",
			rawPath: "/debug/conversations",
			query:   "limit=5",
			want:    "/api/debug/conversations?limit=5",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotPath, gotRawPath, gotQuery string
			next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotRawPath = r.URL.RawPath
				gotQuery = r.URL.RawQuery
			})

			req := httptest.NewRequest(http.MethodGet, "http://example.com"+tt.path, nil)
			req.URL.RawQuery = tt.query
			req.URL.RawPath = tt.rawPath
			rr := httptest.NewRecorder()

			rewritePathPrefix(tt.from, tt.to, next).ServeHTTP(rr, req)

			got := gotPath
			if gotQuery != "" {
				got += "?" + gotQuery
			}
			if got != tt.want {
				t.Fatalf("rewritten path mismatch: got %q want %q", got, tt.want)
			}
			if tt.rawPath != "" && gotRawPath == "" {
				t.Fatalf("expected raw path to be preserved")
			}
		})
	}
}
