package flag

import "go.followtheprocess.codes/cli/flag"

// Config represents the internal configuration of a [Flag].
type Config[T flag.Flaggable] struct {
	// DefaultValue holds the intended default value of the flag.
	DefaultValue T
}
