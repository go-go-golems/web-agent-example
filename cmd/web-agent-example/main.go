package main

import (
	"context"
	"embed"
	"io"
	"net/http"
	"strings"

	clay "github.com/go-go-golems/clay/pkg"
	geppettolayers "github.com/go-go-golems/geppetto/pkg/layers"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	geppettomw "github.com/go-go-golems/geppetto/pkg/inference/middleware"
	rediscfg "github.com/go-go-golems/pinocchio/pkg/redisstream"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"github.com/rs/zerolog/log"

	"github.com/go-go-golems/web-agent-example/pkg/thinkingmode"
)

//go:embed static
var staticFS embed.FS

type Command struct {
	*cmds.CommandDescription
}

func NewCommand() (*Command, error) {
	geLayers, err := geppettolayers.CreateGeppettoLayers()
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
			parameters.NewParameterDefinition("addr", parameters.ParameterTypeString, parameters.WithDefault(":8080"), parameters.WithHelp("HTTP listen address")),
			parameters.NewParameterDefinition("idle-timeout-seconds", parameters.ParameterTypeInteger, parameters.WithDefault(60), parameters.WithHelp("Stop per-conversation reader after N seconds with no sockets (0=disabled)")),
			parameters.NewParameterDefinition("evict-idle-seconds", parameters.ParameterTypeInteger, parameters.WithDefault(300), parameters.WithHelp("Evict conversations after N seconds idle (0=disabled)")),
			parameters.NewParameterDefinition("evict-interval-seconds", parameters.ParameterTypeInteger, parameters.WithDefault(60), parameters.WithHelp("Sweep idle conversations every N seconds (0=disabled)")),
			parameters.NewParameterDefinition("root", parameters.ParameterTypeString, parameters.WithDefault("/"), parameters.WithHelp("Serve the chat UI under a given URL root (e.g., /chat)")),
			parameters.NewParameterDefinition("timeline-dsn", parameters.ParameterTypeString, parameters.WithDefault(""), parameters.WithHelp("SQLite DSN for durable timeline snapshots (enables GET /timeline); preferred over timeline-db")),
			parameters.NewParameterDefinition("timeline-db", parameters.ParameterTypeString, parameters.WithDefault(""), parameters.WithHelp("SQLite DB file path for durable timeline snapshots (enables GET /timeline); DSN is derived with WAL/busy_timeout")),
		),
		cmds.WithLayersList(append(geLayers, redisLayer)...),
	)
	return &Command{CommandDescription: desc}, nil
}

func (c *Command) RunIntoWriter(ctx context.Context, parsed *layers.ParsedLayers, _ io.Writer) error {
	r, err := webchat.NewRouter(ctx, parsed, staticFS,
		webchat.WithEngineFromReqBuilder(newNoCookieEngineFromReqBuilder()),
		webchat.WithEventSinkWrapper(discoSinkWrapper()),
	)
	if err != nil {
		return errors.Wrap(err, "new webchat router")
	}

	r.RegisterMiddleware("webagent-thinking-mode", func(cfg any) geppettomw.Middleware {
		return thinkingmode.NewMiddleware(thinkingmode.ConfigFromAny(cfg))
	})

	r.AddProfile(&webchat.Profile{
		Slug:           "default",
		DefaultPrompt:  "You are a helpful assistant.",
		DefaultMws:     []webchat.MiddlewareUse{{Name: "webagent-thinking-mode", Config: thinkingmode.DefaultConfig()}},
		AllowOverrides: true,
	})

	httpSrv, err := r.BuildHTTPServer()
	if err != nil {
		return errors.Wrap(err, "build http server")
	}

	type serverSettings struct {
		Root string `glazed.parameter:"root"`
	}
	s := &serverSettings{}
	_ = parsed.InitializeStruct(layers.DefaultSlug, s)
	if s.Root != "" && s.Root != "/" {
		parent := http.NewServeMux()
		prefix := s.Root
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		if !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		parent.Handle(prefix, http.StripPrefix(strings.TrimRight(prefix, "/"), r.Handler()))
		httpSrv.Handler = parent
		log.Info().Str("root", prefix).Msg("mounted webchat under custom root")
	}

	srv := webchat.NewFromRouter(ctx, r, httpSrv)
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
	command, err := cli.BuildCobraCommand(c, cli.WithCobraMiddlewaresFunc(geppettolayers.GetCobraCommandGeppettoMiddlewares))
	cobra.CheckErr(err)
	root.AddCommand(command)
	cobra.CheckErr(root.Execute())
}
