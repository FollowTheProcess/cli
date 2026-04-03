package cli_test

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/snapshot"
	"go.followtheprocess.codes/test"
)

func TestCompletionSpec(t *testing.T) {
	tests := []struct {
		name    string
		options []cli.Option
	}{
		{
			name: "no user flags",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "bool flag",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new(bool), "force", 'f', "Force deletion"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "string flag",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new(string), "output", 'o', "Output file"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "count flag",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new(flag.Count), "verbose", 'v', "Verbosity level"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "long-only flag",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new(bool), "dry-run", flag.NoShortHand, "Dry run mode"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "slice flag",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new([]string), "tags", 't', "Tags to apply"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "with subcommands",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.SubCommands(
					cli.CompletionSubCommand(),
					func() (*cli.Command, error) {
						return cli.New("delete",
							cli.Short("Delete a resource"),
							cli.Run(func(_ context.Context, _ *cli.Command) error { return nil }),
						)
					},
					func() (*cli.Command, error) {
						return cli.New("get",
							cli.Short("Get a resource"),
							cli.Run(func(_ context.Context, _ *cli.Command) error { return nil }),
						)
					},
				),
			},
		},
		{
			name: "nested subcommands",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.SubCommands(
					cli.CompletionSubCommand(),
					func() (*cli.Command, error) {
						return cli.New("repo",
							cli.Short("Repository commands"),
							cli.SubCommands(
								func() (*cli.Command, error) {
									return cli.New("clone",
										cli.Short("Clone a repository"),
										cli.Flag(new(string), "branch", 'b', "Branch to clone"),
										cli.Run(func(_ context.Context, _ *cli.Command) error { return nil }),
									)
								},
							),
						)
					},
				),
			},
		},
		{
			name: "three level nesting",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.SubCommands(
					cli.CompletionSubCommand(),
					func() (*cli.Command, error) {
						return cli.New("repo",
							cli.Short("Repository commands"),
							cli.SubCommands(
								func() (*cli.Command, error) {
									return cli.New("pr",
										cli.Short("Pull request commands"),
										cli.SubCommands(
											func() (*cli.Command, error) {
												return cli.New("create",
													cli.Short("Create a pull request"),
													cli.Flag(new(string), "title", 't', "PR title"),
													cli.Run(func(_ context.Context, _ *cli.Command) error { return nil }),
												)
											},
										),
									)
								},
							),
						)
					},
				),
			},
		},
		{
			name: "mixed flag types",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				cli.Flag(new(bool), "force", 'f', "Force the operation"),
				cli.Flag(new(string), "output", 'o', "Output file"),
				cli.Flag(new(flag.Count), "verbose", 'v', "Verbosity level"),
				cli.Flag(new(bool), "dry-run", flag.NoShortHand, "Dry run, make no changes"),
				cli.Flag(new(int), "timeout", flag.NoShortHand, "Timeout in seconds"),
				cli.Flag(new([]string), "tags", flag.NoShortHand, "Tags to apply"),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
		{
			name: "yaml special chars in usage",
			options: []cli.Option{
				cli.Short("My tool"),
				cli.OverrideArgs([]string{"completion"}),
				// Colons, quotes, and parens in usage strings must be safely marshalled.
				cli.Flag(new(bool), "verbose", 'v', "Verbose: show all output"),
				cli.Flag(new(string), "config", 'c', `Config file (default: "~/.config.yaml")`),
				cli.SubCommands(cli.CompletionSubCommand()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := snapshot.New(
				t,
				snapshot.Update(*update),
			)

			stdout := &bytes.Buffer{}

			cmd, err := cli.New("mytool", slices.Concat([]cli.Option{cli.Stdout(stdout)}, tt.options)...)
			test.Ok(t, err)

			err = cmd.Execute(t.Context())
			test.Ok(t, err)

			snap.Snap(stdout.String())
		})
	}
}
