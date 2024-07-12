// Package cli provides a tiny, simple and minimalistic CLI framework for building Go CLI tools.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

// New builds and returns a new [Command].
//
// Without any options passed, the default implementation returns a [Command] with
// no flags, hooked up to [os.Stdin], [os.Stdout] and [os.Stderr], and accepting
// [os.Args] as arguments (with the command path stripped, equivalent to os.Args[1:]).
//
// This default command, when invoked, will print "Hello from {name}\n" to [os.Stdout].
func New(name string, options ...Option) *Command {
	// Default implementation
	cmd := &Command{
		run: func(cmd *Command, args []string) error {
			fmt.Fprintf(cmd.stdout, "Hello from %s\n", name)
			return nil
		},
		flags:  pflag.NewFlagSet(name, pflag.ContinueOnError),
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args[1:],
		name:   name,
	}

	for _, option := range options {
		option(cmd)
	}

	return cmd
}

// Command represents a CLI command. In terms of an example, in the line
// git commit -m <msg>; 'commit' is the command. It can have any number of subcommands
// which themselves can have subcommands etc.
type Command struct {
	// run is the function actually implementing the command, the command and arguments to it, are passed into the function, flags
	// are parsed out before the arguments are passed to Run, so `args` here are the command line arguments minus flags.
	run func(cmd *Command, args []string) error

	// flags is the set of flags for this command.
	flags *pflag.FlagSet

	// stdin is an [io.Reader] from which command input is read.
	//
	// It defaults to [os.Stdin] but can be overridden as desired e.g. for testing.
	stdin io.Reader

	// stdout is an [io.Writer] to which normal command output is written.
	//
	// It defaults to [os.Stdout] but can be overridden as desired e.g. for testing.
	stdout io.Writer

	// stderr is an [io.Writer] to which error command output is written.
	//
	// It defaults to [os.Stderr] but can be overridden as desired e.g. for testing.
	stderr io.Writer

	// name is the name of the command.
	name string

	// short is the one line summary for the command, shown inline in the -h/--help output.
	short string

	// long is the long form description for the command, shown when -h/--help is called on the command itself.
	long string

	// example is examples of how to use the command.
	example []Example

	// args is the arguments passed to the command, default to [os.Args]
	// (excluding the command name, so os.Args[1:]), can be overridden using
	// the [Args] option.
	args []string
}

// Example is a single usage example for a [Command].
//
// The example will be shown in the -h/--help output as follows:
//
//	# Comment
//	$ Command
type Example struct {
	Comment string // The comment for the example.
	Command string // The command string for the example.
}

// String implements [fmt.Stringer] for [Example].
func (e Example) String() string {
	switch {
	case e.Comment == "" && e.Command == "":
		// Empty example, return empty string
		return ""
	case e.Command == "":
		// Empty command, show just the comment
		return fmt.Sprintf("# %s", e.Comment)
	case e.Comment == "":
		// No comment, just show command on it's own
		return fmt.Sprintf("$ %s", e.Command)
	default:
		// Both passed, show the full example
		return fmt.Sprintf("# %s\n$ %s", e.Comment, e.Command)
	}
}

// Execute parses the flags and arguments, and invokes the Command's Run
// function, returning any error.
//
// If the flags fail to parse, an error will be returned and the Run function
// will not be called.
func (c *Command) Execute() error {
	if c == nil {
		return errors.New("Execute called on a nil Command")
	}

	if err := c.Flags().Parse(c.args); err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	argsWithoutFlags := c.flags.Args()

	return c.run(c, argsWithoutFlags)
}

// Flags returns the set of flags for the command.
//
// TODO: Make it so we can add flags via the functional options pattern.
func (c *Command) Flags() *pflag.FlagSet {
	return c.flags
}

// Stdout returns the configured Stdout for the Command.
func (c *Command) Stdout() io.Writer {
	return c.stdout
}

// Stderr returns the configured Stderr for the Command.
func (c *Command) Stderr() io.Writer {
	return c.stderr
}

// Stdin returns the configured Stdin for the Command.
func (c *Command) Stdin() io.Reader {
	return c.stdin
}
