package flag

// Value is an interface representing a Flag value that can be set from the command line.
type Value interface {
	// Name returns the name of the flag.
	Name() string

	// Short returns the shorthand of the flag (or NoShortHand).
	Short() rune

	// Usage returns the usage line for the flag.
	Usage() string

	// String returns the stored value of a flag as a string.
	String() string

	// Default return the default value of a flag as a string.
	//
	// If the flag's default is the zero value for it's type,
	// an empty string is returned.
	Default() string

	// NoArgValue returns astring representation of the value of the flag when no
	// args are passed (e.g --bool implies --bool true).
	NoArgValue() string

	// Type returns the string representation of the flag type e.g. "bool".
	Type() string

	// Set sets the stored value of a flag by parsing the string "str".
	Set(str string) error
}
