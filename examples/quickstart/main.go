package main

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/cli"
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
