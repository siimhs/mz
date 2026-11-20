package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// scrub replaces personal values in recorded files before they are committed.
// The map file has one "real => fake" pair per line; keep it out of git.
// Values in file names are replaced too (the file is renamed). With -check
// nothing is changed and the exit status is 1 if a value is found.
func scrubMain(args []string) {
	fset := flag.NewFlagSet("scrub", flag.ExitOnError)
	mapFile := fset.String("map", "", "file with 'real => fake' lines (required)")
	check := fset.Bool("check", false, "only report, do not change files")
	fset.Parse(args)
	if *mapFile == "" || fset.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: fakemz scrub -map FILE [-check] PATH...")
		os.Exit(2)
	}
	pairs, err := readScrubMap(*mapFile)
	if err != nil {
		log.Fatal(err)
	}
	found := false
	for _, root := range fset.Args() {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out := data
			var hits []string
			for _, p := range pairs {
				if n := bytes.Count(out, p[0]); n > 0 {
					hits = append(hits, fmt.Sprintf("%q x%d", p[0], n))
					out = bytes.ReplaceAll(out, p[0], p[1])
				}
			}
			newName := d.Name()
			for _, p := range pairs {
				newName = strings.ReplaceAll(newName, string(p[0]), string(p[1]))
			}
			if newName != d.Name() {
				hits = append(hits, "file name")
			}
			if len(hits) == 0 {
				return nil
			}
			found = true
			if *check {
				fmt.Printf("%s: %s\n", path, strings.Join(hits, ", "))
				return nil
			}
			fmt.Printf("%s: replaced %s\n", path, strings.Join(hits, ", "))
			if err := os.WriteFile(path, out, 0o644); err != nil {
				return err
			}
			if newName != d.Name() {
				return os.Rename(path, filepath.Join(filepath.Dir(path), newName))
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	if *check && found {
		os.Exit(1)
	}
}

func readScrubMap(path string) ([][2][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var pairs [][2][]byte
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		real, fake, ok := strings.Cut(line, "=>")
		real, fake = strings.TrimSpace(real), strings.TrimSpace(fake)
		if !ok || real == "" {
			return nil, fmt.Errorf("%s:%d: want 'real => fake'", path, n)
		}
		pairs = append(pairs, [2][]byte{[]byte(real), []byte(fake)})
	}
	// Longest first, so a long value is replaced before a shorter value it contains.
	sort.SliceStable(pairs, func(i, j int) bool { return len(pairs[i][0]) > len(pairs[j][0]) })
	return pairs, sc.Err()
}
