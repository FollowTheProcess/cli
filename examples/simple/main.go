package main

import (
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

	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("could not execute root command: %w", err)
	}
	return nil
}
