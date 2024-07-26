package main

import (
	"fmt"
	"strings"

	"github.com/FollowTheProcess/cli"
)

func BuildCLI() *cli.Command {
	demo := cli.New(
		"demo",
		cli.Short("An example CLI to demonstrate the library and play with it for real."),
		cli.Examples(
			cli.Example{
				Comment: "A basic subcommand",
				Command: "demo say hello world",
			},
		),
		cli.Allow(cli.NoArgs),
		cli.SubCommands(buildSayCommand()),
	)

	return demo
}

func buildSayCommand() *cli.Command {
	var (
		shout bool
		count int
		thing string
	)
	say := cli.New(
		"say",
		cli.Short("Print a message"),
		cli.Examples(
			cli.Example{
				Comment: "Say a well known phrase",
				Command: "demo say hello world",
			},
			cli.Example{
				Comment: "Now louder",
				Command: "demo say hello world --shout",
			},
		),
		cli.Run(func(cmd *cli.Command, args []string) error {
			if shout {
				for _, arg := range args {
					fmt.Fprint(cmd.Stdout(), strings.ToUpper(arg), " ")
				}
			} else {
				for _, arg := range args {
					fmt.Fprint(cmd.Stdout(), arg, " ")
				}
			}
			fmt.Printf("Shout: %v\nCount: %v\nThing: %v\n", shout, count, thing)
			fmt.Fprintln(cmd.Stdout())
			return nil
		}),
		cli.Flag(&shout, "shout", "s", false, "Say the message louder"),
		cli.Flag(&count, "count", "c", 0, "Count the things"),
		cli.Flag(&thing, "thing", "t", "", "Name of the thing"),
	)

	return say
}
