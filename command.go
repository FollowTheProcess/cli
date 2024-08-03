// Package cli provides a clean, minimal and simple mechanism for constructing CLI commands.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/FollowTheProcess/cli/internal/table"
)

// New builds and returns a new [Command].
//
// The command can be customised by passing in a number of options enabling you to
// do things like configure stderr and stdout, add or customise help or version output
// add subcommands and run functions etc.
//
// Without any options passed, the default implementation returns a [Command] with no subcommands,
// a -v/--version and a -h/--help flag, hooked up to [os.Stdin], [os.Stdout] and [os.Stderr]
// and accepting arbitrary positional arguments from [os.Args] (with the command path stripped, equivalent to os.Args[1:]).
//
// Options will validate their inputs where possible and return errors which will be bubbled up through New
// to aid debugging invalid configuration.
func New(name string, options ...Option) (*Command, error) {
	// This was actually a nilaway thing, indexing into os.Args without knowing the length
	if len(os.Args) < 1 {
		return nil, fmt.Errorf("bad arguments expected [<command> <args>...], got %v", os.Args)
	}

	// Default implementation
	cfg := config{
		flags:        flag.NewSet(),
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
		args:         os.Args[1:],
		name:         name,
		version:      "dev",
		versionFunc:  defaultVersion,
		short:        "A placeholder for something cool",
		argValidator: AnyArgs(),
	}

	// Ensure we always have at least help and version flags
	defaultOptions := []Option{
		Flag(&cfg.helpCalled, "help", 'h', false, fmt.Sprintf("Show help for %s", name)),
		Flag(&cfg.versionCalled, "version", 'v', false, fmt.Sprintf("Show version info for %s", name)),
	}

	toApply := slices.Concat(options, defaultOptions)

	// Apply the options, gathering up all the validation errors
	// to report in one go. Each option returns only one error
	// so this can be pre-allocated.
	errs := make([]error, 0, len(toApply))
	for _, option := range toApply {
		err := option.apply(&cfg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg.build(), nil
}

// Command represents a CLI command. In terms of an example, in the line
// git commit -m <msg>; 'commit' is the command. It can have any number of subcommands
// which themselves can have subcommands etc. The root command in this example is 'git'.
type Command struct {
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

	// run is the function actually implementing the command, the command and arguments to it, are passed into the function, flags
	// are parsed out before the arguments are passed to Run, so `args` here are the command line arguments minus flags.
	run func(cmd *Command, args []string) error

	// flags is the set of flags for this command.
	flags *flag.Set

	// versionFunc is the function thatgets called when the user calls -v/--version.
	//
	// It can be overridden by the user to customise their version output using
	// the [VersionFunc] [Option].
	versionFunc func(cmd *Command) error

	// parent is the immediate parent of this subcommand. If this command is the root
	// and has no parent, this will be nil.
	parent *Command

	// argValidator is a function that gets called to validate the positional arguments
	// to the command. It defaults to allowing arbitrary arguments, can be overridden using
	// the [AllowArgs] option.
	argValidator ArgValidator

	// name is the name of the command.
	name string

	// short is the one line summary for the command, shown inline in the -h/--help output.
	short string

	// long is the long form description for the command, shown when -h/--help is called on the command itself.
	long string

	// version is the version of this command, shown when -v/--version is called.
	version string

	// examples is examples of how to use the command.
	examples []example

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

	// helpCalled is whether or not the --help flag was used.
	helpCalled bool

	// versionCalled is whether or not the --version flag was used.
	versionCalled bool
}

// example is a single usage example for a [Command].
//
// The example will be shown in the -h/--help output as follows:
//
//	# Comment
//	$ Command
type example struct {
	comment string // The comment for the example.
	command string // The command string for the example.
}

// String implements [fmt.Stringer] for [Example].
func (e example) String() string {
	return fmt.Sprintf("\n# %s\n$ %s\n", e.comment, e.command)
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

	// Regardless of where we call execute, run it only from the root command, this is to ensure
	// that when we use the arguments to go and find the subcommand to run (if needed), then we
	// at the root of the command tree.
	if c.parent != nil {
		return c.root().Execute()
	}

	// Use the raw arguments and the command tree to determine which subcommand (if any)
	// we should be invoking. If it turns out we want to invoke the root command, then
	// cmd here will be c.
	cmd, args := findRequestedCommand(c, c.args)

	if err := cmd.flagSet().Parse(args); err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	// If -h/--help was called, call the defined helpFunc and exit so that
	// the run function is never called.
	helpCalled, ok := cmd.flagSet().Help()
	if !ok {
		// Should never get here as we define a default help
		return errors.New("help flag not defined")
	}
	if helpCalled {
		if err := defaultHelp(cmd); err != nil {
			return fmt.Errorf("help function returned an error: %w", err)
		}
		return nil
	}

	// If -v/--version was called, call the defined versionFunc and exit so that
	// the run function is never called
	versionCalled, ok := cmd.flagSet().Version()
	if !ok {
		// Again, should be unreachable
		return errors.New("version flag not defined")
	}
	if versionCalled {
		if cmd.versionFunc == nil {
			return errors.New("versionFunc was nil")
		}
		if err := cmd.versionFunc(c); err != nil {
			return fmt.Errorf("version function returned an error: %w", err)
		}
		return nil
	}

	// A command cannot have no subcommands and no run function, it must define one or the other
	if cmd.run == nil && len(cmd.subcommands) == 0 {
		return fmt.Errorf(
			"command %s has no subcommands and no run function, a command must either be runnable or have subcommands",
			cmd.name,
		)
	}

	// Validate the arguments using the command's allowedArgs function
	argsWithoutFlags := cmd.flagSet().Args()
	if err := cmd.argValidator(cmd, argsWithoutFlags); err != nil {
		return err
	}

	// If the command is runnable, go and execute its run function
	if cmd.run != nil {
		return cmd.run(cmd, argsWithoutFlags)
	}

	// This basically only happens when we have subcommands defined but pass no args to the root command
	// in which case we'll just show the help text and error
	if err := defaultHelp(cmd); err != nil {
		return err
	}
	return errors.New("invalid arguments")
}

// Flags returns the set of flags for the command.
func (c *Command) flagSet() *flag.Set {
	if c == nil {
		// Only thing to do really, slightly more helpful than a generic
		// nil pointer dereference
		panic("Flags called on a nil Command")
	}
	if c.flags == nil {
		return flag.NewSet()
	}
	return c.flags
}

// Stdout returns the configured Stdout for the Command.
func (c *Command) Stdout() io.Writer {
	return c.root().stdout
}

// Stderr returns the configured Stderr for the Command.
func (c *Command) Stderr() io.Writer {
	return c.root().stderr
}

// Stdin returns the configured Stdin for the Command.
func (c *Command) Stdin() io.Reader {
	return c.root().stdin
}

// ExtraArgs returns any additional arguments following a "--", this is useful for when you want to
// implement argument pass through in your commands.
//
// If there were no extra arguments, it will return nil.
func (c *Command) ExtraArgs() []string {
	return c.flagSet().ExtraArgs()
}

// root returns the root of the command tree.
func (c *Command) root() *Command {
	if c.parent != nil {
		return c.parent.root()
	}
	return c
}

// hasFlag returns whether the command has a flag of the given name defined.
func (c *Command) hasFlag(name string) bool {
	if name == "" {
		return false
	}
	flag, ok := c.flagSet().Get(name)
	if !ok {
		return false
	}

	return flag.DefaultValueNoArg != ""
}

// hasShortFlag returns whether the command has a shorthand flag of the given name defined.
func (c *Command) hasShortFlag(name string) bool {
	if name == "" {
		return false
	}

	char, _ := utf8.DecodeRuneInString(name)

	flag, ok := c.flagSet().GetShort(char)
	if !ok {
		return false
	}

	return flag.DefaultValueNoArg != ""
}

// findRequestedCommand uses the raw arguments and the command tree to determine what
// (if any) subcommand is being requested and return that command along with the arguments
// that were meant for it.
func findRequestedCommand(cmd *Command, args []string) (*Command, []string) {
	// Any arguments without flags could be names of subcommands
	argsWithoutFlags := stripFlags(cmd, args)
	if len(argsWithoutFlags) == 0 {
		// If there are no non-flag arguments, we must already be either at the root command
		// or the correct subcommand
		return cmd, args
	}

	// The next non-flag argument will be the first immediate subcommand
	// e.g. in 'go mod tidy', argsWithoutFlags[0] will be 'mod'
	nextSubCommand := argsWithoutFlags[0]

	// Lookup this immediate subcommand by name and if we find it, recursively call
	// this function so we eventually end up at the end of the command tree with
	// the right arguments
	next := findSubCommand(cmd, nextSubCommand)
	if next != nil {
		return findRequestedCommand(next, argsMinusFirstX(args, nextSubCommand))
	}

	// Found it
	return cmd, args
}

// argsMinusFirstX removes only the first x from args.  Otherwise, commands that look like
// openshift admin policy add-role-to-user admin my-user, lose the admin argument (arg[4]).
func argsMinusFirstX(args []string, x string) []string {
	// Note: this is borrowed from Cobra but ours is a lot simpler because we don't support
	// persistent flags
	for i, arg := range args {
		if arg == x {
			return slices.Delete(args, i, i+1)
		}
	}
	return args
}

// findSubCommand searches the immediate subcommands of cmd by name, looking for next.
//
// If next is not found, it will return nil.
func findSubCommand(cmd *Command, next string) *Command {
	for _, subcommand := range cmd.subcommands {
		if subcommand.name == next {
			return subcommand
		}
	}
	return nil
}

// stripFlags takes a slice of raw command line arguments (including possible flags) and removes
// any arguments that are flags or values passed in to flags e.g. --flag value.
func stripFlags(cmd *Command, args []string) []string {
	if len(args) == 0 {
		return args
	}

	argsWithoutFlags := []string{}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]
		switch {
		case arg == "--":
			// "--" terminates the flags
			return argsWithoutFlags
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "=") && !cmd.hasFlag(arg[2:]):
			// If '--flag arg' then delete arg from args
			fallthrough // (do the same as below)
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !cmd.hasShortFlag(arg[1:]):
			// If '-f arg' then delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				return argsWithoutFlags
			} else {
				args = args[1:]
				continue
			}
		case arg != "" && !strings.HasPrefix(arg, "-"):
			// We have a valid positional arg
			argsWithoutFlags = append(argsWithoutFlags, arg)
		}
	}

	return argsWithoutFlags
}

// defaultHelp is the default for a command's helpFunc.
func defaultHelp(cmd *Command) error {
	if cmd == nil {
		return errors.New("defaultHelp called on a nil Command")
	}
	usage, err := cmd.flagSet().Usage()
	if err != nil {
		return fmt.Errorf("could not write usage: %w", err)
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
	if len(cmd.examples) != 0 {
		s.WriteString("\n\nExamples:")
		for _, example := range cmd.examples {
			s.WriteString(example.String())
		}
	}

	// Now show subcommands
	if len(cmd.subcommands) != 0 {
		s.WriteString("\n\nCommands:\n")
		tab := table.New(s)
		for _, subcommand := range cmd.subcommands {
			tab.Row("  %s\t%s\n", subcommand.name, subcommand.short)
		}
		if err := tab.Flush(); err != nil {
			return fmt.Errorf("could not format subcommands: %w", err)
		}
	}

	// Now options
	if len(cmd.examples) != 0 || len(cmd.subcommands) != 0 {
		// If there were examples or subcommands, the last one would have printed a newline
		s.WriteString("\n")
	} else {
		// If there weren't, we need some more space
		s.WriteString("\n\n")
	}
	s.WriteString("Options:\n")
	s.WriteString(usage)

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
