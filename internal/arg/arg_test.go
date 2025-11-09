package arg_test

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"

	"go.followtheprocess.codes/cli/internal/arg"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/parse"
	"go.followtheprocess.codes/test"
)

// TODO(@FollowTheProcess): Again a LOT of this is a straight copy paste from flag to get some confidence
// in the parsing.
//
// I think we should make an internal parse package or something and shunt all of this stuff in there as it's
// really only testing our Set logic.

func TestArgableTypes(t *testing.T) {
	// We can't do table testing here because Arg[T] is a different type for each test
	// so we can't do a []Arg[T] which is needed to define the test cases
	// so strap in for a bunch of copy pasta
	t.Run("int valid", func(t *testing.T) {
		var i int

		intArg, err := arg.New(&i, "int", "Set an int value", arg.Config[int]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intArg.Type(), "int")
		test.Equal(t, intArg.String(), "42")
		test.Equal(t, intArg.Default(), "")
	})

	t.Run("int invalid", func(t *testing.T) {
		var i int

		intArg, err := arg.New(&i, "int", "Set an int value", arg.Config[int]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int8 valid", func(t *testing.T) {
		var i int8

		intArg, err := arg.New(&i, "int", "Set an int8 value", arg.Config[int8]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int8(42))
		test.Equal(t, intArg.Type(), "int8")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("int8 invalid", func(t *testing.T) {
		var i int8

		intArg, err := arg.New(&i, "int", "Set an int8 value", arg.Config[int8]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int16 valid", func(t *testing.T) {
		var i int16

		intArg, err := arg.New(&i, "int", "Set an int16 value", arg.Config[int16]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int16(42))
		test.Equal(t, intArg.Type(), "int16")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("int16 invalid", func(t *testing.T) {
		var i int16

		intArg, err := arg.New(&i, "int", "Set an int16 value", arg.Config[int16]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int32 valid", func(t *testing.T) {
		var i int32

		intArg, err := arg.New(&i, "int", "Set an int32 value", arg.Config[int32]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int32(42))
		test.Equal(t, intArg.Type(), "int32")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("int32 invalid", func(t *testing.T) {
		var i int32

		intArg, err := arg.New(&i, "int", "Set an int32 value", arg.Config[int32]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int64 valid", func(t *testing.T) {
		var i int64

		intArg, err := arg.New(&i, "int", "Set an int64 value", arg.Config[int64]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int64(42))
		test.Equal(t, intArg.Type(), "int64")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("int64 invalid", func(t *testing.T) {
		var i int64

		intArg, err := arg.New(&i, "int", "Set an int64 value", arg.Config[int64]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint valid", func(t *testing.T) {
		var i uint

		intArg, err := arg.New(&i, "uint", "Set a uint value", arg.Config[uint]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intArg.Type(), "uint")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uint invalid", func(t *testing.T) {
		var i uint

		intArg, err := arg.New(&i, "uint", "Set a uint value", arg.Config[uint]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint8 valid", func(t *testing.T) {
		var i uint8

		intArg, err := arg.New(&i, "uint", "Set a uint8 value", arg.Config[uint8]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint8(42))
		test.Equal(t, intArg.Type(), "uint8")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uint8 invalid", func(t *testing.T) {
		var i uint8

		intArg, err := arg.New(&i, "uint", "Set a uint8 value", arg.Config[uint8]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint16 valid", func(t *testing.T) {
		var i uint16

		intArg, err := arg.New(&i, "uint", "Set a uint16 value", arg.Config[uint16]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint16(42))
		test.Equal(t, intArg.Type(), "uint16")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uint16 invalid", func(t *testing.T) {
		var i uint16

		intArg, err := arg.New(&i, "uint", "Set a uint16 value", arg.Config[uint16]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint32 valid", func(t *testing.T) {
		var i uint32

		intArg, err := arg.New(&i, "uint", "Set a uint32 value", arg.Config[uint32]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint32(42))
		test.Equal(t, intArg.Type(), "uint32")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uint32 invalid", func(t *testing.T) {
		var i uint32

		intArg, err := arg.New(&i, "uint", "Set a uint32 value", arg.Config[uint32]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint64 valid", func(t *testing.T) {
		var i uint64

		intArg, err := arg.New(&i, "uint", "Set a uint64 value", arg.Config[uint64]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint64(42))
		test.Equal(t, intArg.Type(), "uint64")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uint64 invalid", func(t *testing.T) {
		var i uint64

		intArg, err := arg.New(&i, "uint", "Set a uint64 value", arg.Config[uint64]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uintptr valid", func(t *testing.T) {
		var i uintptr

		intArg, err := arg.New(&i, "uintptr", "Set a uintptr value", arg.Config[uintptr]{})
		test.Ok(t, err)

		err = intArg.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uintptr(42))
		test.Equal(t, intArg.Type(), "uintptr")
		test.Equal(t, intArg.String(), "42")
	})

	t.Run("uintptr invalid", func(t *testing.T) {
		var i uintptr

		intArg, err := arg.New(&i, "uintptr", "Set a uintptr value", arg.Config[uintptr]{})
		test.Ok(t, err)

		err = intArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float32 valid", func(t *testing.T) {
		var f float32

		floatArg, err := arg.New(&f, "float", "Set a float32 value", arg.Config[float32]{})
		test.Ok(t, err)

		err = floatArg.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatArg.Type(), "float32")
		test.Equal(t, floatArg.String(), "3.14159")
	})

	t.Run("float32 invalid", func(t *testing.T) {
		var f float32

		floatArg, err := arg.New(&f, "float", "Set a float32 value", arg.Config[float32]{})
		test.Ok(t, err)

		err = floatArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float64 valid", func(t *testing.T) {
		var f float64

		floatArg, err := arg.New(&f, "float", "Set a float64 value", arg.Config[float64]{})
		test.Ok(t, err)

		err = floatArg.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatArg.Type(), "float64")
		test.Equal(t, floatArg.String(), "3.14159")
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64

		floatArg, err := arg.New(&f, "float", "Set a float64 value", arg.Config[float64]{})
		test.Ok(t, err)

		err = floatArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("bool valid", func(t *testing.T) {
		var b bool

		boolArg, err := arg.New(&b, "bool", "Set a bool value", arg.Config[bool]{})
		test.Ok(t, err)

		err = boolArg.Set(format.True)
		test.Ok(t, err)
		test.Equal(t, b, true)
		test.Equal(t, boolArg.Type(), "bool")
		test.Equal(t, boolArg.String(), format.True)
	})

	t.Run("bool invalid", func(t *testing.T) {
		var b bool

		boolArg, err := arg.New(&b, "bool", "Set a bool value", arg.Config[bool]{})
		test.Ok(t, err)

		err = boolArg.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	// No invalid case as all command line args are strings anyway so no real way of
	// getting an error here
	t.Run("string", func(t *testing.T) {
		var str string

		strArg, err := arg.New(&str, "string", "Set a string value", arg.Config[string]{})
		test.Ok(t, err)

		err = strArg.Set("newvalue")
		test.Ok(t, err)
		test.Equal(t, str, "newvalue")
		test.Equal(t, strArg.Type(), "string")
		test.Equal(t, strArg.String(), "newvalue")
	})

	t.Run("byte slice valid", func(t *testing.T) {
		var byt []byte

		byteArg, err := arg.New(&byt, "byte", "Set a byte slice value", arg.Config[[]byte]{})
		test.Ok(t, err)

		err = byteArg.Set("5e")
		test.Ok(t, err)
		test.EqualFunc(t, byt, []byte("^"), bytes.Equal)
		test.Equal(t, byteArg.Type(), "bytesHex")
		test.Equal(t, byteArg.String(), "5e")
	})

	t.Run("byte slice invalid", func(t *testing.T) {
		var byt []byte

		byteArg, err := arg.New(&byt, "byte", "Set a byte slice value", arg.Config[[]byte]{})
		test.Ok(t, err)

		err = byteArg.Set("0xF")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("time.Time valid", func(t *testing.T) {
		var tyme time.Time

		timeArg, err := arg.New(&tyme, "time", "Set a time value", arg.Config[time.Time]{})
		test.Ok(t, err)

		err = timeArg.Set("2024-07-17T07:38:05Z")
		test.Ok(t, err)

		want, err := time.Parse(time.RFC3339, "2024-07-17T07:38:05Z")
		test.Ok(t, err)
		test.Equal(t, tyme, want)
		test.Equal(t, timeArg.Type(), "time")
		test.Equal(t, timeArg.String(), "2024-07-17T07:38:05Z")
	})

	t.Run("time.Time invalid", func(t *testing.T) {
		var tyme time.Time

		timeArg, err := arg.New(&tyme, "time", "Set a time value", arg.Config[time.Time]{})
		test.Ok(t, err)

		err = timeArg.Set("not a time")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("time.Duration valid", func(t *testing.T) {
		var duration time.Duration

		durationArg, err := arg.New(&duration, "duration", "Set a duration value", arg.Config[time.Duration]{})
		test.Ok(t, err)

		err = durationArg.Set("300ms")
		test.Ok(t, err)

		want, err := time.ParseDuration("300ms")
		test.Ok(t, err)
		test.Equal(t, duration, want)
		test.Equal(t, durationArg.Type(), "duration")
		test.Equal(t, durationArg.String(), "300ms")
	})

	t.Run("time.Duration invalid", func(t *testing.T) {
		var duration time.Duration

		durationArg, err := arg.New(&duration, "duration", "Set a duration value", arg.Config[time.Duration]{})
		test.Ok(t, err)

		err = durationArg.Set("not a duration")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("ip valid", func(t *testing.T) {
		var ip net.IP

		ipArg, err := arg.New(&ip, "ip", "Set an IP address", arg.Config[net.IP]{})
		test.Ok(t, err)

		err = ipArg.Set("192.0.2.1")
		test.Ok(t, err)
		test.DiffBytes(t, ip, net.ParseIP("192.0.2.1"))
		test.Equal(t, ipArg.Type(), "ip")
		test.Equal(t, ipArg.String(), "192.0.2.1")
	})

	t.Run("ip invalid", func(t *testing.T) {
		var ip net.IP

		ipArg, err := arg.New(&ip, "ip", "Set an IP address", arg.Config[net.IP]{})
		test.Ok(t, err)

		err = ipArg.Set("not an ip")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})
}
