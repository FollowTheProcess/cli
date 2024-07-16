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
			fmt.Fprintln(cmd.Stdout())
			return nil
		}),
	)

	// In pflag
	say.Flags().BoolVarP(&shout, "shout", "s", false, "Say the message louder")
	say.Flags().IntVarP(&count, "count", "c", 0, "Count the things")
	say.Flags().StringVarP(&thing, "thing", "t", "", "The name of a thing")

	// I'd like a generic version where I could do this... maybe next project
	// flag.New(&shout, "shout", "s", false, "Say the message louder")
	// flag.New(&count, "count", "c", 0, "Count the things")
	// flag.New(&something, "something", "s", "word", "Something is a string")

	return say
}
