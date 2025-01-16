Output ./docs/img/namedargs.gif

Set FontSize 18
Set FontFamily "Geist Mono"
Set Theme "Catppuccin Macchiato"
Set WindowBar Colorful
Set Padding 5
Set Height 800
Set Width 1200
Set Margin 40
Set MarginFill "#7983FF"
Set BorderRadius 10
Set TypingSpeed 100ms

# Setup
Hide
Type "go build -o ./namedargs ./examples/namedargs && clear" Sleep 2s
Enter
Sleep 2s
Show
Sleep 1s

Type "./namedargs --help" Sleep 500ms Enter

Sleep 8s

Type "clear" Sleep 500ms Enter

Sleep 1s

# Cleanup
Hide
Type "rm -f ./namedargs"
Enter
