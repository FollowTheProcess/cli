// Package cancel demonstrates how to handle CTRL+C and cancellation/timeouts
// easily with cli.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cli, err := BuildCLI()
	if err != nil {
		return err
	}

	return cli.Execute(ctx)
}
