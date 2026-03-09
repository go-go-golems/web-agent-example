package main

import (
	"context"
	"embed"
	"io"
	"net/http"
	"strings"

	clay "github.com/go-go-golems/clay/pkg"
	geppettosections "github.com/go-go-golems/geppetto/pkg/sections"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	webhttp "github.com/go-go-golems/pinocchio/pkg/webchat/http"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rediscfg "github.com/go-go-golems/pinocchio/pkg/redisstream"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"github.com/rs/zerolog/log"
)

//go:embed static
var staticFS embed.FS

type Command struct {
	*cmds.CommandDescription
}

func NewCommand() (*Command, error) {
	geSections, err := geppettosections.CreateGeppettoSections()
	if err != nil {
		return nil, errors.Wrap(err, "create geppetto sections")
	}
	redisSection, err := rediscfg.NewParameterLayer()
	if err != nil {
		return nil, err
	}

	desc := cmds.NewCommandDescription(
		"serve",
		cmds.WithShort("Serve the web-agent-example webchat server"),
		cmds.WithFlags(
			fields.New("addr", fields.TypeString, fields.WithDefault(":8080"), fields.WithHelp("HTTP listen address")),
			fields.New("idle-timeout-seconds", fields.TypeInteger, fields.WithDefault(60), fields.WithHelp("Stop per-conversation reader after N seconds with no sockets (0=disabled)")),
			fields.New("evict-idle-seconds", fields.TypeInteger, fields.WithDefault(300), fields.WithHelp("Evict conversations after N seconds idle (0=disabled)")),
			fields.New("evict-interval-seconds", fields.TypeInteger, fields.WithDefault(60), fields.WithHelp("Sweep idle conversations every N seconds (0=disabled)")),
			fields.New("root", fields.TypeString, fields.WithDefault("/"), fields.WithHelp("Serve the chat UI under a given URL root (e.g., /chat)")),
			fields.New("timeline-dsn", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DSN for durable timeline snapshots (enables GET /timeline); preferred over timeline-db")),
			fields.New("timeline-db", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DB file path for durable timeline snapshots (enables GET /timeline); DSN is derived with WAL/busy_timeout")),
			fields.New("turns-dsn", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DSN for durable turn snapshots (enables GET /turns); preferred over turns-db")),
			fields.New("turns-db", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DB file path for durable turn snapshots (enables GET /turns); DSN is derived with WAL/busy_timeout")),
		),
		cmds.WithSections(append(geSections, redisSection)...),
	)
	return &Command{CommandDescription: desc}, nil
}

func (c *Command) RunIntoWriter(ctx context.Context, parsed *values.Values, _ io.Writer) error {
	type serverSettings struct {
		Root string `glazed:"root"`
	}
	s := &serverSettings{}
	if err := parsed.DecodeSectionInto(values.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "decode server settings")
	}

	runtimeComposer, err := newDefaultRuntimeComposer(parsed)
	if err != nil {
		return err
	}
	requestResolver := newStaticRequestResolver()

	srv, err := webchat.NewServer(
		ctx,
		parsed,
		staticFS,
		webchat.WithRuntimeComposer(runtimeComposer),
		webchat.WithEventSinkWrapper(discoSinkWrapper()),
	)
	if err != nil {
		return errors.Wrap(err, "new webchat server")
	}

	chatHandler := webhttp.NewChatHandler(srv.ChatService(), requestResolver)
	wsHandler := webhttp.NewWSHandler(
		srv.StreamHub(),
		requestResolver,
		websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }},
	)
	apiHandler := srv.APIHandler()
	uiHandler := srv.UIHandler()

	appMux := http.NewServeMux()
	appMux.HandleFunc("/chat", chatHandler)
	appMux.HandleFunc("/chat/", chatHandler)
	appMux.HandleFunc("/ws", wsHandler)
	appMux.Handle("/api/", apiHandler)
	registerLegacyAPIAliases(appMux, apiHandler)
	appMux.Handle("/", uiHandler)

	httpSrv := srv.HTTPServer()
	if httpSrv == nil {
		return errors.New("http server is not initialized")
	}
	if s.Root != "" && s.Root != "/" {
		parent := http.NewServeMux()
		prefix := s.Root
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		if !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		parent.Handle(prefix, http.StripPrefix(strings.TrimRight(prefix, "/"), appMux))
		httpSrv.Handler = parent
		log.Info().Str("root", prefix).Msg("mounted webchat under custom root")
	} else {
		httpSrv.Handler = appMux
	}

	return srv.Run(ctx)
}

func registerLegacyAPIAliases(mux *http.ServeMux, apiHandler http.Handler) {
	if mux == nil || apiHandler == nil {
		return
	}

	// Preserve legacy backend paths that local debug tooling still requests directly.
	mux.Handle("/debug", rewritePathPrefix("/debug", "/api/debug", apiHandler))
	mux.Handle("/debug/", rewritePathPrefix("/debug/", "/api/debug/", apiHandler))
	mux.Handle("/timeline", rewritePathPrefix("/timeline", "/api/timeline", apiHandler))
	mux.Handle("/timeline/", rewritePathPrefix("/timeline/", "/api/timeline/", apiHandler))
	mux.Handle("/turns", rewritePathPrefix("/turns", "/api/debug/turns", apiHandler))
	mux.Handle("/turns/", rewritePathPrefix("/turns/", "/api/debug/turns/", apiHandler))
}

func rewritePathPrefix(from, to string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next == nil || r == nil || r.URL == nil {
			http.NotFound(w, r)
			return
		}

		rr := r.Clone(r.Context())
		u := *r.URL
		rr.URL = &u
		rr.URL.Path = to + strings.TrimPrefix(r.URL.Path, from)
		if r.URL.RawPath != "" {
			rr.URL.RawPath = to + strings.TrimPrefix(r.URL.RawPath, from)
		}
		next.ServeHTTP(w, rr)
	})
}

func main() {
	root := &cobra.Command{Use: "web-agent-example", PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return logging.InitLoggerFromCobra(cmd)
	}}

	helpSystem := help.NewHelpSystem()
	help_cmd.SetupCobraRootCommand(helpSystem, root)

	if err := clay.InitGlazed("web-agent-example", root); err != nil {
		cobra.CheckErr(err)
	}

	c, err := NewCommand()
	cobra.CheckErr(err)
	command, err := cli.BuildCobraCommand(c, cli.WithParserConfig(cli.CobraParserConfig{
		AppName: "web-agent-example",
	}))
	cobra.CheckErr(err)
	root.AddCommand(command)
	cobra.CheckErr(root.Execute())
}
