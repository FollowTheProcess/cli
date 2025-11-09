package parse //nolint:testpackage // I need the base and bits values and don't want to export them.

import (
	"errors"
	"strconv"
	"testing"
	"testing/quick"

	"go.followtheprocess.codes/test"
)

// These are basically all just testing that I haven't broken anything
// by wrapping strconv and saves me having to write lots of test cases
// by hand.

func TestInt(t *testing.T) {
	// Check a basic happy path then let quick handle the rest
	got, err := Int("5")
	test.Ok(t, err)
	test.Equal(t, got, 5)

	test := Int

	reference := func(str string) (int, error) {
		val, err := strconv.ParseInt(str, base10, 0)
		return int(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt8(t *testing.T) {
	got, err := Int8("4")
	test.Ok(t, err)
	test.Equal(t, got, 4)

	test := Int8

	reference := func(str string) (int8, error) {
		val, err := strconv.ParseInt(str, base10, bits8)
		return int8(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt16(t *testing.T) {
	got, err := Int16("32")
	test.Ok(t, err)
	test.Equal(t, got, 32)

	test := Int16

	reference := func(str string) (int16, error) {
		val, err := strconv.ParseInt(str, base10, bits16)
		return int16(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt32(t *testing.T) {
	got, err := Int32("12")
	test.Ok(t, err)
	test.Equal(t, got, 12)

	test := Int32

	reference := func(str string) (int32, error) {
		val, err := strconv.ParseInt(str, base10, bits32)
		return int32(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt64(t *testing.T) {
	got, err := Int64("27")
	test.Ok(t, err)
	test.Equal(t, got, 27)

	test := Int64

	reference := func(str string) (int64, error) {
		val, err := strconv.ParseInt(str, base10, bits64)
		return val, err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint(t *testing.T) {
	got, err := Uint("2")
	test.Ok(t, err)
	test.Equal(t, got, 2)

	test := Uint

	reference := func(str string) (uint, error) {
		val, err := strconv.ParseUint(str, base10, 0)
		return uint(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint8(t *testing.T) {
	got, err := Uint8("8")
	test.Ok(t, err)
	test.Equal(t, got, 8)

	test := Uint8

	reference := func(str string) (uint8, error) {
		val, err := strconv.ParseUint(str, base10, bits8)
		return uint8(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint16(t *testing.T) {
	got, err := Uint16("7")
	test.Ok(t, err)
	test.Equal(t, got, 7)

	test := Uint16

	reference := func(str string) (uint16, error) {
		val, err := strconv.ParseUint(str, base10, bits16)
		return uint16(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint32(t *testing.T) {
	got, err := Uint32("47")
	test.Ok(t, err)
	test.Equal(t, got, 47)

	test := Uint32

	reference := func(str string) (uint32, error) {
		val, err := strconv.ParseUint(str, base10, bits32)
		return uint32(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint64(t *testing.T) {
	got, err := Uint64("102")
	test.Ok(t, err)
	test.Equal(t, got, 102)

	test := Uint64

	reference := func(str string) (uint64, error) {
		val, err := strconv.ParseUint(str, base10, bits64)
		return val, err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat32(t *testing.T) {
	got, err := Float32("102.8975")
	test.Ok(t, err)
	test.NearlyEqual(t, got, 102.8975)

	test := Float32

	reference := func(str string) (float32, error) {
		val, err := strconv.ParseFloat(str, bits32)
		return float32(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat64(t *testing.T) {
	got, err := Float64("916156.123")
	test.Ok(t, err)
	test.NearlyEqual(t, got, 916156.123)

	test := Float64

	reference := func(str string) (float64, error) {
		val, err := strconv.ParseFloat(str, bits64)
		return float64(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	got := Error(KindArgument, "test", "blah", "string", errors.New("underlying"))
	want := `parse error: argument "test" received invalid value "blah" (expected string): underlying`

	test.Equal(t, got.Error(), want)
}

func TestErrorSlice(t *testing.T) {
	got := ErrorSlice(KindArgument, "test", "blah", "string", errors.New("underlying"))
	want := `parse error: argument "test" (type string) cannot append element "blah": underlying`

	test.Equal(t, got.Error(), want)
}
