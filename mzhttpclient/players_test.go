package mzhttpclient

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// testdata/players-live-*.html is a raw (mz raw "?p=players") "?p=players" page. It holds account data, so it is
// gitignored; the test skips when it is absent.
func TestScrapePlayersFixture(t *testing.T) {
	b, err := os.ReadFile("../testdata/players-live-2026-11-20.html")
	if err != nil {
		t.Skip("no fixture:", err)
	}
	doc, err := html.Parse(strings.NewReader(string(b)))
	if err != nil {
		t.Fatal(err)
	}
	c := findByTagAndAttr(doc, atom.Div, "id", "players_container")
	if c == nil || c.FirstChild == nil {
		t.Fatal("players_container not found")
	}
	players, err := scrapePlayers(c.FirstChild, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) < 11 {
		t.Fatalf("got %d players, want a full squad", len(players))
	}
	for _, p := range players {
		if p.Name == "" || p.Age < 15 || p.Speed < 0 || p.Speed > 10 {
			t.Errorf("implausible player: %+v", p)
		}
	}
}
