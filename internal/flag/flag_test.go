package flag_test

import (
	"bytes"
	"errors"
	"net"
	"slices"
	"testing"
	"time"

	publicflag "go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/flag"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/parse"
	"go.followtheprocess.codes/test"
)

func TestFlaggableTypes(t *testing.T) {
	// We can't do table testing here because Flag[T] is a different type for each test
	// so we can't do a []Flag[T] which is needed to define the test cases
	// so strap in for a bunch of copy pasta
	t.Run("int valid", func(t *testing.T) {
		var i int

		intFlag, err := flag.New(&i, "int", 'i', "Set an int value", flag.Config[int]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intFlag.Type(), "int")
		test.Equal(t, intFlag.String(), "42")
		test.Equal(t, intFlag.Default(), "42")
	})

	t.Run("int invalid", func(t *testing.T) {
		var i int

		intFlag, err := flag.New(&i, "int", 'i', "Set an int value", flag.Config[int]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int8 valid", func(t *testing.T) {
		var i int8

		intFlag, err := flag.New(&i, "int", 'i', "Set an int8 value", flag.Config[int8]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int8(42))
		test.Equal(t, intFlag.Type(), "int8")
		test.Equal(t, intFlag.String(), "42")
		test.Equal(t, intFlag.Default(), "42")
	})

	t.Run("int8 invalid", func(t *testing.T) {
		var i int8

		intFlag, err := flag.New(&i, "int", 'i', "Set an int8 value", flag.Config[int8]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int16 valid", func(t *testing.T) {
		var i int16

		intFlag, err := flag.New(&i, "int", 'i', "Set an int16 value", flag.Config[int16]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int16(42))
		test.Equal(t, intFlag.Type(), "int16")
		test.Equal(t, intFlag.String(), "42")
		test.Equal(t, intFlag.Default(), "42")
	})

	t.Run("int16 invalid", func(t *testing.T) {
		var i int16

		intFlag, err := flag.New(&i, "int", 'i', "Set an int16 value", flag.Config[int16]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int32 valid", func(t *testing.T) {
		var i int32

		intFlag, err := flag.New(&i, "int", 'i', "Set an int32 value", flag.Config[int32]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int32(42))
		test.Equal(t, intFlag.Type(), "int32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int32 invalid", func(t *testing.T) {
		var i int32

		intFlag, err := flag.New(&i, "int", 'i', "Set an int32 value", flag.Config[int32]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int64 valid", func(t *testing.T) {
		var i int64

		intFlag, err := flag.New(&i, "int", 'i', "Set an int64 value", flag.Config[int64]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, int64(42))
		test.Equal(t, intFlag.Type(), "int64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("int64 invalid", func(t *testing.T) {
		var i int64

		intFlag, err := flag.New(&i, "int", 'i', "Set an int64 value", flag.Config[int64]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("count valid", func(t *testing.T) {
		var c publicflag.Count

		countFlag, err := flag.New(&c, "count", 'c', "Count something", flag.Config[publicflag.Count]{})
		test.Ok(t, err)

		err = countFlag.Set("1")
		test.Ok(t, err)
		test.Equal(t, c, publicflag.Count(1))
		test.Equal(t, countFlag.Type(), "count")
		test.Equal(t, countFlag.String(), "1")

		// Setting it again should increment to 2
		err = countFlag.Set("1")
		test.Ok(t, err)
		test.Equal(t, c, publicflag.Count(2))
		test.Equal(t, countFlag.Type(), "count")
		test.Equal(t, countFlag.String(), "2")

		// Should also be able to set an explicit number e.g. --verbosity=3
		// so should now be 5
		err = countFlag.Set("3")
		test.Ok(t, err)
		test.Equal(t, c, publicflag.Count(5))
		test.Equal(t, countFlag.Type(), "count")
		test.Equal(t, countFlag.String(), "5")
	})

	t.Run("count invalid", func(t *testing.T) {
		var c publicflag.Count

		countFlag, err := flag.New(&c, "count", 'c', "Count something", flag.Config[publicflag.Count]{})
		test.Ok(t, err)

		err = countFlag.Set("a word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint valid", func(t *testing.T) {
		var i uint

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint value", flag.Config[uint]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, 42)
		test.Equal(t, intFlag.Type(), "uint")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint invalid", func(t *testing.T) {
		var i uint

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint value", flag.Config[uint]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint8 valid", func(t *testing.T) {
		var i uint8

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint8 value", flag.Config[uint8]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint8(42))
		test.Equal(t, intFlag.Type(), "uint8")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint8 invalid", func(t *testing.T) {
		var i uint8

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint8 value", flag.Config[uint8]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint16 valid", func(t *testing.T) {
		var i uint16

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint16 value", flag.Config[uint16]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint16(42))
		test.Equal(t, intFlag.Type(), "uint16")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint16 invalid", func(t *testing.T) {
		var i uint16

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint16 value", flag.Config[uint16]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint32 valid", func(t *testing.T) {
		var i uint32

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint32 value", flag.Config[uint32]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint32(42))
		test.Equal(t, intFlag.Type(), "uint32")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint32 invalid", func(t *testing.T) {
		var i uint32

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint32 value", flag.Config[uint32]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint64 valid", func(t *testing.T) {
		var i uint64

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint64 value", flag.Config[uint64]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uint64(42))
		test.Equal(t, intFlag.Type(), "uint64")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uint64 invalid", func(t *testing.T) {
		var i uint64

		intFlag, err := flag.New(&i, "uint", 'i', "Set a uint64 value", flag.Config[uint64]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uintptr valid", func(t *testing.T) {
		var i uintptr

		intFlag, err := flag.New(&i, "uintptr", 'i', "Set a uintptr value", flag.Config[uintptr]{})
		test.Ok(t, err)

		err = intFlag.Set("42")
		test.Ok(t, err)
		test.Equal(t, i, uintptr(42))
		test.Equal(t, intFlag.Type(), "uintptr")
		test.Equal(t, intFlag.String(), "42")
	})

	t.Run("uintptr invalid", func(t *testing.T) {
		var i uintptr

		intFlag, err := flag.New(&i, "uintptr", 'i', "Set a uintptr value", flag.Config[uintptr]{})
		test.Ok(t, err)

		err = intFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float32 valid", func(t *testing.T) {
		var f float32

		floatFlag, err := flag.New(&f, "float", 'f', "Set a float32 value", flag.Config[float32]{})
		test.Ok(t, err)

		err = floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatFlag.Type(), "float32")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float32 invalid", func(t *testing.T) {
		var f float32

		floatFlag, err := flag.New(&f, "float", 'f', "Set a float32 value", flag.Config[float32]{})
		test.Ok(t, err)

		err = floatFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float64 valid", func(t *testing.T) {
		var f float64

		floatFlag, err := flag.New(&f, "float", 'f', "Set a float64 value", flag.Config[float64]{})
		test.Ok(t, err)

		err = floatFlag.Set("3.14159")
		test.Ok(t, err)
		test.Equal(t, f, 3.14159)
		test.Equal(t, floatFlag.Type(), "float64")
		test.Equal(t, floatFlag.String(), "3.14159")
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64

		floatFlag, err := flag.New(&f, "float", 'f', "Set a float64 value", flag.Config[float64]{})
		test.Ok(t, err)

		err = floatFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("bool valid", func(t *testing.T) {
		var b bool

		boolFlag, err := flag.New(&b, "bool", 'b', "Set a bool value", flag.Config[bool]{})
		test.Ok(t, err)

		err = boolFlag.Set(format.True)
		test.Ok(t, err)
		test.Equal(t, b, true)
		test.Equal(t, boolFlag.Type(), "bool")
		test.Equal(t, boolFlag.String(), format.True)
	})

	t.Run("bool invalid", func(t *testing.T) {
		var b bool

		boolFlag, err := flag.New(&b, "bool", 'b', "Set a bool value", flag.Config[bool]{})
		test.Ok(t, err)

		err = boolFlag.Set("word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	// No invalid case as all command line args are strings anyway so no real way of
	// getting an error here
	t.Run("string", func(t *testing.T) {
		var str string

		strFlag, err := flag.New(&str, "string", 's', "Set a string value", flag.Config[string]{})
		test.Ok(t, err)

		err = strFlag.Set("newvalue")
		test.Ok(t, err)
		test.Equal(t, str, "newvalue")
		test.Equal(t, strFlag.Type(), "string")
		test.Equal(t, strFlag.String(), "newvalue")
	})

	t.Run("byte slice valid", func(t *testing.T) {
		var byt []byte

		byteFlag, err := flag.New(&byt, "byte", 'b', "Set a byte slice value", flag.Config[[]byte]{})
		test.Ok(t, err)

		err = byteFlag.Set("5e")
		test.Ok(t, err)
		test.EqualFunc(t, byt, []byte("^"), bytes.Equal)
		test.Equal(t, byteFlag.Type(), "bytesHex")
		test.Equal(t, byteFlag.String(), "5e")
	})

	t.Run("byte slice invalid", func(t *testing.T) {
		var byt []byte

		byteFlag, err := flag.New(&byt, "byte", 'b', "Set a byte slice value", flag.Config[[]byte]{})
		test.Ok(t, err)

		err = byteFlag.Set("0xF")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("time.Time valid", func(t *testing.T) {
		var tyme time.Time

		timeFlag, err := flag.New(&tyme, "time", 't', "Set a time value", flag.Config[time.Time]{})
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

		timeFlag, err := flag.New(&tyme, "time", 't', "Set a time value", flag.Config[time.Time]{})
		test.Ok(t, err)

		err = timeFlag.Set("not a time")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("time.Duration valid", func(t *testing.T) {
		var duration time.Duration

		durationFlag, err := flag.New(
			&duration,
			"duration",
			'd',
			"Set a duration value",
			flag.Config[time.Duration]{},
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
			"Set a duration value",
			flag.Config[time.Duration]{},
		)
		test.Ok(t, err)

		err = durationFlag.Set("not a duration")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("ip valid", func(t *testing.T) {
		var ip net.IP

		ipFlag, err := flag.New(&ip, "ip", 'i', "Set an IP address", flag.Config[net.IP]{})
		test.Ok(t, err)

		err = ipFlag.Set("192.0.2.1")
		test.Ok(t, err)
		test.DiffBytes(t, ip, net.ParseIP("192.0.2.1"))
		test.Equal(t, ipFlag.Type(), "ip")
		test.Equal(t, ipFlag.String(), "192.0.2.1")
	})

	t.Run("ip invalid", func(t *testing.T) {
		var ip net.IP

		ipFlag, err := flag.New(&ip, "ip", 'i', "Set an IP address", flag.Config[net.IP]{})
		test.Ok(t, err)

		err = ipFlag.Set("not an ip")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int slice valid", func(t *testing.T) {
		var slice []int

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]int]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("int slice invalid", func(t *testing.T) {
		var slice []int

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]int]{})
		test.Ok(t, err)

		err = sliceFlag.Set("a word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int8 slice valid", func(t *testing.T) {
		var slice []int8

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]int8]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("int8 slice invalid", func(t *testing.T) {
		var slice []int8

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]int8]{})
		test.Ok(t, err)

		err = sliceFlag.Set("cheese")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int16 slice valid", func(t *testing.T) {
		var slice []int16

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]int16]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("int16 slice invalid", func(t *testing.T) {
		var slice []int16

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]int16]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int32 slice valid", func(t *testing.T) {
		var slice []int32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]int32]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("int32 slice invalid", func(t *testing.T) {
		var slice []int32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]int32]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("int64 slice valid", func(t *testing.T) {
		var slice []int64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]int64]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("int64 slice invalid", func(t *testing.T) {
		var slice []int64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]int64]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint slice valid", func(t *testing.T) {
		var slice []uint

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of uints", flag.Config[[]uint]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("uint slice invalid", func(t *testing.T) {
		var slice []uint

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of unsigned integers", flag.Config[[]uint]{})
		test.Ok(t, err)

		err = sliceFlag.Set("a word")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint16 slice valid", func(t *testing.T) {
		var slice []uint16

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]uint16]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("uint16 slice invalid", func(t *testing.T) {
		var slice []uint16

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]uint16]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint32 slice valid", func(t *testing.T) {
		var slice []uint32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]uint32]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("uint32 slice invalid", func(t *testing.T) {
		var slice []uint32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]uint32]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("uint64 slice valid", func(t *testing.T) {
		var slice []uint64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of ints", flag.Config[[]uint64]{})
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
		test.Equal(t, sliceFlag.String(), "[1, 2]")
	})

	t.Run("uint64 slice invalid", func(t *testing.T) {
		var slice []uint64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of integers", flag.Config[[]uint64]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float32 slice valid", func(t *testing.T) {
		var slice []float32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of floats", flag.Config[[]float32]{})
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
		test.Equal(t, sliceFlag.String(), "[3.14159, 2.7128]")
	})

	t.Run("float32 slice invalid", func(t *testing.T) {
		var slice []float32

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of floats", flag.Config[[]float32]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("float64 slice valid", func(t *testing.T) {
		var slice []float64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of floats", flag.Config[[]float64]{})
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
		test.Equal(t, sliceFlag.String(), "[3.14159, 2.7128]")
	})

	t.Run("float64 slice invalid", func(t *testing.T) {
		var slice []float64

		sliceFlag, err := flag.New(&slice, "slice", 's', "Slice of floats", flag.Config[[]float64]{})
		test.Ok(t, err)

		err = sliceFlag.Set("balls")
		test.Err(t, err)
		test.True(t, errors.Is(err, parse.Err))
	})

	t.Run("string slice valid", func(t *testing.T) {
		// Note: no invalid case for []string because *every* flag value is a string
		// it's impossible to make a bad one
		var slice []string

		sliceFlag, err := flag.New(&slice, "slice", 's', "Append to a slice of strings", flag.Config[[]string]{})
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
			short:    publicflag.NoShortHand,
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
			_, err := flag.New(new(string), tt.flagName, tt.short, "Test me", flag.Config[string]{})
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

		flag, err := flag.New(bang, "bang", 'b', "Nil go bang?", flag.Config[bool]{})
		test.Ok(t, err)

		test.Equal(t, flag.String(), "false")
		test.Equal(t, flag.Type(), "bool")
	})

	t.Run("composite literal", func(t *testing.T) {
		// Users doing naughty things, should still be nil safe
		flag := flag.Flag[bool]{}
		test.Equal(t, flag.String(), "<nil>")
		test.Equal(t, flag.Type(), "<nil>")

		err := flag.Set(format.True)
		test.Err(t, err)
		test.Equal(t, err.Error(), "cannot set value true, flag.value was nil")
	})
}

func BenchmarkFlagSet(b *testing.B) {
	var count int

	flag, err := flag.New(&count, "count", 'c', "Count things", flag.Config[int]{})
	test.Ok(b, err)

	for b.Loop() {
		err := flag.Set("42")
		if err != nil {
			b.Fatalf("flag.Set returned an error: %v", err)
		}
	}
}
