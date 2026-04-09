package flag

import "go.followtheprocess.codes/cli/flag"

// Config represents the internal configuration of a [Flag].
type Config[T flag.Flaggable] struct {
	// DefaultValue holds the intended default value of the flag.
	DefaultValue T
	// EnvVar is the name of an environment variable that may set this flag's value
	// if the flag is not explicitly provided on the command line.
	EnvVar string
}
