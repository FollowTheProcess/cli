package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

func TestExecute(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	testCmd := &cli.Command{
		Run: func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(cmd.Stdout, "Oooh look, it ran, here are some args: %v\n", args)
			return nil
		},
		Stdout: stdout,
		Stderr: stderr,
		Name:   "test",
		Short:  "A simple test command",
		Long:   "Much longer description blah blah blah",
	}

	err := testCmd.Execute([]string{"arg1", "arg2", "arg3"})
	test.Ok(t, err)

	want := "Oooh look, it ran, here are some args: [arg1 arg2 arg3]\n"
	test.Equal(t, stdout.String(), want)
}
