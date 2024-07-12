package cli

import (
	"io"
)

// Option is a functional option for configuring a [Command].
type Option func(*Command)

// Stdin is an [Option] that sets the Stdin for a [Command].
func Stdin(stdin io.Reader) Option {
	return func(cmd *Command) {
		cmd.stdin = stdin
	}
}

// Stdout is an [Option] that sets the Stdout for a [Command].
func Stdout(stdout io.Writer) Option {
	return func(cmd *Command) {
		cmd.stdout = stdout
	}
}

// Stderr is an [Option] that sets the Stderr for a [Command].
func Stderr(stderr io.Writer) Option {
	return func(cmd *Command) {
		cmd.stderr = stderr
	}
}

// Short is an [Option] that sets the one line usage summary for a [Command].
func Short(short string) Option {
	return func(cmd *Command) {
		cmd.short = short
	}
}

// Long is an [Option] that sets the full description for a [Command].
func Long(long string) Option {
	return func(cmd *Command) {
		cmd.long = long
	}
}

// Examples is an [Option] that sets the examples for a [Command].
func Examples(examples ...Example) Option {
	return func(cmd *Command) {
		cmd.example = examples
	}
}

// Run is an [Option] that sets the run function for a [Command].
//
// The run function is the actual implementation of your command i.e. what you
// want it to do when invoked.
func Run(run func(cmd *Command, args []string) error) Option {
	return func(cmd *Command) {
		cmd.run = run
	}
}

// Args is an [Option] that sets the arguments for a [Command].
func Args(args []string) Option {
	return func(cmd *Command) {
		cmd.args = args
	}
}

// HelpFunc is an [Option] that allows for a custom implementation of the -h/--help flag.
//
// A [Command] will have a default implementation of this function that prints a default
// format of the help to [os.Stderr].
func HelpFunc(fn func(cmd *Command) error) Option {
	return func(cmd *Command) {
		cmd.helpFunc = fn
	}
}
