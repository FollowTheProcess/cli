// Package quickstart demonstrates a quickstart example for cli.
package main

import (
	"context"
	"fmt"
	"os"

	"go.followtheprocess.codes/cli"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var count int

	cmd, err := cli.New(
		"quickstart",
		cli.Short("Short description of your command"),
		cli.Long("Much longer text..."),
		cli.Version("v1.2.3"),
		cli.Commit("7bcac896d5ab67edc5b58632c821ec67251da3b8"),
		cli.BuildDate("2024-08-17T10:37:30Z"),
		cli.Stdout(os.Stdout),
		cli.Example("Do a thing", "quickstart something"),
		cli.Example("Count the things", "quickstart something --count 3"),
		cli.Flag(&count, "count", 'c', "Count the things"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Stdout(), "Hello from quickstart!, my args were: %v, count was %d\n", cmd.Args(), count)
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute(context.Background())
}
