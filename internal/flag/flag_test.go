package flag_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/FollowTheProcess/test"
)

func TestFlagValue(t *testing.T) {
	// We can't do table testing here because Flag[T] is a different type for each test
	// so we can't do a []Flag[T] which is needed to define the test cases
	// so strap in for a bunch of copy pasta
	t.Run("int valid", func(t *testing.T) {
		var i int
		intFlag := flag.New(&i, "int", "i", 0, "Set an int value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), 42)
		test.Equal(t, intFlag.Type(), "int")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int invalid", func(t *testing.T) {
		var i int
		intFlag := flag.New(&i, "int", "i", 0, "Set an int value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag int received invalid value "word" (expected int), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int8 valid", func(t *testing.T) {
		var i int8
		intFlag := flag.New(&i, "int8", "i", 0, "Set an int8 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int8(42))
		test.Equal(t, intFlag.Type(), "int8")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int8 invalid", func(t *testing.T) {
		var i int8
		intFlag := flag.New(&i, "int8", "i", 0, "Set an int8 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag int8 received invalid value "word" (expected int8), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int16 valid", func(t *testing.T) {
		var i int16
		intFlag := flag.New(&i, "int16", "i", 0, "Set an int16 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int16(42))
		test.Equal(t, intFlag.Type(), "int16")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int16 invalid", func(t *testing.T) {
		var i int16
		intFlag := flag.New(&i, "int16", "i", 0, "Set an int16 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag int16 received invalid value "word" (expected int16), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int32 valid", func(t *testing.T) {
		var i int32
		intFlag := flag.New(&i, "int32", "i", 0, "Set an int32 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int32(42))
		test.Equal(t, intFlag.Type(), "int32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int32 invalid", func(t *testing.T) {
		var i int32
		intFlag := flag.New(&i, "int32", "i", 0, "Set an int32 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag int32 received invalid value "word" (expected int32), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int64 valid", func(t *testing.T) {
		var i int64
		intFlag := flag.New(&i, "int64", "i", 0, "Set an int64 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), int64(42))
		test.Equal(t, intFlag.Type(), "int64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int64 invalid", func(t *testing.T) {
		var i int64
		intFlag := flag.New(&i, "int64", "i", 0, "Set an int64 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag int64 received invalid value "word" (expected int64), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint valid", func(t *testing.T) {
		var i uint
		intFlag := flag.New(&i, "uint", "i", 0, "Set a uint value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), 42)
		test.Equal(t, intFlag.Type(), "uint")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint invalid", func(t *testing.T) {
		var i uint
		intFlag := flag.New(&i, "uint", "i", 0, "Set a uint value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uint received invalid value "word" (expected uint), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint8 valid", func(t *testing.T) {
		var i uint8
		intFlag := flag.New(&i, "uint8", "i", 0, "Set a uint8 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint8(42))
		test.Equal(t, intFlag.Type(), "uint8")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint8 invalid", func(t *testing.T) {
		var i uint8
		intFlag := flag.New(&i, "uint8", "i", 0, "Set a uint8 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uint8 received invalid value "word" (expected uint8), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint16 valid", func(t *testing.T) {
		var i uint16
		intFlag := flag.New(&i, "uint16", "i", 0, "Set a uint16 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint16(42))
		test.Equal(t, intFlag.Type(), "uint16")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint16 invalid", func(t *testing.T) {
		var i uint16
		intFlag := flag.New(&i, "uint16", "i", 0, "Set a uint16 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uint16 received invalid value "word" (expected uint16), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint32 valid", func(t *testing.T) {
		var i uint32
		intFlag := flag.New(&i, "uint32", "i", 0, "Set a uint32 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint32(42))
		test.Equal(t, intFlag.Type(), "uint32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint32 invalid", func(t *testing.T) {
		var i uint32
		intFlag := flag.New(&i, "uint32", "i", 0, "Set a uint32 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uint32 received invalid value "word" (expected uint32), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint64 valid", func(t *testing.T) {
		var i uint64
		intFlag := flag.New(&i, "uint64", "i", 0, "Set a uint64 value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uint64(42))
		test.Equal(t, intFlag.Type(), "uint64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint64 invalid", func(t *testing.T) {
		var i uint64
		intFlag := flag.New(&i, "uint64", "i", 0, "Set a uint64 value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uint64 received invalid value "word" (expected uint64), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uintptr valid", func(t *testing.T) {
		var i uintptr
		intFlag := flag.New(&i, "uintptr", "i", 0, "Set a uintptr value")
		err := intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, intFlag.Get(), uintptr(42))
		test.Equal(t, intFlag.Type(), "uintptr")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uintptr invalid", func(t *testing.T) {
		var i uintptr
		intFlag := flag.New(&i, "uintptr", "i", 0, "Set a uintptr value")
		err := intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag uintptr received invalid value "word" (expected uintptr), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("float32 valid", func(t *testing.T) {
		var f float32
		floatFlag := flag.New(&f, "float32", "f", 0, "Set a float32 value")
		err := floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, floatFlag.Get(), 3.14159)
		test.Equal(t, floatFlag.Type(), "float32")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float32 invalid", func(t *testing.T) {
		var f float32
		floatFlag := flag.New(&f, "float32", "f", 0, "Set a float32 value")
		err := floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag float32 received invalid value "word" (expected float32), detail: strconv.ParseFloat: parsing "word": invalid syntax`,
		)
	})

	t.Run("float64 valid", func(t *testing.T) {
		var f float64
		floatFlag := flag.New(&f, "float64", "f", 0, "Set a float64 value")
		err := floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, floatFlag.Get(), 3.14159)
		test.Equal(t, floatFlag.Type(), "float64")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64
		floatFlag := flag.New(&f, "float64", "f", 0, "Set a float64 value")
		err := floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag float64 received invalid value "word" (expected float64), detail: strconv.ParseFloat: parsing "word": invalid syntax`,
		)
	})

	t.Run("bool valid", func(t *testing.T) {
		var b bool
		boolFlag := flag.New(&b, "bool", "b", false, "Set a bool value")
		err := boolFlag.Set("true")
		test.Ok(t, err)
		test.Equal(t, boolFlag.Get(), true)
		test.Equal(t, boolFlag.Type(), "bool")
		test.Equal(t, boolFlag.String(), "true")
	})

	t.Run("bool invalid", func(t *testing.T) {
		var b bool
		boolFlag := flag.New(&b, "bool", "b", false, "Set a bool value")
		err := boolFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag bool received invalid value "word" (expected bool), detail: strconv.ParseBool: parsing "word": invalid syntax`,
		)
	})

	// No invalid case as all command line args are strings anyway so no real way of
	// getting an error here
	t.Run("string", func(t *testing.T) {
		var str string
		strFlag := flag.New(&str, "string", "s", "", "Set a string value")
		err := strFlag.Set("newvalue")
		test.Ok(t, err)
		test.Equal(t, strFlag.Get(), "newvalue")
		test.Equal(t, strFlag.Type(), "string")
		test.Equal(t, strFlag.String(), "newvalue")
	})

	t.Run("byte slice valid", func(t *testing.T) {
		var byt []byte
		byteFlag := flag.New(&byt, "byte", "b", []byte(""), "Set a byte slice value")
		err := byteFlag.Set("5e")
		test.Ok(t, err)
		test.EqualFunc(t, byteFlag.Get(), []byte("^"), bytes.Equal)
		test.Equal(t, byteFlag.Type(), "bytesHex")
		test.Equal(t, byteFlag.String(), "5e")
	})

	t.Run("byte slice invalid", func(t *testing.T) {
		var byt []byte
		byteFlag := flag.New(&byt, "byte", "b", []byte(""), "Set a byte slice value")
		err := byteFlag.Set("0xF")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag byte received invalid value "0xF" (expected []uint8), detail: encoding/hex: invalid byte: U+0078 'x'`,
		)
	})

	t.Run("time.Time valid", func(t *testing.T) {
		var tyme time.Time
		timeFlag := flag.New(&tyme, "time", "t", time.Now(), "Set a time value")
		err := timeFlag.Set("2024-07-17T07:38:05Z")
		test.Ok(t, err)

		want, err := time.Parse(time.RFC3339, "2024-07-17T07:38:05Z")
		test.Ok(t, err)
		test.Equal(t, timeFlag.Get(), want)
		test.Equal(t, timeFlag.Type(), "time")
		test.Equal(t, timeFlag.String(), "2024-07-17T07:38:05Z")
	})

	t.Run("time.Time invalid", func(t *testing.T) {
		var tyme time.Time
		timeFlag := flag.New(&tyme, "time", "t", time.Now(), "Set a time value")
		err := timeFlag.Set("not a time")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag time received invalid value "not a time" (expected time.Time), detail: parsing time "not a time" as "2006-01-02T15:04:05Z07:00": cannot parse "not a time" as "2006"`,
		)
	})

	t.Run("time.Duration valid", func(t *testing.T) {
		var duration time.Duration
		durationFlag := flag.New(
			&duration,
			"duration",
			"d",
			time.Duration(0),
			"Set a duration value",
		)
		err := durationFlag.Set("300ms")
		test.Ok(t, err)

		want, err := time.ParseDuration("300ms")
		test.Ok(t, err)
		test.Equal(t, durationFlag.Get(), want)
		test.Equal(t, durationFlag.Type(), "duration")
		test.Equal(t, durationFlag.String(), "300ms")
	})

	t.Run("time.Duration invalid", func(t *testing.T) {
		var duration time.Duration
		durationFlag := flag.New(
			&duration,
			"duration",
			"d",
			time.Duration(0),
			"Set a duration value",
		)
		err := durationFlag.Set("not a duration")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag duration received invalid value "not a duration" (expected time.Duration), detail: time: invalid duration "not a duration"`,
		)
	})

	t.Run("ip valid", func(t *testing.T) {
		var ip net.IP
		ipFlag := flag.New(&ip, "ip", "i", nil, "Set an IP address")
		err := ipFlag.Set("192.0.2.1")
		test.Ok(t, err)
		test.Diff(t, ipFlag.Get(), net.ParseIP("192.0.2.1"))
		test.Equal(t, ipFlag.Type(), "ip")
		test.Equal(t, ipFlag.String(), "192.0.2.1")
	})

	t.Run("ip invalid", func(t *testing.T) {
		var ip net.IP
		ipFlag := flag.New(&ip, "ip", "i", nil, "Set an IP address")
		err := ipFlag.Set("not an ip")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag ip received invalid value "not an ip" (expected net.IP), detail: invalid IP address`,
		)
	})
}
