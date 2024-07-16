# CLI

[![License](https://img.shields.io/github/license/FollowTheProcess/cli)](https://github.com/FollowTheProcess/cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/cli.svg)](https://pkg.go.dev/github.com/FollowTheProcess/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/cli)](https://goreportcard.com/report/github.com/FollowTheProcess/cli)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/cli?logo=github&sort=semver)](https://github.com/FollowTheProcess/cli)
[![CI](https://github.com/FollowTheProcess/cli/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/cli/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/cli)

Tiny, simple, minimal CLI framework for Go

> [!WARNING]
> **CLI is in early development and is not yet ready for use**

## Project Description

`cli` is my attempt at making a very simple, minimalist, yet functional and powerful CLI framework for Go. Inspired by things like <https://github.com/spf13/cobra> and <https://github.com/urfave/cli>.

*"So why make a new one?"*

Good question 🤔

1) I love Cobra and use it in a lot of my own programs, but IMO it's unnecessarily complex. The `Command` struct has lots of fields, has many ways to do things, and a lot of options I never use. It's a fantastic library, but has suffered from success in the sense that everybody wants it to do everything, rather than remaining tightly focussed on the essentials
2) I feel like a lot of time has passed since this was first made and now we have generics and some other cool modern Go idioms, a CLI toolkit could be better designed and have a friendlier API, with less options and one clear way to do things
3) It's fun to implement cool things from scratch and stretch your programming legs on a new problem

> [!NOTE]
> Currently flag parsing is provided by <https://github.com/spf13/pflag> (same as Cobra) while I iterate and make the overall CLI experience as nice as I can. I do plan to add my own flag parsing capability so that `cli` is a completely self contained toolkit.

## Installation

```shell
go get github.com/FollowTheProcess/cli@latest
```

## Quickstart

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
