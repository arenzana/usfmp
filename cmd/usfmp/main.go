package main

import (
	"fmt"
	"os"

	"github.com/arenzana/usfmp/cmd/usfmp/cmd"
)

// Build information set by goreleaser
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Set version information for cobra command
	cmd.SetVersionInfo(version, commit, date, builtBy)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
