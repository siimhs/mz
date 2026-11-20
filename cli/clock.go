package cli

import (
	"fmt"
	"os"
	"time"
)

// now returns the current time, or MZ_NOW (RFC3339) when set. MZ_NOW makes
// date-dependent commands reproducible; the e2e tests rely on it.
func now() time.Time {
	v := os.Getenv("MZ_NOW")
	if v == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "MZ_NOW: %v\n", err)
		os.Exit(2)
	}
	return t.In(time.Local)
}
