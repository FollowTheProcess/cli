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
- [x] Do some benchmarking and see where we can improve
- [x] Remove the `Loop` tag from parsing functions
- [ ] More full test programs as integration tests
- [ ] Write a package example doc
- [ ] Write some fully formed examples under `./examples` for people to see as inspiration
- [x] Make the help output as pretty as possible, see [clap] for inspiration as their help is so nice
- [x] Implement something nice to do the whole `-- <bonus args>` thing, maybe `ExtraArgs`?
- [x] Thin wrapper around tabwriter to keep it consistent
- [ ] Try this on some of my CLI tools to work out the bugs in real world programs
- [ ] Test that shuffles the order of options and ensures the command we get is the same every time
- [x] Use [vhs] to make a nice demo gif

## Ideas

Things that might be good but not committed to doing yet and/or haven't decided if they are worth it.

- [ ] Wrap long description at a sensible limit
- [ ] Named arguments e.g. `git clone REPO` could retrieve an arg by name and would be good for the help text
- [ ] Make a CLI cli (lol) like cobra-cli to generate CLI apps?

[clap]: https://github.com/clap-rs/clap
[vhs]: https://github.com/charmbracelet/vhs
