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
	return sprint(CodeTitle, text)
}

// Bold returns the given text in bold white.
//
// If $NO_COLOR is set, text will be returned unmodified.
func Bold(text string) string {
	return sprint(CodeBold, text)
}

// sprint returns a string with a given colour and the reset code.
//
// It handles checking for NO_COLOR and FORCE_COLOR.
func sprint(code, text string) string {
	// TODO(@FollowTheProcess): I don't like checking *every* time but doing it
	// via e.g. sync.Once means that tests are annoying unless we ensure env vars are
	// set at the process level
	noColor := os.Getenv("NO_COLOR") != ""
	forceColor := os.Getenv("FORCE_COLOR") != ""

	// $FORCE_COLOR overrides $NO_COLOR
	if forceColor {
		return code + text + CodeReset
	}

	// $NO_COLOR is next
	if noColor {
		return text
	}

	// Normal
	return code + text + CodeReset
}
