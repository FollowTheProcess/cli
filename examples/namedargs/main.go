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

func run() error {
	cmd, err := cli.New(
		"copy", // A fictional copy command
		cli.Short("Copy a file from a src to a destination"),
		cli.RequiredArg(
			"src",
			"The file to copy from",
		), // src is required, failure to provide it will error
		cli.OptionalArg("dest", "The destination to copy to", "./dest"), // dest has a default if not provided
		cli.Stdout(os.Stdout),
		cli.Example("Copy a file to somewhere", "copy src.txt ./some/where/else"),
		cli.Example("Use the default destination", "copy src.txt"),
		cli.Run(runCopy()),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func runCopy() func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		// src is required so if not provided will be an error
		// is dest is provided cmd.Arg("dest") will retrieve the value
		// if it's not provided, cmd.Arg("dest") will return the default of "./dest"
		fmt.Fprintf(cmd.Stdout(), "Copying from %s to %s\n", cmd.Arg("src"), cmd.Arg("dest"))
		return nil
	}
}
