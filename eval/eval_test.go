package eval

import (
	"math"
	"strings"
	"testing"
)

const eps = 1e-9

// all returns skills with every 0..10 skill set to v, form 9 and exp 10.
func all(v float64) Skills {
	return Skills{v, v, v, v, v, v, v, v, v, v, v, 10, 9}
}

var (
	keeper  = func() Skills { s := all(3); s.Keeping = 10; s.Intelligence = 7; return s }()
	tackler = func() Skills {
		s := all(4)
		s.Tackling, s.Intelligence, s.Speed, s.Stamina = 10, 8, 7, 7
		return s
	}()
	apWinger = func() Skills {
		s := all(5)
		s.AerialPassing, s.Speed, s.BallControl, s.Stamina = 10, 8, 7, 8
		return s
	}()
	heaStriker = func() Skills {
		s := all(5)
		s.Heading, s.Shooting, s.BallControl = 10, 8, 7
		return s
	}()
)

var allTactics = func() []Tactic {
	var ts []Tactic
	for _, pa := range []Passing{Short, Long, Wings} {
		for _, pr := range []Pressing{Passive, Normal, Committed} {
			for _, op := range []Passing{"", Short, Long, Wings} {
				ts = append(ts, Tactic{pa, pr, op})
			}
		}
	}
	return ts
}()

func TestBaseRowsSumToOne(t *testing.T) {
	w := DefaultWeights()
	for _, p := range Positions {
		sum := 0.0
		for _, v := range w.Base[p] {
			sum += v
		}
		if math.Abs(sum-1) > eps {
			t.Errorf("base %s sums to %v", p, sum)
		}
	}
}

func TestModifiedRowsSumToOne(t *testing.T) {
	m := New()
	for _, tc := range allTactics {
		for _, p := range Positions {
			sum := 0.0
			for _, v := range m.PositionWeights(p, tc) {
				sum += v
			}
			if math.Abs(sum-1) > eps {
				t.Errorf("%s %+v sums to %v", p, tc, sum)
			}
		}
	}
}

func TestModifiersChangeWeights(t *testing.T) {
	m := New()
	short := m.PositionWeights(W, Tactic{Short, Normal, ""})
	wings := m.PositionWeights(W, Tactic{Wings, Normal, ""})
	if wings["AP"] <= short["AP"] {
		t.Errorf("W AP weight: wings %v <= short %v", wings["AP"], short["AP"])
	}
}

func TestBestPosition(t *testing.T) {
	tests := []struct {
		name   string
		skills Skills
		want   []Position // best position must be one of these
	}{
		{"pure keeper", keeper, []Position{GK}},
		{"tackler", tackler, []Position{CD, DM}},
	}
	m := New()
	for _, tt := range tests {
		for _, tc := range allTactics {
			got := m.BestPosition(tt.skills, tc)
			if len(got) != len(Positions) {
				t.Fatalf("%s: %d scores", tt.name, len(got))
			}
			for i := 1; i < len(got); i++ {
				if got[i].Score > got[i-1].Score {
					t.Errorf("%s: not sorted: %v", tt.name, got)
				}
			}
			ok := false
			for _, p := range tt.want {
				ok = ok || got[0].Pos == p
			}
			if !ok {
				t.Errorf("%s %+v: best %s, want one of %v (%v)", tt.name, tc, got[0].Pos, tt.want, got)
			}
		}
	}
}

func TestWingsVsShort(t *testing.T) {
	tests := []struct {
		name   string
		skills Skills
		pos    Position
	}{
		{"high-AP winger", apWinger, W},
		{"high-HEA striker", heaStriker, ST},
	}
	m := New()
	for _, tt := range tests {
		for _, pr := range []Pressing{Passive, Normal, Committed} {
			short := m.Evaluate(tt.skills, tt.pos, Tactic{Short, pr, ""})
			wings := m.Evaluate(tt.skills, tt.pos, Tactic{Wings, pr, ""})
			if wings <= short {
				t.Errorf("%s at %s, %s pressing: wings %.3f <= short %.3f", tt.name, tt.pos, pr, wings, short)
			}
		}
	}
}

func TestCommittedPenalisesLowStamina(t *testing.T) {
	m := New()
	for _, p := range Positions {
		low, high := all(7), all(7)
		low.Stamina, high.Stamina = 4, 10
		ratio := func(pr Pressing) float64 {
			tc := Tactic{Short, pr, ""}
			return m.Evaluate(low, p, tc) / m.Evaluate(high, p, tc)
		}
		if p == GK {
			continue // keeper stamina is capped at 7 and barely matters
		}
		if c, n := ratio(Committed), ratio(Normal); c >= n {
			t.Errorf("%s: low/high stamina committed %.3f >= normal %.3f", p, c, n)
		}
	}
}

func TestCondition(t *testing.T) {
	w := DefaultWeights()
	tests := []struct {
		form, exp, want float64
	}{
		{9, 10, 1},
		{10, 10, 1},
		{7, 10, 0.90},
		{9, 0, 0.80},
		{7, 5, 0.90 * 0.90},
	}
	for _, tt := range tests {
		if got := Condition(w, tt.form, tt.exp); math.Abs(got-tt.want) > eps {
			t.Errorf("condition(form %v, exp %v) = %v, want %v", tt.form, tt.exp, got, tt.want)
		}
	}
}

func TestSkipCondition(t *testing.T) {
	s := all(7)
	s.Form, s.Experience = 5, 2
	m := New()
	tc := Tactic{Short, Normal, ""}
	with := m.Evaluate(s, CM, tc)
	m.SkipCondition = true
	without := m.Evaluate(s, CM, tc)
	if want := without * Condition(m.W, 5, 2); math.Abs(with-want) > eps {
		t.Errorf("with condition %v, want %v", with, want)
	}
}

func TestEvaluateFormula(t *testing.T) {
	// All skills 10, form 9, exp 10: every factor is 1, shape(10) = 10.
	// (A keeper's stamina is capped at 7, so GK stays below 10.)
	for _, p := range Positions[1:] {
		if got := Evaluate(all(10), p, Tactic{Wings, Committed, Wings}); math.Abs(got-10) > eps {
			t.Errorf("%s: all-10 player scores %v, want 10", p, got)
		}
	}
	// shape(5) = 10 * 0.5^0.8; with all skills 5 the score is shape(5) * fatigue.
	m := New()
	shape5 := 10 * math.Pow(0.5, 0.8)
	if got, want := m.Evaluate(all(5), CD, Tactic{Short, Normal, ""}), shape5*(1-0.04*3); math.Abs(got-want) > eps {
		t.Errorf("all-5 CD: %v, want %v", got, want)
	}
	// Keeper floor: Keeping 0 counts as 4.
	a, b := all(5), all(5)
	a.Keeping, b.Keeping = 0, 4
	if x, y := m.Evaluate(a, GK, Tactic{Short, Normal, ""}), m.Evaluate(b, GK, Tactic{Short, Normal, ""}); x != y {
		t.Errorf("keeper floor: %v != %v", x, y)
	}
	// Keeper stamina cap: 7 and 10 are the same.
	a, b = all(5), all(5)
	a.Stamina, b.Stamina = 7, 10
	if x, y := m.Evaluate(a, GK, Tactic{Short, Normal, ""}), m.Evaluate(b, GK, Tactic{Short, Normal, ""}); x != y {
		t.Errorf("keeper stamina cap: %v != %v", x, y)
	}
}

func TestSetPieceScore(t *testing.T) {
	s := Skills{SetPlays: 10, AerialPassing: 6, Shooting: 2}
	if got, want := New().SetPieceScore(s), 0.35*10+0.35*6+0.30*2; math.Abs(got-want) > eps {
		t.Errorf("SetPieceScore = %v, want %v", got, want)
	}
}

func TestParseWeightsRejects(t *testing.T) {
	good := string(DefaultWeightsJSON())
	tests := []struct {
		name, old, new, wantErr string
	}{
		{"unknown skill", `"KEE": 0.58`, `"KEEP": 0.58`, "unknown skill"},
		{"unknown field", `"shapeExponent"`, `"shapeExp"`, "unknown field"},
		{"unknown position", `"W":  {"AP": 1.4`, `"LW": {"AP": 1.4`, "unknown position"},
		{"unknown style", `"committed": {`, `"commited": {`, "unknown style"},
		{"negative", `"KEE": 0.58`, `"KEE": -0.58`, "negative"},
	}
	for _, tt := range tests {
		if !strings.Contains(good, tt.old) {
			t.Fatalf("%s: %q not in weights.json", tt.name, tt.old)
		}
		_, err := ParseWeights([]byte(strings.Replace(good, tt.old, tt.new, 1)))
		if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
			t.Errorf("%s: err %v, want %q", tt.name, err, tt.wantErr)
		}
	}
}

func TestOverrideChangesScores(t *testing.T) {
	good := string(DefaultWeightsJSON())
	w, err := ParseWeights([]byte(strings.Replace(good, `"shapeExponent": 0.8`, `"shapeExponent": 1`, 1)))
	if err != nil {
		t.Fatal(err)
	}
	m := &Model{W: w}
	// Linear shape: an all-6 GK scores 6 (stamina under the cap, at the fatigue threshold).
	if got := m.Evaluate(all(6), GK, Tactic{Short, Normal, ""}); math.Abs(got-6) > eps {
		t.Errorf("all-6 GK with linear shape: %v", got)
	}
}
