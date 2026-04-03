// Package completion demonstrates how to add shell completion support to a CLI
// using [cli.CompletionSubCommand].
package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cmd, err := BuildCLI()
	if err != nil {
		return fmt.Errorf("could not build root command: %w", err)
	}

	if err := cmd.Execute(context.Background()); err != nil {
		return fmt.Errorf("could not execute root command: %w", err)
	}

	return nil
}
