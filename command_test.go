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

func TestExampleString(t *testing.T) {
	tests := []struct {
		name    string
		example cli.Example
		want    string
	}{
		{
			name:    "empty",
			example: cli.Example{},
			want:    "",
		},
		{
			name:    "only command",
			example: cli.Example{Command: "run this program --once"},
			want:    "$ run this program --once",
		},
		{
			name:    "only comment",
			example: cli.Example{Comment: "Run the program once"},
			want:    "# Run the program once",
		},
		{
			name: "both",
			example: cli.Example{
				Comment: "Run the program once",
				Command: "run this program --once",
			},
			want: "# Run the program once\n$ run this program --once",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.example.String(), tt.want)
		})
	}
}
