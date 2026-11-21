package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var positions = []string{"GK", "CD", "FB", "DM", "CM", "W", "ST"}

type evalRow struct {
	ID       int
	Name     string
	Best     string
	Score    float64
	Scores   map[string]float64
	SetPiece float64
}

func TestEvalTable(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("eval", "--passing", "wings")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "eval-wings.txt")
}

func TestEvalJSON(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("eval", "--json")

	res.WantExit(t, 0)
	var rows []evalRow
	if err := json.Unmarshal([]byte(res.Stdout), &rows); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if len(rows) != 46 {
		t.Fatalf("%d players, the page has 46", len(rows))
	}
	if rows[0].Name != "Aurelio Theodore" || rows[0].ID != 224227751 {
		t.Errorf("first player: %+v", rows[0])
	}
	keepers := 0
	for _, r := range rows {
		if len(r.Scores) != len(positions) {
			t.Errorf("%s: %d scores", r.Name, len(r.Scores))
		}
		for _, p := range positions {
			if s := r.Scores[p]; s < 0 || s > 10 || s > r.Score {
				t.Errorf("%s: %s = %v, best %s = %v", r.Name, p, s, r.Best, r.Score)
			}
		}
		if r.Scores[r.Best] != r.Score {
			t.Errorf("%s: best %s scores %v, not %v", r.Name, r.Best, r.Scores[r.Best], r.Score)
		}
		if r.Best == "GK" {
			keepers++
		}
	}
	if keepers == 0 {
		t.Error("nobody's best position is GK")
	}
}

func TestEvalRankedTable(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("eval", "--pos", "ST", "--passing", "wings", "--opponent", "wings")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "eval-st-wings.txt")
	if strings.Contains(res.Stdout, "+0.00") {
		t.Error("a tie is shown as +0.00")
	}
}

func TestEvalRankedForPosition(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("eval", "--pos", "gk", "--pressing", "committed", "--no-condition", "--json")

	res.WantExit(t, 0)
	var rows []struct {
		Rank  int
		Name  string
		Pos   string
		Score float64
	}
	if err := json.Unmarshal([]byte(res.Stdout), &rows); err != nil || len(rows) != 46 {
		t.Fatalf("want 46 rows: %v (%d)", err, len(rows))
	}
	for i, r := range rows {
		if r.Rank != i+1 || r.Pos != "GK" {
			t.Errorf("row %d: %+v", i, r)
		}
		if i > 0 && r.Score > rows[i-1].Score {
			t.Errorf("not sorted at %d: %v > %v", i, r.Score, rows[i-1].Score)
		}
	}
}

func TestEvalBadFlags(t *testing.T) {
	for _, args := range [][]string{
		{"eval", "--passing", "tiki-taka"},
		{"eval", "--pressing", "gegen"},
		{"eval", "--pos", "LW"},
		{"eval", "--weights", "/nonexistent.json"},
	} {
		mz := Start(t) // no routes: must fail before logging in
		res := mz.Run(args...)
		res.WantExit(t, 1)
		res.WantStderrContains(t, args[2])
	}
}

func TestEvalWeightsOverride(t *testing.T) {
	// --dump-weights needs no login; edit its output into an override.
	dump := Start(t).Run("eval", "--dump-weights")
	dump.WantExit(t, 0)
	if !strings.Contains(dump.Stdout, "unvalidated prior") {
		t.Errorf("dump lacks the caveat: %.200s", dump.Stdout)
	}

	// Only Keeping matters, for every position, with no core term and no
	// fatigue: the squad's best keeper ranks first as a striker.
	var w map[string]any
	if err := json.Unmarshal([]byte(dump.Stdout), &w); err != nil {
		t.Fatal(err)
	}
	base := map[string]any{}
	for _, p := range positions {
		base[p] = map[string]float64{"KEE": 1}
	}
	w["base"] = base
	w["blend"].(map[string]any)["position"] = 1.0
	w["blend"].(map[string]any)["core"] = 0.0
	w["fatigue"].(map[string]any)["perPoint"] = 0.0
	b, _ := json.Marshal(w)
	path := filepath.Join(t.TempDir(), "w.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}

	mz := Start(t, Login())
	var players []player
	if err := json.Unmarshal([]byte(mz.Run("players", "--json").Stdout), &players); err != nil {
		t.Fatal(err)
	}
	maxKee := 0.0
	for _, p := range players {
		if p.Keeping > maxKee {
			maxKee = p.Keeping
		}
	}

	mz.Env["MZ_EVAL_WEIGHTS"] = path
	res := mz.Run("eval", "--pos", "ST", "--no-condition", "--json")
	res.WantExit(t, 0)
	var rows []struct {
		Name  string
		Score float64
	}
	if err := json.Unmarshal([]byte(res.Stdout), &rows); err != nil || len(rows) == 0 {
		t.Fatalf("no rows: %v", err)
	}
	var top player
	for _, p := range players {
		if p.Name == rows[0].Name {
			top = p
		}
	}
	if top.Keeping != maxKee {
		t.Errorf("keeping-only ST #1 is %s with Keeping %v; the best keeper has %v", top.Name, top.Keeping, maxKee)
	}

	// A broken override is an error, not silently ignored.
	bad := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(bad, []byte(`{"base": {"GK": {"KEEP": 1}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Start(t).Run("eval", "--weights", bad)
	r.WantExit(t, 1)
	r.WantStderrContains(t, "unknown skill")
}
