package command

import "fmt"

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
