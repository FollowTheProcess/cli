package cli_test

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

func TestExecute(t *testing.T) {
	var force bool
	tests := []struct {
		cmd     *cli.Command // The command under test
		name    string       // Identifier of the test case
		stdout  string       // Desired output to stdout
		stderr  string       // Desired output to stderr
		wantErr bool         // Whether we want an error
	}{
		{
			name: "simple",
			cmd: cli.New(
				"test",
				cli.Args([]string{"arg1", "arg2", "arg3"}),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintln(cmd.Stdout(), "Hello from test")
					return nil
				}),
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
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				cli.Args([]string{"arg1", "arg2", "--force"}),
				cli.Flag(&force, "force", "f", false, "Force something"),
			),
			stdout:  "Oooh look, it ran, here are some args: [arg1 arg2]\n--force was: true\n",
			wantErr: false,
		},
		{
			name: "no run and no subcommands",
			cmd: cli.New(
				"test",
				cli.Args([]string{"arg1", "arg2", "arg3"}),
			),
			wantErr: true,
		},
		{
			name: "bad flag",
			cmd: cli.New(
				"test",
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "Oooh look, it ran, here are some args: %v\n", args)
					fmt.Fprintf(cmd.Stdout(), "--force was: %v\n", force)
					return nil
				}),
				cli.Args([]string{"arg1", "arg2", "arg3", "-]force"}),
				cli.Flag(&force, "force", "f", false, "Force something"),
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

			test.Equal(t, stdout.String(), tt.stdout)
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}

func TestSubCommandExecute(t *testing.T) {
	var (
		force     bool
		something string
		deleteMe  bool
		number    int
	)
	sub1 := cli.New(
		"sub1",
		cli.Run(func(cmd *cli.Command, args []string) error {
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
		cli.Flag(&force, "force", "f", false, "Force for sub1"),
		cli.Flag(&something, "something", "s", "", "Something for sub1"),
	)

	sub2 := cli.New(
		"sub2",
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(
				cmd.Stdout(),
				"Hello from sub2, my args were: %v, delete was %v, number was %d",
				args,
				deleteMe,
				number,
			)
			return nil
		}),
		cli.Flag(&deleteMe, "delete", "d", false, "Delete for sub2"),
		cli.Flag(&number, "number", "n", -1, "Number for sub2"),
	)

	root := cli.New(
		"root",
		cli.SubCommands(sub1, sub2),
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
			cli.Args(tt.args)(root)

			// Test output streams
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			cli.Stderr(stderr)(root)
			cli.Stdout(stdout)(root)

			// Execute the command, we should see the sub commands get executed based on what args we provide
			err := root.Execute()
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
				cli.Example("Do a thing", "test do thing --now"),
				cli.Example("Do a different thing", "test do thing --different"),
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
			test.File(t, stderr.String(), filepath.Join("TestHelp", tt.golden))
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
				cli.VersionFunc(
					func(cmd *cli.Command) error { return errors.New("Uh oh!") },
				),
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
