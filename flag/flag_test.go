package flag_test

import (
	"testing"

	"github.com/FollowTheProcess/cli/flag"
	"github.com/FollowTheProcess/test"
)

func TestFlagValueSet(t *testing.T) {
	// We can't do table testing here because Flag[T] is a different type for each test
	// so we can't do a []Flag[T] which is needed to define the test cases
	// so strap in for a bunch of copy pasta
	t.Run("int valid", func(t *testing.T) {
		var i int
		intFlag := flag.New(&i, "int", "i", 0, "Set an int value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), 42)
	})

	t.Run("int invalid", func(t *testing.T) {
		var i int
		intFlag := flag.New(&i, "int", "i", 0, "Set an int value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag int received invalid value "word" (expected int)`)
	})

	t.Run("int8 valid", func(t *testing.T) {
		var i int8
		intFlag := flag.New(&i, "int8", "i", 0, "Set an int8 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int8(42))
	})

	t.Run("int8 invalid", func(t *testing.T) {
		var i int8
		intFlag := flag.New(&i, "int8", "i", 0, "Set an int8 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag int8 received invalid value "word" (expected int8)`)
	})

	t.Run("int16 valid", func(t *testing.T) {
		var i int16
		intFlag := flag.New(&i, "int16", "i", 0, "Set an int16 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int16(42))
	})

	t.Run("int16 invalid", func(t *testing.T) {
		var i int16
		intFlag := flag.New(&i, "int16", "i", 0, "Set an int16 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag int16 received invalid value "word" (expected int16)`)
	})

	t.Run("int32 valid", func(t *testing.T) {
		var i int32
		intFlag := flag.New(&i, "int32", "i", 0, "Set an int32 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int32(42))
	})

	t.Run("int32 invalid", func(t *testing.T) {
		var i int32
		intFlag := flag.New(&i, "int32", "i", 0, "Set an int32 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag int32 received invalid value "word" (expected int32)`)
	})

	t.Run("int64 valid", func(t *testing.T) {
		var i int64
		intFlag := flag.New(&i, "int64", "i", 0, "Set an int64 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int64(42))
	})

	t.Run("int64 invalid", func(t *testing.T) {
		var i int64
		intFlag := flag.New(&i, "int64", "i", 0, "Set an int64 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag int64 received invalid value "word" (expected int64)`)
	})

	t.Run("uint valid", func(t *testing.T) {
		var i uint
		intFlag := flag.New(&i, "uint", "i", 0, "Set a uint value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), 42)
	})

	t.Run("uint invalid", func(t *testing.T) {
		var i uint
		intFlag := flag.New(&i, "uint", "i", 0, "Set a uint value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uint received invalid value "word" (expected uint)`)
	})

	t.Run("uint8 valid", func(t *testing.T) {
		var i uint8
		intFlag := flag.New(&i, "uint8", "i", 0, "Set a uint8 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint8(42))
	})

	t.Run("uint8 invalid", func(t *testing.T) {
		var i uint8
		intFlag := flag.New(&i, "uint8", "i", 0, "Set a uint8 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uint8 received invalid value "word" (expected uint8)`)
	})

	t.Run("uint16 valid", func(t *testing.T) {
		var i uint16
		intFlag := flag.New(&i, "uint16", "i", 0, "Set a uint16 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint16(42))
	})

	t.Run("uint16 invalid", func(t *testing.T) {
		var i uint16
		intFlag := flag.New(&i, "uint16", "i", 0, "Set a uint16 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uint16 received invalid value "word" (expected uint16)`)
	})

	t.Run("uint32 valid", func(t *testing.T) {
		var i uint32
		intFlag := flag.New(&i, "uint32", "i", 0, "Set a uint32 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint32(42))
	})

	t.Run("uint32 invalid", func(t *testing.T) {
		var i uint32
		intFlag := flag.New(&i, "uint32", "i", 0, "Set a uint32 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uint32 received invalid value "word" (expected uint32)`)
	})

	t.Run("uint64 valid", func(t *testing.T) {
		var i uint64
		intFlag := flag.New(&i, "uint64", "i", 0, "Set a uint64 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint64(42))
	})

	t.Run("uint64 invalid", func(t *testing.T) {
		var i uint64
		intFlag := flag.New(&i, "uint64", "i", 0, "Set a uint64 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uint64 received invalid value "word" (expected uint64)`)
	})

	t.Run("uintptr valid", func(t *testing.T) {
		var i uintptr
		intFlag := flag.New(&i, "uintptr", "i", 0, "Set a uintptr value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uintptr(42))
	})

	t.Run("uintptr invalid", func(t *testing.T) {
		var i uintptr
		intFlag := flag.New(&i, "uintptr", "i", 0, "Set a uintptr value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag uintptr received invalid value "word" (expected uintptr)`)
	})

	t.Run("float32 valid", func(t *testing.T) {
		var f float32
		floatFlag := flag.New(&f, "float32", "f", 0, "Set a float32 value")
		err := floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, floatFlag.Get(), 3.14159)
	})

	t.Run("float32 invalid", func(t *testing.T) {
		var f float32
		floatFlag := flag.New(&f, "float32", "f", 0, "Set a float32 value")
		err := floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag float32 received invalid value "word" (expected float32)`)
	})

	t.Run("float64 valid", func(t *testing.T) {
		var f float64
		floatFlag := flag.New(&f, "float64", "f", 0, "Set a float64 value")
		err := floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, floatFlag.Get(), 3.14159)
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64
		floatFlag := flag.New(&f, "float64", "f", 0, "Set a float64 value")
		err := floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(t, err.Error(), `flag float64 received invalid value "word" (expected float64)`)
	})

	// No invalid case as all command line args are strings anyway so no real way of
	// getting an error here
	t.Run("string", func(t *testing.T) {
		var str string
		strFlag := flag.New(&str, "string", "s", "", "Set a string value")
		err := strFlag.Set("newvalue")
		test.Ok(t, err)
		test.Equal(t, strFlag.Get(), "newvalue")
	})
}
