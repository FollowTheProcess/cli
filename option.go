package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"go.followtheprocess.codes/cli/arg"
	"go.followtheprocess.codes/cli/flag"
	internalarg "go.followtheprocess.codes/cli/internal/arg"
	internalflag "go.followtheprocess.codes/cli/internal/flag"
	"go.followtheprocess.codes/hue"
)

// Option is a functional option for configuring a [Command].
type Option interface {
	apply(cmd *Command) error
}

type stdinOpt struct{ stdin io.Reader }

func (o stdinOpt) apply(cmd *Command) error {
	if o.stdin == nil {
		return errors.New("cannot set Stdin to nil")
	}

	cmd.stdin = o.stdin

	return nil
}

// Stdin is an [Option] that sets the Stdin for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stdin].
//
// Subcommands cannot override the stdin of the root command, when retrieving the
// stdin with [Command.Stdin], the root's stdin will be returned.
//
//	// Set stdin to os.Stdin (the default anyway)
//	cli.New("test", cli.Stdin(os.Stdin))
func Stdin(stdin io.Reader) Option {
	return stdinOpt{stdin: stdin}
}

type stdoutOpt struct{ stdout io.Writer }

func (o stdoutOpt) apply(cmd *Command) error {
	if o.stdout == nil {
		return errors.New("cannot set Stdout to nil")
	}

	cmd.stdout = o.stdout

	return nil
}

// Stdout is an [Option] that sets the Stdout for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stdout].
//
// Subcommands cannot override the stdout of the root command, when retrieving the
// stdout with [Command.Stdout], the root's stdout will be returned.
//
//	// Set stdout to a temporary buffer
//	buf := &bytes.Buffer{}
//	cli.New("test", cli.Stdout(buf))
func Stdout(stdout io.Writer) Option {
	return stdoutOpt{stdout: stdout}
}

type stderrOpt struct{ stderr io.Writer }

func (o stderrOpt) apply(cmd *Command) error {
	if o.stderr == nil {
		return errors.New("cannot set Stderr to nil")
	}

	cmd.stderr = o.stderr

	return nil
}

// Stderr is an [Option] that sets the Stderr for a [Command].
//
// Successive calls will simply overwrite any previous calls. Without this option
// the command will default to [os.Stderr].
//
// Subcommands cannot override the stderr of the root command, when retrieving the
// stdin with [Command.Stderr], the root's stderr will be returned.
//
//	// Set stderr to a temporary buffer
//	buf := &bytes.Buffer{}
//	cli.New("test", cli.Stderr(buf))
func Stderr(stderr io.Writer) Option {
	return stderrOpt{stderr: stderr}
}

type noColourOpt struct{ noColour bool }

func (o noColourOpt) apply(_ *Command) error {
	hue.Enabled(!o.noColour)

	return nil
}

// NoColour is an [Option] that disables all colour output from the [Command].
//
// CLI respects the values of $NO_COLOR and $FORCE_COLOR automatically so this need
// not be set for most applications.
//
// Setting this option takes precedence over all other colour configuration.
func NoColour(noColour bool) Option {
	return noColourOpt{noColour: noColour}
}

type shortOpt struct{ short string }

func (o shortOpt) apply(cmd *Command) error {
	if o.short == "" {
		return errors.New("cannot set command short description to an empty string")
	}

	cmd.short = strings.TrimSpace(o.short)

	return nil
}

// Short is an [Option] that sets the one line usage summary for a [Command].
//
// The one line usage will appear in the help text as well as alongside
// subcommands when they are listed.
//
// For consistency of formatting, all leading and trailing whitespace is stripped.
//
// Successive calls will simply overwrite any previous calls.
//
//	cli.New("rm", cli.Short("Delete files and directories"))
func Short(short string) Option {
	return shortOpt{short: short}
}

type longOpt struct{ long string }

func (o longOpt) apply(cmd *Command) error {
	if o.long == "" {
		return errors.New("cannot set command long description to an empty string")
	}

	cmd.long = strings.TrimSpace(o.long)

	return nil
}

// Long is an [Option] that sets the full description for a [Command].
//
// The long description will appear in the help text for a command. Users
// are responsible for wrapping the text at a sensible width.
//
// For consistency of formatting, all leading and trailing whitespace is stripped.
//
// Successive calls will simply overwrite any previous calls.
//
//	cli.New("rm", cli.Long("... lots of text here"))
func Long(long string) Option {
	return longOpt{long: long}
}

type exampleOpt struct {
	comment string
	command string
}

func (o exampleOpt) apply(cmd *Command) error {
	if o.comment == "" {
		return errors.New("example comment cannot be empty")
	}

	if o.command == "" {
		return errors.New("example command cannot be empty")
	}

	cmd.examples = append(cmd.examples, example(o))

	return nil
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
	return exampleOpt{comment: comment, command: command}
}

type runOpt struct {
	run func(ctx context.Context, cmd *Command) error
}

func (o runOpt) apply(cmd *Command) error {
	if o.run == nil {
		return errors.New("cannot set Run to nil")
	}

	cmd.run = o.run

	return nil
}

// Run is an [Option] that sets the run function for a [Command].
//
// The run function is the actual implementation of your command i.e. what you
// want it to do when invoked.
//
// Successive calls overwrite previous ones.
func Run(run func(ctx context.Context, cmd *Command) error) Option {
	return runOpt{run: run}
}

type overrideArgsOpt struct{ args []string }

func (o overrideArgsOpt) apply(cmd *Command) error {
	if o.args == nil {
		return errors.New("cannot set Args to nil")
	}

	cmd.rawArgs = o.args

	return nil
}

// OverrideArgs is an [Option] that sets the arguments for a [Command], overriding
// any arguments parsed from the command line.
//
// Without this option, the command will default to os.Args[1:], this option is particularly
// useful for testing.
//
// Successive calls override previous ones.
//
//	// Override arguments for testing
//	cli.New("test", cli.OverrideArgs([]string{"test", "me"}))
func OverrideArgs(args []string) Option {
	return overrideArgsOpt{args: args}
}

type versionOpt struct{ version string }

func (o versionOpt) apply(cmd *Command) error {
	cmd.version = o.version

	return nil
}

// Version is an [Option] that sets the version for a [Command].
//
// Without this option, the command defaults to a version of "dev".
//
//	cli.New("test", cli.Version("v1.2.3"))
func Version(version string) Option {
	return versionOpt{version: version}
}

type commitOpt struct{ commit string }

func (o commitOpt) apply(cmd *Command) error {
	cmd.commit = o.commit

	return nil
}

// Commit is an [Option] that sets the commit hash for a binary built with CLI. It is particularly
// useful for embedding rich version info into a binary using [ldflags].
//
// Without this option, the commit hash is simply omitted from the version info
// shown when -v/--version is called.
//
// If set to a non empty string, the commit hash will be shown.
//
//	cli.New("test", cli.Commit("b43fd2c"))
//
// [ldflags]: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
func Commit(commit string) Option {
	return commitOpt{commit: commit}
}

type buildDateOpt struct{ date string }

func (o buildDateOpt) apply(cmd *Command) error {
	cmd.buildDate = o.date

	return nil
}

// BuildDate is an [Option] that sets the build date for a binary built with CLI. It is particularly
// useful for embedding rich version info into a binary using [ldflags]
//
// Without this option, the build date is simply omitted from the version info
// shown when -v/--version is called.
//
// If set to a non empty string, the build date will be shown.
//
//	cli.New("test", cli.BuildDate("2024-07-06T10:37:30Z"))
//
// [ldflags]: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
func BuildDate(date string) Option {
	return buildDateOpt{date: date}
}

type subCommandsOpt struct{ builders []Builder }

// In Cobra the AddCommand method has to protect against a command adding itself
// as a subcommand, this is impossible in cli due to the functional options pattern, the
// root command will not exist as a variable inside the call to cli.New.
func (o subCommandsOpt) apply(cmd *Command) error {
	for _, builder := range o.builders {
		subcommand, err := builder()
		if err != nil {
			return fmt.Errorf("could not build subcommand: %w", err)
		}

		cmd.subcommands = append(cmd.subcommands, subcommand)
	}

	if name, found := anyDuplicates(cmd.subcommands...); found {
		return fmt.Errorf("subcommand %q already defined", name)
	}

	return nil
}

// SubCommands is an [Option] that attaches 1 or more subcommands to the command being configured.
//
// Sub commands must have unique names, any duplicates will result in an error.
//
// This option is additive and can be called as many times as desired, subcommands are
// effectively appended on every call.
func SubCommands(builders ...Builder) Option {
	return subCommandsOpt{builders: builders}
}

type flagOpt[T flag.Flaggable] struct {
	target  *T
	name    string
	usage   string
	options []FlagOption[T]
	short   rune
}

func (o flagOpt[T]) apply(cmd *Command) error {
	if _, ok := cmd.flags.Get(o.name); ok {
		return fmt.Errorf("flag %q already defined", o.name)
	}

	var flagCfg internalflag.Config[T]

	for _, option := range o.options {
		if err := option.apply(&flagCfg); err != nil {
			return fmt.Errorf("could not apply flag option: %w", err)
		}
	}

	f, err := internalflag.New(o.target, o.name, o.short, o.usage, flagCfg)
	if err != nil {
		return err
	}

	if err := internalflag.AddToSet(cmd.flags, f); err != nil {
		return fmt.Errorf("could not add flag %q to command %q: %w", o.name, cmd.name, err)
	}

	return nil
}

// Flag is an [Option] that adds a typed flag to a [Command], storing its value in a variable via its
// pointer 'target'.
//
// The variable is set when the flag is parsed during command execution. By default, the flag
// will assume the zero value for its type, the default may be provided explicitly using
// the [FlagDefault] option.
//
// If the default value is not the zero value for the type T, the flags usage message will
// show the default value in the commands help text.
//
// To add a long flag only (e.g. --delete with no -d option), pass [NoShortHand] for "short".
//
// Flags linked to slice values (e.g. []string) work by appending the passed values to the slice
// so multiple values may be given by repeat usage of the flag e.g. --items "one" --items "two".
//
//	// Add a force flag
//	var force bool
//	cli.New("rm", cli.Flag(&force, "force", 'f', "Force deletion without confirmation"))
func Flag[T flag.Flaggable](target *T, name string, short rune, usage string, options ...FlagOption[T]) Option {
	return flagOpt[T]{target: target, name: name, short: short, usage: usage, options: options}
}

type argOpt[T arg.Argable] struct {
	target  *T
	name    string
	usage   string
	options []ArgOption[T]
}

func (o argOpt[T]) apply(cmd *Command) error {
	var argCfg internalarg.Config[T]

	for _, option := range o.options {
		if err := option.apply(&argCfg); err != nil {
			return fmt.Errorf("could not apply arg option: %w", err)
		}
	}

	a, err := internalarg.New(o.target, o.name, o.usage, argCfg)
	if err != nil {
		return err
	}

	cmd.args = append(cmd.args, a)

	return nil
}

// Arg is an [Option] that adds a typed argument to a [Command], storing its value in a variable via its
// pointer 'target'.
//
// The variable is set when the argument is parsed during command execution.
//
// Args linked to slice values (e.g. []string) must be defined last as they eagerly consume
// all remaining command line arguments.
//
// The argument may be given a default value with the [ArgDefault] option. Without this option
// the argument will be required, i.e. failing to provide it on the command line is an error, but
// when a default is given and the value omitted on the command line, the default is used in
// its place.
//
//	// Add an int arg that defaults to 1
//	var number int
//	cli.New("add", cli.Arg(&number, "number", "Add a number", cli.ArgDefault(1)))
func Arg[T arg.Argable](p *T, name, usage string, options ...ArgOption[T]) Option {
	return argOpt[T]{target: p, name: name, usage: usage, options: options}
}

type argDefaultOpt[T arg.Argable] struct{ value T }

//nolint:unused // Satisfies the unexported ArgOption.apply method, staticcheck can't see across the interface.
func (o argDefaultOpt[T]) apply(cfg *internalarg.Config[T]) error {
	cfg.DefaultValue = &o.value
	return nil
}

// ArgDefault is a [cli.ArgOption] that sets the default value for a positional argument.
//
// By default, positional arguments are required, but by providing a default value
// via this option, you mark the argument as not required.
//
// If a default is given and the argument is not provided via the command line, the
// default is used in its place.
func ArgDefault[T arg.Argable](value T) ArgOption[T] {
	return argDefaultOpt[T]{value: value}
}

type envOpt[T flag.Flaggable] struct{ name string }

//nolint:unused // Satisfies the unexported FlagOption.apply method, staticcheck can't see across the interface.
func (o envOpt[T]) apply(cfg *internalflag.Config[T]) error {
	if o.name == "" {
		return errors.New("env var name cannot be empty")
	}

	cfg.EnvVar = o.name

	return nil
}

// Env is a [FlagOption] that associates an environment variable with a flag.
//
// When the flag is not explicitly set on the command line, CLI checks the named
// environment variable. If it is set and non-empty, its value is parsed using the
// same mechanism as command-line values. If it is not set or is empty, the flag
// retains its default value.
//
// For scalar flags, command-line values always take priority over environment variables.
// For slice and count flags, the environment variable provides a base value and any
// CLI flags accumulate on top.
//
// Slice flags accept comma-separated values:
//
//	MYTOOL_ITEMS='one,two,three'
//
//	var noApprove bool
//	cli.Flag(&noApprove, "no-approve", cli.NoShortHand, "Skip approval", cli.Env[bool]("MYTOOL_NO_APPROVE"))
func Env[T flag.Flaggable](name string) FlagOption[T] {
	return envOpt[T]{name: name}
}

type flagDefaultOpt[T flag.Flaggable] struct{ value T }

//nolint:unused // Satisfies the unexported FlagOption.apply method, staticcheck can't see across the interface.
func (o flagDefaultOpt[T]) apply(cfg *internalflag.Config[T]) error {
	cfg.DefaultValue = o.value
	return nil
}

// FlagDefault is a [cli.FlagOption] that sets the default value for command line flag.
//
// By default, a flag's default value is the zero value for its type. But using this
// option, you may set a non-zero default value that the flag should inherit if not
// provided on the command line.
func FlagDefault[T flag.Flaggable](value T) FlagOption[T] {
	return flagDefaultOpt[T]{value: value}
}

// anyDuplicates checks the list of commands for ones with duplicate names, if a duplicate
// is found, it's name and true are returned, else "", false.
func anyDuplicates(cmds ...*Command) (string, bool) {
	seen := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		if slices.Contains(seen, cmd.name) {
			return cmd.name, true
		}

		seen = append(seen, cmd.name)
	}

	return "", false
}

// ArgOption is a functional option for configuring an [Arg].
type ArgOption[T arg.Argable] interface {
	apply(cfg *internalarg.Config[T]) error
}

// FlagOption is a functional option for configuring a [Flag].
type FlagOption[T flag.Flaggable] interface {
	apply(cfg *internalflag.Config[T]) error
}
