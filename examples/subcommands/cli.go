package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/cli/flag"
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
		cli.SubCommands(buildSayCommand, buildDoCommand),
	)
}

type sayOptions struct {
	thing string
	items []string
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
		cli.Flag(&options.shout, "shout", 's', false, "Say the message louder"),
		cli.Flag(&options.count, "count", 'c', 0, "Count the things"),
		cli.Flag(&options.thing, "thing", 't', "", "Name of the thing"),
		cli.Flag(&options.items, "item", 'i', nil, "Items to add to a list"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			if options.shout {
				for _, arg := range cmd.Args() {
					fmt.Fprintln(cmd.Stdout(), strings.ToUpper(arg), " ")
				}
			} else {
				for _, arg := range cmd.Args() {
					fmt.Fprintln(cmd.Stdout(), arg, " ")
				}
			}

			fmt.Printf(
				"Shout: %v\nCount: %v\nThing: %v\nItems: %v\n",
				options.shout,
				options.count,
				options.thing,
				options.items,
			)
			fmt.Fprintln(cmd.Stdout())

			return nil
		}),
	)
}

type doOptions struct {
	count     int
	fast      bool
	verbosity flag.Count
	duration  time.Duration
}

func buildDoCommand() (*cli.Command, error) {
	var options doOptions

	var thing *url.URL

	return cli.New(
		"do",
		cli.Short("Do a thing"),
		cli.Example("Do something", "demo do something --fast"),
		cli.Example("Do it 3 times", "demo do something --count 3"),
		cli.Example("Do it for a specific duration", "demo do something --duration 1m30s"),
		cli.Version("do version"),
		cli.Arg(&thing, "thing", "Thing to do"),
		cli.Flag(&options.count, "count", 'c', 1, "Number of times to do the thing"),
		cli.Flag(&options.fast, "fast", 'f', false, "Do the thing quickly"),
		cli.Flag(&options.verbosity, "verbosity", 'v', 0, "Increase the verbosity level"),
		cli.Flag(&options.duration, "duration", 'd', 1*time.Second, "Do the thing for a specific duration"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			if options.fast {
				fmt.Fprintf(
					cmd.Stdout(),
					"Doing %s %d times, but faster! (will still take %v)\n",
					thing,
					options.count,
					options.duration,
				)
			} else {
				fmt.Fprintf(cmd.Stdout(), "Doing %s %d times for %v\n", thing, options.count, options.duration)
			}

			fmt.Fprintf(cmd.Stdout(), "Verbosity level was %d\n", options.verbosity)

			return nil
		}),
	)
}
