# TODOs

Things I want to do to make this library as good as it can be and a better, simpler, more intuitive alternative to e.g. Cobra:

- [x] Replace `pflag.FlagSet` with my own implementation including parsing
- [x] Make sure `cli.Option` is hidden and a `Command` can only be modified prior to building
- [x] Ensure/double check that `cli.Option` applications are order independent
- [x] Do some validation and maybe return an error from `cli.New` (TBC, might not need it)
- [x] Better document the functional options, e.g. what happens if we call each one multiple times
- [x] Change `cli.Example` to be one at a time and append rather than pass a slice
- [x] Rename `cmd.example` to `cmd.examples`
- [x] Clean up/rewrite some of the functions borrowed from Cobra e.g. `argsMinusFirstX`
- [ ] Clean up and optimise the parsing logic
- [x] Remove the `Loop` tag from parsing functions
- [ ] More full test programs as integration tests
- [ ] Write a package example doc
- [ ] Make the help output as pretty as possible, see [clap] for inspiration as their help is so nice
- [x] Implement something nice to do the whole `-- <bonus args>` thing, maybe `ExtraArgs`?
- [x] Thin wrapper around tabwriter to keep it consistent

[clap]: https://github.com/clap-rs/clap
