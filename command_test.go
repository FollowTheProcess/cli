package cli_test

import (
	"bytes"
	goflag "flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"slices"
	"testing"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/snapshot"
	"go.followtheprocess.codes/test"
)

var (
	debug  = goflag.Bool("debug", false, "Print debug output during tests")
	update = goflag.Bool("update", false, "Update golden files")
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
				cli.OverrideArgs([]string{"hello", "there"}),
				cli.Arg(new(string), "first", "The first word"), // Expect the positional arguments
				cli.Arg(new(string), "second", "The second word"),
				cli.Stdin(os.Stdin), // Set stdin for the lols
			},
			wantErr: false,
		},
		{
			name:   "simple with flag",
			stdout: "My arguments were: [hello there]\nForce was: true\n",
			stderr: "",
			options: []cli.Option{
				cli.OverrideArgs([]string{"hello", "there", "--force"}),
				cli.Arg(new(string), "first", "The first word"), // Expect the positional arguments
				cli.Arg(new(string), "second", "The second word"),
			},
			wantErr: false,
		},
		{
			name:   "bad flag",
			stdout: "",
			stderr: "",
			options: []cli.Option{
				cli.OverrideArgs([]string{"arg1", "arg2", "arg3", "-]force"}),
				cli.Arg(new(string), "first", "The first arg"), // Expect the positional arguments
				cli.Arg(new(string), "second", "The second arg"),
				cli.Arg(new(string), "third", "The third arg"),
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
				cli.Run(func(cmd *cli.Command) error {
					fmt.Fprintf(cmd.Stdout(), "My arguments were: %v\nForce was: %v\n", cmd.Args(), force)

					return nil
				}),
				cli.Flag(&force, "force", 'f', false, "Force something"),
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
		name        string       // Test case name
		stdout      string       // Expected stdout
		stderr      string       // Expected stderr
		args        []string     // Args passed to root command
		extra       []string     // Extra args after "--" if present
		sub1Options []cli.Option // Options for subcommand 1
		sub2Options []cli.Option // Options for subcommand 1
		wantErr     bool         // Whether or not we wanted an error
	}{
		{
			name:   "invoke sub1 no flags",
			stdout: "Hello from sub1, my args were: [my subcommand args], force was false, something was <empty>, extra args: []",
			stderr: "",
			args:   []string{"sub1", "my", "subcommand", "args"},
			sub1Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
			},
			wantErr: false,
		},
		{
			name:   "invoke sub2 no flags",
			stdout: "Hello from sub2, my args were: [my different args], delete was false, number was -1",
			stderr: "",
			args:   []string{"sub2", "my", "different", "args"},
			sub2Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
			},
			wantErr: false,
		},
		{
			name:   "invoke sub1 with flags",
			stdout: "Hello from sub1, my args were: [my subcommand args], force was true, something was here, extra args: []",
			stderr: "",
			args:   []string{"sub1", "my", "subcommand", "args", "--force", "--something", "here"},
			sub1Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
			},
			wantErr: false,
		},
		{
			name:   "invoke sub1 with arg terminator",
			stdout: "Hello from sub1, my args were: [my subcommand args more args here], force was true, something was here, extra args: [more args here]",
			stderr: "",
			args: []string{
				"sub1",
				"my",
				"subcommand",
				"args",
				"--force",
				"--something",
				"here",
				"--",
				"more",
				"args",
				"here",
			},
			sub1Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
				cli.Arg(new(string), "four", "Fourth arg"),
				cli.Arg(new(string), "five", "Fifth arg"),
				cli.Arg(new(string), "six", "Sixth arg"),
			},
			wantErr: false,
		},
		{
			name:   "invoke sub1 with sub1 in the arg list",
			stdout: "Hello from sub1, my args were: [my sub1 args sub1 more args here], force was true, something was here, extra args: []",
			stderr: "",
			args: []string{
				"sub1",
				"my",
				"sub1",
				"args",
				"sub1",
				"--force",
				"--something",
				"here",
				"more",
				"args",
				"here",
			},
			sub1Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
				cli.Arg(new(string), "four", "Fourth arg"),
				cli.Arg(new(string), "five", "Fifth arg"),
				cli.Arg(new(string), "six", "Sixth arg"),
				cli.Arg(new(string), "seven", "Seventh arg"),
			},
			wantErr: false,
		},
		{
			name:   "invoke sub1 with sub1 as a flag value",
			stdout: "Hello from sub1, my args were: [my subcommand args more args here], force was true, something was sub2, extra args: []",
			stderr: "",
			args: []string{
				"sub1",
				"my",
				"subcommand",
				"args",
				"--force",
				"--something",
				"sub2",
				"more",
				"args",
				"here",
			},
			sub1Options: []cli.Option{
				cli.Arg(new(string), "one", "First arg"),
				cli.Arg(new(string), "two", "Second arg"),
				cli.Arg(new(string), "three", "Third arg"),
				cli.Arg(new(string), "four", "Fourth arg"),
				cli.Arg(new(string), "five", "Fifth arg"),
				cli.Arg(new(string), "six", "Sixth arg"),
			},
			wantErr: false,
		},
		{
			name:    "invoke root with no args",
			stdout:  "",
			stderr:  "",
			args:    []string{},
			wantErr: true,
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

			sub1 := func() (*cli.Command, error) {
				defaultOpts := []cli.Option{
					cli.Flag(&force, "force", 'f', false, "Force for sub1"),
					cli.Flag(&something, "something", 's', "", "Something for sub1"),
					cli.Run(func(cmd *cli.Command) error {
						if something == "" {
							something = "<empty>"
						}

						extra, ok := cmd.ExtraArgs()
						if !ok {
							extra = []string{}
						}

						fmt.Fprintf(
							cmd.Stdout(),
							"Hello from sub1, my args were: %v, force was %v, something was %s, extra args: %v",
							cmd.Args(),
							force,
							something,
							extra,
						)

						return nil
					}),
				}
				opts := slices.Concat(defaultOpts, tt.sub1Options)

				return cli.New(
					"sub1",
					opts...,
				)
			}

			sub2 := func() (*cli.Command, error) {
				defaultOpts := []cli.Option{
					cli.Run(func(cmd *cli.Command) error {
						fmt.Fprintf(
							cmd.Stdout(),
							"Hello from sub2, my args were: %v, delete was %v, number was %d",
							cmd.Args(),
							deleteMe,
							number,
						)

						return nil
					}),
					cli.Flag(&deleteMe, "delete", 'd', false, "Delete for sub2"),
					cli.Flag(&number, "number", 'n', -1, "Number for sub2"),
				}

				opts := slices.Concat(defaultOpts, tt.sub2Options)

				return cli.New(
					"sub2",
					opts...,
				)
			}

			root, err := cli.New(
				"root",
				cli.SubCommands(sub1, sub2),
				cli.Stdout(stdoutBuf),
				cli.Stderr(stderrBuf),
				cli.OverrideArgs(tt.args),
			)

			test.Ok(t, err)

			// Execute the command, we should see the sub commands get executed based on what args we provide
			err = root.Execute()
			test.WantErr(t, err, tt.wantErr)

			if !tt.wantErr {
				test.Equal(t, stdoutBuf.String(), tt.stdout)
				test.Equal(t, stderrBuf.String(), tt.stderr)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	sub1 := func() (*cli.Command, error) {
		return cli.New(
			"sub1",
			cli.Short("Do one thing"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub1")

				return nil
			}))
	}

	sub2 := func() (*cli.Command, error) {
		return cli.New(
			"sub2",
			cli.Short("Do another thing"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub2")

				return nil
			}),
		)
	}

	sub3 := func() (*cli.Command, error) {
		return cli.New(
			"very-long-subcommand",
			cli.Short("Wow so long"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub3")

				return nil
			}),
		)
	}

	tests := []struct {
		name    string       // Identifier of the test case
		options []cli.Option // Options to apply to the command
		wantErr bool         // Whether we want an error
	}{
		{
			name: "default long",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "default short",
			options: []cli.Option{
				cli.OverrideArgs([]string{"-h"}),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with examples",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Example("Do a thing", "test do thing --now"),
				cli.Example("Do a different thing", "test do thing --different"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with named arguments",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Arg(new(string), "src", "The file to copy"),                                       // This one is required
				cli.Arg(new(string), "dest", "Destination to copy to", cli.ArgDefault("default.txt")), // This one is optional
				cli.Arg(new(int), "other", "Something else", cli.ArgDefault(0)),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with verbosity count",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Arg(new(string), "src", "The file to copy"),                                           // This one is required
				cli.Arg(new(string), "dest", "Destination to copy to", cli.ArgDefault("destination.txt")), // This one is optional
				cli.Flag(new(flag.Count), "verbosity", 'v', 0, "Increase the verbosity level"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with full description",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "full description strip whitespace",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Short("  \t\n A cool CLI to do things   \n "),
				cli.Long("  \t\n\n A longer, probably multiline description \t\n\n "),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with no description",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "with subcommands",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2),
			},
			wantErr: false,
		},
		{
			name: "subcommands different lengths",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2, sub3),
			},
			wantErr: false,
		},
		{
			name: "with subcommands and flags",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2),
				cli.Flag(new(bool), "delete", 'd', false, "Delete something"),
				cli.Flag(new(int), "count", flag.NoShortHand, -1, "Count something"),
				cli.Flag(new([]string), "things", flag.NoShortHand, nil, "Names of things"),
				cli.Flag(new([]string), "more", flag.NoShortHand, []string{"one", "two"}, "Names of things with a default"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := snapshot.New(t, snapshot.Update(*update), snapshot.Color(os.Getenv("CI") == ""))

			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			// Test specific overrides to the options in the table
			options := []cli.Option{
				cli.Stdout(stdout),
				cli.Stderr(stderr),
				cli.NoColour(true),
			}

			cmd, err := cli.New("test", slices.Concat(options, tt.options)...)

			test.Ok(t, err)

			err = cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			if *debug {
				fmt.Printf("DEBUG\n_____\n\n%s\n", stderr.String())
			}

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// --help output should be as per the snapshot
			snap.Snap(stderr.String())
		})
	}
}

func TestVersion(t *testing.T) {
	sub1 := func() (*cli.Command, error) {
		return cli.New(
			"sub1",
			cli.Short("Do one thing"),
			// No version set on sub1
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub1")

				return nil
			}))
	}

	sub2 := func() (*cli.Command, error) {
		return cli.New(
			"sub2",
			cli.Short("Do another thing"),
			cli.Version("sub2 version text"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub2")

				return nil
			}),
		)
	}

	tests := []struct {
		name    string       // Name of the test case
		stderr  string       // Expected output to stderr
		options []cli.Option // Options to apply to the test command
		wantErr bool         // Whether we want an error or not
	}{
		{
			name: "default long",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: dev\n",
			wantErr: false,
		},
		{
			name: "default short",
			options: []cli.Option{
				cli.OverrideArgs([]string{"-V"}),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: dev\n",
			wantErr: false,
		},
		{
			name: "with version",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.Version("v3.1.7"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: v3.1.7\n",
			wantErr: false,
		},
		{
			name: "with commit",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.Commit("eedb45b"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: dev\nCommit: eedb45b\n",
			wantErr: false,
		},
		{
			name: "with build date",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.BuildDate("2024-04-11T02:23:42Z"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: dev\nBuildDate: 2024-04-11T02:23:42Z\n",
			wantErr: false,
		},
		{
			name: "with version and commit",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.Version("v8.17.6"),
				cli.Commit("b9aaafd"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: v8.17.6\nCommit: b9aaafd\n",
			wantErr: false,
		},
		{
			name: "with all",
			options: []cli.Option{
				cli.OverrideArgs([]string{"--version"}),
				cli.Version("v8.17.6"),
				cli.Commit("b9aaafd"),
				cli.BuildDate("2024-08-17T10:37:30Z"),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			stderr:  "version-test\n\nVersion: v8.17.6\nCommit: b9aaafd\nBuildDate: 2024-08-17T10:37:30Z\n",
			wantErr: false,
		},
		{
			name: "call on subcommand with no version",
			options: []cli.Option{
				cli.OverrideArgs([]string{"sub1", "--version"}),
				cli.Version("v8.17.6"),
				cli.Commit("b9aaafd"),
				cli.BuildDate("2024-08-17T10:37:30Z"),
				cli.SubCommands(sub1, sub2),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			// Should show the root commands version info
			stderr:  "sub1\n\nVersion: v8.17.6\nCommit: b9aaafd\nBuildDate: 2024-08-17T10:37:30Z\n",
			wantErr: false,
		},
		{
			name: "call on subcommand with version",
			options: []cli.Option{
				cli.OverrideArgs([]string{"sub2", "--version"}),
				cli.Version("v8.17.6"),
				cli.Commit("b9aaafd"),
				cli.BuildDate("2024-08-17T10:37:30Z"),
				cli.SubCommands(sub1, sub2),
				cli.Run(func(cmd *cli.Command) error { return nil }),
			},
			// Should show sub2's version text
			stderr:  "sub2\n\nVersion: sub2 version text\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			// Test specific overrides to the options in the table
			options := []cli.Option{
				cli.Stdout(stdout),
				cli.Stderr(stderr),
				cli.NoColour(true),
			}

			cmd, err := cli.New("version-test", slices.Concat(options, tt.options)...)
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
		errMsg  string       // Expected error message
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
			name:    "nil override args",
			options: []cli.Option{cli.OverrideArgs(nil)},
			errMsg:  "cannot set Args to nil",
		},
		{
			name:    "nil Run",
			options: []cli.Option{cli.Run(nil)},
			errMsg:  "cannot set Run to nil",
		},
		{
			name: "flag already exists",
			options: []cli.Option{
				cli.Flag(new(int), "count", 'c', 0, "Count something"),
				cli.Flag(new(int), "count", 'c', 0, "Count something (again)"),
			},
			errMsg: `flag "count" already defined`,
		},
		{
			name: "flag short already exists",
			options: []cli.Option{
				cli.Flag(new(int), "count", 'c', 0, "Count something"),
				cli.Flag(new(string), "config", 'c', "", "Path to config file"),
			},
			errMsg: `could not add flag "config" to command "test": shorthand "c" already in use for flag "count"`,
		},
		{
			name:    "example comment empty",
			options: []cli.Option{cli.Example("", "command here")},
			errMsg:  "example comment cannot be empty",
		},
		{
			name:    "example command empty",
			options: []cli.Option{cli.Example("comment here", "")},
			errMsg:  "example command cannot be empty",
		},
		{
			name:    "empty short description",
			options: []cli.Option{cli.Short("")},
			errMsg:  "cannot set command short description to an empty string",
		},
		{
			name:    "empty long description",
			options: []cli.Option{cli.Long("")},
			errMsg:  "cannot set command long description to an empty string",
		},
		{
			name:    "empty arg name",
			options: []cli.Option{cli.Arg(new(string), "", "empty required arg")},
			errMsg:  `invalid arg name "": must not be empty`,
		},
		{
			name:    "arg name contains whitespace",
			options: []cli.Option{cli.Arg(new(string), "a space", "some space things")},
			errMsg:  `invalid arg name "a space": cannot contain whitespace`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.New("test", tt.options...)
			test.Err(t, err)                      // Invalid option should have triggered an error
			test.Equal(t, err.Error(), tt.errMsg) // Error message was not as expected
		})
	}
}

func TestDuplicateSubCommands(t *testing.T) {
	sub1 := func() (*cli.Command, error) {
		return cli.New(
			"sub1",
			cli.Run(func(cmd *cli.Command) error { return nil }),
		)
	}

	sub2 := func() (*cli.Command, error) {
		return cli.New(
			"sub2",
			cli.Run(func(cmd *cli.Command) error { return nil }),
		)
	}

	sub1Again := func() (*cli.Command, error) {
		return cli.New(
			"sub1",
			cli.Run(func(cmd *cli.Command) error { return nil }),
		)
	}

	_, err := cli.New(
		"root",
		cli.SubCommands(sub1, sub2, sub1Again), // This should cause the error
	)

	test.Err(t, err)

	if err != nil {
		test.Equal(t, err.Error(), `subcommand "sub1" already defined`)
	}
}

func TestCommandNoRunNoSub(t *testing.T) {
	_, err := cli.New(
		"root",
		cli.OverrideArgs([]string{}),
		// Run function missing and no subcommand
	)
	test.Err(t, err)
}

func TestExecuteNilCommand(t *testing.T) {
	var cmd *cli.Command

	err := cmd.Execute()
	test.Err(t, err)

	if err != nil {
		test.Equal(t, err.Error(), "Execute called on a nil Command")
	}
}

// The order in which we apply options shouldn't matter, this test
// shuffles the order of the options and asserts the Command we get
// out behaves the same as a baseline.
func TestCommandOptionOrder(t *testing.T) {
	baseLineStdout := &bytes.Buffer{}
	baseLineStderr := &bytes.Buffer{}

	shuffleStdout := &bytes.Buffer{}
	shuffleStderr := &bytes.Buffer{}

	var (
		f     bool
		count int
	)

	sub := func() (*cli.Command, error) {
		return cli.New(
			"sub",
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub")

				return nil
			}),
		)
	}

	options := []cli.Option{
		cli.OverrideArgs([]string{"some", "args", "here", "--flag", "--count", "10"}),
		cli.Short("Short description"),
		cli.Long("Long description"),
		cli.Example("Do a thing", "demo run something --flag"),
		cli.Run(func(cmd *cli.Command) error {
			fmt.Fprintf(cmd.Stdout(), "args: %v, flag: %v, count: %v\n", cmd.Args(), f, count)

			return nil
		}),
		cli.Arg(new(string), "first", "First arg"), // It just needs *some* args
		cli.Arg(new(string), "second", "Second arg"),
		cli.Arg(new(string), "third", "Third arg"),
		cli.Version("v1.2.3"),
		cli.SubCommands(sub),
		cli.Flag(&f, "flag", 'f', false, "Set a bool flag"),
		cli.Flag(&count, "count", 'c', 0, "Count a thing"),
	}

	baseLineOptions := slices.Concat(
		options,
		[]cli.Option{
			cli.Stderr(baseLineStderr), // Set output streams specific to the baseline
			cli.Stdout(baseLineStdout),
		})

	baseline, err := cli.New("baseline", baseLineOptions...)
	test.Ok(t, err)

	err = baseline.Execute()
	test.Ok(t, err)

	// Make sure the baseline is behaving as expected
	test.Equal(t, baseLineStdout.String(), "args: [some args here], flag: true, count: 10\n")
	test.Equal(t, baseLineStderr.String(), "")

	// Shuffley shuffle, 100 permutations should do it
	for range 100 {
		shuffled := shuffle(options)

		// Set output streams specific to the shuffled command
		shuffleOptions := slices.Concat(
			shuffled,
			[]cli.Option{
				cli.Stderr(shuffleStderr),
				cli.Stdout(shuffleStdout),
			},
		)

		// Make a Command with the randomly ordered options
		shuffle, err := cli.New("shuffle", shuffleOptions...)
		test.Ok(t, err)

		// The two commands should behave equivalently
		err = shuffle.Execute()
		test.Ok(t, err)

		test.Equal(t, shuffleStdout.String(), baseLineStdout.String())
		test.Equal(t, shuffleStderr.String(), baseLineStderr.String())

		// Clear the buffers for the next loop
		shuffleStderr.Reset()
		shuffleStdout.Reset()
	}
}

func BenchmarkExecuteHelp(b *testing.B) {
	sub1 := func() (*cli.Command, error) {
		return cli.New(
			"sub1",
			cli.Short("Do one thing"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub1")

				return nil
			}),
		)
	}

	sub2 := func() (*cli.Command, error) {
		return cli.New(
			"sub2",
			cli.Short("Do another thing"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub2")

				return nil
			}),
		)
	}

	sub3 := func() (*cli.Command, error) {
		return cli.New(
			"very-long-subcommand",
			cli.Short("Wow so long"),
			cli.Run(func(cmd *cli.Command) error {
				fmt.Fprintln(cmd.Stdout(), "Hello from sub3")

				return nil
			}),
		)
	}

	cmd, err := cli.New(
		"bench-help",
		cli.Short("A helpful benchmark command"),
		cli.Long("Much longer text..."),
		cli.Example("Do a thing", "bench-help very-long-subcommand --flag"),
		cli.SubCommands(sub1, sub2, sub3),
		cli.OverrideArgs([]string{"--help"}),
		cli.Stderr(io.Discard),
		cli.Stdout(io.Discard),
	)
	test.Ok(b, err)

	for b.Loop() {
		err := cmd.Execute()
		if err != nil {
			b.Fatalf("Execute returned an error: %v", err)
		}
	}
}

// Benchmarks calling New to build a typical CLI.
func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_, err := cli.New(
			"benchy",
			cli.Short("A typical CLI to benchmark calling cli.New"),
			cli.Version("dev"),
			cli.Commit("dfdddaf"),
			cli.Example("An example", "bench --help"),
			cli.Flag(new(bool), "force", 'f', false, "Force something"),
			cli.Flag(new(string), "name", 'n', "", "The name of something"),
			cli.Flag(new(int), "count", 'c', 1, "Count something"),
			cli.Run(func(cmd *cli.Command) error { return nil }),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarks calling the --version flag.
func BenchmarkVersion(b *testing.B) {
	cmd, err := cli.New(
		"bench-version",
		cli.Version("v1.2.3"),
		cli.Commit("f0d129cbda5cac5f57ffe091481bfb35fc6d5ee4"),
		cli.BuildDate("2025-02-16T07:36:04Z"),
		cli.OverrideArgs([]string{"--version"}),
		cli.Stderr(io.Discard),
		cli.Stdout(io.Discard),
		cli.Run(func(cmd *cli.Command) error { return nil }),
	)
	test.Ok(b, err)

	for b.Loop() {
		err := cmd.Execute()
		if err != nil {
			b.Fatalf("Execute returned an error: %v", err)
		}
	}
}

// shuffle returns a randomly ordered copy of items.
func shuffle[T any](items []T) []T {
	shuffled := slices.Clone(items)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
