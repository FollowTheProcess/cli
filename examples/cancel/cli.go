package main

import (
	"context"
	"fmt"
	"time"

	"go.followtheprocess.codes/cli"
)

func BuildCLI() (*cli.Command, error) {
	return cli.New(
		"cancel",
		cli.Short("Cancel me!"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second):
					fmt.Fprintln(cmd.Stdout(), "working...")
				}
			}
		}),
	)
}
