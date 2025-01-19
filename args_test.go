package cli_test

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

func TestArgValidators(t *testing.T) {
	tests := []struct {
		name    string       // Identifier of the test case
		stdout  string       // Desired output to stdout
		stderr  string       // Desired output to stderr
		errMsg  string       // If we wanted an error, what should it say
		options []cli.Option // Options to apply to the command
		wantErr bool         // Whether we want an error
	}{
		{
			name: "anyargs",
			options: []cli.Option{
				cli.OverrideArgs([]string{"some", "args", "here"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from anyargs")
					return nil
				}),
				cli.Allow(cli.AnyArgs()),
			},
			wantErr: false,
			stdout:  "Hello from anyargs\n",
		},
		{
			name: "noargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				cli.Allow(cli.NoArgs()),
			},
			wantErr: false,
			stdout:  "Hello from noargs\n",
		},
		{
			name: "noargs fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"some", "args", "here"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				cli.Allow(cli.NoArgs()),
			},
			wantErr: true,
			errMsg:  "command test accepts no arguments but got [some args here]",
		},
		{
			name: "noargs subcommand",
			options: []cli.Option{
				cli.OverrideArgs([]string{"subb", "args", "here"}), // Note: subb is typo of sub
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from noargs")
					return nil
				}),
				cli.Allow(cli.NoArgs()),
				cli.SubCommands(
					func() (*cli.Command, error) {
						return cli.New(
							"sub",
							cli.Run(func(_ *cli.Command, _ []string) error { return nil }),
						)
					},
				),
			},
			wantErr: true,
			errMsg:  `unknown subcommand "subb" for command "test", available subcommands: [sub]`,
		},
		{
			name: "minargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"loads", "more", "args", "here"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				cli.Allow(cli.MinArgs(3)),
			},
			wantErr: false,
			stdout:  "Hello from minargs\n",
		},
		{
			name: "minargs fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"only", "two"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from minargs")
					return nil
				}),
				cli.Allow(cli.MinArgs(3)),
			},
			wantErr: true,
			errMsg:  "command test requires at least 3 arguments, but got 2: [only two]",
		},
		{
			name: "maxargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				cli.Allow(cli.MaxArgs(2)),
			},
			wantErr: false,
			stdout:  "Hello from maxargs\n",
		},
		{
			name: "maxargs fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"loads", "of", "args", "here", "wow", "so", "many"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from maxargs")
					return nil
				}),
				cli.Allow(cli.MaxArgs(3)),
			},
			wantErr: true,
			errMsg:  "command test has a limit of 3 argument(s), but got 7: [loads of args here wow so many]",
		},
		{
			name: "exactargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				cli.Allow(cli.ExactArgs(2)),
			},
			wantErr: false,
			stdout:  "Hello from exactargs\n",
		},
		{
			name: "exactargs fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"not", "three", "but", "four"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from exactargs")
					return nil
				}),
				cli.Allow(cli.ExactArgs(3)),
			},
			wantErr: true,
			errMsg:  "command test requires exactly 3 arguments, but got 4: [not three but four]",
		},
		{
			name: "betweenargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"two", "args"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(1, 4)),
			},
			wantErr: false,
			stdout:  "Hello from betweenargs\n",
		},
		{
			name: "betweenargs fail high",
			options: []cli.Option{
				cli.OverrideArgs([]string{"not", "three", "but", "more", "than", "four"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(1, 4)),
			},
			wantErr: true,
			errMsg:  "command test requires between 1 and 4 arguments, but got 6: [not three but more than four]",
		},
		{
			name: "betweenargs fail low",
			options: []cli.Option{
				cli.OverrideArgs([]string{"not", "three"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from betweenargs")
					return nil
				}),
				cli.Allow(cli.BetweenArgs(3, 5)),
			},
			wantErr: true,
			errMsg:  "command test requires between 3 and 5 arguments, but got 2: [not three]",
		},
		{
			name: "validargs pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"valid", "args", "only"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				cli.Allow(cli.ValidArgs([]string{"only", "valid", "args"})), // Order doesn't matter
			},

			wantErr: false,
			stdout:  "Hello from validargs\n",
		},
		{
			name: "validargs fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"valid", "args", "only", "bad"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from validargs")
					return nil
				}),
				cli.Allow(cli.ValidArgs([]string{"only", "valid", "args"})),
			},
			wantErr: true,
			errMsg:  "command test got an invalid argument bad, expected one of [only valid args]",
		},
		{
			name: "combine pass",
			options: []cli.Option{
				cli.OverrideArgs([]string{"four", "args", "all", "valid"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				cli.Allow(
					cli.Combine(
						cli.ExactArgs(4),
						cli.ValidArgs([]string{"valid", "all", "four", "args"}),
					),
				),
			},
			wantErr: false,
			stdout:  "Hello from combine\n",
		},
		{
			name: "combine fail",
			options: []cli.Option{
				cli.OverrideArgs([]string{"valid", "args", "only", "bad", "five"}),
				cli.Run(func(cmd *cli.Command, _ []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from combine")
					return nil
				}),
				cli.Allow(
					cli.Combine(
						cli.BetweenArgs(1, 4),
						cli.ValidArgs([]string{"only", "valid", "args", "here"}),
					),
				),
			},
			wantErr: true,
			errMsg:  "command test requires between 1 and 4 arguments, but got 5: [valid args only bad five]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			// Test specific overrides to the options in the table
			options := []cli.Option{cli.Stdout(stdout), cli.Stderr(stderr)}

			cmd, err := cli.New("test", slices.Concat(tt.options, options)...)
			test.Ok(t, err)

			err = cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			if tt.wantErr {
				test.Equal(t, err.Error(), tt.errMsg)
			}

			test.Equal(t, stdout.String(), tt.stdout)
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}
