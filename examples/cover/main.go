package main

import (
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
		cli.Allow(cli.MinArgs(1)),
		cli.Stdout(os.Stdout),
		cli.Example("Do a thing", "demo thing --count"),
		cli.Flag(&count, "count", 'c', 0, "Count the thing"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(cmd.Stdout(), "Hello from demo")
			return nil
		}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
