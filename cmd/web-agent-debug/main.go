package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "chat", "ws", "timeline", "run", "serve":
		if err := runCommand(os.Args[1], os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "web-agent-debug: CLI harness for webchat (/chat + /ws)\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  web-agent-debug <chat|ws|timeline|run|serve> [flags]\n\n")
}

func runCommand(cmd string, args []string) error {
	switch cmd {
	case "chat":
		return runChat(args)
	case "ws":
		return runWS(args)
	case "timeline":
		return runTimeline(args)
	case "run":
		return runRun(args)
	case "serve":
		return runServe(args)
	default:
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		if err := fs.Parse(args); err != nil {
			return err
		}
		return fmt.Errorf("command %s not implemented yet", cmd)
	}
}
