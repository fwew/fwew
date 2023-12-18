module fwew

go 1.20

require (
	github.com/c-bata/go-prompt v0.2.3
	github.com/fwew/fwew-lib/v5 v5.7.1-dev.0.20231215061223-07bc3664c516
)

require (
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v0.0.0-20200520122047-c3ffed290a03 // indirect
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
)

//for testing on a local machine's fwew-lib
//replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
