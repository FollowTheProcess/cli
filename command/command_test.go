package command_test

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/FollowTheProcess/cli/command"
	"github.com/FollowTheProcess/test"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		cmd        *command.Command                         // The command under test
		customiser func(t *testing.T, cmd *command.Command) // An optional function to customise the command e.g. add flags
		name       string                                   // Identifier of the test case
		stdout     string                                   // Desired output to stdout
		stderr     string                                   // Desired output to stderr
		wantErr    bool                                     // Whether we want an error
	}{
		{
			name: "simple",
			cmd: command.New(
				"test",
				command.Args([]string{"arg1", "arg2", "arg3"}),
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from test")
					return nil
				}),
			),
			stdout:  "Hello from test\n",
			wantErr: false,
		},
		{
			name: "simple with flags",
			cmd: command.New(
				"test",
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "Oooh look, it ran, here are some args: %v\n", args)
					force, err := cmd.Flags().GetBool("force")
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				command.Args([]string{"arg1", "arg2", "--force"}),
			),
			customiser: func(t *testing.T, cmd *command.Command) {
				// Set flags in the customiser
				t.Helper()
				cmd.Flags().BoolP("force", "f", false, "Force something")
			},
			stdout:  "Oooh look, it ran, here are some args: [arg1 arg2]\n--force was: true\n",
			wantErr: false,
		},
		{
			name: "no run and no subcommands",
			cmd: command.New(
				"test",
				command.Args([]string{"arg1", "arg2", "arg3"}),
			),
			wantErr: true,
		},
		{
			name: "bad flag",
			cmd: command.New(
				"test",
				command.Run(func(cmd *command.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "Oooh look, it ran, here are some args: %v\n", args)
					force, err := cmd.Flags().GetBool("force")
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				command.Args([]string{"arg1", "arg2", "arg3", "-]force"}),
			),
			customiser: func(t *testing.T, cmd *command.Command) {
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

			command.Stderr(stderr)(tt.cmd)
			command.Stdout(stdout)(tt.cmd)

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

func TestSubCommandExecute(t *testing.T) {
	sub1 := command.New(
		"sub1",
		command.Run(func(cmd *command.Command, args []string) error {
			force, err := cmd.Flags().GetBool("force")
			if err != nil {
				return err
			}
			something, err := cmd.Flags().GetString("something")
			if err != nil {
				return err
			}
			if something == "" {
				something = "<empty>"
			}
			fmt.Fprintf(
				cmd.Stdout(),
				"Hello from sub1, my args were: %v, force was %v, something was %s",
				args,
				force,
				something,
			)
			return nil
		}),
	)
	sub1.Flags().BoolP("force", "f", false, "Force for sub1")
	sub1.Flags().StringP("something", "s", "", "Something for sub1")

	sub2 := command.New(
		"sub2",
		command.Run(func(cmd *command.Command, args []string) error {
			deleteFlag, err := cmd.Flags().GetBool("delete")
			if err != nil {
				return err
			}
			number, err := cmd.Flags().GetInt("number")
			if err != nil {
				return err
			}
			fmt.Fprintf(
				cmd.Stdout(),
				"Hello from sub2, my args were: %v, delete was %v, number was %d",
				args,
				deleteFlag,
				number,
			)
			return nil
		}),
	)
	sub2.Flags().BoolP("delete", "d", false, "Delete for sub2")
	sub2.Flags().IntP("number", "n", -1, "Number for sub2")

	root := command.New(
		"root",
		command.SubCommands(sub1, sub2),
	)

	tests := []struct {
		name    string   // Test case name
		stdout  string   // Expected stdout
		stderr  string   // Expected stderr
		args    []string // Args passed to root command
		wantErr bool     // Whether or not we wanted an error
	}{
		{
			name:    "invoke sub1 no flags",
			stdout:  "Hello from sub1, my args were: [my subcommand args], force was false, something was <empty>",
			stderr:  "",
			args:    []string{"sub1", "my", "subcommand", "args"},
			wantErr: false,
		},
		{
			name:    "invoke sub2 no flags",
			stdout:  "Hello from sub2, my args were: [my different args], delete was false, number was -1",
			stderr:  "",
			args:    []string{"sub2", "my", "different", "args"},
			wantErr: false,
		},
		{
			name:    "invoke sub1 with flags",
			stdout:  "Hello from sub1, my args were: [my subcommand args], force was true, something was here",
			stderr:  "",
			args:    []string{"sub1", "my", "subcommand", "args", "--force", "--something", "here"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the args on the root command
			command.Args(tt.args)(root)

			// Test output streams
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			command.Stderr(stderr)(root)
			command.Stdout(stdout)(root)

			// Execute the command, we should see the sub commands get executed based on what args we provide
			err := root.Execute()
			test.Ok(t, err)

			test.Equal(t, stdout.String(), tt.stdout)
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}

func TestHelp(t *testing.T) {
	tests := []struct {
		cmd     *command.Command // The command under test
		name    string           // Identifier of the test case
		golden  string           // The name of the file relative to testdata containing to expected output
		wantErr bool             // Whether we want an error
		debug   bool             // Whether or not to print the produced help text to stderr (useful for debugging)
	}{
		{
			name: "default long",
			cmd: command.New(
				"test",
				command.Args([]string{"--help"}),
			),
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "default short",
			cmd: command.New(
				"test",
				command.Args([]string{"-h"}),
			),
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "with examples",
			cmd: command.New(
				"test",
				command.Args([]string{"--help"}),
				command.Examples(
					command.Example{Comment: "Do a thing", Command: "test do thing --now"},
					command.Example{
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
			cmd: command.New(
				"test",
				command.Args([]string{"--help"}),
				command.Short("A cool CLI to do things"),
				command.Long("A longer, probably multiline description"),
			),
			golden:  "full.txt",
			wantErr: false,
		},
		{
			name: "with no description",
			cmd: command.New(
				"test",
				command.Args([]string{"--help"}),
				command.Short(""),
				command.Long(""),
			),
			golden:  "no-about.txt",
			wantErr: false,
		},
		{
			name: "with subcommands",
			cmd: command.New(
				"test",
				command.Args([]string{"--help"}),
				command.Short("A cool CLI to do things"),
				command.Long("A longer, probably multiline description"),
				command.SubCommands(
					command.New("sub1", command.Short("Do one thing")),
					command.New("sub2", command.Short("Do another thing")),
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

			command.Stderr(stderr)(tt.cmd)
			command.Stdout(stdout)(tt.cmd)

			err := tt.cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// Show the help output, can aid debugging
			if tt.debug {
				fmt.Print(stderr.String())
			}

			// --help output should be as per the golden file
			test.File(t, stderr.String(), filepath.Join("TestHelp", tt.golden))
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string           // Name of the test case
		cmd     *command.Command // Command under test
		stderr  string           // Expected output to stderr
		wantErr bool             // Whether we want an error or not
	}{
		{
			name:    "default long",
			cmd:     command.New("version-test", command.Args([]string{"--version"})),
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name:    "default short",
			cmd:     command.New("version-test", command.Args([]string{"-v"})),
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "custom version",
			cmd: command.New(
				"version-test",
				command.Args([]string{"--version"}),
				command.Version("v1.2.3"),
			),
			stderr:  "version-test, version: v1.2.3\n",
			wantErr: false,
		},
		{
			name: "custom versionFunc",
			cmd: command.New(
				"version-test",
				command.Args([]string{"--version"}),
				command.VersionFunc(func(cmd *command.Command) error {
					fmt.Fprintln(cmd.Stderr(), "Do something custom here")
					return nil
				}),
			),
			stderr:  "Do something custom here\n",
			wantErr: false,
		},
		{
			name: "return error",
			cmd: command.New(
				"version-test",
				command.Args([]string{"--version"}),
				command.VersionFunc(
					func(cmd *command.Command) error { return errors.New("Uh oh!") },
				),
			),
			wantErr: true,
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
		example command.Example
		want    string
	}{
		{
			name:    "empty",
			example: command.Example{},
			want:    "",
		},
		{
			name:    "only command",
			example: command.Example{Command: "run this program --once"},
			want:    "\n$ run this program --once\n",
		},
		{
			name:    "only comment",
			example: command.Example{Comment: "Run the program once"},
			want:    "\n# Run the program once\n",
		},
		{
			name: "both",
			example: command.Example{
				Comment: "Run the program once",
				Command: "run this program --once",
			},
			want: "\n# Run the program once\n$ run this program --once\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.example.String(), tt.want)
		})
	}
}
