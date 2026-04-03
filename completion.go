package cli

import (
	"context"
	"fmt"
	"slices"
	"strings"

	publicflag "go.followtheprocess.codes/cli/flag"
	"go.yaml.in/yaml/v4"
)

const completionLong = `
"Outputs a carapace-spec YAML document describing this command's flags and subcommands.

Redirect the output to register completions with carapace-bin:

  mytool completion > ~/.config/carapace/specs/mytool.yaml",
`

// CompletionSubCommand returns a [Builder] that constructs a "completion" subcommand.
//
// When run, it outputs a carapace-spec YAML document to stdout describing the full
// command tree, all flags, and all subcommands.
//
// Wire it into your root command:
//
//	cli.New("mytool",
//	    cli.SubCommands(cli.CompletionSubCommand()),
//	    ...
//	)
//
// Users register completions by running:
//
//	mytool completion > ~/.config/carapace/specs/mytool.yaml
//
// carapace-bin then provides completions across bash, zsh, fish, nushell,
// powershell, and more. See https://github.com/carapace-sh/carapace-spec
// for how to extend the generated YAML with semantic completion hints.
func CompletionSubCommand() Builder {
	return func() (*Command, error) {
		return New(
			"completion",
			Short("Output a carapace-spec YAML document to stdout"),
			Long(completionLong),
			Run(func(_ context.Context, cmd *Command) error {
				data, err := marshalCompletionSpec(cmd.root())
				if err != nil {
					return fmt.Errorf("generating carapace spec: %w", err)
				}

				_, err = cmd.Stdout().Write(data)

				return err
			}),
		)
	}
}

// specCommand is the internal representation of a carapace-spec YAML command node.
//
// See https://github.com/carapace-sh/carapace-spec for the full schema.
type specCommand struct {
	Name            string            `yaml:"name"`
	Description     string            `yaml:"description,omitempty"`
	Flags           map[string]string `yaml:"flags,omitempty"`
	PersistentFlags map[string]string `yaml:"persistentflags,omitempty"`
	Commands        []specCommand     `yaml:"commands,omitempty"`
}

// marshalCompletionSpec generates a carapace-spec YAML document for cmd and its
// full subcommand tree.
func marshalCompletionSpec(cmd *Command) ([]byte, error) {
	spec := buildSpecTree(cmd)

	// Help and version are on every command; declare them once as persistentflags.
	spec.PersistentFlags = persistentFlagsFrom(cmd)

	data, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("marshalling carapace spec: %w", err)
	}

	return data, nil
}

// isSystemFlag reports whether name is an automatically-injected flag.
// These flags appear on every command and are emitted as persistentflags
// on the root rather than repeated in every command's flags map.
func isSystemFlag(name string) bool {
	return name == "help" || name == "version"
}

// buildSpecTree recursively builds a specCommand tree from cmd.
func buildSpecTree(cmd *Command) specCommand {
	spec := specCommand{
		Name:        cmd.name,
		Description: cmd.short,
		Flags:       specFlagsFrom(cmd),
	}

	if len(cmd.subcommands) > 0 {
		sorted := make([]*Command, len(cmd.subcommands))
		copy(sorted, cmd.subcommands)
		slices.SortFunc(sorted, func(a, b *Command) int {
			return strings.Compare(a.name, b.name)
		})

		spec.Commands = make([]specCommand, len(sorted))
		for i, sub := range sorted {
			spec.Commands[i] = buildSpecTree(sub)
		}
	}

	return spec
}

// specFlagsFrom builds the carapace-spec flags map for cmd, excluding system
// flags (help, version) which are emitted as persistentflags on the root.
func specFlagsFrom(cmd *Command) map[string]string {
	flags := make(map[string]string)

	for name, fl := range cmd.flagSet().All() {
		if isSystemFlag(name) {
			continue
		}

		key := specFlagKey(name, fl.Short(), fl.Type(), fl.NoArgValue())
		flags[key] = fl.Usage()
	}

	if len(flags) == 0 {
		return nil
	}

	return flags
}

// persistentFlagsFrom builds the carapace-spec persistentflags map from the
// system flags (help, version) that are automatically added to every command.
func persistentFlagsFrom(cmd *Command) map[string]string {
	flags := make(map[string]string)

	for name, fl := range cmd.flagSet().All() {
		if !isSystemFlag(name) {
			continue
		}

		key := specFlagKey(name, fl.Short(), fl.Type(), fl.NoArgValue())
		flags[key] = fl.Usage()
	}

	if len(flags) == 0 {
		return nil
	}

	return flags
}

// specFlagKey encodes a flag name, shorthand, and type into a carapace-spec
// flag map key.
//
// Modifier selection (in order):
//
//	bool                    → no suffix    ("--force" or "-f, --force")
//	count                   → *            ("--verbose*" or "-v, --verbose*")
//	noArgValue != ""        → ?            (optional-value flag)
//	everything else         → =            ("--output=" or "-o, --output=")
func specFlagKey(name string, short rune, typ, noArgValue string) string {
	var base string
	if short != publicflag.NoShortHand {
		base = fmt.Sprintf("-%s, --%s", string(short), name)
	} else {
		base = "--" + name
	}

	switch {
	case typ == "bool":
		return base
	case typ == "count":
		return base + "*"
	case noArgValue != "":
		return base + "?"
	default:
		return base + "="
	}
}
