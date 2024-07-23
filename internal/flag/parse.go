package flag

import (
	"fmt"
	"strconv"
)

// signed is the same as constraints.Signed but we don't have to depend
// on golang/x/exp.
type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// unsigned is the same as constraints.Unsigned but we don't hve to depend
// on golang/x/exp.
type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// errParse is a helper to quickly return a consistent error in the face of flag
// value parsing errors.
func errParse[T any](name, str string, typ T, err error) error {
	return fmt.Errorf(
		"flag %s received invalid value %q (expected %T), detail: %w",
		name,
		str,
		typ,
		err,
	)
}

// parseInt is a generic helper to parse all signed integers, given a bit size.
//
// It returns the parsed value or an error.
func parseInt[T signed](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseInt(str, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// parseUint is a generic helper to parse all signed integers, given a bit size.
//
// It returns the parsed value or an error.
func parseUint[T unsigned](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseUint(str, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// parseFloat is a generic helper to parse floating point numbers, given a bit size.
//
// It returns the parsed value or an error.
func parseFloat[T ~float32 | ~float64](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseFloat(str, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// formatInt is a generic helper to return a string representation of any signed integer.
func formatInt[T signed](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

// formatUint is a generic helper to return a string representation of any unsigned integer.
func formatUint[T unsigned](in T) string {
	return strconv.FormatUint(uint64(in), 10)
}

// formatFloat is a generic helper to return a string representation of any floating point digit.
func formatFloat[T ~float32 | ~float64](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatFloat(float64(in), 'g', -1, bits)
	}
}