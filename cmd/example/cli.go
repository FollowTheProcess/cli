package main

import (
	"fmt"
	"strings"

	"github.com/FollowTheProcess/cli/command"
)

func BuildCLI() *command.Command {
	demo := command.New(
		"demo",
		command.Short("An example CLI to demonstrate the library and play with it for real."),
		command.Examples(
			command.Example{
				Comment: "A basic subcommand",
				Command: "demo say hello world",
			},
		),
		command.Allow(command.NoArgs),
		command.SubCommands(buildSayCommand()),
	)

	return demo
}

func buildSayCommand() *command.Command {
	var (
		shout bool
		count int
		thing string
	)
	say := command.New(
		"say",
		command.Short("Print a message"),
		command.Examples(
			command.Example{
				Comment: "Say a well known phrase",
				Command: "demo say hello world",
			},
			command.Example{
				Comment: "Now louder",
				Command: "demo say hello world --shout",
			},
		),
		command.Run(func(cmd *command.Command, args []string) error {
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

	// With my version (not hooked up to can't actually be parsed yet)
	// flag.New(&shout, "shout", "s", false, "Say the message louder")
	// flag.New(&count, "count", "c", 0, "Count the things")
	// flag.New(&something, "something", "s", "word", "Something is a string")

	return say
}
