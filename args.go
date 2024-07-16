package cli

import (
	"fmt"
	"slices"
)

// AnyArgs is a positional argument validator that allows any arbitrary args,
// it never returns an error.
//
// This is the default argument validator on a [Command] instantiated with cli.New.
func AnyArgs(cmd *Command, args []string) error {
	return nil
}

// NoArgs is a positional argument validator that does not allow any arguments,
// it returns an error if there are any arguments.
func NoArgs(cmd *Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("command %s accepts no arguments but got %v", cmd.name, args)
	}
	return nil
}

// MinArgs is a positional argument validator that requires at least n arguments.
func MinArgs(n int) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf(
				"command %s requires at least %d arguments, but got %d: %v",
				cmd.name,
				n,
				len(args),
				args,
			)
		}
		return nil
	}
}

// MaxArgs is a positional argument validator that returns an error if there are more than n arguments.
func MaxArgs(n int) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		if len(args) > n {
			return fmt.Errorf(
				"command %s has a limit of %d arguments, but got %d: %v",
				cmd.name,
				n,
				len(args),
				args,
			)
		}
		return nil
	}
}

// ExactArgs is a positional argument validator that allows exactly n args, any more
// or less will return an error.
func ExactArgs(n int) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf(
				"command %s requires exactly %d arguments, but got %d: %v",
				cmd.name,
				n,
				len(args),
				args,
			)
		}
		return nil
	}
}

// BetweenArgs is a positional argument validator that allows between min and max arguments (inclusive),
// any outside that range will return an error.
func BetweenArgs(min, max int) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		nArgs := len(args)
		if nArgs < min || nArgs > max {
			return fmt.Errorf(
				"command %s requires between %d and %d arguments, but got %d: %v",
				cmd.name,
				min,
				max,
				nArgs,
				args,
			)
		}
		return nil
	}
}

// ValidArgs is a positional argument validator that only allows arguments that are contained in
// the valid slice. If any non-valid arguments are seen, an error will be returned.
func ValidArgs(valid []string) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		for _, arg := range args {
			if !slices.Contains(valid, arg) {
				return fmt.Errorf(
					"command %s got an invalid argument %s, expected one of %v",
					cmd.name,
					arg,
					valid,
				)
			}
		}
		return nil
	}
}

// Combine allows multiple positional argument validators to be composed together.
//
// The first validator to fail will be the one that returns the error.
func Combine(
	validators ...func(cmd *Command, args []string) error,
) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		for _, validator := range validators {
			if err := validator(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
