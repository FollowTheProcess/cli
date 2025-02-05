package flag_test

import (
	"bytes"
	"net"
	"slices"
	"testing"
	"time"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/FollowTheProcess/test"
)

func TestFlaggableTypes(t *testing.T) {
	// We can't do table testing here because Flag[T] is a different type for each test
	// so we can't do a []Flag[T] which is needed to define the test cases
	// so strap in for a bunch of copy pasta
	t.Run("int valid", func(t *testing.T) {
		var i int
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intFlag.Type(), "int")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int invalid", func(t *testing.T) {
		var i int
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "int" received invalid value "word" (expected int), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int8 valid", func(t *testing.T) {
		var i int8
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int8 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int8(42))
		test.Equal(t, intFlag.Type(), "int8")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int8 invalid", func(t *testing.T) {
		var i int8
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int8 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "int" received invalid value "word" (expected int8), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int16 valid", func(t *testing.T) {
		var i int16
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int16 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int16(42))
		test.Equal(t, intFlag.Type(), "int16")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int16 invalid", func(t *testing.T) {
		var i int16
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int16 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "int" received invalid value "word" (expected int16), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int32 valid", func(t *testing.T) {
		var i int32
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int32 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int32(42))
		test.Equal(t, intFlag.Type(), "int32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int32 invalid", func(t *testing.T) {
		var i int32
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int32 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "int" received invalid value "word" (expected int32), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("int64 valid", func(t *testing.T) {
		var i int64
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int64 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int64(42))
		test.Equal(t, intFlag.Type(), "int64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int64 invalid", func(t *testing.T) {
		var i int64
		intFlag, err := flag.New(&i, "int", 'i', 0, "Set an int64 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "int" received invalid value "word" (expected int64), detail: strconv.ParseInt: parsing "word": invalid syntax`,
		)
	})

	t.Run("count valid", func(t *testing.T) {
		var c flag.Count
		countFlag, err := flag.New(&c, "count", 'c', 0, "Count something")
		test.Ok(t, err)

		err = countFlag.Set("1")
		test.Ok(t, err)
		test.Equal(t, c, flag.Count(1))
		test.Equal(t, countFlag.Type(), "count")
		test.Equal(t, countFlag.String(), "1")

		// Setting it again should increment to 2
		err = countFlag.Set("1")
		test.Ok(t, err)
		test.Equal(t, c, flag.Count(2))
		test.Equal(t, countFlag.Type(), "count")
		test.Equal(t, countFlag.String(), "2")
	})

	t.Run("uint valid", func(t *testing.T) {
		var i uint
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intFlag.Type(), "uint")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint invalid", func(t *testing.T) {
		var i uint
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uint" received invalid value "word" (expected uint), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint8 valid", func(t *testing.T) {
		var i uint8
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint8 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint8(42))
		test.Equal(t, intFlag.Type(), "uint8")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint8 invalid", func(t *testing.T) {
		var i uint8
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint8 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uint" received invalid value "word" (expected uint8), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint16 valid", func(t *testing.T) {
		var i uint16
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint16 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint16(42))
		test.Equal(t, intFlag.Type(), "uint16")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint16 invalid", func(t *testing.T) {
		var i uint16
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint16 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uint" received invalid value "word" (expected uint16), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint32 valid", func(t *testing.T) {
		var i uint32
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint32 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint32(42))
		test.Equal(t, intFlag.Type(), "uint32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint32 invalid", func(t *testing.T) {
		var i uint32
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint32 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uint" received invalid value "word" (expected uint32), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uint64 valid", func(t *testing.T) {
		var i uint64
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint64 value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint64(42))
		test.Equal(t, intFlag.Type(), "uint64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint64 invalid", func(t *testing.T) {
		var i uint64
		intFlag, err := flag.New(&i, "uint", 'i', 0, "Set a uint64 value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uint" received invalid value "word" (expected uint64), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("uintptr valid", func(t *testing.T) {
		var i uintptr
		intFlag, err := flag.New(&i, "uintptr", 'i', 0, "Set a uintptr value")
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uintptr(42))
		test.Equal(t, intFlag.Type(), "uintptr")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uintptr invalid", func(t *testing.T) {
		var i uintptr
		intFlag, err := flag.New(&i, "uintptr", 'i', 0, "Set a uintptr value")
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "uintptr" received invalid value "word" (expected uintptr), detail: strconv.ParseUint: parsing "word": invalid syntax`,
		)
	})

	t.Run("float32 valid", func(t *testing.T) {
		var f float32
		floatFlag, err := flag.New(&f, "float", 'f', 0, "Set a float32 value")
		test.Ok(t, err)

		err = floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatFlag.Type(), "float32")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float32 invalid", func(t *testing.T) {
		var f float32
		floatFlag, err := flag.New(&f, "float", 'f', 0, "Set a float32 value")
		test.Ok(t, err)

		err = floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "float" received invalid value "word" (expected float32), detail: strconv.ParseFloat: parsing "word": invalid syntax`,
		)
	})

	t.Run("float64 valid", func(t *testing.T) {
		var f float64
		floatFlag, err := flag.New(&f, "float", 'f', 0, "Set a float64 value")
		test.Ok(t, err)

		err = floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatFlag.Type(), "float64")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64
		floatFlag, err := flag.New(&f, "float", 'f', 0, "Set a float64 value")
		test.Ok(t, err)

		err = floatFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "float" received invalid value "word" (expected float64), detail: strconv.ParseFloat: parsing "word": invalid syntax`,
		)
	})

	t.Run("bool valid", func(t *testing.T) {
		var b bool
		boolFlag, err := flag.New(&b, "bool", 'b', false, "Set a bool value")
		test.Ok(t, err)

		err = boolFlag.Set("true")
		test.Ok(t, err)
		test.Equal(t, b, true)
		test.Equal(t, boolFlag.Type(), "bool")
		test.Equal(t, boolFlag.String(), "true")
	})

	t.Run("bool invalid", func(t *testing.T) {
		var b bool
		boolFlag, err := flag.New(&b, "bool", 'b', false, "Set a bool value")
		test.Ok(t, err)

		err = boolFlag.Set("word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "bool" received invalid value "word" (expected bool), detail: strconv.ParseBool: parsing "word": invalid syntax`,
		)
	})

	// No invalid case as all command line args are strings anyway so no real way of
	// getting an error here
	t.Run("string", func(t *testing.T) {
		var str string
		strFlag, err := flag.New(&str, "string", 's', "", "Set a string value")
		test.Ok(t, err)

		err = strFlag.Set("newvalue")
		test.Ok(t, err)
		test.Equal(t, str, "newvalue")
		test.Equal(t, strFlag.Type(), "string")
		test.Equal(t, strFlag.String(), "newvalue")
	})

	t.Run("byte slice valid", func(t *testing.T) {
		var byt []byte
		byteFlag, err := flag.New(&byt, "byte", 'b', []byte(""), "Set a byte slice value")
		test.Ok(t, err)

		err = byteFlag.Set("5e")
		test.Ok(t, err)
		test.EqualFunc(t, byt, []byte("^"), bytes.Equal)
		test.Equal(t, byteFlag.Type(), "bytesHex")
		test.Equal(t, byteFlag.String(), "5e")
	})

	t.Run("byte slice invalid", func(t *testing.T) {
		var byt []byte
		byteFlag, err := flag.New(&byt, "byte", 'b', []byte(""), "Set a byte slice value")
		test.Ok(t, err)

		err = byteFlag.Set("0xF")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "byte" received invalid value "0xF" (expected []uint8), detail: encoding/hex: invalid byte: U+0078 'x'`,
		)
	})

	t.Run("time.Time valid", func(t *testing.T) {
		var tyme time.Time
		timeFlag, err := flag.New(&tyme, "time", 't', time.Now(), "Set a time value")
		test.Ok(t, err)

		err = timeFlag.Set("2024-07-17T07:38:05Z")
		test.Ok(t, err)

		want, err := time.Parse(time.RFC3339, "2024-07-17T07:38:05Z")
		test.Ok(t, err)
		test.Equal(t, tyme, want)
		test.Equal(t, timeFlag.Type(), "time")
		test.Equal(t, timeFlag.String(), "2024-07-17T07:38:05Z")
	})

	t.Run("time.Time invalid", func(t *testing.T) {
		var tyme time.Time
		timeFlag, err := flag.New(&tyme, "time", 't', time.Now(), "Set a time value")
		test.Ok(t, err)

		err = timeFlag.Set("not a time")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "time" received invalid value "not a time" (expected time.Time), detail: parsing time "not a time" as "2006-01-02T15:04:05Z07:00": cannot parse "not a time" as "2006"`,
		)
	})

	t.Run("time.Duration valid", func(t *testing.T) {
		var duration time.Duration
		durationFlag, err := flag.New(
			&duration,
			"duration",
			'd',
			time.Duration(0),
			"Set a duration value",
		)
		test.Ok(t, err)

		err = durationFlag.Set("300ms")
		test.Ok(t, err)

		want, err := time.ParseDuration("300ms")
		test.Ok(t, err)
		test.Equal(t, duration, want)
		test.Equal(t, durationFlag.Type(), "duration")
		test.Equal(t, durationFlag.String(), "300ms")
	})

	t.Run("time.Duration invalid", func(t *testing.T) {
		var duration time.Duration
		durationFlag, err := flag.New(
			&duration,
			"duration",
			'd',
			time.Duration(0),
			"Set a duration value",
		)
		test.Ok(t, err)

		err = durationFlag.Set("not a duration")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "duration" received invalid value "not a duration" (expected time.Duration), detail: time: invalid duration "not a duration"`,
		)
	})

	t.Run("ip valid", func(t *testing.T) {
		var ip net.IP
		ipFlag, err := flag.New(&ip, "ip", 'i', nil, "Set an IP address")
		test.Ok(t, err)

		err = ipFlag.Set("192.0.2.1")
		test.Ok(t, err)
		test.DiffBytes(t, ip, net.ParseIP("192.0.2.1"))
		test.Equal(t, ipFlag.Type(), "ip")
		test.Equal(t, ipFlag.String(), "192.0.2.1")
	})

	t.Run("ip invalid", func(t *testing.T) {
		var ip net.IP
		ipFlag, err := flag.New(&ip, "ip", 'i', nil, "Set an IP address")
		test.Ok(t, err)

		err = ipFlag.Set("not an ip")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "ip" received invalid value "not an ip" (expected net.IP), detail: invalid IP address`,
		)
	})

	t.Run("int slice valid", func(t *testing.T) {
		var slice []int
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("int slice invalid", func(t *testing.T) {
		var slice []int
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("a word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []int) cannot append element "a word": strconv.ParseInt: parsing "a word": invalid syntax`,
		)
	})

	t.Run("int8 slice valid", func(t *testing.T) {
		var slice []int8
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int8{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int8")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int8{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int8")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("int8 slice invalid", func(t *testing.T) {
		var slice []int8
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("cheese")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []int8) cannot append element "cheese": strconv.ParseInt: parsing "cheese": invalid syntax`,
		)
	})

	t.Run("int16 slice valid", func(t *testing.T) {
		var slice []int16
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int16{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int16")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int16{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int16")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("int16 slice invalid", func(t *testing.T) {
		var slice []int16
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []int16) cannot append element "balls": strconv.ParseInt: parsing "balls": invalid syntax`,
		)
	})

	t.Run("int32 slice valid", func(t *testing.T) {
		var slice []int32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int32{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int32")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int32{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int32")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("int32 slice invalid", func(t *testing.T) {
		var slice []int32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []int32) cannot append element "balls": strconv.ParseInt: parsing "balls": invalid syntax`,
		)
	})

	t.Run("int64 slice valid", func(t *testing.T) {
		var slice []int64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int64{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int64")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []int64{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]int64")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("int64 slice invalid", func(t *testing.T) {
		var slice []int64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []int64) cannot append element "balls": strconv.ParseInt: parsing "balls": invalid syntax`,
		)
	})

	t.Run("uint slice valid", func(t *testing.T) {
		var slice []uint
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of uints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("uint slice invalid", func(t *testing.T) {
		var slice []uint
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of unsigned integers")
		test.Ok(t, err)

		err = sliceFlag.Set("a word")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []uint) cannot append element "a word": strconv.ParseUint: parsing "a word": invalid syntax`,
		)
	})

	t.Run("uint16 slice valid", func(t *testing.T) {
		var slice []uint16
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint16{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint16")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint16{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint16")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("uint16 slice invalid", func(t *testing.T) {
		var slice []uint16
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []uint16) cannot append element "balls": strconv.ParseUint: parsing "balls": invalid syntax`,
		)
	})

	t.Run("uint32 slice valid", func(t *testing.T) {
		var slice []uint32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint32{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint32")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint32{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint32")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("uint32 slice invalid", func(t *testing.T) {
		var slice []uint32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []uint32) cannot append element "balls": strconv.ParseUint: parsing "balls": invalid syntax`,
		)
	})

	t.Run("uint64 slice valid", func(t *testing.T) {
		var slice []uint64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of ints")
		test.Ok(t, err)

		err = sliceFlag.Set("1") // Append 1 to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint64{1}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint64")
		test.Equal(t, sliceFlag.String(), "[1]")

		err = sliceFlag.Set("2") // Now 2
		test.Ok(t, err)

		test.EqualFunc(t, slice, []uint64{1, 2}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]uint64")
		test.Equal(t, sliceFlag.String(), "[1 2]")
	})

	t.Run("uint64 slice invalid", func(t *testing.T) {
		var slice []uint64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of integers")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []uint64) cannot append element "balls": strconv.ParseUint: parsing "balls": invalid syntax`,
		)
	})

	t.Run("float32 slice valid", func(t *testing.T) {
		var slice []float32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of floats")
		test.Ok(t, err)

		err = sliceFlag.Set("3.14159") // Append pi to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []float32{3.14159}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]float32")
		test.Equal(t, sliceFlag.String(), "[3.14159]")

		err = sliceFlag.Set("2.7128") // Now e
		test.Ok(t, err)

		test.EqualFunc(t, slice, []float32{3.14159, 2.7128}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]float32")
		test.Equal(t, sliceFlag.String(), "[3.14159 2.7128]")
	})

	t.Run("float32 slice invalid", func(t *testing.T) {
		var slice []float32
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of floats")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []float32) cannot append element "balls": strconv.ParseFloat: parsing "balls": invalid syntax`,
		)
	})

	t.Run("float64 slice valid", func(t *testing.T) {
		var slice []float64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of floats")
		test.Ok(t, err)

		err = sliceFlag.Set("3.14159") // Append pi to the slice
		test.Ok(t, err)

		test.EqualFunc(t, slice, []float64{3.14159}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]float64")
		test.Equal(t, sliceFlag.String(), "[3.14159]")

		err = sliceFlag.Set("2.7128") // Now e
		test.Ok(t, err)

		test.EqualFunc(t, slice, []float64{3.14159, 2.7128}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]float64")
		test.Equal(t, sliceFlag.String(), "[3.14159 2.7128]")
	})

	t.Run("float64 slice invalid", func(t *testing.T) {
		var slice []float64
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Slice of floats")
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.Equal(
			t,
			err.Error(),
			`flag "slice" (type []float64) cannot append element "balls": strconv.ParseFloat: parsing "balls": invalid syntax`,
		)
	})

	t.Run("string slice valid", func(t *testing.T) {
		// Note: no invalid case for []string because *every* flag value is a string
		// it's impossible to make a bad one
		var slice []string
		sliceFlag, err := flag.New(&slice, "slice", 's', nil, "Append to a slice of strings")
		test.Ok(t, err)

		err = sliceFlag.Set("a string")
		test.Ok(t, err)

		test.EqualFunc(t, slice, []string{"a string"}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]string")
		test.Equal(t, sliceFlag.String(), `["a string"]`)

		err = sliceFlag.Set("another string")
		test.Ok(t, err)

		test.EqualFunc(t, slice, []string{"a string", "another string"}, slices.Equal)
		test.Equal(t, sliceFlag.Type(), "[]string")
		test.Equal(t, sliceFlag.String(), `["a string", "another string"]`)
	})
}

func TestFlagValidation(t *testing.T) {
	tests := []struct {
		name     string // Name of the test case
		flagName string // Input flag name
		errMsg   string // If we wanted an error, what should it say
		short    rune   // Flag shorthand
		wantErr  bool   // Whether we want an error
	}{
		{
			name:     "short is uppercase",
			flagName: "delete",
			short:    'D',
			wantErr:  false,
		},
		{
			name:     "valid short",
			flagName: "delete",
			short:    'd',
			wantErr:  false,
		},
		{
			name:     "no shorthand",
			flagName: "delete",
			short:    flag.NoShortHand,
			wantErr:  false,
		},
		{
			name:     "hyphen separated",
			flagName: "dry-run",
			short:    'd',
			wantErr:  false,
		},
		{
			name:     "empty name",
			flagName: "",
			wantErr:  true,
			errMsg:   `invalid flag name "": must not be empty`,
		},
		{
			name:     "whitespace",
			flagName: " ",
			wantErr:  true,
			errMsg:   `invalid flag name " ": cannot contain whitespace`,
		},
		{
			name:     "mixed case",
			flagName: "heLlO",
			wantErr:  true,
			errMsg:   `invalid flag name "heLlO": contains upper case character "L"`,
		},
		{
			name:     "underscore",
			flagName: "set_default",
			wantErr:  true,
			errMsg:   `invalid flag name "set_default": contains non ascii letter: "_"`,
		},
		{
			name:     "digits",
			flagName: "some-06digit",
			wantErr:  true,
			errMsg:   `invalid flag name "some-06digit": contains non ascii letter: "0"`,
		},
		{
			name:     "just hyphen",
			flagName: "-",
			wantErr:  true,
			errMsg:   `invalid flag name "-": trailing hyphen`,
		},
		{
			name:     "leading hyphen",
			flagName: "-something",
			wantErr:  true,
			errMsg:   `invalid flag name "-something": leading hyphen`,
		},
		{
			name:     "trailing hyphen",
			flagName: "something-",
			wantErr:  true,
			errMsg:   `invalid flag name "something-": trailing hyphen`,
		},
		{
			name:     "non ascii",
			flagName: "語ç日ð本",
			wantErr:  true,
			errMsg:   `invalid flag name "語ç日ð本": contains non ascii character: "語"`,
		},
		{
			name:     "short is digit",
			flagName: "delete",
			short:    '7',
			wantErr:  true,
			errMsg:   `invalid shorthand for flag "delete": invalid character, must be a single ASCII letter, got "7"`,
		},
		{
			name:     "short is non ascii",
			flagName: "delete",
			short:    '本',
			wantErr:  true,
			errMsg:   `invalid shorthand for flag "delete": invalid character, must be a single ASCII letter, got "本"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := flag.New(new(string), tt.flagName, tt.short, "", "Test me")
			test.WantErr(t, err, tt.wantErr)

			if err != nil {
				test.Equal(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestFlagNilSafety(t *testing.T) {
	t.Run("with new", func(t *testing.T) {
		// Should be impossible to make a nil pointer dereference when using .New
		var bang *bool

		flag, err := flag.New(bang, "bang", 'b', false, "Nil go bang?")
		test.Ok(t, err)

		test.Equal(t, flag.String(), "false")
		test.Equal(t, flag.Type(), "bool")
	})

	t.Run("composite literal", func(t *testing.T) {
		// Users doing naughty things, should still be nil safe
		flag := flag.Flag[bool]{}
		test.Equal(t, flag.String(), "<nil>")
		test.Equal(t, flag.Type(), "<nil>")

		err := flag.Set("true")
		test.Err(t, err)
		test.Equal(t, err.Error(), "cannot set value true, flag.value was nil")
	})
}

func BenchmarkFlagSet(b *testing.B) {
	var count int
	flag, err := flag.New(&count, "count", 'c', 0, "Count things")
	test.Ok(b, err)

	b.ResetTimer()

	for range b.N {
		err := flag.Set("42")
		if err != nil {
			b.Fatalf("flag.Set returned an error: %v", err)
		}
	}
}
