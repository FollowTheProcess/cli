// Package cli provides a tiny, simple and minimalistic CLI framework for building Go CLI tools.
package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/pflag"
)

// Command represents a CLI command.
type Command struct {
	// Run is the function actually implementing the command, the command and arguments to it, are passed into the function, flags
	// are parsed out before the arguments are passed to Run, so `args` here are the command line arguments minus flags.
	Run func(cmd *Command, args []string) error

	// flags is the set of flags for this command.
	flags *pflag.FlagSet

	// Stdin is an [io.Reader] from which command input is read.
	//
	// It defaults to [os.Stdin] but can be overridden as desired e.g. for testing.
	Stdin io.Reader

	// Stdout is an [io.Writer] to which normal command output is written.
	//
	// It defaults to [os.Stdout] but can be overridden as desired e.g. for testing.
	Stdout io.Writer

	// Stderr is an [io.Writer] to which error command output is written.
	//
	// It defaults to [os.Stderr] but can be overridden as desired e.g. for testing.
	Stderr io.Writer

	// Name is the name of the command.
	Name string

	// Short is the one line summary for the command, shown inline in the -h/--help output.
	Short string

	// Long is the long form description for the command, shown when -h/--help is called on the command itself.
	Long string

	// Example is examples of how to use the command.
	Example []Example
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
// The arguments should not include the command name.
//
// If the flags fail to parse, an error will be returned and the Run function
// will not be called.
//
//	err := cmd.Execute(os.Args[1:])
func (c *Command) Execute(args []string) error {
	if c == nil {
		return errors.New("Execute called on a nil Command")
	}

	if err := c.Flags().Parse(args); err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	argsWithoutFlags := c.flags.Args()

	return c.Run(c, argsWithoutFlags)
}

// Flags returns the set of flags for the command.
func (c *Command) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = pflag.NewFlagSet(c.Name, pflag.ContinueOnError)
	}
	return c.flags
}
