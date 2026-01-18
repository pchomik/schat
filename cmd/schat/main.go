package main

import (
	"flag"
	"fmt"
	"os"
	"schat/internal/schat"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of schat:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "schat [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nSupported providers:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  opencode-cli\tDefault provider for OpenCode AI\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  cursor\t\tProvider for Cursor AI\n")
	}

	provider := flag.String("provider", "opencode-cli", "AI provider to use (options: opencode-cli, cursor)")
	flag.Parse()

	flag.Usage = nil

	validProviders := map[string]bool{
		"opencode-cli": true,
		"cursor":       true,
	}

	if !validProviders[*provider] {
		fmt.Fprintf(os.Stderr, "Invalid provider '%s'. Valid options are: opencode-cli, cursor\n", *provider)
		os.Exit(1)
	}

	schat.Run(*provider)
}
