package format //nolint:testpackage // I need the base and bits values and don't want to export them.

import (
	"strconv"
	"testing"
	"testing/quick"

	"go.followtheprocess.codes/test"
)

func TestInt(t *testing.T) {
	//nolint:gocritic // It wants me to "unlambda" this but it's generic so I can't
	test := func(n int) string {
		return Int(n)
	}

	reference := func(n int) string {
		return strconv.FormatInt(int64(n), base10)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint(t *testing.T) {
	//nolint:gocritic // It wants me to "unlambda" this but it's generic so I can't
	test := func(n uint) string {
		return Uint(n)
	}

	reference := func(n uint) string {
		return strconv.FormatUint(uint64(n), base10)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat32(t *testing.T) {
	test := Float32

	reference := func(f float32) string {
		return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits32)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat64(t *testing.T) {
	test := Float64

	reference := func(f float64) string {
		return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits64)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestSlice(t *testing.T) {
	oneString := []string{"one"}
	twoStrings := []string{"one", "two"}
	strings := []string{"one", "two", "three"}
	ints := []int{1, 2, 3}
	floats := []float64{1.0, 2.0, 3.0}
	bools := []bool{true, true, false}

	test.Equal(t, Slice(oneString), `["one"]`)
	test.Equal(t, Slice(twoStrings), `["one", "two"]`)
	test.Equal(t, Slice(strings), `["one", "two", "three"]`, test.Context("strings"))
	test.Equal(t, Slice(ints), "[1, 2, 3]", test.Context("ints"))
	test.Equal(t, Slice(floats), "[1, 2, 3]", test.Context("floats"))
	test.Equal(t, Slice(bools), "[true, true, false]", test.Context("bools"))
}
