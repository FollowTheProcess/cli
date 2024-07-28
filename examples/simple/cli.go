package main

import (
	"fmt"
	"strings"

	"github.com/FollowTheProcess/cli"
)

func BuildCLI() (*cli.Command, error) {
	say, err := buildSayCommand()
	if err != nil {
		return nil, err
	}

	demo, err := cli.New(
		"demo",
		cli.Short("An example CLI to demonstrate the library and play with it for real."),
		cli.Example("A basic subcommand", "demo say hello world"),
		cli.Allow(cli.NoArgs()),
		cli.SubCommands(say),
	)
	if err != nil {
		return nil, err
	}

	return demo, nil
}

type sayOptions struct {
	thing string
	count int
	shout bool
}

func buildSayCommand() (*cli.Command, error) {
	var options sayOptions
	say, err := cli.New(
		"say",
		cli.Short("Print a message"),
		cli.Example("Say a well known phrase", "demo say hello world"),
		cli.Example("Now louder", "demo say hello world --shout"),
		cli.Run(runSay(&options)),
		cli.Flag(&options.shout, "shout", "s", false, "Say the message louder"),
		cli.Flag(&options.count, "count", "c", 0, "Count the things"),
		cli.Flag(&options.thing, "thing", "t", "", "Name of the thing"),
	)
	if err != nil {
		return nil, err
	}

	return say, nil
}

func runSay(options *sayOptions) func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		if options.shout {
			for _, arg := range args {
				fmt.Fprint(cmd.Stdout(), strings.ToUpper(arg), " ")
			}
		} else {
			for _, arg := range args {
				fmt.Fprint(cmd.Stdout(), arg, " ")
			}
		}
		fmt.Printf("Shout: %v\nCount: %v\nThing: %v\n", options.shout, options.count, options.thing)
		fmt.Fprintln(cmd.Stdout())
		return nil
	}
}
