// Package hockey scores ice hockey players per position, line tactic and
// slot.
//
// It implements the model in docs/hockey-eval/RESULTS.md. Every weight and
// constant is an unvalidated prior and lives in weights.json (embedded; an
// override file can replace it). The package is pure: no HTTP and no
// ManagerZone types.
//
// Ice time is a lineup property (it depends on every line-mate's stamina),
// so it is not part of a player's score; only the number of regular lines
// filled is, through the Stamina weight.
package hockey

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Skills are a player's skills on MZ's 0..10 scale; fractions count.
type Skills struct {
	Power, Stamina, Skating, Quickness, Shooting, Checking,
	PuckControl, PlayIntelligence, Passing, Keeping float64
	Form float64
	// Trained lists skill codes trained on match day (optional); each plays
	// as s + trainedBonus*(10 - s).
	Trained []string
}

// Skill codes used in weights.json, in vector order.
var skillCodes = []string{"POW", "STA", "SKA", "QUI", "SHO", "CHK", "PC", "PI", "PAS", "KEE"}

// SkillCodes returns the skill codes in vector order.
func SkillCodes() []string { return append([]string(nil), skillCodes...) }

var skillIndex = func() map[string]int {
	m := map[string]int{}
	for i, c := range skillCodes {
		m[c] = i
	}
	return m
}()

const (
	iPOW = iota
	iSTA
	iSKA
	iQUI
	iSHO
	iCHK
	iPC
	iPI
	iPAS
	iKEE
	nSkills
)

type vec [nSkills]float64

func (s Skills) vec() vec {
	return vec{s.Power, s.Stamina, s.Skating, s.Quickness, s.Shooting, s.Checking,
		s.PuckControl, s.PlayIntelligence, s.Passing, s.Keeping}
}

func (r Row) vec() vec {
	var v vec
	for k, x := range r {
		v[skillIndex[k]] = x
	}
	return v
}

// Position is where a player plays. Left and right are cosmetic.
type Position string

const (
	GK Position = "GK" // goalkeeper
	D  Position = "D"  // defender (LD/RD)
	C  Position = "C"  // centre
	W  Position = "W"  // winger (LW/RW)
)

// Positions in rink order.
var Positions = []Position{GK, D, C, W}

func (p Position) valid() bool { return p == GK || p == D || p == C || p == W }

// ParsePosition accepts GK, D, C or W in any case.
func ParsePosition(s string) (Position, error) {
	p := Position(strings.ToUpper(s))
	if !p.valid() {
		return "", fmt.Errorf("unknown hockey position %q (want GK, D, C or W)", s)
	}
	return p, nil
}

// Mentality is a line's mentality.
type Mentality string

const (
	Defensive Mentality = "defensive"
	Balanced  Mentality = "balanced"
	Attacking Mentality = "attacking"
)

func (m Mentality) valid() bool { return m == Defensive || m == Balanced || m == Attacking }

// ParseMentality accepts defensive, balanced or attacking.
func ParseMentality(s string) (Mentality, error) {
	m := Mentality(strings.ToLower(s))
	if !m.valid() {
		return "", fmt.Errorf("unknown mentality %q (want defensive, balanced or attacking)", s)
	}
	return m, nil
}

// Pressing is a line's pressing.
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

// SlotKind is the kind of line: regular, powerplay or boxplay.
type SlotKind string

const (
	Regular   SlotKind = "L"
	Powerplay SlotKind = "PP"
	Boxplay   SlotKind = "BP"
)

func (k SlotKind) valid() bool { return k == Regular || k == Powerplay || k == Boxplay }

// Slot is a line in the tactic: L1..L4, PP1..PP2 or BP1..BP2.
type Slot struct {
	Kind SlotKind
	N    int
}

func (s Slot) String() string { return string(s.Kind) + strconv.Itoa(s.N) }

// L1 is the first regular line.
var L1 = Slot{Regular, 1}

// ParseSlot accepts L1..L4, PP1..PP2 or BP1..BP2 in any case.
func ParseSlot(s string) (Slot, error) {
	u := strings.ToUpper(s)
	for _, k := range []struct {
		kind SlotKind
		max  int
	}{{Powerplay, 2}, {Boxplay, 2}, {Regular, 4}} {
		if rest, ok := strings.CutPrefix(u, string(k.kind)); ok {
			if n, err := strconv.Atoi(rest); err == nil && n >= 1 && n <= k.max {
				return Slot{k.kind, n}, nil
			}
		}
	}
	return Slot{}, fmt.Errorf("unknown slot %q (want L1-L4, PP1-PP2 or BP1-BP2)", s)
}

// Line limits for Line.Lines.
const (
	MinLines = 2
	MaxLines = 4
)

// Line is where a player is evaluated: the line's tactic, its slot, and how
// many regular lines the lineup fills. Zero values mean balanced, normal,
// L1 and 4 lines.
type Line struct {
	Mentality Mentality
	Pressing  Pressing
	Slot      Slot
	Lines     int
}

func (l Line) norm() Line {
	if l.Mentality == "" {
		l.Mentality = Balanced
	}
	if l.Pressing == "" {
		l.Pressing = Normal
	}
	if l.Slot.Kind == "" {
		l.Slot = L1
	}
	if l.Lines == 0 {
		l.Lines = MaxLines
	}
	return l
}

// Model evaluates players with a set of weights.
type Model struct {
	W *Weights
	// SkipCondition drops condition(form): the long-term or transfer view.
	SkipCondition bool
}

// New returns a model with the embedded weights.
func New() *Model { return &Model{W: DefaultWeights()} }

// PositionWeights returns the normalised weights for pos on line l (sum 1).
func (m *Model) PositionWeights(pos Position, l Line) Row {
	r := Row{}
	for i, x := range m.weights(pos, l) {
		if x != 0 {
			r[skillCodes[i]] = x
		}
	}
	return r
}

func (m *Model) weights(pos Position, l Line) vec {
	l = l.norm()
	w := m.W.Base[pos].vec()
	apply := func(r Row) {
		for k, f := range r {
			w[skillIndex[k]] *= f
		}
	}
	apply(m.W.Modifiers.Mentality[l.Mentality][pos])
	apply(m.W.Modifiers.Pressing[l.Pressing][pos])
	apply(m.W.Modifiers.Slot[l.Slot.Kind][pos])
	if contains(m.W.Lines.Positions, pos) {
		if f, ok := m.W.Lines.StaminaFactor[strconv.Itoa(l.Lines)]; ok {
			w[iSTA] *= f
		}
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

// effective returns the skills as played: trained skills get the bonus.
func (m *Model) effective(s Skills) vec {
	v := s.vec()
	for _, k := range s.Trained {
		if i, ok := skillIndex[strings.ToUpper(k)]; ok {
			v[i] += m.W.TrainedBonus * (10 - v[i])
		}
	}
	return v
}

func (m *Model) shape(x float64) float64 {
	x = math.Max(0, math.Min(10, x))
	return 10 * math.Pow(x/10, m.W.ShapeExponent)
}

// Evaluate scores a player at pos on line l, about 0..10.
func (m *Model) Evaluate(s Skills, pos Position, l Line) float64 {
	v := m.effective(s)
	w := m.weights(pos, l)
	score := 0.0
	for i := range w {
		score += w[i] * m.shape(v[i])
	}
	if pf := m.W.PowerFloor; contains(pf.Positions, pos) {
		score *= math.Max(0, 1-pf.PerPoint*math.Max(0, pf.Threshold-v[iPOW]))
	}
	return score * m.condition(s)
}

func (m *Model) condition(s Skills) float64 {
	if m.SkipCondition {
		return 1
	}
	c := m.W.Condition
	return math.Max(0, 1-c.FormPerPoint*math.Max(0, c.FormTarget-s.Form))
}

// Score is a player's score at one position.
type Score struct {
	Pos   Position `json:"pos"`
	Score float64  `json:"score"`
}

// BestPosition scores every position, best first. Ties keep rink order.
func (m *Model) BestPosition(s Skills, l Line) []Score {
	out := make([]Score, 0, len(Positions))
	for _, p := range Positions {
		out = append(out, Score{p, m.Evaluate(s, p, l)})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

// LineShape is the five skaters of a line.
var LineShape = []Position{C, W, W, D, D}

// Member is one skater picked for a line.
type Member struct {
	Player int      `json:"player"` // index into the players passed in
	Pos    Position `json:"pos"`
	Score  float64  `json:"score"`
}

// LineUp is one filled line and its strength (the sum of its five skater
// scores). It says nothing about ice time.
type LineUp struct {
	Slot     Slot     `json:"-"`
	Members  []Member `json:"members"`
	Strength float64  `json:"strength"`
}

// LineStrength sums the five skaters' scores.
func LineStrength(scores []float64) float64 {
	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	return sum
}

// PickLines fills regular lines L1..Ln greedily (a suggestion, not an
// optimum): for each line in turn, the highest remaining (player, position)
// score takes a free C, W or D place, until the line has C, 2 W and 2 D.
// A player is used once. Lines stop early when skaters run out. Every line
// is scored with l's mentality and pressing and n lines.
func (m *Model) PickLines(players []Skills, n int, l Line) []LineUp {
	type cand struct {
		i     int
		pos   Position
		score float64
	}
	used := make([]bool, len(players))
	var out []LineUp
	for k := 1; k <= n; k++ {
		ll := l
		ll.Slot, ll.Lines = Slot{Regular, k}, n
		var cs []cand
		for i, s := range players {
			if used[i] {
				continue
			}
			for _, p := range []Position{C, W, D} {
				cs = append(cs, cand{i, p, m.Evaluate(s, p, ll)})
			}
		}
		sort.SliceStable(cs, func(a, b int) bool { return cs[a].score > cs[b].score })
		free := map[Position]int{C: 1, W: 2, D: 2}
		line := LineUp{Slot: ll.Slot}
		for _, c := range cs {
			if used[c.i] || free[c.pos] == 0 {
				continue
			}
			used[c.i] = true
			free[c.pos]--
			line.Members = append(line.Members, Member{c.i, c.pos, c.score})
		}
		if len(line.Members) < len(LineShape) {
			break
		}
		// Shape order: C, W, W, D, D.
		order := map[Position]int{C: 0, W: 1, D: 2}
		sort.SliceStable(line.Members, func(a, b int) bool { return order[line.Members[a].Pos] < order[line.Members[b].Pos] })
		scores := make([]float64, len(line.Members))
		for j, mb := range line.Members {
			scores[j] = mb.Score
		}
		line.Strength = LineStrength(scores)
		out = append(out, line)
	}
	return out
}

func contains(ps []Position, p Position) bool {
	for _, q := range ps {
		if p == q {
			return true
		}
	}
	return false
}
