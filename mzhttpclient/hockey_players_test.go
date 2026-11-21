package mzhttpclient

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// The recorded, scrubbed hockey players page used by the e2e tests.
func TestScrapeHockeyPlayersFixture(t *testing.T) {
	b, err := os.ReadFile("../e2e/testdata/mz/hockey/players.html")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(strings.NewReader(string(b)))
	if err != nil {
		t.Fatal(err)
	}
	players, err := scrapeHockeyPlayers(doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 71 {
		t.Fatalf("got %d players, the page has 71", len(players))
	}
	// First player, checked by hand against the page: Form 9 + 4px bar,
	// PI 9, Power 1.
	p := players[0]
	if p.Name != "Joosep Vanamets" || p.ID != 34807132 || p.Age != 25 || p.Shoots != "Left" ||
		p.Height != 168 || p.Weight != 63 || p.BirthSeason != 67 ||
		p.Value != 760324 || p.ValueCurrency != "EUR" || p.Salary != 8582 || p.SalaryCurrency != "EUR" {
		t.Errorf("first player: %+v", p)
	}
	if p.Form != 9.5 || p.PlayIntelligence != 9 || p.Power != 1 {
		t.Errorf("first player skills: Form %v PI %v Power %v", p.Form, p.PlayIntelligence, p.Power)
	}
	fractions, youth := 0, 0
	for _, p := range players {
		for _, v := range []float32{p.PlayIntelligence, p.Power, p.Skating, p.Passing, p.Quickness,
			p.Shooting, p.Keeping, p.PuckControl, p.Checking, p.Stamina, p.Form} {
			if v < 0 || v > 10.75 {
				t.Errorf("%s: skill %v out of range", p.Name, v)
			}
			if v != float32(int(v)) {
				fractions++
			}
		}
		if p.Salary == 0 {
			youth++
		}
		if p.Name == "" || p.ID == 0 || p.Age < 15 {
			t.Errorf("implausible player: %+v", p)
		}
	}
	if fractions == 0 {
		t.Error("no fractional skill read")
	}
	if youth == 0 || youth == len(players) {
		t.Errorf("%d players without salary", youth)
	}
}
