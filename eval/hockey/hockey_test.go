package hockey

import (
	"math"
	"strings"
	"testing"
)

const eps = 1e-9

// all returns skills with every skill set to v and form 9 (no condition penalty).
func all(v float64) Skills {
	return Skills{v, v, v, v, v, v, v, v, v, v, 9, nil}
}

func shape(x float64) float64 { return 10 * math.Pow(x/10, 0.8) }

var allLines = func() []Line {
	var ls []Line
	for _, me := range []Mentality{Defensive, Balanced, Attacking} {
		for _, pr := range []Pressing{Passive, Normal, Committed} {
			for _, sl := range []string{"L1", "L4", "PP1", "PP2", "BP1", "BP2"} {
				s, _ := ParseSlot(sl)
				for n := MinLines; n <= MaxLines; n++ {
					ls = append(ls, Line{me, pr, s, n})
				}
			}
		}
	}
	return ls
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
	for _, l := range allLines {
		for _, p := range Positions {
			sum := 0.0
			for _, v := range m.PositionWeights(p, l) {
				sum += v
			}
			if math.Abs(sum-1) > eps {
				t.Errorf("%s %+v sums to %v", p, l, sum)
			}
		}
	}
}

func TestModifiersByHand(t *testing.T) {
	m := New()
	// D, attacking, powerplay: SHO .06*1.6*1.5, QUI .11*1.1, STA .15*1.1,
	// CHK .18*0.9; the row then sums to 1.092.
	w := m.PositionWeights(D, Line{Attacking, Normal, Slot{Powerplay, 1}, 4})
	sum := 0.28 + 0.165 + 0.12 + 0.121 + 0.144 + 0.162 + 0.03 + 0.05 + 0.02
	if math.Abs(w["SHO"]-0.144/sum) > eps || math.Abs(w["POW"]-0.28/sum) > eps {
		t.Errorf("D attacking PP: SHO %v POW %v, want %v %v", w["SHO"], w["POW"], 0.144/sum, 0.28/sum)
	}
	// W, 2 lines, committed: STA .15*1.6*1.25 = .30; the rest unchanged.
	w = m.PositionWeights(W, Line{Balanced, Committed, L1, 2})
	if want := 0.30 / 1.15; math.Abs(w["STA"]-want) > eps {
		t.Errorf("W 2 lines committed: STA %v, want %v", w["STA"], want)
	}
	// The GK has no modifiers, and the line count does not touch his Stamina.
	base := m.PositionWeights(GK, Line{})
	for _, l := range allLines {
		for k, v := range m.PositionWeights(GK, l) {
			if math.Abs(v-base[k]) > eps {
				t.Errorf("GK %s changed under %+v: %v vs %v", k, l, v, base[k])
			}
		}
	}
}

func TestScoresByHand(t *testing.T) {
	m := New()
	cases := []struct {
		name string
		s    Skills
		pos  Position
		want float64
	}{
		{"all 10", all(10), C, 10},
		{"all 8, at the power floor", all(8), W, shape(8)},
		{"power 6 costs 8%", func() Skills { s := all(8); s.Power = 6; return s }(), D,
			(0.28*shape(6) + 0.72*shape(8)) * 0.92},
		{"form 7 costs 10%", func() Skills { s := all(10); s.Form = 7; return s }(), C, 9},
		{"keeper: no power floor", func() Skills { s := all(5); s.Keeping = 10; s.Power = 0; return s }(), GK,
			0.62*10 + 0.38*shape(5)},
		{"trained skill", func() Skills { s := all(9); s.Trained = []string{"pow"}; return s }(), D,
			0.28*shape(9.1) + 0.72*shape(9)},
	}
	for _, c := range cases {
		if got := m.Evaluate(c.s, c.pos, Line{}); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
	skip := &Model{W: DefaultWeights(), SkipCondition: true}
	s := all(10)
	s.Form = 3
	if got := skip.Evaluate(s, W, Line{}); math.Abs(got-10) > eps {
		t.Errorf("SkipCondition: %v, want 10", got)
	}
}

func TestDirections(t *testing.T) {
	m := New()
	checker := all(7)
	checker.Checking = 10
	shooter := all(7)
	shooter.Shooting = 10
	// Defensive wingers value Checking more; powerplay values Shooting.
	if m.Evaluate(checker, W, Line{Mentality: Defensive}) <= m.Evaluate(checker, W, Line{}) {
		t.Error("defensive W does not favour a checker")
	}
	if m.Evaluate(shooter, D, Line{Slot: Slot{Powerplay, 1}}) <= m.Evaluate(shooter, D, Line{}) {
		t.Error("PP D does not favour a shooter")
	}
	// Fewer lines: low stamina hurts more.
	tired := all(9)
	tired.Stamina = 4
	if m.Evaluate(tired, C, Line{Lines: 2}) >= m.Evaluate(tired, C, Line{Lines: 4}) {
		t.Error("2 lines do not penalise low stamina more than 4")
	}
	best := m.BestPosition(func() Skills { s := all(3); s.Keeping = 10; return s }(), Line{})
	if best[0].Pos != GK {
		t.Errorf("keeper's best position is %s", best[0].Pos)
	}
	best = m.BestPosition(checker, Line{})
	if best[0].Pos != D {
		t.Errorf("checker's best position is %s", best[0].Pos)
	}
}

func TestParse(t *testing.T) {
	for in, want := range map[string]Slot{"l1": {Regular, 1}, "L4": {Regular, 4}, "pp2": {Powerplay, 2}, "BP1": {Boxplay, 1}} {
		got, err := ParseSlot(in)
		if err != nil || got != want {
			t.Errorf("ParseSlot(%q) = %v, %v", in, got, err)
		}
		if strings.ToUpper(in) != got.String() {
			t.Errorf("%v.String() = %q", got, got.String())
		}
	}
	for _, bad := range []string{"L5", "L0", "PP3", "BP", "X1", ""} {
		if _, err := ParseSlot(bad); err == nil {
			t.Errorf("ParseSlot(%q) accepted", bad)
		}
	}
	if _, err := ParsePosition("LW"); err == nil {
		t.Error("LW accepted")
	}
	if p, err := ParsePosition("d"); err != nil || p != D {
		t.Errorf("ParsePosition(d) = %v, %v", p, err)
	}
	if _, err := ParseMentality("ultra"); err == nil {
		t.Error("ultra mentality accepted")
	}
}

func TestPickLines(t *testing.T) {
	m := New()
	var players []Skills
	for i := 0; i < 12; i++ {
		players = append(players, all(float64(10-i%6)))
	}
	lines := m.PickLines(players, 3, Line{})
	if len(lines) != 2 {
		t.Fatalf("%d lines from 12 skaters, want 2", len(lines))
	}
	seen := map[int]bool{}
	for k, l := range lines {
		if l.Slot != (Slot{Regular, k + 1}) {
			t.Errorf("line %d slot %v", k, l.Slot)
		}
		for j, mb := range l.Members {
			if mb.Pos != LineShape[j] {
				t.Errorf("line %d member %d is %s", k, j, mb.Pos)
			}
			if seen[mb.Player] {
				t.Errorf("player %d used twice", mb.Player)
			}
			seen[mb.Player] = true
		}
	}
	if lines[0].Strength <= lines[1].Strength {
		t.Errorf("L1 %v not stronger than L2 %v", lines[0].Strength, lines[1].Strength)
	}
}

func TestParseWeightsRejectsTypos(t *testing.T) {
	for _, bad := range []string{
		`{"base": {"GK": {"KEEP": 1}}}`,
		`{"bases": {}}`,
		strings.Replace(string(defaultWeightsJSON), `"PP": {`, `"PQ": {`, 1),
		strings.Replace(string(defaultWeightsJSON), `"attacking": {`, `"offensive": {`, 1),
		strings.Replace(string(defaultWeightsJSON), `"2": 1.6`, `"5": 1.6`, 1),
	} {
		if _, err := ParseWeights([]byte(bad)); err == nil {
			t.Errorf("accepted: %.80s", bad)
		}
	}
}
