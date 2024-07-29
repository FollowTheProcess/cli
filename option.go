package cli

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/spf13/pflag"
)

// Option is a functional option for configuring a [Command].
type Option interface {
	// Apply the option to the config, returning an error if the
	// option cannot be applied for whatever reason.
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
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
	run          func(cmd *Command, args []string) error
	flags        *pflag.FlagSet
	versionFunc  func(cmd *Command) error
	parent       *Command
	argValidator ArgValidator
	name         string
	short        string
	long         string
	version      string
	examples     []example
	args         []string
	subcommands  []*Command
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
		argValidator: c.argValidator,
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
		if stdin == nil {
			return errors.New("cannot set Stdin to nil")
		}
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
		if stdout == nil {
			return errors.New("cannot set Stdout to nil")
		}
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
		if stderr == nil {
			return errors.New("cannot set Stderr to nil")
		}
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
		if short == "" {
			return errors.New("cannot set command short description to an empty string")
		}
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
		if long == "" {
			return errors.New("cannot set command long description to an empty string")
		}
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
	// TODO: Make sure both comment and command are not empty, then can ditch the
	// complexity in example.String()
	f := func(cfg *config) error {
		errs := make([]error, 0, 2) //nolint:mnd // 2 here is because we have two arguments
		if comment == "" {
			errs = append(errs, errors.New("example comment cannot be empty"))
		}
		if command == "" {
			errs = append(errs, errors.New("example command cannot be empty"))
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
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
		if run == nil {
			return errors.New("cannot set Run to nil")
		}
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
		if args == nil {
			return errors.New("cannot set Args to nil")
		}
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
		if version == "" {
			return errors.New("cannot set Version to an empty string")
		}
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
		if fn == nil {
			return errors.New("cannot set VersionFunc to nil")
		}
		cfg.versionFunc = fn
		return nil
	}
	return option(f)
}

// SubCommands is an [Option] that attaches 1 or more subcommands to the command being configured.
//
// Sub commands must have unique names, any duplicates will result in an error.
//
// This option is additive and can be called as many times as desired, subcommands are
// effectively appended on every call.
func SubCommands(subcommands ...*Command) Option {
	// Note: In Cobra the AddCommand method has to protect against a command adding itself
	// as a subcommand, this is impossible in cli due to the functional options pattern, the
	// root command will not exist as a variable inside the call to cli.New.

	f := func(cfg *config) error {
		// Add the subcommands to the command this is being called on
		cfg.subcommands = append(cfg.subcommands, subcommands...)

		// Any duplicates in the list of subcommands (by name) is an error
		if name, found := anyDuplicates(cfg.subcommands...); found {
			return fmt.Errorf("subcommand %q already defined", name)
		}
		return nil
	}
	return option(f)
}

// Allow is an [Option] that allows for validating positional arguments to a [Command].
//
// You provide a validator function that returns an error if it encounters invalid arguments, and it will
// be run for you, passing in the non-flag arguments to the [Command] that was called.
//
// Successive calls overwrite previous ones, use [Combine] to compose multiple validators.
//
//	// No positional arguments allowed
//	cli.New("test", cli.Allow(cli.NoArgs()))
func Allow(validator ArgValidator) Option {
	f := func(cfg *config) error {
		if validator == nil {
			return errors.New("cannot set Allow to a nil ArgValidator")
		}
		cfg.argValidator = validator
		return nil
	}
	return option(f)
}

// Flag is an [Option] that adds a flag to a [Command], storing its value in a variable.
//
// The variable is set when the flag is parsed during command execution. The value provided
// in the call to [Flag] is used as the default value.
//
// To add a long flag only (e.g. --delete with no -d option), simply pass "" for short.
//
//	// Add a force flag
//	var force bool
//	cli.New("rm", cli.Flag(&force, "force", "f", false, "Force deletion without confirmation"))
func Flag[T flag.Flaggable](p *T, name string, short string, value T, usage string) Option {
	f := func(cfg *config) error {
		if name == "" {
			return errors.New("flag names must not be empty")
		}
		if err := validateFlagName(name); err != nil {
			return fmt.Errorf("invalid flag name: %w", err)
		}
		if cfg.flags.Lookup(name) != nil {
			return fmt.Errorf("flag %q already defined", name)
		}

		// len(short) > 1 means an error, shorthand must be a single character
		if length := utf8.RuneCountInString(short); length > 1 {
			return fmt.Errorf("shorthand for flag %q must be a single ASCII letter, got %q which has %d letters", name, short, length)
		}

		if short != "" {
			// Shorthand must be a valid ASCII letter
			char, _ := utf8.DecodeRuneInString(short)
			if char == utf8.RuneError || char > unicode.MaxASCII || !unicode.IsLetter(char) {
				return fmt.Errorf(
					"shorthand for flag %q is an invalid character, must be a single ASCII letter, got %q",
					name,
					string(char),
				)
			}
		}

		// Short is now either "" or a single letter
		flag := flag.New(p, name, short, value, usage)
		var defVal string
		if flag.Type() == "bool" {
			defVal = "true"
		}

		// Annoyingly pflag does the same checks we've done above but will panic on error, hopefully
		// the above checks will come in handy when I implement my own flag parsing
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

// validateFlagName ensures a flag name is valid, returning an error if it's not.
//
// Flags names must be all lower case ASCII letters, a hypen separator is allowed e.g. "set-default"
// but this must be in between letters, not leading or trailing.
func validateFlagName(name string) error {
	before, after, found := strings.Cut(name, "-")

	// Hyphen must be in between "words" like "set-default"
	// we can't have "-default" or "default-"
	if found && after == "" {
		return fmt.Errorf("trailing hyphen: %q", name)
	}

	if found && before == "" {
		return fmt.Errorf("leading hyphen: %q", name)
	}
	for _, char := range name {
		// Only ASCII characters allowed
		if char > unicode.MaxASCII {
			return fmt.Errorf("non ascii character: %q", string(char))
		}
		// Only non-letter character allowed is a hyphen
		if !unicode.IsLetter(char) && char != '-' {
			return fmt.Errorf("not ascii letter: %q", string(char))
		}
		// Any upper case letters are not allowed
		if unicode.IsLetter(char) && !unicode.IsLower(char) {
			return fmt.Errorf("upper case character %q", string(char))
		}
	}

	return nil
}
