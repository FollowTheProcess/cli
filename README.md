# CLI

[![License](https://img.shields.io/github/license/FollowTheProcess/cli)](https://github.com/FollowTheProcess/cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/cli.svg)](https://pkg.go.dev/github.com/FollowTheProcess/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/cli)](https://goreportcard.com/report/github.com/FollowTheProcess/cli)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/cli?logo=github&sort=semver)](https://github.com/FollowTheProcess/cli)
[![CI](https://github.com/FollowTheProcess/cli/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/cli/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/cli)

Tiny, simple, but powerful CLI framework for modern Go 🚀

> [!WARNING]
> **CLI is in early development and is not yet ready for use**

![caution](img/caution.png)

- [CLI](#cli)
  - [Project Description](#project-description)
    - [Core Principles](#core-principles)
      - [😱 Well behaved libraries don't panic](#-well-behaved-libraries-dont-panic)
      - [🧘🏻 Keep it Simple](#-keep-it-simple)
      - [👨🏻‍🔬 Use Modern Techniques](#-use-modern-techniques)
      - [🥹 A Beautiful API](#-a-beautiful-api)
      - [🔐 Immutable State](#-immutable-state)
      - [🚧 Good Libraries are Hard to Misuse](#-good-libraries-are-hard-to-misuse)
  - [Installation](#installation)
  - [Quickstart](#quickstart)
    - [Credits](#credits)

## Project Description

`cli` is a simple, minimalist, yet functional and powerful CLI framework for Go. Inspired by things like [spf13/cobra] and [urfave/cli], but building on lessons learned and using modern Go techniques and idioms.

### Core Principles

#### 😱 Well behaved libraries don't panic

`cli` validates heavily and returns errors for you to handle. By contrast [spf13/cobra] (and [spf13/pflag]) panic in a number of conditions including:

- Duplicate subcommand being added
- Duplicate flag added
- Invalid shorthand flag letter

The design of `cli` is such that commands are instantiated with `cli.New` and a number of [functional options]. These options are in charge of configuring your command and each will perform validation prior to applying the setting.

These errors are joined and bubbled up to you in one go via `cli.New` so you don't have to play error whack-a-mole.

#### 🧘🏻 Keep it Simple

`cli` has a tiny public interface and gives you only what you need to build amazing CLI apps, no more confusing options and hundreds of struct fields.

There is one and only one way to do things (and that is *usually* to use an option in `cli.New`)

#### 👨🏻‍🔬 Use Modern Techniques

The dominant Go CLI toolkits were mostly built many years (and many versions of Go) ago. They are reliable and battle hardened but because of their high number of users, they have had to be very conservative with changes.

`cli` has none of these constraints and can use bang up to date Go techniques and idioms.

One example is generics, consider how you define a flag:

```go
var force bool
cli.New("demo", cli.Flag(&force, "force", 'f', false, "Force something"))
```

Note the type `bool` is inferred by `cli.Flag`. No more `flag.BoolStringSliceVarP` 🎉

#### 🥹 A Beautiful API

`cli` heavily leverages the [functional options] pattern to create a delightful experience building a CLI tool. It almost reads like plain english:

```go
var count int
cmd, err := cli.New(
    "test",
    cli.Short("Short description of your command"),
    cli.Long("Much longer text..."),
    cli.Version("v1.2.3"),
    cli.Allow(cli.MinArgs(1)),
    cli.Stdout(os.Stdout),
    cli.Example("Do a thing", "test run thing --now"),
    cli.Flag(&count, "count", 'c', 0, "Count the things"),
)
```

#### 🔐 Immutable State

Typically, these sorts of things are implemented with a big struct with lots of fields. `cli` is no different in this regard.

What *is* different though is that this large struct can **only** be configured with `cli.New`. Once you've built your command, it can't be modified.

This eliminates a whole class of bugs and prevents misconfiguration and footguns 🔫

#### 🚧 Good Libraries are Hard to Misuse

Everything in `cli` is (hopefully) clear, intuitive, and well-documented. There's a tonne of strict validation in a bunch of places and wherever possible, misuse results in a compilation error.

Consider the following example of a bad shorthand value:

```go
var delete bool

// Note: bad shorthand, it's two letters
cli.New("demo", cli.Flag(&delete, "delete", "de", false, "Delete something"))
```

In `cli` this is impossible as we use `rune` as the type for a flag shorthand, so the above example would not compile. Instead you must specify a valid rune:

```go
var delete bool

// Ahhh, that's better
cli.New("demo", cli.Flag(&delete, "delete", 'd', false, "Delete something"))
```

And if you don't want a shorthand? i.e. just `--delete` with no `-d` option:

```go
var delete bool
cli.New("demo", cli.Flag(&delete, "delete", flag.NoShortHand, false, "Delete something"))
```

## Installation

Convinced? I hope so...

```shell
go get github.com/FollowTheProcess/cli@latest
```

## Quickstart

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[spf13/cobra]: https://github.com/spf13/cobra
[spf13/pflag]: https://github.com/spf13/pflag
[urfave/cli]: https://github.com/urfave/cli
[functional options]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
