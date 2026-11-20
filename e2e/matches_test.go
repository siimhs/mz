package e2e

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

type match struct {
	Id, HomeTeamId, AwayTeamId, HomeGoals, AwayGoals, MatchType, Status string
	Date                                                                time.Time
}

func TestMatchesUpcoming(t *testing.T) {
	mz := Start(t, Login(), UpcomingMatches())

	res := mz.Run("matches")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "matches-upcoming.txt")
}

func TestMatchesOfAnotherTeam(t *testing.T) {
	mz := Start(t, Login(),
		Route("GET", "/xml/team_matchlist.php?team_id=849080&match_status=2", File("matches-upcoming.xml")))

	res := mz.Run("matches", "849080")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "matches-upcoming.txt") // same list, asked for by team id
	if n := len(mz.Find("GET", "/xml/team_matchlist.php*")); n != 1 {
		t.Errorf("%d match list requests, want 1", n)
	}
}

func TestMatchesPlayedFriendlies(t *testing.T) {
	mz := Start(t, Login(), PlayedMatches())

	res := mz.Run("matches", "--played", "--type", "friendly", "--json")

	res.WantExit(t, 0)
	var ms []match
	if err := json.Unmarshal([]byte(res.Stdout), &ms); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, res.Stdout)
	}
	xml, err := os.ReadFile("testdata/mz/matches-played.xml")
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Count(string(xml), `type="friendly"`); len(ms) != want {
		t.Errorf("%d friendlies, the list has %d", len(ms), want)
	}
	for _, m := range ms {
		if m.MatchType != "friendly" || m.Status != "played" || (m.HomeTeamId != myTeam && m.AwayTeamId != myTeam) {
			t.Errorf("unexpected match: %+v", m)
		}
		if m.HomeGoals == "" || m.AwayGoals == "" {
			t.Errorf("played match without score: %+v", m)
		}
	}
	// Dates are MZ server time, shown in the local zone.
	if len(ms) > 0 && ms[0].Date.Format(time.RFC3339) != "2026-11-20T20:00:00+02:00" {
		t.Errorf("first match date %v", ms[0].Date)
	}
}

func TestRaw(t *testing.T) {
	mz := Start(t, Login(), UpcomingMatches())

	res := mz.Run("raw", "xml/team_matchlist.php?team_id="+myTeam+"&match_status=2")

	res.WantExit(t, 0)
	want, err := os.ReadFile("testdata/mz/matches-upcoming.xml")
	if err != nil {
		t.Fatal(err)
	}
	if res.Stdout != string(want) {
		t.Errorf("raw output is not the page byte for byte:\n%s", diff(string(want), res.Stdout))
	}
}
