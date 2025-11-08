// Package namedargs demonstrates how to use named positional arguments in cli.
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

type runArgs struct {
	src   string
	dest  string
	count int
}

func run() error {
	var arguments runArgs

	cmd, err := cli.New(
		"copy", // A fictional copy command
		cli.Short("Copy a file from a src to a destination"),
		cli.Stdout(os.Stdout),
		cli.Arg(&arguments.src, "src", "The file to copy from"),
		cli.Arg(&arguments.dest, "dest", "The file to copy to"),
		cli.Arg(&arguments.count, "count", "The number of things"),
		cli.Example("Copy a file to somewhere", "copy src.txt ./some/where/else"),
		cli.Example("Use the default destination", "copy src.txt"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(cmd.Stdout(), "Copying from %s to %s %d times\n", arguments.src, arguments.dest, arguments.count)
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
