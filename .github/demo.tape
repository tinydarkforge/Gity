# vhs tape — generates docs/demo.gif
#
# Install vhs:  brew install vhs
# Record:       vhs .github/demo.tape
#
# Pre-req: ollama running, intake_REPO set to a sandbox repo with ~5 open issues.

Output docs/demo.gif

Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Dracula"
Set TypingSpeed 45ms
Set PlaybackSpeed 1.0

Type "intake"
Enter
Sleep 2s

# Navigate issues pane
Tab
Sleep 500ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 500ms

# Open create flow
Type "c"
Sleep 1s

# Pick template
Down
Sleep 300ms
Enter
Sleep 500ms

# Title
Type "Mobile perf regression on login"
Enter
Sleep 500ms

# Paste context
Type "users report 3-4s login delay on Android after the 2.1.0 release, seeing high JS thread time in Chrome DevTools Performance panel"
Sleep 800ms
Ctrl+s
Sleep 4s

# Preview
Sleep 2s
Escape
Sleep 500ms
Type "q"
Sleep 500ms
