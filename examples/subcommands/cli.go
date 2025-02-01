package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/FollowTheProcess/cli"
)

func BuildCLI() (*cli.Command, error) {
	return cli.New(
		"demo",
		cli.Short("An example CLI to demonstrate the library and play with it for real."),
		cli.Version("dev"),
		cli.Commit("5ecb7caefacb12a75db4deebaa3cafcfd9d5c7c2"),
		cli.BuildDate("2024-08-17T10:37:30Z"),
		cli.Example("A basic subcommand", "demo say hello world"),
		cli.Example("Can do things", "demo do something --count 3"),
		cli.Allow(cli.NoArgs()),
		cli.SubCommands(buildSayCommand, buildDoCommand),
	)
}

type sayOptions struct {
	thing string
	count int
	shout bool
}

func buildSayCommand() (*cli.Command, error) {
	var options sayOptions
	return cli.New(
		"say",
		cli.Short("Print a message"),
		cli.Example("Say a well known phrase", "demo say hello world"),
		cli.Example("Now louder", "demo say hello world --shout"),
		cli.Run(runSay(&options)),
		cli.Flag(&options.shout, "shout", 's', false, "Say the message louder"),
		cli.Flag(&options.count, "count", 'c', 0, "Count the things"),
		cli.Flag(&options.thing, "thing", 't', "", "Name of the thing"),
	)
}

type doOptions struct {
	count     int
	fast      bool
	verbosity cli.FlagCount
	duration  time.Duration
}

func buildDoCommand() (*cli.Command, error) {
	var options doOptions
	return cli.New(
		"do",
		cli.Short("Do a thing"),
		cli.Example("Do something", "demo do something --fast"),
		cli.Example("Do it 3 times", "demo do something --count 3"),
		cli.Example("Do it for a specific duration", "demo do something --duration 1m30s"),
		cli.Allow(cli.ExactArgs(1)), // Only allowed to do one thing
		cli.Flag(&options.count, "count", 'c', 1, "Number of times to do the thing"),
		cli.Flag(&options.fast, "fast", 'f', false, "Do the thing quickly"),
		cli.Flag(&options.verbosity, "verbosity", 'v', 0, "Increase the verbosity level"),
		cli.Flag(&options.duration, "duration", 'd', 1*time.Second, "Do the thing for a specific duration"),
		cli.Run(runDo(&options)),
	)
}

func runSay(options *sayOptions) func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		if options.shout {
			for _, arg := range args {
				fmt.Fprintln(cmd.Stdout(), strings.ToUpper(arg), " ")
			}
		} else {
			for _, arg := range args {
				fmt.Fprintln(cmd.Stdout(), arg, " ")
			}
		}
		fmt.Printf("Shout: %v\nCount: %v\nThing: %v\n", options.shout, options.count, options.thing)
		fmt.Fprintln(cmd.Stdout())
		return nil
	}
}

func runDo(options *doOptions) func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		if options.fast {
			fmt.Fprintf(
				cmd.Stdout(),
				"Doing %s %d times, but faster! (will still take %v)\n",
				args[0],
				options.count,
				options.duration,
			)
		} else {
			fmt.Fprintf(cmd.Stdout(), "Doing %s %d times for %v\n", args[0], options.count, options.duration)
		}

		fmt.Fprintf(cmd.Stdout(), "Verbosity level was %d\n", options.verbosity)

		return nil
	}
}
