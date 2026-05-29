//glazedclilint:file-ignore legacy debug Cobra command uses raw flags; migrate to Glazed fields in a follow-up
package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//go:embed static
var staticFS embed.FS

type serveOptions struct {
	Addr    string
	Backend string
}

func runServe(args []string) error {
	opts := &serveOptions{}
	flagSet := flag.NewFlagSet("serve", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.StringVar(&opts.Addr, "addr", ":8090", "HTTP listen address for debug UI")
	flagSet.StringVar(&opts.Backend, "backend", "http://localhost:8080", "Backend origin for /debug, /ws, /chat, /timeline, /turns, /hydrate, /api proxying")
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	target, err := parseBackendURL(opts.Backend)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Proxy backend endpoints so the UI can be served from a dedicated port.
	for _, p := range []string{
		"/debug",
		"/debug/",
		"/ws",
		"/ws/",
		"/chat",
		"/chat/",
		"/timeline",
		"/timeline/",
		"/turns",
		"/turns/",
		"/hydrate",
		"/hydrate/",
		"/api",
		"/api/",
	} {
		mux.Handle(p, proxy)
	}

	// Expose built asset files directly when present.
	if assetsFS, err := fs.Sub(staticFS, "static/dist/assets"); err == nil {
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	}

	mux.Handle("/", newSPAHandler())

	srv := &http.Server{
		Addr:              opts.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "serving web-agent-debug UI on %s (proxy backend: %s)\n", opts.Addr, target.String())
	return srv.ListenAndServe()
}

func parseBackendURL(raw string) (*url.URL, error) {
	backend := normalizeBackend(raw)
	u, err := url.Parse(backend)
	if err != nil {
		return nil, fmt.Errorf("invalid --backend URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported --backend scheme %q (expected http or https)", u.Scheme)
	}
	if strings.TrimSpace(u.Host) == "" {
		return nil, errors.New("invalid --backend URL: missing host")
	}
	return u, nil
}

func newSPAHandler() http.Handler {
	distFS, distErr := fs.Sub(staticFS, "static/dist")
	var distFileServer http.Handler
	if distErr == nil {
		distFileServer = http.FileServer(http.FS(distFS))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean("/" + r.URL.Path)
		trimmed := strings.TrimPrefix(cleanPath, "/")
		if trimmed == "." {
			trimmed = ""
		}

		if distErr == nil && trimmed != "" && fileExists(distFS, trimmed) {
			rr := *r
			rr.URL = copyURL(r.URL)
			rr.URL.Path = "/" + trimmed
			distFileServer.ServeHTTP(w, &rr)
			return
		}

		if distErr == nil && fileExists(distFS, "index.html") {
			if err := serveFileFromFS(w, distFS, "index.html"); err == nil {
				return
			}
		}

		if fileExists(staticFS, "static/index.html") {
			if err := serveFileFromFS(w, staticFS, "static/index.html"); err == nil {
				return
			}
		}

		http.Error(w, "ui build not found; run `go generate ./cmd/web-agent-debug`", http.StatusInternalServerError)
	})
}

func copyURL(u *url.URL) *url.URL {
	if u == nil {
		return &url.URL{}
	}
	u2 := *u
	return &u2
}

func fileExists(fsys fs.FS, name string) bool {
	if fsys == nil {
		return false
	}
	info, err := fs.Stat(fsys, filepath.Clean(name))
	return err == nil && !info.IsDir()
}

func serveFileFromFS(w http.ResponseWriter, fsys fs.FS, name string) error {
	data, err := fs.ReadFile(fsys, filepath.Clean(name))
	if err != nil {
		return err
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ctype := mime.TypeByExtension(ext); ctype != "" {
		w.Header().Set("Content-Type", ctype)
	} else if ext == ".html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	_, err = w.Write(data)
	return err
}
