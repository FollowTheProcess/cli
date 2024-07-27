package cli_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
)

var (
	debug  = flag.Bool("debug", false, "Print debug output during tests")
	update = flag.Bool("update", false, "Update golden files")
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string       // Name of the test case
		stdout  string       // Expected output to stdout
		stderr  string       // Expected output to stderr
		options []cli.Option // Options to apply to the test command
		wantErr bool         // Whether we want an error or not
	}{
		{
			name:   "simple",
			stdout: "My arguments were: [hello there]\nForce was: false\n",
			stderr: "",
			options: []cli.Option{
				cli.Args([]string{"hello", "there"}),
				cli.Stdin(os.Stdin), // Set stdin for the lols
			},
			wantErr: false,
		},
		{
			name:   "simple with flag",
			stdout: "My arguments were: [hello there]\nForce was: true\n",
			stderr: "",
			options: []cli.Option{
				cli.Args([]string{"hello", "there", "--force"}),
			},
			wantErr: false,
		},
		{
			name:   "bad flag",
			stdout: "",
			stderr: "",
			options: []cli.Option{
				cli.Args([]string{"arg1", "arg2", "arg3", "-]force"}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var force bool

			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			// Test specific overrides to the options in the table
			options := []cli.Option{
				cli.Stdout(stdout),
				cli.Stderr(stderr),
				cli.Run(func(cmd *cli.Command, args []string) error {
					fmt.Fprintf(cmd.Stdout(), "My arguments were: %v\nForce was: %v\n", args, force)
					return nil
				}),
				cli.Flag(&force, "force", "f", false, "Force something"),
			}

			cmd, err := cli.New("test", slices.Concat(options, tt.options)...)
			test.Ok(t, err)

			err = cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			test.Equal(t, stdout.String(), tt.stdout)
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}

func TestSubCommandExecute(t *testing.T) {
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
			var (
				force     bool
				something string
				deleteMe  bool
				number    int
				stdoutBuf = &bytes.Buffer{}
				stderrBuf = &bytes.Buffer{}
			)

			sub1, err := cli.New(
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

			test.Ok(t, err)

			sub2, err := cli.New(
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

			test.Ok(t, err)

			root, err := cli.New(
				"root",
				cli.SubCommands(sub1, sub2),
				cli.Stdout(stdoutBuf),
				cli.Stderr(stderrBuf),
				cli.Args(tt.args),
			)

			test.Ok(t, err)

			// Execute the command, we should see the sub commands get executed based on what args we provide
			err = root.Execute()
			test.WantErr(t, err, tt.wantErr)

			test.Equal(t, stdoutBuf.String(), tt.stdout)
			test.Equal(t, stderrBuf.String(), tt.stderr)
		})
	}
}

func TestHelp(t *testing.T) {
	sub1, err := cli.New("sub1", cli.Short("Do one thing"))
	test.Ok(t, err)

	sub2, err := cli.New("sub2", cli.Short("Do another thing"))
	test.Ok(t, err)
	tests := []struct {
		name    string       // Identifier of the test case
		golden  string       // Name of the file containing expected output
		options []cli.Option // Options to apply to the command
		wantErr bool         // Whether we want an error
	}{
		{
			name:    "default long",
			options: []cli.Option{cli.Args([]string{"--help"})},
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name:    "default.short",
			options: []cli.Option{cli.Args([]string{"-h"})},
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "with examples",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Example("Do a thing", "test do thing --now"),
				cli.Example("Do a different thing", "test do thing --different"),
			},
			golden:  "with-examples.txt",
			wantErr: false,
		},
		{
			name: "with full description",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
			},
			golden:  "full.txt",
			wantErr: false,
		},
		{
			name: "with no description",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short(""),
				cli.Long(""),
			},
			golden:  "no-about.txt",
			wantErr: false,
		},
		{
			name: "with subcommands",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2),
			},
			golden:  "subcommands.txt",
			wantErr: false,
		},
		{
			name: "with subcommands and flags",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2),
				cli.Flag(new(bool), "delete", "d", false, "Delete something"),
				cli.Flag(new(int), "count", "", -1, "Count something"),
			},
			golden:  "subcommands-flags.txt",
			wantErr: false,
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

			if *debug {
				fmt.Printf("DEBUG\n_____\n\n%s\n", stderr.String())
			}

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// --help output should be as per the golden file
			test.File(t, stderr.String(), filepath.Join("TestHelp", tt.golden))

			if *update {
				err := os.WriteFile(filepath.Join("TestHelp", tt.golden), stderr.Bytes(), os.ModePerm)
				test.Ok(t, err)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string       // Name of the test case
		stderr  string       // Expected output to stderr
		options []cli.Option // Options to apply to the test command
		wantErr bool         // Whether we want an error or not
	}{
		{
			name: "default long",
			options: []cli.Option{
				cli.Args([]string{"--version"}),
			},
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "default short",
			options: []cli.Option{
				cli.Args([]string{"-v"}),
			},
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "custom version",
			options: []cli.Option{
				cli.Args([]string{"--version"}),
				cli.Version("v1.2.3"),
			},
			stderr:  "version-test, version: v1.2.3\n",
			wantErr: false,
		},
		{
			name: "custom versionFunc",
			options: []cli.Option{
				cli.Args([]string{"--version"}),
				cli.VersionFunc(func(cmd *cli.Command) error {
					fmt.Fprintln(cmd.Stderr(), "Do something custom here")
					return nil
				}),
			},
			stderr:  "Do something custom here\n",
			wantErr: false,
		},
		{
			name: "return error",
			options: []cli.Option{
				cli.Args([]string{"--version"}),
				cli.VersionFunc(func(cmd *cli.Command) error {
					return errors.New("Uh oh!")
				}),
			},
			stderr:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			// Test specific overrides to the options in the table
			options := []cli.Option{cli.Stdout(stdout), cli.Stderr(stderr)}

			cmd, err := cli.New("version-test", slices.Concat(tt.options, options)...)
			test.Ok(t, err)

			err = cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// --version output should be as desired
			test.Equal(t, stderr.String(), tt.stderr)
		})
	}
}

func TestOptionValidation(t *testing.T) {
	tests := []struct {
		name    string       // Name of the test case
		errMsg  string       // If we wanted an error, what should it say
		options []cli.Option // Options to apply to the command
	}{
		{
			name:    "nil stdin",
			options: []cli.Option{cli.Stdin(nil)},
			errMsg:  "cannot set Stdin to nil",
		},
		{
			name:    "nil stdout",
			options: []cli.Option{cli.Stdout(nil)},
			errMsg:  "cannot set Stdout to nil",
		},
		{
			name:    "nil stderr",
			options: []cli.Option{cli.Stderr(nil)},
			errMsg:  "cannot set Stderr to nil",
		},
		{
			name:    "nil all three",
			options: []cli.Option{cli.Stdout(nil), cli.Stderr(nil), cli.Stdin(nil)},
			errMsg:  "cannot set Stdout to nil\ncannot set Stderr to nil\ncannot set Stdin to nil",
		},
		{
			name:    "nil args",
			options: []cli.Option{cli.Args(nil)},
			errMsg:  "cannot set Args to nil",
		},
		{
			name:    "empty version",
			options: []cli.Option{cli.Version("")},
			errMsg:  "cannot set Version to an empty string",
		},
		{
			name:    "nil VersionFunc",
			options: []cli.Option{cli.VersionFunc(nil)},
			errMsg:  "cannot set VersionFunc to nil",
		},
		{
			name:    "nil Run",
			options: []cli.Option{cli.Run(nil)},
			errMsg:  "cannot set Run to nil",
		},
		{
			name:    "nil ArgValidator",
			options: []cli.Option{cli.Allow(nil)},
			errMsg:  "cannot set Allow to a nil ArgValidator",
		},
		{
			name: "flag already exists",
			options: []cli.Option{
				cli.Flag(new(int), "count", "c", 0, "Count something"),
				cli.Flag(new(int), "count", "c", 0, "Count something (again)"),
			},
			errMsg: `flag "count" already defined`,
		},
		{
			name:    "short too long",
			options: []cli.Option{cli.Flag(new(bool), "short", "word", false, "Set something")},
			errMsg:  `shorthand for flag "short" must be a single ASCII letter, got "word" which has 4 letters`,
		},
		{
			name:    "short is digit",
			options: []cli.Option{cli.Flag(new(bool), "short", "7", false, "Set something")},
			errMsg:  `shorthand for flag "short" is an invalid character, must be a single ASCII letter, got "7"`,
		},
		{
			name:    "short is non ascii",
			options: []cli.Option{cli.Flag(new(bool), "short", "本", false, "Set something")},
			errMsg:  `shorthand for flag "short" is an invalid character, must be a single ASCII letter, got "本"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.New("test", tt.options...)
			test.Err(t, err)
			test.Equal(t, err.Error(), tt.errMsg)
		})
	}
}
