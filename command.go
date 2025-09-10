// Package cli provides a clean, minimal and simple mechanism for constructing CLI commands.
package cli // import "go.followtheprocess.codes/cli"

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"go.followtheprocess.codes/cli/internal/colour"
	"go.followtheprocess.codes/cli/internal/flag"
	"go.followtheprocess.codes/cli/internal/table"
)

const (
	helpBufferSize    = 1024 // helpBufferSize is sufficient to hold most command --help text.
	versionBufferSize = 256  // versionBufferSize is sufficient to hold all the --version text.
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
		flags:        flag.NewSet(),
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
		args:         os.Args[1:],
		name:         name,
		version:      "dev",
		short:        "A placeholder for something cool",
		argValidator: AnyArgs(),
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

	// commit is the commit hash of the binary, if passed shown when -v/--version is called
	commit string

	// buildDate is the date the release binary was built on, if passed shown when -v/--version is called
	buildDate string

	// examples is examples of how to use the command.
	examples []example

	// args are the raw arguments passed to the command prior to any parsing, defaulting to [os.Args]
	// (excluding the command name, so os.Args[1:]), can be overridden using
	// the [OverrideArgs] option for e.g. testing.
	args []string

	// positionalArgs are the named positional arguments to the command, positional arguments
	// may be retrieved from within command logic by name and this also significantly
	// enhances the help message.
	positionalArgs []positionalArg

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
	return fmt.Sprintf("\n  # %s\n  $ %s\n", e.comment, e.command)
}

// Execute parses the flags and arguments, and invokes the Command's Run
// function, returning any error.
//
// If the flags fail to parse, an error will be returned and the Run function
// will not be called.
func (cmd *Command) Execute() error {
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
	cmd, args := findRequestedCommand(cmd, cmd.args)

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
		if err := defaultVersion(cmd); err != nil {
			return fmt.Errorf("version function returned an error: %w", err)
		}

		return nil
	}

	// Validate the arguments using the command's allowedArgs function
	argsWithoutFlags := cmd.flagSet().Args()
	if err := cmd.argValidator(cmd, argsWithoutFlags); err != nil {
		return err
	}

	// Now we have the actual positional arguments to the command, we can use our
	// named arguments to assign the given values (or the defaults) to the arguments
	// so they may be retrieved by name.
	//
	// We're modifying the slice in place here, hence not using a range loop as it
	// would take a copy of the c.positionalArgs slice
	for i := range len(cmd.positionalArgs) {
		if i >= len(argsWithoutFlags) {
			arg := cmd.positionalArgs[i]

			// If we've fallen off the end of argsWithoutFlags and the positionalArg at this
			// index does not have a default, it means the arg was required but not provided
			if arg.defaultValue == requiredArgMarker {
				return fmt.Errorf("missing required argument %q, expected at position %d", arg.name, i)
			}
			// It does have a default, so use that instead
			cmd.positionalArgs[i].value = arg.defaultValue
		} else {
			// We are in a valid index in both slices which means the named positional
			// argument at this index was provided on the command line, so all we need
			// to do is set its value
			cmd.positionalArgs[i].value = argsWithoutFlags[i]
		}
	}

	// If the command is runnable, go and execute its run function
	if cmd.run != nil {
		return cmd.run(cmd, argsWithoutFlags)
	}

	// The only way we get here is if the command has subcommands defined but got no arguments given to it
	// so just show the usage and error
	if err := defaultHelp(cmd); err != nil {
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

// Arg looks up a named positional argument by name.
//
// If the argument was defined with a default, and it was not provided on the command line
// then the value returned will be the default value.
//
// If no named argument exists with the given name, it will return "".
func (cmd *Command) Arg(name string) string {
	for _, arg := range cmd.positionalArgs {
		if arg.name == name {
			// arg.value will have been set to the default already during command line parsing
			// if the arg was not provided
			return arg.value
		}
	}

	return ""
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

// subcommandNames returns a list of all the names of the current command's registered subcommands.
func (cmd *Command) subcommandNames() []string {
	names := make([]string, 0, len(cmd.subcommands))
	for _, sub := range cmd.subcommands {
		names = append(names, sub.name)
	}

	return names
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

	s.WriteString(colour.Title("Usage"))
	s.WriteString(": ")
	s.WriteString(colour.Bold(cmd.name))

	if len(cmd.subcommands) == 0 {
		// We don't have any subcommands so usage will be:
		// "Usage: {name} [OPTIONS] ARGS..."
		s.WriteString(" [OPTIONS] ")

		if len(cmd.positionalArgs) > 0 {
			// If we have named args, use the names in the help text
			writePositionalArgs(cmd, s)
		} else {
			// We have no named arguments so do the best we can
			// TODO(@FollowTheProcess): Can we detect if cli.NoArgs was used in which case
			// omit this
			s.WriteString("ARGS...")
		}
	} else {
		// We do have subcommands, so usage will instead be:
		// "Usage: {name} [OPTIONS] COMMAND"
		s.WriteString(" [OPTIONS] COMMAND")
	}

	// If we have named arguments, list them explicitly and use their descriptions
	if len(cmd.positionalArgs) != 0 {
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
	if len(cmd.examples) != 0 || len(cmd.subcommands) != 0 || len(cmd.positionalArgs) != 0 {
		// If there were examples or subcommands or named arguments, the last one would have printed a newline
		s.WriteString("\n")
	} else {
		// If there weren't, we need some more space
		s.WriteString("\n\n")
	}

	s.WriteString(colour.Title("Options"))
	s.WriteString(":\n\n")
	s.WriteString(usage)

	// Subcommand help
	if len(cmd.subcommands) != 0 {
		writeFooter(cmd, s)
	}

	fmt.Fprint(cmd.Stderr(), s.String())

	return nil
}

// writePositionalArgs writes any positional arguments in the correct
// format for the top level usage string in the help text string builder.
func writePositionalArgs(cmd *Command, s *strings.Builder) {
	for _, arg := range cmd.positionalArgs {
		displayName := strings.ToUpper(arg.name)

		if arg.defaultValue != requiredArgMarker {
			// If it has a default, it's an optional argument so wrap it
			// in brackets e.g. [FILE]
			s.WriteString("[")
			s.WriteString(displayName)
			s.WriteString("]")
		} else {
			// It's required, so just FILE
			s.WriteString(displayName)
		}

		s.WriteString(" ")
	}
}

// writeArgumentsSection writes the entire positional arguments block to the help
// text string builder.
func writeArgumentsSection(cmd *Command, s *strings.Builder) error {
	s.WriteString("\n\n")
	s.WriteString(colour.Title("Arguments"))
	s.WriteString(":\n")
	tab := table.New(s)

	for _, arg := range cmd.positionalArgs {
		switch arg.defaultValue {
		case requiredArgMarker:
			tab.Row("  %s\t%s\t[required]\n", colour.Bold(arg.name), arg.description)
		case "":
			tab.Row("  %s\t%s\t[default %q]\n", colour.Bold(arg.name), arg.description, arg.defaultValue)
		default:
			tab.Row("  %s\t%s\t[default %s]\n", colour.Bold(arg.name), arg.description, arg.defaultValue)
		}
	}

	if err := tab.Flush(); err != nil {
		return fmt.Errorf("could not format arguments: %w", err)
	}

	return nil
}

// writeExamples writes the examples block to the help text string builder.
func writeExamples(cmd *Command, s *strings.Builder) {
	// If there were positional args, the last one would have printed a newline
	if len(cmd.positionalArgs) != 0 {
		s.WriteString("\n")
	} else {
		// If not, we need a bit more space
		s.WriteString("\n\n")
	}

	s.WriteString(colour.Title("Examples"))
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

	s.WriteString(colour.Title("Commands"))
	s.WriteByte(':')
	s.WriteByte('\n')

	tab := table.New(s)
	for _, subcommand := range cmd.subcommands {
		tab.Row("  %s\t%s\n", colour.Bold(subcommand.name), subcommand.short)
	}

	if err := tab.Flush(); err != nil {
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

// defaultVersion is the default for a command's versionFunc.
func defaultVersion(cmd *Command) error {
	if cmd == nil {
		return errors.New("defaultVersion called on a nil Command")
	}

	s := &strings.Builder{}
	s.Grow(versionBufferSize)
	s.WriteString(colour.Title(cmd.name))
	s.WriteString("\n\n")
	s.WriteString(colour.Bold("Version:"))
	s.WriteString(" ")
	s.WriteString(cmd.version)
	s.WriteString("\n")

	if cmd.commit != "" {
		s.WriteString(colour.Bold("Commit:"))
		s.WriteString(" ")
		s.WriteString(cmd.commit)
		s.WriteString("\n")
	}

	if cmd.buildDate != "" {
		s.WriteString(colour.Bold("BuildDate:"))
		s.WriteString(" ")
		s.WriteString(cmd.buildDate)
		s.WriteString("\n")
	}

	fmt.Fprint(cmd.stderr, s.String())

	return nil
}
