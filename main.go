package main

import (
	// Embed the timezone database so the static binary works on systems
	// without /usr/share/zoneinfo (e.g. minimal containers).
	_ "time/tzdata"

	"github.com/siimhs/mz/cli"
)

func main() {
	cli.Execute()
}
