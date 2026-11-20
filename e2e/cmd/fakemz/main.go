// fakemz is the fixture tool for mz's e2e tests: a recording proxy to the real
// site and a scrubber for recorded files. It shares no code with mz.
package main

import (
	"fmt"
	"os"
)

const usage = `usage: fakemz <command> [flags]

commands:
  record  proxy to the real site and save responses as a case
  scrub   replace personal values in recorded files
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "record":
		recordMain(os.Args[2:])
	case "scrub":
		scrubMain(os.Args[2:])
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
}
