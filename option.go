package cli

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/FollowTheProcess/cli/internal/flag"
)

// NoShortHand should be passed as the "short" argument to [Flag] if the desired flag
// should be the long hand version only e.g. --count, not -c/--count.
const NoShortHand = flag.NoShortHand

// Flaggable is a type constraint that defines any type capable of being parsed as a command line flag.
//
// It's worth noting that the complete set of supported types is wider than this constraint appears
// as e.g. a [time.Duration] is actually just an int64 underneath, likewise a [net.IP] is actually just []byte.
type Flaggable flag.Flaggable

// FlagCount is a type used for a flag who's job is to increment a counter, e.g. a "verbosity"
// flag may be passed "-vvv" which should increase the verbosity level to 3.
type FlagCount flag.Count

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
	stdin         io.Reader
	stdout        io.Writer
	stderr        io.Writer
	run           func(cmd *Command, args []string) error
	flags         *flag.Set
	versionFunc   func(cmd *Command) error
	parent        *Command
	argValidator  ArgValidator
	name          string
	short         string
	long          string
	version       string
	commit        string
	buildDate     string
	examples      []example
	args          []string
	subcommands   []*Command
	helpCalled    bool
	versionCalled bool
}

// build builds an returns a Command from the config, applying validation
// to the whole thing.
func (c *config) build() *Command {
	cmd := &Command{
		stdin:         c.stdin,
		stdout:        c.stdout,
		stderr:        c.stderr,
		run:           c.run,
		flags:         c.flags,
		versionFunc:   c.versionFunc,
		parent:        c.parent,
		argValidator:  c.argValidator,
		name:          c.name,
		short:         c.short,
		long:          c.long,
		version:       c.version,
		commit:        c.commit,
		buildDate:     c.buildDate,
		examples:      c.examples,
		args:          c.args,
		subcommands:   c.subcommands,
		helpCalled:    c.helpCalled,
		versionCalled: c.versionCalled,
	}

	// Loop through each subcommand and set this command as their immediate parent
	for _, subcommand := range cmd.subcommands {
		subcommand.parent = cmd
	}

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
// The long description will appear in the help text for a command. Users
// are responsible for wrapping the text at a sensible width.
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
	f := func(cfg *config) error {
		if comment == "" {
			return errors.New("example comment cannot be empty")
		}
		if command == "" {
			return errors.New("example command cannot be empty")
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
	f := func(cfg *config) error {
		if commit == "" {
			return errors.New("cannot set Commit to an empty string")
		}
		cfg.commit = commit
		return nil
	}
	return option(f)
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
	f := func(cfg *config) error {
		if date == "" {
			return errors.New("cannot set BuildDate to an empty string")
		}
		cfg.buildDate = date
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
func SubCommands(builders ...Builder) Option {
	// Note: In Cobra the AddCommand method has to protect against a command adding itself
	// as a subcommand, this is impossible in cli due to the functional options pattern, the
	// root command will not exist as a variable inside the call to cli.New.

	f := func(cfg *config) error {
		// Add the subcommands to the command this is being called on
		for _, builder := range builders {
			subcommand, err := builder()
			if err != nil {
				return fmt.Errorf("could not build subcommand: %w", err)
			}
			cfg.subcommands = append(cfg.subcommands, subcommand)
		}

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
// To add a long flag only (e.g. --delete with no -d option), pass [NoShortHand] for "short".
//
//	// Add a force flag
//	var force bool
//	cli.New("rm", cli.Flag(&force, "force", 'f', false, "Force deletion without confirmation"))
func Flag[T Flaggable](p *T, name string, short rune, value T, usage string) Option {
	f := func(cfg *config) error {
		if _, ok := cfg.flags.Get(name); ok {
			return fmt.Errorf("flag %q already defined", name)
		}

		f, err := flag.New(p, name, short, value, usage)
		if err != nil {
			return err
		}

		if err := flag.AddToSet(cfg.flags, f); err != nil {
			return fmt.Errorf("could not add flag %q to command %q: %w", name, cfg.name, err)
		}

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
