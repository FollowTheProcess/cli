package cli_test

import (
	"bytes"
	"errors"
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
				// Set flags in the customiser
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
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}

func TestHelp(t *testing.T) {
	tests := []struct {
		cmd     *cli.Command // The command under test
		name    string       // Identifier of the test case
		golden  string       // The name of the file relative to testdata containing to expected output
		wantErr bool         // Whether we want an error
		debug   bool         // Whether or not to print the produced help text to stderr (useful for debugging)
	}{
		{
			name: "default long",
			cmd: cli.New(
				"test",
				cli.Args([]string{"--help"}),
			),
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "default short",
			cmd: cli.New(
				"test",
				cli.Args([]string{"-h"}),
			),
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "with examples",
			cmd: cli.New(
				"test",
				cli.Args([]string{"--help"}),
				cli.Examples(
					cli.Example{Comment: "Do a thing", Command: "test do thing --now"},
					cli.Example{
						Comment: "Do a different thing",
						Command: "test do thing --different",
					},
				),
			),
			golden:  "with-examples.txt",
			wantErr: false,
		},
		{
			name: "with full description",
			cmd: cli.New(
				"test",
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
			),
			golden:  "full.txt",
			wantErr: false,
		},
		{
			name: "with no description",
			cmd: cli.New(
				"test",
				cli.Args([]string{"--help"}),
				cli.Short(""),
				cli.Long(""),
			),
			golden:  "no-about.txt",
			wantErr: false,
		},
		{
			name: "with subcommands",
			cmd: cli.New(
				"test",
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(
					cli.New("sub1", cli.Short("Do one thing")),
					cli.New("sub2", cli.Short("Do another thing")),
				),
			),
			golden:  "subcommands.txt",
			wantErr: false,
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

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// Show the help output, can aid debugging
			if tt.debug {
				fmt.Print(stderr.String())
			}

			// --help output should be as per the golden file
			test.File(t, stderr.String(), tt.golden)
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string       // Name of the test case
		cmd     *cli.Command // Command under test
		stderr  string       // Expected output to stderr
		wantErr bool         // Whether we want an error or not
	}{
		{
			name:    "default long",
			cmd:     cli.New("version-test", cli.Args([]string{"--version"})),
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name:    "default short",
			cmd:     cli.New("version-test", cli.Args([]string{"-v"})),
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "custom version",
			cmd: cli.New(
				"version-test",
				cli.Args([]string{"--version"}),
				cli.Version("v1.2.3"),
			),
			stderr:  "version-test, version: v1.2.3\n",
			wantErr: false,
		},
		{
			name: "custom versionFunc",
			cmd: cli.New(
				"version-test",
				cli.Args([]string{"--version"}),
				cli.VersionFunc(func(cmd *cli.Command) error {
					fmt.Fprintln(cmd.Stderr(), "Do something custom here")
					return nil
				}),
			),
			stderr:  "Do something custom here\n",
			wantErr: false,
		},
		{
			name: "return error",
			cmd: cli.New(
				"version-test",
				cli.Args([]string{"--version"}),
				cli.VersionFunc(func(cmd *cli.Command) error { return errors.New("Uh oh!") }),
			),
			wantErr: true,
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

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// --version output should be as desired
			test.Equal(t, stderr.String(), tt.stderr)
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
			want:    "\n  $ run this program --once\n",
		},
		{
			name:    "only comment",
			example: cli.Example{Comment: "Run the program once"},
			want:    "\n  # Run the program once\n",
		},
		{
			name: "both",
			example: cli.Example{
				Comment: "Run the program once",
				Command: "run this program --once",
			},
			want: "\n  # Run the program once\n  $ run this program --once\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.example.String(), tt.want)
		})
	}
}
