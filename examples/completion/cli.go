package main

import (
	"context"
	"fmt"
	"strings"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/cli/flag"
)

// BuildCLI constructs the root command with two subcommands and a completion
// subcommand wired in via [cli.CompletionSubCommand].
//
// Running "mytool completion" outputs a carapace-spec YAML document that
// describes the full command tree. Redirect it once to register completions
// with carapace-bin:
//
//	mytool completion > ~/.config/carapace/specs/mytool.yaml
//
// carapace-bin then provides completions across bash, zsh, fish, nushell,
// PowerShell, and more — no shell-specific scripts required.
func BuildCLI() (*cli.Command, error) {
	return cli.New(
		"mytool",
		cli.Short("An example tool with shell completion support"),
		cli.Version("dev"),
		cli.Example("Say hello", "mytool hello --name Alice"),
		cli.Example("Say goodbye loudly", "mytool goodbye --name Bob --shout"),
		cli.Example("Register completions (run once)", "mytool completion > ~/.config/carapace/specs/mytool.yaml"),
		cli.SubCommands(
			buildHelloCommand,
			buildGoodbyeCommand,
			cli.CompletionSubCommand(), // wire in completion support
		),
	)
}

type helloOptions struct {
	name      string
	verbosity flag.Count
}

func buildHelloCommand() (*cli.Command, error) {
	var opts helloOptions

	return cli.New(
		"hello",
		cli.Short("Greet someone"),
		cli.Example("Greet a person", "mytool hello --name Alice"),
		cli.Example("Greet verbosely", "mytool hello --name Alice -vvv"),
		cli.Flag(&opts.name, "name", 'n', "Name of the person to greet", cli.FlagDefault("World")),
		cli.Flag(&opts.verbosity, "verbose", 'v', "Increase output verbosity"),
		cli.Run(func(_ context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Stdout(), "Hello, %s!\n", opts.name)

			if opts.verbosity > 0 {
				fmt.Fprintf(cmd.Stdout(), "(verbosity level: %d)\n", opts.verbosity)
			}

			return nil
		}),
	)
}

type goodbyeOptions struct {
	name  string
	shout bool
}

func buildGoodbyeCommand() (*cli.Command, error) {
	var opts goodbyeOptions

	return cli.New(
		"goodbye",
		cli.Short("Bid someone farewell"),
		cli.Example("Say goodbye", "mytool goodbye --name Bob"),
		cli.Example("Say it louder", "mytool goodbye --name Bob --shout"),
		cli.Flag(&opts.name, "name", 'n', "Name of the person to farewell", cli.FlagDefault("World")),
		cli.Flag(&opts.shout, "shout", 's', "Say the farewell in uppercase"),
		cli.Run(func(_ context.Context, cmd *cli.Command) error {
			msg := fmt.Sprintf("Goodbye, %s!", opts.name)
			if opts.shout {
				msg = strings.ToUpper(msg)
			}

			fmt.Fprintln(cmd.Stdout(), msg)

			return nil
		}),
	)
}
