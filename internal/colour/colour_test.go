package colour_test

import (
	"testing"

	"github.com/FollowTheProcess/cli/internal/colour"
	"github.com/FollowTheProcess/test"
)

func TestColour(t *testing.T) {
	tests := []struct {
		name    string                   // Name of the test case
		text    string                   // Text to colour
		fn      func(text string) string // Printer function to return the coloured version of text
		want    string                   // Expected result containing ANSI escape codes
		noColor bool                     // Whether to set the $NO_COLOR env var
	}{
		{
			name: "bold",
			text: "hello bold",
			fn:   colour.Bold,
			want: colour.CodeBold + "hello bold" + colour.CodeReset,
		},
		{
			name:    "bold no color",
			text:    "hello bold",
			fn:      colour.Bold,
			noColor: true,
			want:    "hello bold",
		},
		{
			name: "title",
			text: "Section",
			fn:   colour.Title,
			want: colour.CodeTitle + "Section" + colour.CodeReset,
		},
		{
			name:    "title no color",
			text:    "Section",
			fn:      colour.Title,
			noColor: true,
			want:    "Section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.noColor {
				t.Setenv("NO_COLOR", "true")
			}
			got := tt.fn(tt.text)
			test.Equal(t, got, tt.want)
		})
	}
}

func TestCodesAllSameLength(t *testing.T) {
	test.True(t, len(colour.CodeBold) == len(colour.CodeReset))
	test.True(t, len(colour.CodeBold) == len(colour.CodeTitle))
	test.True(t, len(colour.CodeReset) == len(colour.CodeTitle))
}

func BenchmarkBold(b *testing.B) {
	for range b.N {
		colour.Bold("Some bold text here")
	}
}

func BenchmarkTitle(b *testing.B) {
	for range b.N {
		colour.Title("Some title here")
	}
}
