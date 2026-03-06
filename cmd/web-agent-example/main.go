package main

import (
	"context"
	"embed"
	"io"
	"net/http"
	"strings"

	clay "github.com/go-go-golems/clay/pkg"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	geppettosections "github.com/go-go-golems/geppetto/pkg/sections"
	rediscfg "github.com/go-go-golems/pinocchio/pkg/redisstream"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	webhttp "github.com/go-go-golems/pinocchio/pkg/webchat/http"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

//go:embed static
var staticFS embed.FS

type Command struct {
	*cmds.CommandDescription
}

func NewCommand() (*Command, error) {
	geLayers, err := geppettosections.CreateGeppettoSections()
	if err != nil {
		return nil, errors.Wrap(err, "create geppetto layers")
	}
	redisLayer, err := rediscfg.NewParameterLayer()
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
			fields.New("timeline-dsn", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DSN for durable timeline snapshots (enables GET /api/timeline); preferred over timeline-db")),
			fields.New("timeline-db", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DB file path for durable timeline snapshots (enables GET /api/timeline); DSN is derived with WAL/busy_timeout")),
			fields.New("turns-dsn", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DSN for durable turn snapshots (enables GET /api/debug/turns); preferred over turns-db")),
			fields.New("turns-db", fields.TypeString, fields.WithDefault(""), fields.WithHelp("SQLite DB file path for durable turn snapshots (enables GET /api/debug/turns); DSN is derived with WAL/busy_timeout")),
		),
		cmds.WithSections(append(geLayers, redisLayer)...),
	)
	return &Command{CommandDescription: desc}, nil
}

func (c *Command) RunIntoWriter(ctx context.Context, parsed *values.Values, _ io.Writer) error {
	srv, err := webchat.NewServer(ctx, parsed, staticFS,
		webchat.WithRuntimeComposer(newWebAgentRuntimeComposer(parsed)),
		webchat.WithEventSinkWrapper(discoSinkWrapper()),
	)
	if err != nil {
		return errors.Wrap(err, "new webchat server")
	}

	resolver := newNoCookieRequestResolver()
	chatHandler := webhttp.NewChatHandler(srv.ChatService(), resolver)
	wsHandler := webhttp.NewWSHandler(
		srv.StreamHub(),
		resolver,
		websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }},
	)
	timelineHandler := webhttp.NewTimelineHandler(
		srv.TimelineService(),
		log.With().Str("component", "web-agent-example").Str("route", "/api/timeline").Logger(),
	)

	type serverSettings struct {
		Root string `glazed:"root"`
	}
	s := &serverSettings{}
	if err := parsed.DecodeSectionInto(values.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "decode server settings")
	}

	appMux := http.NewServeMux()
	appMux.HandleFunc("/chat", chatHandler)
	appMux.HandleFunc("/chat/", chatHandler)
	appMux.HandleFunc("/ws", wsHandler)
	appMux.HandleFunc("/api/timeline", timelineHandler)
	appMux.HandleFunc("/api/timeline/", timelineHandler)
	appMux.Handle("/api/", srv.APIHandler())
	appMux.Handle("/", srv.UIHandler())

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
	command, err := cli.BuildCobraCommand(c, cli.WithCobraMiddlewaresFunc(geppettosections.GetCobraCommandGeppettoMiddlewares))
	cobra.CheckErr(err)
	root.AddCommand(command)
	cobra.CheckErr(root.Execute())
}
