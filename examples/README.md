# Examples

This directory contains a bunch of example programs built with `cli` to show you how the library works and how you might implement common patterns.

- [Examples](#examples)
  - [`./cover`](#cover)
  - [`./quickstart`](#quickstart)
  - [`./subcommands`](#subcommands)

## `./cover`

Not really an example, but holds the source code used to generate the cover image with [silicon]

## `./quickstart`

Implements the [quickstart] command from the main project README

## `./subcommands`

A CLI with multiple subcommands, each with their own flags and expected arguments. Shows how to easily store parsed flag values in an options struct and pass them around your program.

[quickstart]: https://github.com/FollowTheProcess/cli#quickstart
[silicon]: https://github.com/Aloxaf/silicon
