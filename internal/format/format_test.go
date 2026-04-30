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
	tests := []struct {
		got  func() string
		name string
		want string
	}{
		{
			name: "nil string slice",
			got:  func() string { return Slice([]string(nil)) },
			want: "[]",
		},
		{
			name: "empty int slice",
			got:  func() string { return Slice([]int{}) },
			want: "[]",
		},
		{
			name: "one string",
			got:  func() string { return Slice([]string{"one"}) },
			want: `["one"]`,
		},
		{
			name: "two strings",
			got:  func() string { return Slice([]string{"one", "two"}) },
			want: `["one", "two"]`,
		},
		{
			name: "three strings",
			got:  func() string { return Slice([]string{"one", "two", "three"}) },
			want: `["one", "two", "three"]`,
		},
		{
			name: "strings with escapes",
			got:  func() string { return Slice([]string{"hi\nthere", "tab\there", `quote"here`}) },
			want: `["hi\nthere", "tab\there", "quote\"here"]`,
		},
		{
			name: "empty string element",
			got:  func() string { return Slice([]string{""}) },
			want: `[""]`,
		},
		{
			name: "ints",
			got:  func() string { return Slice([]int{1, 2, 3}) },
			want: "[1, 2, 3]",
		},
		{
			name: "negative ints",
			got:  func() string { return Slice([]int{-1, 0, 1}) },
			want: "[-1, 0, 1]",
		},
		{
			name: "int8s",
			got:  func() string { return Slice([]int8{-128, 0, 127}) },
			want: "[-128, 0, 127]",
		},
		{
			name: "int64s",
			got:  func() string { return Slice([]int64{-1 << 62, 0, 1 << 62}) },
			want: "[-4611686018427387904, 0, 4611686018427387904]",
		},
		{
			name: "uints",
			got:  func() string { return Slice([]uint{1, 2, 3}) },
			want: "[1, 2, 3]",
		},
		{
			name: "uint64s",
			got:  func() string { return Slice([]uint64{0, 1, 1 << 63}) },
			want: "[0, 1, 9223372036854775808]",
		},
		{
			name: "floats",
			got:  func() string { return Slice([]float64{1.0, 2.0, 3.0}) },
			want: "[1, 2, 3]",
		},
		{
			name: "floats with decimals",
			got:  func() string { return Slice([]float64{1.5, -2.25, 3.125}) },
			want: "[1.5, -2.25, 3.125]",
		},
		{
			name: "float32s",
			got:  func() string { return Slice([]float32{1.5, -2.25, 3.125}) },
			want: "[1.5, -2.25, 3.125]",
		},
		{
			name: "bools",
			got:  func() string { return Slice([]bool{true, true, false}) },
			want: "[true, true, false]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.got(), tt.want)
		})
	}
}

func BenchmarkSlice(b *testing.B) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8}
	int64s := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	uints := []uint{1, 2, 3, 4, 5, 6, 7, 8}
	floats := []float64{1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5}
	strs := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta"}
	bools := []bool{true, false, true, false, true, false, true, false}

	b.Run("ints", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(ints)
		}
	})

	b.Run("int64s", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(int64s)
		}
	})

	b.Run("uints", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(uints)
		}
	})

	b.Run("floats", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(floats)
		}
	})

	b.Run("strings", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(strs)
		}
	})

	b.Run("bools", func(b *testing.B) {
		for b.Loop() {
			_ = Slice(bools)
		}
	})

	b.Run("empty", func(b *testing.B) {
		var s []int

		for b.Loop() {
			_ = Slice(s)
		}
	})
}
