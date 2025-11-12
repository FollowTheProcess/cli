// Package cli provides a clean, minimal and simple mechanism for constructing CLI commands.
package cli // import "go.followtheprocess.codes/cli"

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"go.followtheprocess.codes/cli/internal/arg"
	"go.followtheprocess.codes/cli/internal/flag"
	"go.followtheprocess.codes/cli/internal/style"
)

const (
	helpBufferSize    = 1024                               // helpBufferSize is sufficient to hold most command --help text.
	versionBufferSize = 256                                // versionBufferSize is sufficient to hold all the --version text.
	defaultVersion    = "dev"                              // defaultVersion is the version shown in --version when the user has not provided one.
	defaultShort      = "A placeholder for something cool" // defaultShort is the default value for cli.Short.
)

// Builder is a function that constructs and returns a [Command], it makes constructing
// complex command trees easier as they can be passed directly to the [SubCommands] option.
type Builder func() (*Command, error)

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
		flags:   flag.NewSet(),
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		rawArgs: os.Args[1:],
		name:    name,
		version: defaultVersion,
		short:   defaultShort,
	}

	// Apply the options, gathering up all the validation errors
	// to report in one go
	var errs error
	for _, option := range options {
		errs = errors.Join(errs, option.apply(&cfg))
	}

	// Ensure we always have at least help and version flags
	err := Flag(&cfg.helpCalled, "help", 'h', false, "Show help for "+name).apply(&cfg)
	errs = errors.Join(errs, err) // nil errors are discarded in join

	err = Flag(&cfg.versionCalled, "version", 'V', false, "Show version info for "+name).apply(&cfg)

	errs = errors.Join(errs, err)
	if errs != nil {
		return nil, errs
	}

	// Additional validation that can't be done per-option
	// A command cannot have no subcommands and no run function, it must define one or the other
	if cfg.run == nil && len(cfg.subcommands) == 0 {
		return nil, fmt.Errorf(
			"command %s has no subcommands and no run function, a command must either be runnable or have subcommands",
			cfg.name,
		)
	}

	return cfg.build(), nil
}

// Command represents a CLI command.
//
// Commands in cli are recursive, that means that the root command
// is not different or special compared to any of its subcommands, this structure
// makes defining complex command trees as simple as creating a single command.
//
// In the command line 'git commit -m "Message"' both 'git' and 'commit'
// would be a clio command, '-m' would be a Flag taking a string argument.
//
// Commands are constructed with the [New] function and customised by
// providing a number of functional options to layer different settings
// and functionality.
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

	// run is the function actually implementing the command, the command is passed into the function for access
	// to things like cmd.Stdout().
	run func(ctx context.Context, cmd *Command) error

	// flags is the set of flags for this command.
	flags *flag.Set

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

	// commit is the commit hash of the binary, if passed shown when -v/--version is called
	commit string

	// buildDate is the date the release binary was built on, if passed shown when -v/--version is called
	buildDate string

	// examples is examples of how to use the command.
	examples []example

	// rawArgs are the raw arguments passed to the command prior to any parsing, defaulting to [os.Args]
	// (excluding the command name, so os.Args[1:]), can be overridden using
	// the [OverrideArgs] option for e.g. testing.
	rawArgs []string

	// args are the command line arguments declared by the user using the [cli.Args] option.
	args []arg.Value

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

// Execute parses the flags and arguments, and invokes the Command's Run
// function, returning any error.
//
// If the flags fail to parse, an error will be returned and the Run function
// will not be called.
func (cmd *Command) Execute(ctx context.Context) error {
	if cmd == nil {
		return errors.New("Execute called on a nil Command")
	}

	// Regardless of where we call execute, run it only from the root command, this is to ensure
	// that when we use the arguments to go and find the subcommand to run (if needed), then we
	// at the root of the command tree.
	if cmd.parent != nil {
		return fmt.Errorf("Execute must be called on the root of the command tree, was called on %s", cmd.name)
	}

	// Use the raw arguments and the command tree to determine which subcommand (if any)
	// we should be invoking and swap that into 'cmd'.
	//
	// Slightly magical trick but it simplifies a lot of stuff below.
	cmd, args := findRequestedCommand(cmd, cmd.rawArgs)

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
		if err := showHelp(cmd); err != nil {
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
		if err := showVersion(cmd); err != nil {
			return fmt.Errorf("could not show version: %w", err)
		}

		return nil
	}

	nonExtraArgs := cmd.flagSet().Args()
	terminatorIndex := slices.Index(nonExtraArgs, "--")

	if terminatorIndex != -1 {
		nonExtraArgs = nonExtraArgs[:terminatorIndex]
	}

	for i, argument := range cmd.args {
		var str string
		// The argument has been provided
		if len(nonExtraArgs) > i {
			str = nonExtraArgs[i]
		} else {
			// It hasn't, use the default
			str = argument.Default()
			if str == "" {
				return fmt.Errorf("argument %q is required and no value was provided", argument.Name())
			}
		}

		if err := argument.Set(str); err != nil {
			return fmt.Errorf("could not parse argument %q from provided input %q: %w", argument.Name(), str, err)
		}
	}

	// If the command is runnable, go and execute its run function
	if cmd.run != nil {
		return cmd.run(ctx, cmd)
	}

	// The only way we get here is if the command has subcommands defined but got no arguments given to it
	// so just show the usage and error
	if err := showHelp(cmd); err != nil {
		return err
	}

	return fmt.Errorf("command %q expected arguments (subcommands) but got none", cmd.name)
}

// Stdout returns the configured Stdout for the Command.
func (cmd *Command) Stdout() io.Writer {
	return cmd.root().stdout
}

// Stderr returns the configured Stderr for the Command.
func (cmd *Command) Stderr() io.Writer {
	return cmd.root().stderr
}

// Stdin returns the configured Stdin for the Command.
func (cmd *Command) Stdin() io.Reader {
	return cmd.root().stdin
}

// Args returns the positional arguments passed to the command.
func (cmd *Command) Args() []string {
	return cmd.flagSet().Args()
}

// ExtraArgs returns any additional arguments following a "--", and a boolean indicating
// whether or not they were present. This is useful for when you want to implement argument
// pass through in your commands.
//
// If there were no extra arguments, it will return nil, false.
func (cmd *Command) ExtraArgs() (args []string, ok bool) {
	extra := cmd.flagSet().ExtraArgs()
	if len(extra) > 0 {
		return extra, true
	}

	return nil, false
}

// Flags returns the set of flags for the command.
func (cmd *Command) flagSet() *flag.Set {
	if cmd == nil {
		// Only thing to do really, slightly more helpful than a generic
		// nil pointer dereference
		panic("flagSet called on a nil Command")
	}

	if cmd.flags == nil {
		return flag.NewSet()
	}

	return cmd.flags
}

// root returns the root of the command tree.
func (cmd *Command) root() *Command {
	if cmd.parent != nil {
		return cmd.parent.root()
	}

	return cmd
}

// hasFlag returns whether the command has a flag of the given name defined.
func (cmd *Command) hasFlag(name string) bool {
	flag, ok := cmd.flagSet().Get(name)
	if !ok {
		return false
	}

	return flag.NoArgValue() != ""
}

// hasShortFlag returns whether the command has a shorthand flag of the given name defined.
func (cmd *Command) hasShortFlag(name string) bool {
	if name == "" {
		return false
	}

	char, _ := utf8.DecodeRuneInString(name)

	flag, ok := cmd.flagSet().GetShort(char)
	if !ok {
		return false
	}

	return flag.NoArgValue() != ""
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
			}

			args = args[1:]

			continue

		case arg != "" && !strings.HasPrefix(arg, "-"):
			// We have a valid positional arg
			argsWithoutFlags = append(argsWithoutFlags, arg)
		}
	}

	return argsWithoutFlags
}

// showHelp is the default for a command's helpFunc.
func showHelp(cmd *Command) error {
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
	s.Grow(helpBufferSize)

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

	s.WriteString(style.Title.Text("Usage"))
	s.WriteString(": ")
	s.WriteString(style.Bold.Text(cmd.name))

	if len(cmd.subcommands) == 0 {
		// We don't have any subcommands so usage will be:
		// "Usage: {name} [OPTIONS] ARGS..."
		s.WriteString(" [OPTIONS]")

		if len(cmd.args) > 0 {
			// If we have named args, use the names in the help text
			writePositionalArgs(cmd, s)
		} else {
			// Otherwise, the command accepts arbitrary arguments
			s.WriteString(" ARGS...")
		}
	} else {
		// We do have subcommands, so usage will instead be:
		// "Usage: {name} [OPTIONS] COMMAND"
		s.WriteString(" [OPTIONS] COMMAND")
	}

	// If we have defined, list them explicitly and use their descriptions
	if len(cmd.args) != 0 {
		if err := writeArgumentsSection(cmd, s); err != nil {
			return err
		}
	}

	// If the user defined some examples, show those
	if len(cmd.examples) != 0 {
		writeExamples(cmd, s)
	}

	// Now show subcommands
	if len(cmd.subcommands) != 0 {
		if err := writeSubcommands(cmd, s); err != nil {
			return err
		}
	}

	// Now options
	if len(cmd.examples) != 0 || len(cmd.subcommands) != 0 || len(cmd.args) != 0 {
		// If there were examples or subcommands or named arguments, the last one would have printed a newline
		s.WriteString("\n")
	} else {
		// If there weren't, we need some more space
		s.WriteString("\n\n")
	}

	s.WriteString(style.Title.Text("Options"))
	s.WriteString(":\n\n")
	s.WriteString(usage)

	// Subcommand help
	if len(cmd.subcommands) != 0 {
		writeFooter(cmd, s)
	}

	// Note: It's important to use cmd.Stderr() here over cmd.stderr
	// as it resolves to the root's stderr
	fmt.Fprint(cmd.Stderr(), s.String())

	return nil
}

// writePositionalArgs writes any positional arguments in the correct
// format for the top level usage string in the help text string builder.
func writePositionalArgs(cmd *Command, s *strings.Builder) {
	for _, arg := range cmd.args {
		s.WriteString(" ")

		displayName := strings.ToUpper(arg.Name())

		if arg.Default() != "" {
			// It has a default so is not required
			s.WriteString("[")
			s.WriteString(displayName)
			s.WriteString("]")
		} else {
			// It is required
			s.WriteString(displayName)
		}
	}
}

// writeArgumentsSection writes the entire positional arguments block to the help
// text string builder.
func writeArgumentsSection(cmd *Command, s *strings.Builder) error {
	s.WriteString("\n\n")
	s.WriteString(style.Title.Text("Arguments"))
	s.WriteString(":\n\n")
	tw := style.Tabwriter(s)

	for _, arg := range cmd.args {
		switch arg.Default() {
		case "":
			// It's required
			fmt.Fprintf(tw, "  %s\t%s\t%s\t[required]\n", style.Bold.Text(arg.Name()), arg.Type(), arg.Usage())
		default:
			// It has a default
			fmt.Fprintf(tw, "  %s\t%s\t%s\t[default %s]\n", style.Bold.Text(arg.Name()), arg.Type(), arg.Usage(), arg.Default())
		}
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("could not format arguments: %w", err)
	}

	return nil
}

// writeExamples writes the examples block to the help text string builder.
func writeExamples(cmd *Command, s *strings.Builder) {
	// If there were positional args, the last one would have printed a newline
	if len(cmd.args) != 0 {
		s.WriteString("\n")
	} else {
		// If not, we need a bit more space
		s.WriteString("\n\n")
	}

	s.WriteString(style.Title.Text("Examples"))
	s.WriteByte(':')
	s.WriteString("\n\n")

	s.WriteString("  # ")
	s.WriteString(cmd.examples[0].comment)
	s.WriteByte('\n')
	s.WriteString("  $ ")
	s.WriteString(cmd.examples[0].command)

	for _, example := range cmd.examples[1:] {
		s.WriteString("\n\n")
		s.WriteString("  # ")
		s.WriteString(example.comment)
		s.WriteByte('\n')
		s.WriteString("  $ ")
		s.WriteString(example.command)
	}

	s.WriteByte('\n')
}

// writeSubcommands writes the subcommand block to the help text string builder.
func writeSubcommands(cmd *Command, s *strings.Builder) error {
	// If there were examples, the last one would have printed a newline
	if len(cmd.examples) != 0 {
		s.WriteByte('\n')
	} else {
		s.WriteString("\n\n")
	}

	s.WriteString(style.Title.Text("Commands"))
	s.WriteByte(':')
	s.WriteString("\n\n")

	tw := style.Tabwriter(s)
	for _, subcommand := range cmd.subcommands {
		fmt.Fprintf(tw, "  %s\t%s\n", style.Bold.Text(subcommand.name), subcommand.short)
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("could not format subcommands: %w", err)
	}

	return nil
}

// writeFooter writes the footer to the help text string builder.
func writeFooter(cmd *Command, s *strings.Builder) {
	s.WriteByte('\n')
	s.WriteString(`Use "`)
	s.WriteString(cmd.name)
	s.WriteString(" [command] -h/--help")
	s.WriteString(`" `)
	s.WriteString("for more information about a command.")
	s.WriteByte('\n')
}

// showVersion is the default implementation of the --version flag.
func showVersion(cmd *Command) error {
	if cmd == nil {
		return errors.New("defaultVersion called on a nil Command")
	}

	name := cmd.name // Incase we need to show the subcommand name

	if cmd.version == defaultVersion {
		// User has not set a version for this command, so we show the root version info
		cmd = cmd.root()
	}

	s := &strings.Builder{}
	s.Grow(versionBufferSize)
	s.WriteString(style.Title.Text(name))
	s.WriteString("\n\n")
	s.WriteString(style.Bold.Text("Version:"))
	s.WriteString(" ")
	s.WriteString(cmd.version)
	s.WriteString("\n")

	if cmd.commit != "" {
		s.WriteString(style.Bold.Text("Commit:"))
		s.WriteString(" ")
		s.WriteString(cmd.commit)
		s.WriteString("\n")
	}

	if cmd.buildDate != "" {
		s.WriteString(style.Bold.Text("BuildDate:"))
		s.WriteString(" ")
		s.WriteString(cmd.buildDate)
		s.WriteString("\n")
	}

	// Note: It's important to use cmd.Stderr() here over cmd.stderr
	// as it resolves to the root's stderr
	fmt.Fprint(cmd.Stderr(), s.String())

	return nil
}
