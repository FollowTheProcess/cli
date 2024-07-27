package cli

import (
	"fmt"
	"io"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/spf13/pflag"
)

// TODO: We now have the ability to handle config errors in options so we should validate where possible

// Option is a functional option for configuring a [Command].
type Option interface {
	// Apply the option to the config
	apply(*config) error
}

// option is a function adapter implementing the Option interface, analogous
// to http.HandlerFunc.
type option func(*config) error

// apply implements the Option interface for option.
func (o option) apply(cfg *config) error {
	return o(cfg)
}

// config represents the internal configuration of a [Command].
type config struct {
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	run         func(cmd *Command, args []string) error
	flags       *pflag.FlagSet
	versionFunc func(cmd *Command) error
	parent      *Command
	allowArgs   func(cmd *Command, args []string) error
	name        string
	short       string
	long        string
	version     string
	examples    []example
	args        []string
	subcommands []*Command
}

// build builds an returns a Command from the config.
func (c *config) build() *Command {
	cmd := &Command{
		stdin:        c.stdin,
		stdout:       c.stdout,
		stderr:       c.stderr,
		run:          c.run,
		flags:        c.flags,
		versionFunc:  c.versionFunc,
		parent:       c.parent,
		argValidator: c.allowArgs,
		name:         c.name,
		short:        c.short,
		long:         c.long,
		version:      c.version,
		examples:     c.examples,
		args:         c.args,
		subcommands:  c.subcommands,
	}

	// Loop through each subcommand and set this command as their immediate parent
	for _, subcommand := range cmd.subcommands {
		subcommand.parent = cmd
	}

	// Add the help and version flags
	cmd.flagSet().BoolP("help", "h", false, fmt.Sprintf("Show help for %s", cmd.name))
	cmd.flagSet().BoolP("version", "v", false, fmt.Sprintf("Show version info for %s", cmd.name))

	return cmd
}

// Stdin is an [Option] that sets the Stdin for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stdin].
//
//	// Set stdin to os.Stdin (the default anyway)
//	cli.New("test", cli.Stdin(os.Stdin))
func Stdin(stdin io.Reader) Option {
	f := func(cfg *config) error {
		cfg.stdin = stdin
		return nil
	}
	return option(f)
}

// Stdout is an [Option] that sets the Stdout for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stdout].
//
//	// Set stdout to a temporary buffer
//	buf := &bytes.Buffer{}
//	cli.New("test", cli.Stdout(buf))
func Stdout(stdout io.Writer) Option {
	f := func(cfg *config) error {
		cfg.stdout = stdout
		return nil
	}
	return option(f)
}

// Stderr is an [Option] that sets the Stderr for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stderr].
//
//	// Set stderr to a temporary buffer
//	buf := &bytes.Buffer{}
//	cli.New("test", cli.Stderr(buf))
func Stderr(stderr io.Writer) Option {
	f := func(cfg *config) error {
		cfg.stderr = stderr
		return nil
	}
	return option(f)
}

// Short is an [Option] that sets the one line usage summary for a [Command].
//
// The one line usage will appear in the help text as well as alongside
// subcommands when they are listed.
//
// Successive calls will simply overwrite any previous calls.
//
//	cli.New("rm", cli.Short("Delete files and directories"))
func Short(short string) Option {
	f := func(cfg *config) error {
		cfg.short = short
		return nil
	}
	return option(f)
}

// Long is an [Option] that sets the full description for a [Command].
//
// The long description will appear in the help text for a command.
//
// Successive calls will simply overwrite any previous calls.
//
//	cli.New("rm", cli.Long("... lots of text here"))
func Long(long string) Option {
	f := func(cfg *config) error {
		cfg.long = long
		return nil
	}
	return option(f)
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
// An arbitrary number of examples can be added to a [Command], and calls to [Example] are additive.
func Example(comment, command string) Option {
	f := func(cfg *config) error {
		cfg.examples = append(cfg.examples, example{comment: comment, command: command})
		return nil
	}
	return option(f)
}

// Run is an [Option] that sets the run function for a [Command].
//
// The run function is the actual implementation of your command i.e. what you
// want it to do when invoked.
//
// Successive calls overwrite previous ones.
func Run(run func(cmd *Command, args []string) error) Option {
	f := func(cfg *config) error {
		cfg.run = run
		return nil
	}
	return option(f)
}

// Args is an [Option] that sets the arguments for a [Command].
//
// Without this option, the command will default to os.Args[1:], this option is particularly
// useful for testing.
//
// Successive calls override previous ones.
//
//	// Override arguments for testing
//	cli.New("test", cli.Args([]string{"test", "me"}))
func Args(args []string) Option {
	f := func(cfg *config) error {
		cfg.args = args
		return nil
	}
	return option(f)
}

// Version is an [Option] that sets the version for a [Command].
//
// Without this option, the command defaults to a version of "dev".
//
//	cli.New("test", cli.Version("v1.2.3"))
func Version(version string) Option {
	f := func(cfg *config) error {
		cfg.version = version
		return nil
	}
	return option(f)
}

// VersionFunc is an [Option] that allows for a custom implementation of the -v/--version flag.
//
// A [Command] will have a default implementation of this function that prints a default
// format of the version info to [os.Stderr].
//
// This option is particularly useful if you want to inject ldflags in at build time for
// e.g commit hash.
func VersionFunc(fn func(cmd *Command) error) Option {
	f := func(cfg *config) error {
		cfg.versionFunc = fn
		return nil
	}
	return option(f)
}

// SubCommands is an [Option] that attaches 1 or more subcommands to the command being configured.
func SubCommands(subcommands ...*Command) Option {
	// TODO: We need to handle some potential misconfigurations here
	// 1) The user could attempt to add the command itself as a subcommand
	// 2) The user could try and add a duplicate subcommand, it shouldn't appear in the slice twice
	f := func(cfg *config) error {
		// Add the subcommands to the command this is being called on
		cfg.subcommands = append(cfg.subcommands, subcommands...)
		return nil
	}
	return option(f)
}

// Allow is an [Option] that allows for validating positional arguments to a [Command].
//
// You provide a validator function that returns an error if it encounters invalid arguments, and it will
// be run for you, passing in the non-flag arguments to the [Command] that was called.
//
// Successive calls overwrite previous ones.
//
//	// No positional arguments allowed
//	cli.New("test", cli.Allow(cli.NoArgs))
func Allow(validator func(cmd *Command, args []string) error) Option {
	f := func(cfg *config) error {
		cfg.allowArgs = validator
		return nil
	}
	return option(f)
}

// Flag is an [Option] that adds a flag to a [Command], storing its value in a variable.
//
// The variable is set when the flag is parsed during command execution.
//
//	// Add a force flag
//	var force bool
//	cli.New("rm", cli.Flag(&force, "force", "f", false, "Force deletion without confirmation"))
func Flag[T flag.Flaggable](p *T, name string, short string, value T, usage string) Option {
	// TODO: Some potential errors here although pflag just panics on most of them
	// when replacing with my own version, I need to handle the errors properly here like
	// duplicate flags etc.

	f := func(cfg *config) error {
		flag := flag.New(p, name, short, value, usage)
		var defVal string
		if flag.Type() == "bool" {
			defVal = "true"
		}
		cfg.flags.AddFlag(&pflag.Flag{
			Name:        name,
			Shorthand:   short,
			Usage:       usage,
			Value:       flag,
			DefValue:    flag.String(),
			NoOptDefVal: defVal,
		})
		return nil
	}
	return option(f)
}
