# Examples

This directory contains a bunch of example programs built with `cli` to show you how the library works and how you might implement common patterns.

- [Examples](#examples)
  - [`./cover`](#cover)
  - [`./quickstart`](#quickstart)
  - [`./subcommands`](#subcommands)
  - [`./namedargs`](#namedargs)
  - [`./cancel`](#cancel)
  - [`./completion`](#completion)
    - [TODO](#todo)

## `./cover`

Not really an example, but holds the source code used to generate the cover image with [freeze]

## `./quickstart`

Implements the [quickstart] command from the main project README

![quickstart](../docs/img/quickstart.gif)

## `./subcommands`

A CLI with multiple subcommands, each with their own flags and expected arguments. Shows how to easily store parsed flag values in an options struct and pass them around your program.

![subcommands](../docs/img/subcommands.gif)

## `./namedargs`

A CLI with named positional arguments that may or may not have default values. Shows how to retrieve these arguments by name and use them without having to care if they were provided or you're using the default.

![namedargs](../docs/img/namedargs.gif)

## `./cancel`

This examples shows how `cli` requiring you to pass a `context.Context` to your run functions leads to elegant and resilient cancellation and `CTRL+C` handling.

![cancel](../docs/img/cancel.gif)

## `./completion`

Demonstrates how to add shell completion support to any CLI by wiring in `cli.CompletionSubCommand()`. Running `mytool completion` outputs a [carapace-spec] YAML document describing the full command tree. Redirect it once to register completions with [carapace-bin]:

```shell
mytool completion > ~/.config/carapace/specs/mytool.yaml
```

carapace-bin then provides completions across bash, zsh, fish, nushell, PowerShell, and more — no shell-specific scripts required.

### TODO

- Replicate one or two well known CLI tools as an example
  - Docker
  - Cargo

[quickstart]: <https://github.com/FollowTheProcess/cli#quickstart>
[freeze]: <https://github.com/charmbracelet/freeze>
[carapace-spec]: <https://github.com/carapace-sh/carapace-spec>
[carapace-bin]: <https://github.com/carapace-sh/carapace-bin>
