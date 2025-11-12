// Package cover is a simple CLI demonstrating the core features of this library.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.followtheprocess.codes/cli"
)

func main() {
	var count int

	cmd, err := cli.New(
		"demo",
		cli.Short("Short description of your command"),
		cli.Long("Much longer text..."),
		cli.Version("v1.2.3"),
		cli.Stdout(os.Stdout),
		cli.Example("Do a thing", "demo thing --count"),
		cli.Flag(&count, "count", 'c', 0, "Count the thing"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from demo, my arguments were: ", cmd.Args())
			return nil
		}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Good command line tools allow for timeouts, cancellations etc.
	// so in cli, you pass a context.Context to your root command, and
	// it gets passed down to your Run function.
	if err := cmd.Execute(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
