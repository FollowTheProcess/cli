package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

func TestArgValidators(t *testing.T) {
	tests := []struct {
		cmd     *cli.Command // The command undkzer test
		name    string       // Identifier of the test case
		stdout  string       // Desired output to stdout
		stderr  string       // Desired output to stderr
		errMsg  string       // If we wanted an error, what should it say
		wantErr bool         // Whether we want an error
	}{
		{
			name: "anyargs",
			cmd: cli.New(
				"anyargs",
				cli.Args([]string{"some", "args", "here"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from anyargs")
					return nil
				}),
				cli.Allow(cli.AnyArgs),
			),
			wantErr: false,
			stdout:  "Hello from anyargs\n",
		},
		{
			name: "noargs pass",
			cmd: cli.New(
				"noargs",
				cli.Args([]string{}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				cli.Allow(cli.NoArgs),
			),
			wantErr: false,
			stdout:  "Hello from noargs\n",
		},
		{
			name: "noargs fail",
			cmd: cli.New(
				"noargs",
				cli.Args([]string{"arg1", "arg2", "arg3"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				cli.Allow(cli.NoArgs),
			),
			wantErr: true,
			errMsg:  "command noargs accepts no arguments but got [arg1 arg2 arg3]",
		},
		{
			name: "minargs pass",
			cmd: cli.New(
				"minargs",
				cli.Args([]string{"loads", "more", "args", "here"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				cli.Allow(cli.MinArgs(3)),
			),
			wantErr: false,
			stdout:  "Hello from minargs\n",
		},
		{
			name: "minargs fail",
			cmd: cli.New(
				"minargs",
				cli.Args([]string{"only", "two"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				cli.Allow(cli.MinArgs(3)),
			),
			wantErr: true,
			errMsg:  "command minargs requires at least 3 arguments, but got 2: [only two]",
		},
		{
			name: "maxargs pass",
			cmd: cli.New(
				"maxargs",
				cli.Args([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				cli.Allow(cli.MaxArgs(2)),
			),
			wantErr: false,
			stdout:  "Hello from maxargs\n",
		},
		{
			name: "maxargs fail",
			cmd: cli.New(
				"maxargs",
				cli.Args([]string{"loads", "of", "args", "here", "wow", "so", "many"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				cli.Allow(cli.MaxArgs(3)),
			),
			wantErr: true,
			errMsg:  "command maxargs has a limit of 3 arguments, but got 7: [loads of args here wow so many]",
		},
		{
			name: "exactargs pass",
			cmd: cli.New(
				"exactargs",
				cli.Args([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				cli.Allow(cli.ExactArgs(2)),
			),
			wantErr: false,
			stdout:  "Hello from exactargs\n",
		},
		{
			name: "exactargs fail",
			cmd: cli.New(
				"exactargs",
				cli.Args([]string{"not", "three", "but", "four"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				cli.Allow(cli.ExactArgs(3)),
			),
			wantErr: true,
			errMsg:  "command exactargs requires exactly 3 arguments, but got 4: [not three but four]",
		},
		{
			name: "betweenargs pass",
			cmd: cli.New(
				"betweenargs",
				cli.Args([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(1, 4)),
			),
			wantErr: false,
			stdout:  "Hello from betweenargs\n",
		},
		{
			name: "betweenargs fail high",
			cmd: cli.New(
				"betweenargs",
				cli.Args([]string{"not", "three", "but", "more", "than", "four"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(1, 4)),
			),
			wantErr: true,
			errMsg:  "command betweenargs requires between 1 and 4 arguments, but got 6: [not three but more than four]",
		},
		{
			name: "betweenargs fail low",
			cmd: cli.New(
				"betweenargs",
				cli.Args([]string{"not", "three"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(3, 5)),
			),
			wantErr: true,
			errMsg:  "command betweenargs requires between 3 and 5 arguments, but got 2: [not three]",
		},
		{
			name: "validargs pass",
			cmd: cli.New(
				"validargs",
				cli.Args([]string{"valid", "args", "only"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				cli.Allow(
					cli.ValidArgs([]string{"only", "valid", "args"}),
				), // Order doesn't matter
			),
			wantErr: false,
			stdout:  "Hello from validargs\n",
		},
		{
			name: "validargs fail",
			cmd: cli.New(
				"validargs",
				cli.Args([]string{"valid", "args", "only", "bad"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				cli.Allow(cli.ValidArgs([]string{"only", "valid", "args"})),
			),
			wantErr: true,
			errMsg:  "command validargs got an invalid argument bad, expected one of [only valid args]",
		},
		{
			name: "combine pass",
			cmd: cli.New(
				"combine",
				cli.Args([]string{"four", "args", "all", "valid"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				cli.Allow(
					cli.Combine(
						cli.ExactArgs(4),
						cli.ValidArgs([]string{"valid", "all", "four", "args"}),
					),
				),
			),
			wantErr: false,
			stdout:  "Hello from combine\n",
		},
		{
			name: "combine fail",
			cmd: cli.New(
				"combine",
				cli.Args([]string{"valid", "args", "only", "bad", "five"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				cli.Allow(
					cli.Combine(
						cli.BetweenArgs(1, 4),
						cli.ValidArgs([]string{"only", "valid", "args", "here"}),
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

			cli.Stderr(stderr)(tt.cmd)
			cli.Stdout(stdout)(tt.cmd)

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
