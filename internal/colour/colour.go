// Package colour implements basic text colouring for cli's limited needs.
//
// In particular, it's not expected to provide every ANSI code, just the ones we need. The codes have also been padded so that they are
// the same length, which means [text/tabwriter] will correctly calculate alignment as long as styles are not mixed within a table.
package colour

import "os"

// ANSI codes for coloured output, they are all the same length so as not to throw off
// alignment of [text/tabwriter].
const (
	CodeReset = "\x1b[000000m" // Reset all attributes
	CodeTitle = "\x1b[1;37;4m" // Bold, white & underlined
	CodeBold  = "\x1b[1;0037m" // Bold & white
)

// Title returns the given text in a title style, bold white and underlined.
//
// If $NO_COLOR is set, text will be returned unmodified.
func Title(text string) string {
	if noColour() {
		return text
	}
	return CodeTitle + text + CodeReset
}

// Bold returns the given text in bold white.
//
// If $NO_COLOR is set, text will be returned unmodified.
func Bold(text string) string {
	if noColour() {
		return text
	}
	return CodeBold + text + CodeReset
}

// noColour returns whether the $NO_COLOR env var was set.
func noColour() bool {
	return os.Getenv("NO_COLOR") != ""
}
