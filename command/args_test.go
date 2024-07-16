package command_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/FollowTheProcess/cli/command"
	"github.com/FollowTheProcess/test"
)

func TestArgValidators(t *testing.T) {
	tests := []struct {
		cmd     *command.Command // The command undkzer test
		name    string           // Identifier of the test case
		stdout  string           // Desired output to stdout
		stderr  string           // Desired output to stderr
		errMsg  string           // If we wanted an error, what should it say
		wantErr bool             // Whether we want an error
	}{
		{
			name: "anyargs",
			cmd: command.New(
				"anyargs",
				command.Args([]string{"some", "args", "here"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from anyargs")
					return nil
				}),
				command.Allow(command.AnyArgs),
			),
			wantErr: false,
			stdout:  "Hello from anyargs\n",
		},
		{
			name: "noargs pass",
			cmd: command.New(
				"noargs",
				command.Args([]string{}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				command.Allow(command.NoArgs),
			),
			wantErr: false,
			stdout:  "Hello from noargs\n",
		},
		{
			name: "noargs fail",
			cmd: command.New(
				"noargs",
				command.Args([]string{"arg1", "arg2", "arg3"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				command.Allow(command.NoArgs),
			),
			wantErr: true,
			errMsg:  "command noargs accepts no arguments but got [arg1 arg2 arg3]",
		},
		{
			name: "minargs pass",
			cmd: command.New(
				"minargs",
				command.Args([]string{"loads", "more", "args", "here"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				command.Allow(command.MinArgs(3)),
			),
			wantErr: false,
			stdout:  "Hello from minargs\n",
		},
		{
			name: "minargs fail",
			cmd: command.New(
				"minargs",
				command.Args([]string{"only", "two"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				command.Allow(command.MinArgs(3)),
			),
			wantErr: true,
			errMsg:  "command minargs requires at least 3 arguments, but got 2: [only two]",
		},
		{
			name: "maxargs pass",
			cmd: command.New(
				"maxargs",
				command.Args([]string{"two", "args"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				command.Allow(command.MaxArgs(2)),
			),
			wantErr: false,
			stdout:  "Hello from maxargs\n",
		},
		{
			name: "maxargs fail",
			cmd: command.New(
				"maxargs",
				command.Args([]string{"loads", "of", "args", "here", "wow", "so", "many"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				command.Allow(command.MaxArgs(3)),
			),
			wantErr: true,
			errMsg:  "command maxargs has a limit of 3 arguments, but got 7: [loads of args here wow so many]",
		},
		{
			name: "exactargs pass",
			cmd: command.New(
				"exactargs",
				command.Args([]string{"two", "args"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				command.Allow(command.ExactArgs(2)),
			),
			wantErr: false,
			stdout:  "Hello from exactargs\n",
		},
		{
			name: "exactargs fail",
			cmd: command.New(
				"exactargs",
				command.Args([]string{"not", "three", "but", "four"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				command.Allow(command.ExactArgs(3)),
			),
			wantErr: true,
			errMsg:  "command exactargs requires exactly 3 arguments, but got 4: [not three but four]",
		},
		{
			name: "betweenargs pass",
			cmd: command.New(
				"betweenargs",
				command.Args([]string{"two", "args"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				command.Allow(command.BetweenArgs(1, 4)),
			),
			wantErr: false,
			stdout:  "Hello from betweenargs\n",
		},
		{
			name: "betweenargs fail high",
			cmd: command.New(
				"betweenargs",
				command.Args([]string{"not", "three", "but", "more", "than", "four"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				command.Allow(command.BetweenArgs(1, 4)),
			),
			wantErr: true,
			errMsg:  "command betweenargs requires between 1 and 4 arguments, but got 6: [not three but more than four]",
		},
		{
			name: "betweenargs fail low",
			cmd: command.New(
				"betweenargs",
				command.Args([]string{"not", "three"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				command.Allow(command.BetweenArgs(3, 5)),
			),
			wantErr: true,
			errMsg:  "command betweenargs requires between 3 and 5 arguments, but got 2: [not three]",
		},
		{
			name: "validargs pass",
			cmd: command.New(
				"validargs",
				command.Args([]string{"valid", "args", "only"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				command.Allow(
					command.ValidArgs([]string{"only", "valid", "args"}),
				), // Order doesn't matter
			),
			wantErr: false,
			stdout:  "Hello from validargs\n",
		},
		{
			name: "validargs fail",
			cmd: command.New(
				"validargs",
				command.Args([]string{"valid", "args", "only", "bad"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				command.Allow(command.ValidArgs([]string{"only", "valid", "args"})),
			),
			wantErr: true,
			errMsg:  "command validargs got an invalid argument bad, expected one of [only valid args]",
		},
		{
			name: "combine pass",
			cmd: command.New(
				"combine",
				command.Args([]string{"four", "args", "all", "valid"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				command.Allow(
					command.Combine(
						command.ExactArgs(4),
						command.ValidArgs([]string{"valid", "all", "four", "args"}),
					),
				),
			),
			wantErr: false,
			stdout:  "Hello from combine\n",
		},
		{
			name: "combine fail",
			cmd: command.New(
				"combine",
				command.Args([]string{"valid", "args", "only", "bad", "five"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				command.Allow(
					command.Combine(
						command.BetweenArgs(1, 4),
						command.ValidArgs([]string{"only", "valid", "args", "here"}),
					),
				),
			),
			wantErr: true,
			errMsg:  "command combine requires between 1 and 4 arguments, but got 5: [valid args only bad five]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			command.Stderr(stderr)(tt.cmd)
			command.Stdout(stdout)(tt.cmd)

			err := tt.cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			if tt.wantErr {
				test.Equal(t, err.Error(), tt.errMsg)
			}

			test.Equal(t, stdout.String(), tt.stdout)
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}
