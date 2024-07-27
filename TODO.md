# TODOs

Things I want to do to make this library as good as it can be and a better, simpler, more intuitive alternative to e.g. Cobra:

- [ ] Replace `pflag.FlagSet` with my own implementation including parsing
- [ ] Make sure `cli.Option` is hidden and a `Command` can only be modified prior to building
- [ ] Ensure/double check that `cli.Option` applications are order independent
- [ ] Do some validation and maybe return an error from `cli.New` (TBC, might not need it)
- [ ] Better document the functional options, e.g. what happens if we call each one multiple times
- [x] Change `cli.Example` to be one at a time and append rather than pass a slice
- [x] Rename `cmd.example` to `cmd.examples`
- [ ] Clean up/rewrite some of the functions borrowed from Cobra e.g. `argsMinusFirstX`
- [ ] Clean up and optimise the parsing logic
- [ ] Remove the `Loop` tag from parsing functions
- [ ] More full test programs as integration tests
- [ ] Write a package example doc
- [ ] Make the help output as pretty as possible, see [clap] for inspiration as their help is so nice

[clap]: https://github.com/clap-rs/clap
