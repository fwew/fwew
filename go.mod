module github.com/fwew/fwew/v5

go 1.25.0

require (
	github.com/c-bata/go-prompt v0.2.6
	github.com/fwew/fwew-lib/v5 v5.27.2
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.20 // indirect
	github.com/mattn/go-tty v0.0.7 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

//for testing on a local machine's fwew-lib
//replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
