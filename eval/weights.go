package eval

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed weights.json
var defaultWeightsJSON []byte

// DefaultWeightsJSON returns the embedded weights file, e.g. as a starting
// point for an override file.
func DefaultWeightsJSON() []byte { return bytes.Clone(defaultWeightsJSON) }

// Row maps a skill code (SPD, STA, PI, PAS, SHO, HEA, KEE, BC, TAC, AP, SP) to
// a weight or a factor.
type Row map[string]float64

// Weights is the whole model as stored in weights.json. All numbers are
// placeholder priors (docs/player-eval/RESULTS.md).
type Weights struct {
	Comment string `json:"_comment,omitempty"`

	// Base weights per position. Re-normalised after modifiers.
	Base map[Position]Row `json:"base"`

	Modifiers struct {
		Passing         map[Passing]map[Position]Row  `json:"passing"`
		Pressing        map[Pressing]map[Position]Row `json:"pressing"`
		OpponentPassing map[Passing]map[Position]Row  `json:"opponentPassing"`
	} `json:"modifiers"`

	GK struct {
		KeepingFloor float64 `json:"keepingFloor"`
		StaminaCap   float64 `json:"staminaCap"`
	} `json:"gk"`

	// shape(x) = 10 * (x/10)^ShapeExponent
	ShapeExponent float64 `json:"shapeExponent"`

	// Outfield: Position*posScore + Core*mean(shape(CoreSkills)).
	Blend struct {
		Position   float64  `json:"position"`
		Core       float64  `json:"core"`
		CoreSkills []string `json:"coreSkills"`
	} `json:"blend"`

	// fatigue = 1 - PerPoint * max(0, Threshold[pos] (+CommittedExtra) - STA)
	Fatigue struct {
		PerPoint       float64              `json:"perPoint"`
		CommittedExtra float64              `json:"committedExtra"`
		Threshold      map[Position]float64 `json:"threshold"`
	} `json:"fatigue"`

	// condition = (1 - FormPerPoint*max(0, FormTarget-form)) * (ExpBase + ExpPerPoint*exp)
	Condition struct {
		FormTarget   float64 `json:"formTarget"`
		FormPerPoint float64 `json:"formPerPoint"`
		ExpBase      float64 `json:"expBase"`
		ExpPerPoint  float64 `json:"expPerPoint"`
	} `json:"condition"`

	// Weighted sum of raw skills.
	SetPiece Row `json:"setPiece"`
}

// ParseWeights decodes and checks a weights file. Unknown fields, positions,
// tactics and skill codes are errors, so a typo in an override cannot be
// silently ignored.
func ParseWeights(b []byte) (*Weights, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var w Weights
	if err := dec.Decode(&w); err != nil {
		return nil, fmt.Errorf("weights: %w", err)
	}
	if err := w.validate(); err != nil {
		return nil, fmt.Errorf("weights: %w", err)
	}
	return &w, nil
}

// LoadWeights reads a weights file from disk.
func LoadWeights(path string) (*Weights, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	w, err := ParseWeights(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return w, nil
}

// DefaultWeights returns a fresh copy of the embedded weights.
func DefaultWeights() *Weights {
	w, err := ParseWeights(defaultWeightsJSON)
	if err != nil {
		panic("embedded " + err.Error())
	}
	return w
}

func (w *Weights) validate() error {
	checkRow := func(where string, r Row) error {
		for k, v := range r {
			if _, ok := skillIndex[k]; !ok {
				return fmt.Errorf("%s: unknown skill %q", where, k)
			}
			if v < 0 {
				return fmt.Errorf("%s: %s is negative", where, k)
			}
		}
		return nil
	}
	checkPositions := func(where string, m map[Position]Row) error {
		for p, r := range m {
			if !p.valid() {
				return fmt.Errorf("%s: unknown position %q", where, p)
			}
			if err := checkRow(fmt.Sprintf("%s.%s", where, p), r); err != nil {
				return err
			}
		}
		return nil
	}

	if err := checkPositions("base", w.Base); err != nil {
		return err
	}
	for _, p := range Positions {
		r, ok := w.Base[p]
		if !ok {
			return fmt.Errorf("base: missing position %s", p)
		}
		sum := 0.0
		for _, v := range r {
			sum += v
		}
		if sum <= 0 {
			return fmt.Errorf("base.%s: weights sum to 0", p)
		}
		if _, ok := w.Fatigue.Threshold[p]; !ok {
			return fmt.Errorf("fatigue.threshold: missing position %s", p)
		}
	}
	for t, m := range w.Modifiers.Passing {
		if !t.valid() {
			return fmt.Errorf("modifiers.passing: unknown style %q", t)
		}
		if err := checkPositions("modifiers.passing."+string(t), m); err != nil {
			return err
		}
	}
	for t, m := range w.Modifiers.Pressing {
		if !t.valid() {
			return fmt.Errorf("modifiers.pressing: unknown style %q", t)
		}
		if err := checkPositions("modifiers.pressing."+string(t), m); err != nil {
			return err
		}
	}
	for t, m := range w.Modifiers.OpponentPassing {
		if !t.valid() {
			return fmt.Errorf("modifiers.opponentPassing: unknown style %q", t)
		}
		if err := checkPositions("modifiers.opponentPassing."+string(t), m); err != nil {
			return err
		}
	}
	for p := range w.Fatigue.Threshold {
		if !p.valid() {
			return fmt.Errorf("fatigue.threshold: unknown position %q", p)
		}
	}
	for _, k := range w.Blend.CoreSkills {
		if _, ok := skillIndex[k]; !ok {
			return fmt.Errorf("blend.coreSkills: unknown skill %q", k)
		}
	}
	if w.Blend.Core > 0 && len(w.Blend.CoreSkills) == 0 {
		return fmt.Errorf("blend: core weight without coreSkills")
	}
	if w.ShapeExponent <= 0 {
		return fmt.Errorf("shapeExponent must be > 0")
	}
	return checkRow("setPiece", w.SetPiece)
}
