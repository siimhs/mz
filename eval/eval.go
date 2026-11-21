// Package eval scores football players per position and tactic.
//
// It implements the model in docs/player-eval/RESULTS.md. Every weight and
// constant is an unvalidated prior and lives in weights.json (embedded; an
// override file can replace it), so the model can be tuned without code
// changes. The package is pure: no HTTP and no ManagerZone types.
package eval

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Skills are a player's skills on MZ's 0..10 scale; fractions count.
type Skills struct {
	Speed, Stamina, Intelligence, Passing, Shooting, Heading, Keeping,
	BallControl, Tackling, AerialPassing, SetPlays float64
	Experience, Form float64
}

// Skill codes used in weights.json, in vector order.
var skillCodes = []string{"SPD", "STA", "PI", "PAS", "SHO", "HEA", "KEE", "BC", "TAC", "AP", "SP"}

var skillIndex = func() map[string]int {
	m := map[string]int{}
	for i, c := range skillCodes {
		m[c] = i
	}
	return m
}()

const (
	iSPD = iota
	iSTA
	iPI
	iPAS
	iSHO
	iHEA
	iKEE
	iBC
	iTAC
	iAP
	iSP
	nSkills
)

type vec [nSkills]float64

func (s Skills) vec() vec {
	return vec{s.Speed, s.Stamina, s.Intelligence, s.Passing, s.Shooting, s.Heading,
		s.Keeping, s.BallControl, s.Tackling, s.AerialPassing, s.SetPlays}
}

func (r Row) vec() vec {
	var v vec
	for k, x := range r {
		v[skillIndex[k]] = x
	}
	return v
}

// Position is where a player plays.
type Position string

const (
	GK Position = "GK" // goalkeeper
	CD Position = "CD" // central defender
	FB Position = "FB" // wide defender
	DM Position = "DM" // holding midfielder
	CM Position = "CM" // central/attacking midfielder
	W  Position = "W"  // wide midfielder or forward
	ST Position = "ST" // central forward
)

// Positions in field order.
var Positions = []Position{GK, CD, FB, DM, CM, W, ST}

func (p Position) valid() bool {
	for _, q := range Positions {
		if p == q {
			return true
		}
	}
	return false
}

// ParsePosition accepts a position code in any case.
func ParsePosition(s string) (Position, error) {
	p := Position(strings.ToUpper(s))
	if !p.valid() {
		return "", fmt.Errorf("unknown position %q (want GK, CD, FB, DM, CM, W or ST)", s)
	}
	return p, nil
}

// Passing is a passing style.
type Passing string

const (
	Short Passing = "short"
	Long  Passing = "long"
	Wings Passing = "wings"
)

func (p Passing) valid() bool { return p == Short || p == Long || p == Wings }

// ParsePassing accepts short, long or wings.
func ParsePassing(s string) (Passing, error) {
	p := Passing(strings.ToLower(s))
	if !p.valid() {
		return "", fmt.Errorf("unknown passing style %q (want short, long or wings)", s)
	}
	return p, nil
}

// Pressing is a pressing (marking) style.
type Pressing string

const (
	Passive   Pressing = "passive"
	Normal    Pressing = "normal"
	Committed Pressing = "committed"
)

func (p Pressing) valid() bool { return p == Passive || p == Normal || p == Committed }

// ParsePressing accepts passive, normal or committed.
func ParsePressing(s string) (Pressing, error) {
	p := Pressing(strings.ToLower(s))
	if !p.valid() {
		return "", fmt.Errorf("unknown pressing %q (want passive, normal or committed)", s)
	}
	return p, nil
}

// Tactic is our tactic plus, optionally, the opponent's passing style.
type Tactic struct {
	Passing  Passing
	Pressing Pressing
	// OpponentPassing is "" when unknown.
	OpponentPassing Passing
}

// Model evaluates players with a set of weights.
type Model struct {
	W *Weights
	// SkipCondition drops condition(form, exp): the long-term or transfer view.
	SkipCondition bool
}

// New returns a model with the embedded weights.
func New() *Model { return &Model{W: DefaultWeights()} }

// Evaluate scores a player at a position with the default model.
func Evaluate(s Skills, pos Position, t Tactic) float64 { return New().Evaluate(s, pos, t) }

// PositionWeights returns the normalised weights for pos under t (sum 1).
func (m *Model) PositionWeights(pos Position, t Tactic) Row {
	w := m.weights(pos, t)
	r := Row{}
	for i, x := range w {
		if x != 0 {
			r[skillCodes[i]] = x
		}
	}
	return r
}

func (m *Model) weights(pos Position, t Tactic) vec {
	w := m.W.Base[pos].vec()
	apply := func(r Row) {
		for k, f := range r {
			w[skillIndex[k]] *= f
		}
	}
	apply(m.W.Modifiers.Passing[t.Passing][pos])
	apply(m.W.Modifiers.Pressing[t.Pressing][pos])
	if t.OpponentPassing != "" {
		apply(m.W.Modifiers.OpponentPassing[t.OpponentPassing][pos])
	}
	sum := 0.0
	for _, x := range w {
		sum += x
	}
	if sum > 0 {
		for i := range w {
			w[i] /= sum
		}
	}
	return w
}

// Evaluate scores a player at a position, about 0..10.
func (m *Model) Evaluate(s Skills, pos Position, t Tactic) float64 {
	return m.skill(s, pos, t) * m.fatigue(s, pos, t.Pressing) * m.condition(s)
}

// effective returns the skills as the model sees them at pos.
func (m *Model) effective(s Skills, pos Position) vec {
	v := s.vec()
	if pos == GK {
		v[iKEE] = math.Max(v[iKEE], m.W.GK.KeepingFloor)
		v[iSTA] = math.Min(v[iSTA], m.W.GK.StaminaCap)
	}
	return v
}

func (m *Model) shape(x float64) float64 {
	x = math.Max(0, math.Min(10, x))
	return 10 * math.Pow(x/10, m.W.ShapeExponent)
}

func (m *Model) skill(s Skills, pos Position, t Tactic) float64 {
	v := m.effective(s, pos)
	w := m.weights(pos, t)
	posScore := 0.0
	for i := range w {
		posScore += w[i] * m.shape(v[i])
	}
	if pos == GK {
		return posScore
	}
	b := m.W.Blend
	core := 0.0
	for _, k := range b.CoreSkills {
		core += m.shape(v[skillIndex[k]])
	}
	if len(b.CoreSkills) > 0 {
		core /= float64(len(b.CoreSkills))
	}
	return b.Position*posScore + b.Core*core
}

func (m *Model) fatigue(s Skills, pos Position, pr Pressing) float64 {
	f := m.W.Fatigue
	t := f.Threshold[pos]
	if pr == Committed {
		t += f.CommittedExtra
	}
	sta := m.effective(s, pos)[iSTA]
	return math.Max(0, 1-f.PerPoint*math.Max(0, t-sta))
}

func (m *Model) condition(s Skills) float64 {
	if m.SkipCondition {
		return 1
	}
	return Condition(m.W, s.Form, s.Experience)
}

// Condition is the present-strength factor: 1 at form >= target and exp 10.
func Condition(w *Weights, form, exp float64) float64 {
	c := w.Condition
	mf := 1 - c.FormPerPoint*math.Max(0, c.FormTarget-form)
	me := c.ExpBase + c.ExpPerPoint*exp
	return math.Max(0, mf) * math.Max(0, me)
}

// Score is a player's score at one position.
type Score struct {
	Pos   Position `json:"pos"`
	Score float64  `json:"score"`
}

// BestPosition scores every position, best first. Ties keep field order.
func (m *Model) BestPosition(s Skills, t Tactic) []Score {
	out := make([]Score, 0, len(Positions))
	for _, p := range Positions {
		out = append(out, Score{p, m.Evaluate(s, p, t)})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

// SetPieceScore rates a set-piece taker (low confidence; not part of any
// position score): a weighted sum of raw skills, 0.35 SP + 0.35 AP + 0.30 SHO
// by default.
func (m *Model) SetPieceScore(s Skills) float64 {
	v := s.vec()
	w := m.W.SetPiece.vec()
	sum := 0.0
	for i := range w {
		sum += w[i] * v[i]
	}
	return sum
}
