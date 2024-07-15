// Package cli provides a tiny, simple and minimalistic CLI framework for building Go CLI tools.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/pflag"
)

// TableWriter config, used for showing subcommands in help.
const (
	minWidth = 0    // Min cell width
	tabWidth = 8    // Tab width in spaces
	padding  = 1    // Padding
	padChar  = '\t' // Char to pad with
)

// New builds and returns a new [Command].
//
// The command can be customised by passing in a number of options enabling you to
// do things like configure stderr and stdout, add or customise help or version output
// add subcommands and run functions etc.
//
// Without any options passed, the default implementation returns a [Command] with no subcommands,
// a -v/--version and a -h/--help flag, hooked up to [os.Stdin], [os.Stdout] and [os.Stderr]
// and accepting [os.Args] as arguments (with the command path stripped, equivalent to os.Args[1:]).
func New(name string, options ...Option) *Command {
	// Default implementation
	cmd := &Command{
		flags:       pflag.NewFlagSet(name, pflag.ContinueOnError),
		stdin:       os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		args:        os.Args[1:],
		name:        name,
		version:     "dev",
		versionFunc: defaultVersion,
		short:       "A placeholder for something cool",
	}

	// Add the built in help and version flags
	cmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Show help for %s", name))
	cmd.Flags().BoolP("version", "v", false, fmt.Sprintf("Show version info for %s", name))

	// Apply the options
	for _, option := range options {
		option(cmd)
	}

	return cmd
}

// Command represents a CLI command. In terms of an example, in the line
// git commit -m <msg>; 'commit' is the command. It can have any number of subcommands
// which themselves can have subcommands etc. The root command in this example is 'git'.
type Command struct {
	// run is the function actually implementing the command, the command and arguments to it, are passed into the function, flags
	// are parsed out before the arguments are passed to Run, so `args` here are the command line arguments minus flags.
	run func(cmd *Command, args []string) error

	// flags is the set of flags for this command.
	flags *pflag.FlagSet

	// versionFunc is the function thatgets called when the user calls -v/--version.
	//
	// It can be overridden by the user to customise their version output using
	// the [VersionFunc] [Option].
	versionFunc func(cmd *Command) error

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

	// parent is the immediate parent of this subcommand. If this command is the root
	// and has no parent, this will be nil.
	parent *Command

	// name is the name of the command.
	name string

	// short is the one line summary for the command, shown inline in the -h/--help output.
	short string

	// long is the long form description for the command, shown when -h/--help is called on the command itself.
	long string

	// version is the version of this command, shown when -v/--version is called.
	version string

	// example is examples of how to use the command.
	example []Example

	// args is the arguments passed to the command, default to [os.Args]
	// (excluding the command name, so os.Args[1:]), can be overridden using
	// the [Args] option for e.g. testing.
	args []string

	// subcommands is the list of subcommands this command has directly underneath it,
	// these may have any number of subcommands under them, this is how we form nested
	// command structures.
	//
	// If the command has no subcommands, this slice will be nil.
	subcommands []*Command
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
		return fmt.Sprintf("\n  # %s\n", e.Comment)
	case e.Comment == "":
		// No comment, just show command on it's own
		return fmt.Sprintf("\n  $ %s\n", e.Command)
	default:
		// Both passed, show the full example
		return fmt.Sprintf("\n  # %s\n  $ %s\n", e.Comment, e.Command)
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

	// Regardless of where we call execute, run it only from the root command
	if c.parent != nil {
		return c.root().Execute()
	}

	if err := c.Flags().Parse(c.args); err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	// Check if we should be responding to -h/--help
	helpCalled, err := c.Flags().GetBool("help")
	if err != nil {
		// We shouldn't ever get here because we define a default for help
		return fmt.Errorf("could not parse help flag: %w", err)
	}

	// If -h/--help was called, call the defined helpFunc and exit so that
	// the run function is never called.
	if helpCalled {
		if err = defaultHelp(c); err != nil {
			return fmt.Errorf("help function returned an error: %w", err)
		}
		return nil
	}

	// Check if we should be responding to -v/--version
	versionCalled, err := c.Flags().GetBool("version")
	if err != nil {
		// Again, shouldn't ever get here
		return fmt.Errorf("could not parse version flag: %w", err)
	}

	// If -v/--version was called, call the defined versionFunc and exit so that
	// the run function is never called
	if versionCalled {
		if c.versionFunc == nil {
			return errors.New("versionFunc was nil")
		}
		if err := c.versionFunc(c); err != nil {
			return fmt.Errorf("version function returned an error: %w", err)
		}
		return nil
	}

	// Not all commands are runnable, e.g. if this command is the root of a subcommand
	// it will define subcommands but no run function itself. We must decide here what to do when
	// this command is executed.

	// A command cannot have
	// no subcommands and no run function, it must define one or the other
	if c.run == nil && len(c.subcommands) == 0 {
		return fmt.Errorf(
			"command %s has no subcommands and no run function, a command must either be runnable or have subcommands",
			c.name,
		)
	}

	// If the command is runnable, go and execute it's run function
	if c.run != nil {
		return c.run(c, c.Flags().Args())
	}

	// TODO: If the command defines subcommands, we need to parse the args and determine which subcommand to go and run

	return nil
}

// Flags returns the set of flags for the command.
//
// TODO: Make it so we can add flags via the functional options pattern.
func (c *Command) Flags() *pflag.FlagSet {
	if c == nil {
		// Only thing to do really, slightly more helpful than a generic
		// nil pointer dereference
		panic("Flags called on a nil Command")
	}
	if c.flags == nil {
		return pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	}
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

// root returns the root of the command tree.
func (c *Command) root() *Command {
	if c.parent != nil {
		return c.parent.root()
	}
	return c
}

// defaultHelp is the default for a command's helpFunc.
func defaultHelp(cmd *Command) error {
	if cmd == nil {
		return errors.New("defaultHelp called on a nil Command")
	}
	// Note: The decision to not use text/template here is intentional, template calls
	// reflect.Value.MethodByName() and/or reflect.Type.MethodByName() which disables dead
	// code elimination in the compiler, meaning any application that uses cli for it's
	// command line interface will not be run through dead code elimination which could cause
	// significant increase in memory consumption and disk space.
	// See https://github.com/spf13/cobra/issues/2015
	s := &strings.Builder{}

	// If we have a short description, write that
	if cmd.short != "" {
		s.WriteString(cmd.short)
		s.WriteString("\n\n")
	}

	// If we have a long description, write that
	if cmd.long != "" {
		s.WriteString(cmd.long)
		s.WriteString("\n\n")
	}

	// TODO: See if we can be clever about dynamically generating the syntax for e.g. variadic args
	// required args, flags etc. based on what the command has defined.
	if len(cmd.subcommands) == 0 {
		// We don't have any subcommands so usage will be:
		// "Usage: {name} [OPTIONS] ARGS..."
		s.WriteString("Usage: ")
		s.WriteString(cmd.name)
		s.WriteString(" [OPTIONS] ARGS...")
	} else {
		// We do have subcommands, so usage will instead be:
		// "Usage: {name} [OPTIONS] COMMAND"
		s.WriteString("Usage: ")
		s.WriteString(cmd.name)
		s.WriteString(" [OPTIONS] COMMAND")
	}

	// If the user defined some examples, show those
	if len(cmd.example) != 0 {
		s.WriteString("\n\nExamples:")
		for _, example := range cmd.example {
			s.WriteString(example.String())
		}
	}

	// Now show subcommands
	if len(cmd.subcommands) != 0 {
		s.WriteString("\n\nCommands:\n")
		tab := tabwriter.NewWriter(s, minWidth, tabWidth, padding, padChar, tabwriter.AlignRight)
		for _, subcommand := range cmd.subcommands {
			fmt.Fprintf(tab, "  %s\t%s\n", subcommand.name, subcommand.short)
		}
		if err := tab.Flush(); err != nil {
			return fmt.Errorf("could not format subcommands: %w", err)
		}
	}

	// Now options
	if len(cmd.example) != 0 || len(cmd.subcommands) != 0 {
		// If there were examples or subcommands, the last one would have printed a newline
		s.WriteString("\n")
	} else {
		// If there weren't, we need some more space
		s.WriteString("\n\n")
	}
	s.WriteString("Options:\n")
	s.WriteString(cmd.Flags().FlagUsages())

	// Subcommand help
	s.WriteString("\n")
	s.WriteString(`Use "`)
	s.WriteString(cmd.name)
	s.WriteString(" [command] -h/--help")
	s.WriteString(`" `)
	s.WriteString("for more information about a command.")
	s.WriteString("\n")

	fmt.Fprint(cmd.Stderr(), s.String())

	return nil
}

// defaultVersion is the default for a command's versionFunc.
func defaultVersion(cmd *Command) error {
	if cmd == nil {
		return errors.New("defaultVersion called on a nil Command")
	}
	fmt.Fprintf(cmd.Stderr(), "%s, version: %s\n", cmd.name, cmd.version)
	return nil
}
