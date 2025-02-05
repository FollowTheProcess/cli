// Package colour implements basic text colouring for cli's limited needs.
//
// In particular, it's not expected to provide every ANSI code, just the ones we need. The codes have also been padded so that they are
// the same length, which means [text/tabwriter] will correctly calculate alignment as long as styles are not mixed within a table.
package colour

import (
	"os"
	"sync"
	"sync/atomic"
)

// ANSI codes for coloured output, they are all the same length so as not to throw off
// alignment of [text/tabwriter].
const (
	CodeReset = "\x1b[000000m" // Reset all attributes
	CodeTitle = "\x1b[1;39;4m" // Bold, white & underlined
	CodeBold  = "\x1b[1;0039m" // Bold & white
)

// Disable is a flag that disables all colour text, it overrides both
// $FORCE_COLOR and $NO_COLOR, setting it to true will always make this
// package return plain text and not check any other config.
var Disable atomic.Bool

// getColourOnce is a [sync.OnceValues] function that returns the state of
// $NO_COLOR and $FORCE_COLOR, once and only once to avoid us calling
// os.Getenv on every call to a colour function.
var getColourOnce = sync.OnceValues(getColour)

// getColour returns whether $NO_COLOR and $FORCE_COLOR were set.
func getColour() (noColour, forceColour bool) {
	no := os.Getenv("NO_COLOR") != ""
	force := os.Getenv("FORCE_COLOR") != ""

	return no, force
}

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
// It handles checking for NO_COLOR and FORCE_COLOR. If the global var
// [Disable] is true then nothing else is checked and plain text is returned.
func sprint(code, text string) string {
	// Our global variable is above all else
	if Disable.Load() {
		return text
	}

	noColor, forceColor := getColourOnce()

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
