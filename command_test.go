package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		cmd        *cli.Command                         // The command under test
		customiser func(t *testing.T, cmd *cli.Command) // An optional function to customise the command e.g. add flags
		name       string                               // Identifier of the test case
		stdout     string                               // Desired output to stdout
		stderr     string                               // Desired output to stderr
		wantErr    bool                                 // Whether we want an error
	}{
		{
			name: "simple",
			cmd: cli.New(
				"test",
				cli.Args([]string{"arg1", "arg2", "arg3"}),
			),
			stdout:  "Hello from test\n",
			wantErr: false,
		},
		{
			name: "simple with flags",
			cmd: cli.New(
				"test",
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "Oooh look, it ran, here are some args: %v\n", args)
					force, err := cmd.Flags().GetBool("force")
					test.Ok(t, err)
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				cli.Args([]string{"arg1", "arg2", "--force"}),
			),
			customiser: func(t *testing.T, cmd *cli.Command) {
				t.Helper()
				cmd.Flags().BoolP("force", "f", false, "Force something")
			},
			stdout:  "Oooh look, it ran, here are some args: [arg1 arg2]\n--force was: true\n",
			wantErr: false,
		},
		{
			name: "bad flag",
			cmd: cli.New(
				"test",
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "Oooh look, it ran, here are some args: %v\n", args)
					force, err := cmd.Flags().GetBool("force")
					test.Ok(t, err)
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				cli.Args([]string{"arg1", "arg2", "arg3", "-]force"}),
			),
			customiser: func(t *testing.T, cmd *cli.Command) {
				t.Helper()
				cmd.Flags().BoolP("force", "f", false, "Force something")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			cli.Stderr(stderr)(tt.cmd)
			cli.Stdout(stdout)(tt.cmd)

			// Customise if it's set
			if tt.customiser != nil {
				tt.customiser(t, tt.cmd)
			}

			err := tt.cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			test.Equal(t, stdout.String(), tt.stdout)
		})
	}
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
