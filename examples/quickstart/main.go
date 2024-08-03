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
	cmd, err := cli.New(
		"quickstart",
		cli.Allow(cli.AnyArgs()),
		cli.Short("quick demo CLI to show the library"),
		cli.Run(runQuickstart),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func runQuickstart(cmd *cli.Command, args []string) error {
	fmt.Fprintln(cmd.Stdout(), "Hello from quickstart!")
	return nil
}
