package cli_test

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/test"
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
		name    string   // Test case name
		stdout  string   // Expected stdout
		stderr  string   // Expected stderr
		args    []string // Args passed to root command
		extra   []string // Extra args after "--" if present
		wantErr bool     // Whether or not we wanted an error
	}{
		{
			name:    "invoke sub1 no flags",
			stdout:  "Hello from sub1, my args were: [my subcommand args], force was false, something was <empty>, extra args: []",
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
			stdout:  "Hello from sub1, my args were: [my subcommand args], force was true, something was here, extra args: []",
			stderr:  "",
			args:    []string{"sub1", "my", "subcommand", "args", "--force", "--something", "here"},
			wantErr: false,
		},
		{
			name:    "invoke sub1 with arg terminator",
			stdout:  "Hello from sub1, my args were: [my subcommand args more args here], force was true, something was here, extra args: [more args here]",
			stderr:  "",
			args:    []string{"sub1", "my", "subcommand", "args", "--force", "--something", "here", "--", "more", "args", "here"},
			wantErr: false,
		},
		{
			name:    "invoke sub1 with sub1 in the arg list",
			stdout:  "Hello from sub1, my args were: [my sub1 args sub1 more args here], force was true, something was here, extra args: []",
			stderr:  "",
			args:    []string{"sub1", "my", "sub1", "args", "sub1", "--force", "--something", "here", "more", "args", "here"},
			wantErr: false,
		},
		{
			name:    "invoke sub1 with sub1 as a flag value",
			stdout:  "Hello from sub1, my args were: [my subcommand args more args here], force was true, something was sub2, extra args: []",
			stderr:  "",
			args:    []string{"sub1", "my", "subcommand", "args", "--force", "--something", "sub2", "more", "args", "here"},
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

			sub1, err := cli.New(
				"sub1",
				cli.Run(func(cmd *cli.Command, args []string) error {
					if something == "" {
						something = "<empty>"
					}
					fmt.Fprintf(
						cmd.Stdout(),
						"Hello from sub1, my args were: %v, force was %v, something was %s, extra args: %v",
						args,
						force,
						something,
						cmd.ExtraArgs(),
					)
					return nil
				}),

				cli.Flag(&force, "force", 'f', false, "Force for sub1"),
				cli.Flag(&something, "something", 's', "", "Something for sub1"),
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
				cli.Flag(&deleteMe, "delete", 'd', false, "Delete for sub2"),
				cli.Flag(&number, "number", 'n', -1, "Number for sub2"),
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

			if !tt.wantErr {
				test.Equal(t, stdoutBuf.String(), tt.stdout)
				test.Equal(t, stderrBuf.String(), tt.stderr)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	sub1, err := cli.New(
		"sub1",
		cli.Short("Do one thing"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub1")
			return nil
		}),
	)
	test.Ok(t, err)

	sub2, err := cli.New(
		"sub2",
		cli.Short("Do another thing"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub2")
			return nil
		}),
	)
	test.Ok(t, err)

	sub3, err := cli.New(
		"very-long-subcommand",
		cli.Short("Wow so long"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub3")
			return nil
		}),
	)
	test.Ok(t, err)

	tests := []struct {
		name    string       // Identifier of the test case
		golden  string       // Name of the file containing expected output
		options []cli.Option // Options to apply to the command
		wantErr bool         // Whether we want an error
	}{
		{
			name: "default long",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
			},
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "default.short",
			options: []cli.Option{
				cli.Args([]string{"-h"}),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
			},
			golden:  "default-help.txt",
			wantErr: false,
		},
		{
			name: "with examples",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Example("Do a thing", "test do thing --now"),
				cli.Example("Do a different thing", "test do thing --different"),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
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
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
			},
			golden:  "full.txt",
			wantErr: false,
		},
		{
			name: "with no description",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
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
			name: "subcommands different lengths",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2, sub3),
			},
			golden:  "subcommands-different-lengths.txt",
			wantErr: false,
		},
		{
			name: "with subcommands and flags",
			options: []cli.Option{
				cli.Args([]string{"--help"}),
				cli.Short("A cool CLI to do things"),
				cli.Long("A longer, probably multiline description"),
				cli.SubCommands(sub1, sub2),
				cli.Flag(new(bool), "delete", 'd', false, "Delete something"),
				cli.Flag(new(int), "count", cli.NoShortHand, -1, "Count something"),
			},
			golden:  "subcommands-flags.txt",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Force no colour in tests
			t.Setenv("NO_COLOR", "true")

			stderr := &bytes.Buffer{}
			stdout := &bytes.Buffer{}

			golden := filepath.Join(test.Data(t), "TestHelp", tt.golden)

			// Test specific overrides to the options in the table
			options := []cli.Option{cli.Stdout(stdout), cli.Stderr(stderr)}
			options = append(options, tt.options...)

			cmd, err := cli.New("test", options...)

			test.Ok(t, err)

			err = cmd.Execute()
			test.WantErr(t, err, tt.wantErr)

			if *debug {
				fmt.Printf("DEBUG\n_____\n\n%s\n", stderr.String())
			}

			if *update {
				t.Logf("Updating %s\n", golden)
				err := os.WriteFile(golden, stderr.Bytes(), os.ModePerm)
				test.Ok(t, err)
			}

			// Should have no output to stdout
			test.Equal(t, stdout.String(), "")

			// --help output should be as per the golden file
			test.File(t, stderr.String(), golden)
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
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
			},
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "default short",
			options: []cli.Option{
				cli.Args([]string{"-v"}),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
			},
			stderr:  "version-test, version: dev\n",
			wantErr: false,
		},
		{
			name: "custom version",
			options: []cli.Option{
				cli.Args([]string{"--version"}),
				cli.Version("v1.2.3"),
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
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
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
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
				cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
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
				cli.Flag(new(int), "count", 'c', 0, "Count something"),
				cli.Flag(new(int), "count", 'c', 0, "Count something (again)"),
			},
			errMsg: `flag "count" already defined`,
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
			name:    "example both empty",
			options: []cli.Option{cli.Example("", "")},
			errMsg:  "example comment cannot be empty\nexample command cannot be empty",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.New("test", tt.options...)
			test.Err(t, err)
			test.Equal(t, err.Error(), tt.errMsg)
		})
	}
}

func TestDuplicateSubCommands(t *testing.T) {
	sub1, err := cli.New(
		"sub1",
		cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
	)
	test.Ok(t, err)

	sub2, err := cli.New(
		"sub2",
		cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
	)
	test.Ok(t, err)

	sub1Again, err := cli.New(
		"sub1",
		cli.Run(func(cmd *cli.Command, args []string) error { return nil }),
	)
	test.Ok(t, err) // Shouldn't error at this point as it's not joined up

	_, err = cli.New(
		"root",
		cli.SubCommands(sub1, sub2, sub1Again), // This should cause the error
	)

	test.Err(t, err)
	test.Equal(t, err.Error(), `subcommand "sub1" already defined`)
}

func TestCommandNoRunNoSub(t *testing.T) {
	_, err := cli.New(
		"root",
		cli.Args([]string{}),
		// Run function missing and no subcommand
	)
	test.Err(t, err)
}

func TestExecuteNonRootCommand(t *testing.T) {
	sub, err := cli.New(
		"sub",
		cli.Args([]string{"hello"}),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub")
			return nil
		}),
	)
	test.Ok(t, err)

	_, err = cli.New(
		"root",
		cli.Args([]string{"sub"}),
		cli.SubCommands(sub),
	)
	test.Ok(t, err)

	// Call sub's Execute, we should get an error
	err = sub.Execute()
	test.Err(t, err)
	if err != nil {
		test.Equal(t, err.Error(), "Execute must be called on the root of the command tree, was called on sub")
	}
}

func BenchmarkExecuteHelp(b *testing.B) {
	sub1, err := cli.New(
		"sub1",
		cli.Short("Do one thing"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub1")
			return nil
		}),
	)
	test.Ok(b, err)

	sub2, err := cli.New(
		"sub2",
		cli.Short("Do another thing"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub2")
			return nil
		}),
	)
	test.Ok(b, err)

	sub3, err := cli.New(
		"very-long-subcommand",
		cli.Short("Wow so long"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintln(cmd.Stdout(), "Hello from sub3")
			return nil
		}),
	)
	test.Ok(b, err)

	cmd, err := cli.New(
		"bench-help",
		cli.Short("A helpful benchmark command"),
		cli.Long("Much longer text..."),
		cli.Example("Do a thing", "bench-help very-long-subcommand --flag"),
		cli.SubCommands(sub1, sub2, sub3),
		cli.Args([]string{"--help"}),
		cli.Stderr(io.Discard),
		cli.Stdout(io.Discard),
	)
	test.Ok(b, err)

	b.ResetTimer()
	for range b.N {
		err := cmd.Execute()
		if err != nil {
			b.Fatalf("Execute returned an error: %v", err)
		}
	}
}
