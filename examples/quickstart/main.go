package main

import (
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
		cli.Allow(cli.MinArgs(1)), // Must have at least one argument
		cli.Stdout(os.Stdout),
		cli.Example("Do a thing", "quickstart something"),
		cli.Example("Count the things", "quickstart something --count 3"),
		cli.Flag(&count, "count", 'c', 0, "Count the things"),
		cli.Run(runQuickstart(&count)),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func runQuickstart(count *int) func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		fmt.Fprintf(cmd.Stdout(), "Hello from quickstart!, my args were: %v, count was %d\n", args, *count)
		return nil
	}
}
