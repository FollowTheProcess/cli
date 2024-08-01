package flag_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/FollowTheProcess/test"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string                            // The name of the test case
		errMsg  string                            // If we did get an error, what should it say
		newSet  func(t *testing.T) *flag.Set      // Function to create a set for testing
		test    func(t *testing.T, set *flag.Set) // Function to invoke to test the set's functionality
		args    []string                          // The arguments provided to parse
		wantErr bool                              // Whether we want a parse error
	}{
		{
			name: "empty set no args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "empty set args no flags",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), []string{"some", "args", "here", "no", "flags"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "no", "flags"},
			wantErr: false,
		},
		{
			name: "empty set with flags",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("flag")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), []string{"some", "args", "here"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "--flag", "-s"},
			wantErr: true,
			errMsg:  "unrecognised flag: --flag",
		},
		{
			name: "undefined flag long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u`,
		},
		{
			name: "undefined flag long with value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown", "value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short with value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u", "value"},
			wantErr: true,
			errMsg:  "unrecognised shorthand flag: -u",
		},
		{
			name: "undefined flag long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown=value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u=value"},
			wantErr: true,
			errMsg:  "unrecognised shorthand flag: -u",
		},
		{
			name: "undefined flag shortvalue",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-uvalue"},
			wantErr: true,
			errMsg:  "unrecognised shorthand flag: -u",
		},
		{
			name: "bad syntax long empty name",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--"},
			wantErr: true,
			errMsg:  `invalid flag name "": must not be empty`,
		},
		{
			name: "bad syntax short empty name",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-"},
			wantErr: true,
			errMsg:  `invalid flag name "": must not be empty`,
		},
		{
			name: "bad syntax long extra hyphen",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"---"},
			wantErr: true,
			errMsg:  `invalid flag name "-": trailing hyphen`,
		},
		{
			name: "bad syntax long leading whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-- delete"},
			wantErr: true,
			errMsg:  `invalid flag name " delete": cannot contain whitespace`,
		},
		{
			name: "bad syntax short leading whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"- d"},
			wantErr: true,
			errMsg:  `invalid flag name " d": cannot contain whitespace`,
		},
		{
			name: "bad syntax long trailing whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--delete "},
			wantErr: true,
			errMsg:  `invalid flag name "delete ": cannot contain whitespace`,
		},
		{
			name: "bad syntax short trailing whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-d "},
			wantErr: true,
			errMsg:  `invalid flag name "d ": cannot contain whitespace`,
		},
		{
			name: "bad syntax long internal whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--de lete"},
			wantErr: true,
			errMsg:  `invalid flag name "de lete": cannot contain whitespace`,
		},
		{
			name: "valid long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, value.Type(), "bool")
				test.Equal(t, value.String(), "true")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--delete"},
			wantErr: false,
		},
		{
			name: "valid long with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, value.Type(), "bool")
				test.Equal(t, value.String(), "true")

				test.EqualFunc(t, set.Args(), []string{"some", "subcommand"}, slices.Equal)
			},
			args:    []string{"some", "subcommand", "--delete"},
			wantErr: false,
		},
		{
			name: "valid short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, value.Type(), "bool")
				test.Equal(t, value.String(), "true")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-d"},
			wantErr: false,
		},
		{
			name: "valid short with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, value.Type(), "bool")
				test.Equal(t, value.String(), "true")

				test.EqualFunc(t, set.Args(), []string{"some", "arg"}, slices.Equal)
			},
			args:    []string{"some", "arg", "-d"},
			wantErr: false,
		},
		{
			name: "valid long value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--count", "1"},
			wantErr: false,
		},
		{
			name: "valid long value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"some", "arg", "more", "args"}, slices.Equal)
			},
			args:    []string{"some", "arg", "--count", "1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid long value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, value.Type(), "uint")
				test.Equal(t, value.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-c", "1"},
			wantErr: false,
		},
		{
			name: "valid short value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "-c", "1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid short value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, value.Type(), "uint")
				test.Equal(t, value.String(), "0") // Shouldn't have been set
			},
			args:    []string{"-n", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")
			},
			args:    []string{"--count=1"},
			wantErr: false,
		},
		{
			name: "valid long equals value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "--count=1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, value.Type(), "uint")
				test.Equal(t, value.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number=-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")
			},
			args:    []string{"-c=1"},
			wantErr: false,
		},
		{
			name: "valid short equals value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, value.Type(), "int")
				test.Equal(t, value.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "-c=1", "more", "args"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.newSet(t)
			err := set.Parse(tt.args)
			test.WantErr(t, err, tt.wantErr)

			if err != nil {
				test.Equal(t, err.Error(), tt.errMsg)
			}

			if tt.test != nil {
				tt.test(t, set)
			}
		})
	}
}

func TestFlagSet(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set      // Function to build the flag set under test
		test   func(t *testing.T, set *flag.Set) // Function to test the set
		name   string                            // Name of the test case
	}{
		{
			name: "empty",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe get",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe add",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(bool), "force", 'f', false, "Force something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), "cannot add flag to a nil set")
				}
			},
		},
		{
			name: "nil safe parse",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				err := set.Parse([]string{"args", "here", "doesn't", "matter"})
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), "Parse called on a nil set")
				}
			},
		},
		{
			name: "duplicate flag added",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				// Add the same flag again
				err = flag.AddToSet(set, f)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), `flag "count" already defined`)
				}
			},
		},
		{
			name: "different flag same name added",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				f2, err := flag.New(new(uint), "count", 'C', 0, "Count something 2")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				// Add the different flag with the same na,e
				err = flag.AddToSet(set, f2)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), `flag "count" already defined`)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.newSet(t)
			tt.test(t, set)
		})
	}
}
