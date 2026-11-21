package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

var hockeyPositions = []string{"GK", "D", "C", "W"}

// The recorded hockey squad (testdata/mz/hockey/players.html) has 71 players.
const hockeySquad = 71

func TestHockeyPlayersTable(t *testing.T) {
	mz := Start(t, HockeyLogin())

	res := mz.Run("--sport", "hockey", "players")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "hockey-players.txt")
}

func TestHockeyEvalTable(t *testing.T) {
	mz := Start(t, HockeyLogin())

	res := mz.Run("--sport", "hockey", "eval", "--mentality", "defensive", "--lines", "3")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "hockey-eval-defensive-3.txt")
	if strings.Contains(res.Stdout, "SP") {
		t.Error("hockey table has a set-piece column")
	}
}

func TestHockeyEvalJSON(t *testing.T) {
	mz := Start(t, HockeyLogin())
	mz.Env["MZ_SPORT"] = "hockey"

	res := mz.Run("eval", "--json")

	res.WantExit(t, 0)
	var rows []map[string]any
	if err := json.Unmarshal([]byte(res.Stdout), &rows); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if len(rows) != hockeySquad {
		t.Fatalf("%d players, the page has %d", len(rows), hockeySquad)
	}
	keepers := 0
	for _, r := range rows {
		if _, ok := r["setPiece"]; ok {
			t.Errorf("%v: hockey row has setPiece", r["name"])
		}
		scores := r["scores"].(map[string]any)
		if len(scores) != len(hockeyPositions) {
			t.Errorf("%v: %d scores", r["name"], len(scores))
		}
		best, score := r["best"].(string), r["score"].(float64)
		for _, p := range hockeyPositions {
			if s := scores[p].(float64); s < 0 || s > 10 || s > score {
				t.Errorf("%v: %s = %v, best %s = %v", r["name"], p, s, best, score)
			}
		}
		if scores[best].(float64) != score {
			t.Errorf("%v: best %s is not the top score", r["name"], best)
		}
		if best == "GK" {
			keepers++
		}
	}
	if keepers == 0 {
		t.Error("nobody's best position is GK")
	}
}

func TestHockeyEvalRanked(t *testing.T) {
	mz := Start(t, HockeyLogin())

	res := mz.Run("--sport", "hockey", "eval", "--pos", "c", "--slot", "pp1", "--pressing", "committed")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "hockey-eval-c-pp1.txt")
	if strings.Contains(res.Stdout, "+0.00") {
		t.Error("a tie is shown as +0.00")
	}
}

func TestHockeyEvalTeam(t *testing.T) {
	mz := Start(t, HockeyLogin())

	res := mz.Run("--sport", "hockey", "eval", "--team", "--lines", "4")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "hockey-eval-team.txt")

	js := mz.Run("--sport", "hockey", "eval", "--team", "--json")
	js.WantExit(t, 0)
	var tv struct {
		Goalkeeper struct{ ID int }
		Lines      []struct {
			Slot     string
			Strength float64
			Members  []struct {
				ID    int
				Pos   string
				Score float64
			}
		}
	}
	if err := json.Unmarshal([]byte(js.Stdout), &tv); err != nil {
		t.Fatal(err)
	}
	if len(tv.Lines) != 4 || tv.Goalkeeper.ID == 0 {
		t.Fatalf("want a keeper and 4 lines: %+v", tv)
	}
	seen := map[int]bool{tv.Goalkeeper.ID: true}
	for _, l := range tv.Lines {
		sum := 0.0
		for _, m := range l.Members {
			if seen[m.ID] {
				t.Errorf("%s: player %d used twice", l.Slot, m.ID)
			}
			seen[m.ID] = true
			sum += m.Score
		}
		if len(l.Members) != 5 || sum-l.Strength > 0.03 || l.Strength-sum > 0.03 {
			t.Errorf("%s: %d members summing to %v, strength %v", l.Slot, len(l.Members), sum, l.Strength)
		}
	}
}

func TestEvalSportFlags(t *testing.T) {
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--sport", "hockey", "eval", "--passing", "wings"}, "--passing does not apply to hockey"},
		{[]string{"--sport", "hockey", "eval", "--opponent", "long"}, "--opponent does not apply to hockey"},
		{[]string{"eval", "--mentality", "attacking"}, "--mentality does not apply to football"},
		{[]string{"eval", "--slot", "PP1"}, "--slot does not apply to football"},
		{[]string{"eval", "--lines", "3"}, "--lines does not apply to football"},
		{[]string{"eval", "--team"}, "--team does not apply to football"},
		{[]string{"--sport", "hockey", "eval", "--pos", "ST"}, `"ST"`},
		{[]string{"--sport", "hockey", "eval", "--mentality", "ultra"}, `"ultra"`},
		{[]string{"--sport", "hockey", "eval", "--slot", "L5"}, `"L5"`},
		{[]string{"--sport", "hockey", "eval", "--lines", "5"}, "--lines 5"},
		{[]string{"--sport", "hockey", "eval", "--slot", "L4", "--lines", "3"}, "not filled"},
		{[]string{"--sport", "hockey", "eval", "--team", "--pos", "C"}, "--team"},
	} {
		mz := Start(t) // no routes: must fail before logging in
		res := mz.Run(c.args...)
		res.WantExit(t, 1)
		res.WantStderrContains(t, c.want)
	}
}

func TestHockeyEvalDumpWeights(t *testing.T) {
	res := Start(t).Run("--sport", "hockey", "eval", "--dump-weights")

	res.WantExit(t, 0)
	for _, want := range []string{"unvalidated prior", `"powerFloor"`, `"PP"`} {
		if !strings.Contains(res.Stdout, want) {
			t.Errorf("dump lacks %s: %.200s", want, res.Stdout)
		}
	}
	// The football override variable does not leak into hockey.
	mz := Start(t)
	mz.Env["MZ_EVAL_WEIGHTS"] = "/nonexistent.json"
	mz.Run("--sport", "hockey", "eval", "--dump-weights").WantExit(t, 0)
}
