package flag_test

import (
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
			args:    []string{},
			wantErr: false,
		},
		{
			name: "empty set args no flags",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"some", "args", "here", "no", "flags"},
			wantErr: false,
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
			name: "bad syntax empty name",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--"},
			wantErr: true,
			errMsg:  `invalid flag name "": must not be empty`,
		},
		{
			name: "bad syntax extra hyphen",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"---"},
			wantErr: true,
			errMsg:  `invalid flag name "-": trailing hyphen`,
		},
		{
			name: "bad syntax leading whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-- delete"},
			wantErr: true,
			errMsg:  `invalid flag name " delete": cannot contain whitespace`,
		},
		{
			name: "bad syntax trailing whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--delete "},
			wantErr: true,
			errMsg:  `invalid flag name "delete ": cannot contain whitespace`,
		},
		{
			name: "bad syntax internal whitespace",
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

				err = set.Add("delete", 'd', "Delete something", f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				value, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, value.Type(), "bool")
				test.Equal(t, value.String(), "true")
			},
			args:    []string{"--delete"},
			wantErr: false,
		},
		{
			name: "valid long value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = set.Add("count", 'c', "Count something", f)
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
			args:    []string{"--count", "1"},
			wantErr: false,
		},
		{
			name: "valid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = set.Add("count", 'c', "Count something", f)
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
