package flag

// Value is an interface representing a Flag value that can be set from the command line.
type Value interface {
	// Return the name of the flag
	Name() string

	// Return the shorthand of the flag (or NoShortHand)
	Short() rune

	// Return the usage line for the flag
	Usage() string

	// Print the stored value of a flag
	String() string

	// String representation of the value of the flag when no args are passed (e.g --bool implies --bool true)
	NoArgValue() string

	// Return the string representation of the flag type e.g. "bool"
	Type() string

	// Set the stored value of a flag by parsing the string "str"
	Set(str string) error
}
