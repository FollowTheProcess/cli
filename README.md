# CLI

[![License](https://img.shields.io/github/license/FollowTheProcess/cli)](https://github.com/FollowTheProcess/cli)
[![Go Reference](https://pkg.go.dev/badge/go.followtheprocess.codes/cli.svg)](https://pkg.go.dev/go.followtheprocess.codes/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/cli)](https://goreportcard.com/report/github.com/FollowTheProcess/cli)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/cli?logo=github&sort=semver)](https://github.com/FollowTheProcess/cli)
[![CI](https://github.com/FollowTheProcess/cli/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/cli/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/cli)

Tiny, simple, but powerful CLI framework for modern Go 🚀

<p align="center">
<img src="https://github.com/FollowTheProcess/cli/raw/main/docs/img/demo.png" alt="demo">
</p>

> [!WARNING]
> **CLI is still in development and is not yet stable**

- [CLI](#cli)
  - [Project Description](#project-description)
  - [Installation](#installation)
  - [Quickstart](#quickstart)
  - [Usage](#usage)
    - [Commands](#commands)
    - [Sub Commands](#sub-commands)
    - [Flags](#flags)
    - [Arguments](#arguments)
  - [Core Principles](#core-principles)
    - [😱 Well behaved libraries don't panic](#-well-behaved-libraries-dont-panic)
    - [🧘🏻 Keep it Simple](#-keep-it-simple)
    - [👨🏻‍🔬 Use Modern Techniques](#-use-modern-techniques)
    - [🥹 A Beautiful API](#-a-beautiful-api)
    - [🔐 Immutable State](#-immutable-state)
    - [🚧 Good Libraries are Hard to Misuse](#-good-libraries-are-hard-to-misuse)
  - [In the Wild](#in-the-wild)

## Project Description

`cli` is a simple, minimalist, yet functional and powerful CLI framework for Go. Inspired by things like [spf13/cobra] and [urfave/cli], but building on lessons learned and using modern Go techniques and idioms.

## Installation

```shell
go get go.followtheprocess.codes/cli@latest
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "os"

    "go.followtheprocess.codes/cli"
)

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func run() error {
    var count int

    cmd, err := cli.New(
        "quickstart",
        cli.Short("Short description of your command"),
        cli.Long("Much longer text..."),
        cli.Version("v1.2.3"),
        cli.Commit("7bcac896d5ab67edc5b58632c821ec67251da3b8"),
        cli.BuildDate("2024-08-17T10:37:30Z"),
        cli.Stdout(os.Stdout),
        cli.Example("Do a thing", "quickstart something"),
        cli.Example("Count the things", "quickstart something --count 3"),
        cli.Flag(&count, "count", 'c', 0, "Count the things"),
        cli.Run(func(ctx context.Context, cmd *cli.Command) error {
            fmt.Fprintf(cmd.Stdout(), "Hello from quickstart!, my args were: %v, count was %d\n", cmd.Args(), count)
            return nil
        }),
    )
    if err != nil {
        return err
    }

    return cmd.Execute(context.Background())
}
```

Will get you the following:

![quickstart](https://github.com/FollowTheProcess/cli/raw/main/docs/img/quickstart.gif)

> [!TIP]
> See usage section below and more examples under [`./examples`](https://github.com/FollowTheProcess/cli/tree/main/examples)

## Usage

### Commands

To create CLI commands, you simply call `cli.New`:

```go
cmd, err := cli.New(
    "name", // The name of your command
    cli.Short("A new command") // Shown in the help
    cli.Run(func(ctx context.Context, cmd *cli.Command) error {
        // This function is what your command does
        fmt.Printf("name called with args: %v\n", cmd.Args())
        return nil
    })
)
```

> [!TIP]
> The command can be customised by applying any number of [functional options] for setting the help text, describing the arguments or flags it takes, adding subcommands etc. see <https://pkg.go.dev/github.com/FollowTheProcess/cli#Option>

### Sub Commands

To add a subcommand underneath the command you've just created, it's again `cli.New`:

```go
// Best to abstract it into a function
func buildSubcommand(ctx context.Context) (*cli.Command, error) {
    return cli.New(
        "sub", // Name of the sub command e.g. 'clone' for 'git clone'
        cli.Short("A sub command"),
        // etc..
    )
}
```

And add it to your parent command:

```go
ctx := context.Background()

// From the example above
cmd, err := cli.New(
    "name", // The name of your command
    // ...
    cli.SubCommands(buildSubcommand(ctx)),
)
```

This pattern can be repeated recursively to create complex command structures.

### Flags

Flags in `cli` are generic, that is, there is *one* way to add a flag to your command, and that's with the `cli.Flag` option to `cli.New`

```go
// These will get set at command line parse time
// based on your flags
type options struct {
    name string
    force bool
    size uint
    items []string
}

func buildCmd() (*cli.Command, error) {
    var opts options
    return cli.New(
        // ...
        // Signature is cli.Flag(*T, name, shorthand, default, description)
        cli.Flag(&options.name, "name", 'n', "", "The name of something"),
        cli.Flag(&options.force, "force", cli.NoShortHand, false, "Force delete without confirmation"),
        cli.Flag(&options.size, "size", 's', 0, "Size of something"),
        cli.Flag(&options.items, "items", 'i', nil, "Items to include"),
        // ...
    )
}
```

The types are all inferred automatically! No more `BoolSliceVarP` ✨

The types you can use for flags currently are:

- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `uintptr`
- `float32`
- `float64`
- `string`
- `bool`
- `[]byte` (interpreted as a hex string)
- `Count` (special type for flags that count things e.g. a `--verbosity` flag may be used like `-vvv` to increase verbosity to 3)
- `time.Time`
- `time.Duration`
- `net.IP`
- `[]int`
- `[]int8`
- `[]int16`
- `[]int32`
- `[]int64`
- `[]uint`
- `[]uint16`
- `[]uint32`
- `[]uint64`
- `[]float32`
- `[]float64`
- `[]string`

> [!NOTE]
> You basically can't get this wrong, if you try and use an unsupported type, the Go compiler will yell at you

### Arguments

There are two approaches to positional arguments in `cli`, you can either just get the raw arguments yourself with `cmd.Args()` and do whatever you want with them:

```go
cli.New(
    "my-command",
    // ...
    cli.Run(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Fprintf(cmd.Stdout(), "Hello! My arguments were: %v\n", cmd.Args())
        return nil
    })
)
```

This will return a `[]string` containing all the positional arguments to your command (not flags, they've already been parsed elsewhere!)

Or, if you want to get smarter 🧠 `cli` allows you to define *type safe* representations of your arguments, with or without default values! This follows a similar
idea to [Flags](#flags)

That works like this:

```go
// Define a struct to hold your arguments
type myArgs struct {
    name     string
    age      int
    employed bool
}

// Instantiate it
var args myArgs

// Tell cli about your arguments with the Arg option
cli.New(
    "my-command",
    // ... other options here
    cli.Arg(&args.name, "name", "The name of a person"),
    cli.Arg(&args.age, "age", "How old the person is"),
    cli.Arg(&args.employed, "employed", "Whether they are employed", cli.ArgDefault(true))
)
```

And just like [Flags](#flags), your argument types are all inferred and parsed automatically ✨, and you get nicer `--help` output too!

Have a look at the [`./examples`](https://github.com/FollowTheProcess/cli/tree/main/examples) to see more!

Other CLI libraries that do this sort of thing rely heavily on reflection and struct tags, `cli` does things differently! We prefer Go code
to be 100% type safe and statically checkable, so all this is implemented using generics and compile time checks 🚀

> [!NOTE]
> Just like flags, you can't really get this wrong. The types you can use for arguments are part of a generic constraint so using the wrong type results in a compiler error

The types you can currently use for positional args are:

- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `uintptr`
- `float32`
- `float64`
- `string`
- `*url.URL`
- `bool`
- `[]byte` (interpreted as a hex string)
- `time.Time`
- `time.Duration`
- `net.IP`

> [!WARNING]
> Slice types are not supported (yet), for those you need to use the `cmd.Args()` method to get the arguments manually. I'm working on this!

## Core Principles

When designing and implementing `cli`, I had some core goals and guiding principles for implementation.

### 😱 Well behaved libraries don't panic

`cli` validates heavily and returns errors for you to handle. By contrast [spf13/cobra] (and by extension [spf13/pflag]) panic in a number of (IMO unnecessary) conditions including:

- Duplicate subcommand
- Command adding itself as a subcommand
- Duplicate flag
- Invalid shorthand flag letter

The design of `cli` is such that commands are instantiated with `cli.New` and a number of [functional options]. These options are in charge of configuring your command and each will perform validation prior to applying the setting.

These errors are joined and bubbled up to you in one go via `cli.New` so you don't have to play error whack-a-mole, and more importantly your application won't panic!

### 🧘🏻 Keep it Simple

`cli` has an intentionally small public interface and gives you only what you need to build amazing CLI apps:

- No huge structs with hundreds of fields
- No confusing or conflicting options
- Customisation in areas where it makes sense, sensible opinionated defaults everywhere else
- No reflection or struct tags

There is one and only one way to do things (and that is *usually* to use an option in `cli.New`)

### 👨🏻‍🔬 Use Modern Techniques

The dominant Go CLI toolkits were mostly built many years (and many versions of Go) ago. They are reliable and battle hardened but because of their high number of users, they have had to be very conservative with changes.

`cli` has none of these constraints and can use bang up to date Go techniques and idioms.

One example is generics, consider how you define a flag:

```go
var force bool
cli.New("demo", cli.Flag(&force, "force", 'f', false, "Force something"))
```

Note the type `bool` is inferred by `cli.Flag`. This will work with any type allowed by the `Flaggable` generic constraint so you'll get compile time feedback if you've got it wrong. No more `flag.BoolStringSliceVarP` 🎉

### 🥹 A Beautiful API

`cli` heavily leverages the [functional options] pattern to create a delightful experience building a CLI tool. It almost reads like plain english:

```go
var count int
cmd, err := cli.New(
    "test",
    cli.Short("Short description of your command"),
    cli.Long("Much longer text..."),
    cli.Version("v1.2.3"),
    cli.Stdout(os.Stdout),
    cli.Example("Do a thing", "test run thing --now"),
    cli.Flag(&count, "count", 'c', 0, "Count the things"),
)
```

### 🔐 Immutable State

Typically, commands are implemented as a big struct with lots of fields. `cli` is no different in this regard.

What *is* different though is that this large struct can **only** be configured with `cli.New`. Once you've built your command, it can't be modified.

This eliminates a whole class of bugs and prevents misconfiguration and footguns 🔫

### 🚧 Good Libraries are Hard to Misuse

Everything in `cli` is (hopefully) clear, intuitive, and well-documented. There's a tonne of strict validation in a bunch of places and wherever possible, misuse results in a compilation error.

Consider the following example of a bad shorthand value:

```go
var delete bool

// Note: "de" is a bad shorthand, it's two letters
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
cli.New("demo", cli.Flag(&delete, "delete", cli.NoShortHand, false, "Delete something"))
```

## In the Wild

I built `cli` for my own uses really, so I've quickly adopted it across a number of tools. See the following projects for some working examples in real code:

- <https://github.com/FollowTheProcess/txtract>
- <https://github.com/FollowTheProcess/tag>
- <https://github.com/FollowTheProcess/gowc>
- <https://github.com/FollowTheProcess/zap>

[spf13/cobra]: https://github.com/spf13/cobra
[spf13/pflag]: https://github.com/spf13/pflag
[urfave/cli]: https://github.com/urfave/cli
[functional options]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
