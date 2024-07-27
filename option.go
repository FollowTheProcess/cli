package cli

import (
	"io"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/spf13/pflag"
)

// Option is a functional option for configuring a [Command].
type Option func(*Command)

// Stdin is an [Option] that sets the Stdin for a [Command].
func Stdin(stdin io.Reader) Option {
	return func(cmd *Command) {
		cmd.stdin = stdin
	}
}

// Stdout is an [Option] that sets the Stdout for a [Command].
func Stdout(stdout io.Writer) Option {
	return func(cmd *Command) {
		cmd.stdout = stdout
	}
}

// Stderr is an [Option] that sets the Stderr for a [Command].
func Stderr(stderr io.Writer) Option {
	return func(cmd *Command) {
		cmd.stderr = stderr
	}
}

// Short is an [Option] that sets the one line usage summary for a [Command].
func Short(short string) Option {
	return func(cmd *Command) {
		cmd.short = short
	}
}

// Long is an [Option] that sets the full description for a [Command].
func Long(long string) Option {
	return func(cmd *Command) {
		cmd.long = long
	}
}

// Example is an [Option] that adds an example to a [Command].
//
// Examples take the form of an explanatory comment and a command
// showing the command to the CLI, these will show up in the help text.
//
// For example, a program called "myrm" that deletes files and directories
// might have an example declared as follows:
//
//	cli.Example("Delete a folder recursively without confirmation", "myrm ./dir --recursive --force")
//
// Which would show up in the help text like so:
//
//	Examples:
//	# Delete a folder recursively without confirmation
//	$ myrm ./dir --recursive --force
//
// An arbitrary number of examples can be added to a [Command], each call to Example
// will add another example.
func Example(comment, command string) Option {
	return func(cmd *Command) {
		cmd.examples = append(cmd.examples, example{comment: comment, command: command})
	}
}

// Run is an [Option] that sets the run function for a [Command].
//
// The run function is the actual implementation of your command i.e. what you
// want it to do when invoked.
func Run(run func(cmd *Command, args []string) error) Option {
	return func(cmd *Command) {
		cmd.run = run
	}
}

// Args is an [Option] that sets the arguments for a [Command].
func Args(args []string) Option {
	return func(cmd *Command) {
		cmd.args = args
	}
}

// Version is an [Option] that sets the version for a [Command].
func Version(version string) Option {
	return func(cmd *Command) {
		cmd.version = version
	}
}

// VersionFunc is an [Option] that allows for a custom implementation of the -v/--version flag.
//
// A [Command] will have a default implementation of this function that prints a default
// format of the version info to [os.Stderr].
func VersionFunc(fn func(cmd *Command) error) Option {
	return func(cmd *Command) {
		cmd.versionFunc = fn
	}
}

// SubCommands is an [Option] that attaches 1 or more subcommands to the command being configured.
func SubCommands(subcommands ...*Command) Option {
	return func(cmd *Command) {
		// Add the subcommands to the command this is being called on
		cmd.subcommands = append(cmd.subcommands, subcommands...)

		// Loop through each subcommand and set this command as their immediate parent
		for _, subcommand := range subcommands {
			subcommand.parent = cmd
		}
	}
}

// Allow is an [Option] that allows for validating positional arguments to a [Command].
//
// You provide a validator function that returns an error if it encounters invalid arguments, and it will
// be run for you, passing in the non-flag arguments to the [Command] that was called.
func Allow(validator func(cmd *Command, args []string) error) Option {
	return func(cmd *Command) {
		cmd.allowArgs = validator
	}
}

// Flag is an [Option] that adds a flag to a [Command].
func Flag[T flag.Flaggable](p *T, name string, short string, value T, usage string) Option {
	return func(cmd *Command) {
		flag := flag.New(p, name, short, value, usage)
		var defVal string
		if flag.Type() == "bool" {
			defVal = "true"
		}
		cmd.flagSet().AddFlag(&pflag.Flag{
			Name:        name,
			Shorthand:   short,
			Usage:       usage,
			Value:       flag,
			DefValue:    flag.String(),
			NoOptDefVal: defVal,
		})
	}
}
