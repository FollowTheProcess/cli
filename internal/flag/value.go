package flag

// Value is an interface representing a Flag value that can be set from the command line.
type Value interface {
	String() string       // Print the stored value of a flag
	Type() string         // Return the string representation of the flag type e.g. "bool"
	Set(str string) error // Set the stored value of a flag by parsing the string "str"
}
