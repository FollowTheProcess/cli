package flag

// Value represents any value that can be set by a command line flag.
type Value interface {
	Set(str string) error // Set the value for a flag from the command line
}
