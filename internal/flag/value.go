package flag

// Value is an interface representing a Flag value that can be set from the command line.
type Value interface {
	Name() string         // Return the name of the flag
	Short() rune          // Return the shorthand of the flag (or NoShortHand)
	Usage() string        // Return the usage line for the flag
	String() string       // Print the stored value of a flag
	NoArgValue() string   // String representation of the value of the flag when no args are passed (e.g --bool implies --bool true)
	Type() string         // Return the string representation of the flag type e.g. "bool"
	Set(str string) error // Set the stored value of a flag by parsing the string "str"
}
