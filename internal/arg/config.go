package arg

import "go.followtheprocess.codes/cli/arg"

// Config represents internal configuration of an [Arg].
type Config[T arg.Argable] struct {
	// DefaultValue holds the intended default value of the argument.
	//
	// If it is nil, the argument is required.
	//
	// A non-nil value indicates the argument is not required and if not
	// provided on the command line, will assume the value DefaultValue points to.
	DefaultValue *T
}
