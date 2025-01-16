Output ./docs/img/subcommands.gif

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
Type "go build -o ./subcommands ./examples/subcommands && clear" Sleep 2s
Enter
Sleep 2s
Show
Sleep 1s

Type "./subcommands --help" Sleep 500ms  Enter

Sleep 8s

Type "clear" Sleep 500ms Enter

Sleep 1s

Type "./subcommands say --help" Sleep 500ms  Enter

Sleep 4s

Type "clear" Sleep 500ms Enter

Sleep 1s

Type "./subcommands do --help" Sleep 500ms  Enter

Sleep 4s

Type "clear" Sleep 500ms Enter

Sleep 1s

Type "./subcommands say 'hello world' --shout --thing demo --count 42" Sleep 500ms Enter

Sleep 4s

Type "clear" Sleep 500ms Enter

Sleep 1s

Type "./subcommands do 'something' --fast --count 3 --duration 5m27s -vvv" Sleep 500ms Enter

Sleep 4s

Type "clear" Sleep 500ms Enter

Sleep 1s

# Cleanup
Hide
Type "rm -f ./subcommands"
Enter
