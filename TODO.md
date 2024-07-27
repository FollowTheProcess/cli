# TODOs

Things I want to do to make this library as good as it can be and a better, simpler, more intuitive alternative to e.g. Cobra:

- [ ] Replace `pflag.FlagSet` with my own implementation including parsing
- [ ] Ensure/double check that `cli.Option` applications are order independent
- [ ] Do some validation and maybe return an error from `cli.New` (TBC, might not need it)
- [ ] Better document the functional options, e.g. what happens if we call each one multiple times
- [ ] Change `cli.Example` to be one at a time and append rather than pass a slice
- [ ] Rename `cmd.example` to `cmd.examples`
- [ ] Clean up/rewrite some of the functions borrowed from Cobra e.g. `argsMinusFirstX`
- [ ] Remove the `Loop` tag from parsing functions
- [ ] More full test programs as integration tests
- [ ] Write a package example doc
