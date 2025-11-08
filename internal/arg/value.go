package arg

// Value is an interface representing an Arg value that can be set from the command line.
type Value interface {
	// Name returns the name of the argument.
	Name() string

	// Usage returns the usage line for the argument.
	Usage() string

	// String returns the stored value of the argument as a string.
	String() string

	// Type returns the string representation of the argument type e.g. "bool".
	Type() string

	// Set sets the stored value of an arg by parsing the string "str".
	Set(str string) error

	// Default returns the default value as a string, or "" if the argument
	// is required.
	Default() string
}
