package hockey

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

//go:embed weights.json
var defaultWeightsJSON []byte

// DefaultWeightsJSON returns the embedded weights file.
func DefaultWeightsJSON() []byte { return bytes.Clone(defaultWeightsJSON) }

// Row maps a skill code (POW STA SKA QUI SHO CHK PC PI PAS KEE) to a weight
// or a factor.
type Row map[string]float64

// Weights is the whole hockey model as stored in weights.json. All numbers
// are placeholder priors (docs/hockey-eval/RESULTS.md).
type Weights struct {
	Comment string `json:"_comment,omitempty"`

	// Base weights per position. Re-normalised after modifiers.
	Base map[Position]Row `json:"base"`

	Modifiers struct {
		Mentality map[Mentality]map[Position]Row `json:"mentality"`
		Pressing  map[Pressing]map[Position]Row  `json:"pressing"`
		Slot      map[SlotKind]map[Position]Row  `json:"slot"`
	} `json:"modifiers"`

	// STA *= StaminaFactor[nLines] for Positions.
	Lines struct {
		Comment       string             `json:"_comment,omitempty"`
		StaminaFactor map[string]float64 `json:"staminaFactor"`
		Positions     []Position         `json:"positions"`
	} `json:"lines"`

	// shape(x) = 10 * (x/10)^ShapeExponent
	ShapeExponent float64 `json:"shapeExponent"`

	// score *= 1 - PerPoint * max(0, Threshold - POW), for Positions.
	PowerFloor struct {
		Threshold float64    `json:"threshold"`
		PerPoint  float64    `json:"perPoint"`
		Positions []Position `json:"positions"`
	} `json:"powerFloor"`

	// condition = 1 - FormPerPoint * max(0, FormTarget - form)
	Condition struct {
		FormTarget   float64 `json:"formTarget"`
		FormPerPoint float64 `json:"formPerPoint"`
	} `json:"condition"`

	// A skill trained on match day plays as s + TrainedBonus*(10 - s).
	TrainedBonus float64 `json:"trainedBonus"`
}

// ParseWeights decodes and checks a weights file. Unknown fields, positions,
// settings and skill codes are errors.
func ParseWeights(b []byte) (*Weights, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var w Weights
	if err := dec.Decode(&w); err != nil {
		return nil, fmt.Errorf("hockey weights: %w", err)
	}
	if err := w.validate(); err != nil {
		return nil, fmt.Errorf("hockey weights: %w", err)
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
	checkList := func(where string, ps []Position) error {
		for _, p := range ps {
			if !p.valid() {
				return fmt.Errorf("%s: unknown position %q", where, p)
			}
		}
		return nil
	}

	if err := checkPositions("base", w.Base); err != nil {
		return err
	}
	for _, p := range Positions {
		sum := 0.0
		for _, v := range w.Base[p] {
			sum += v
		}
		if sum <= 0 {
			return fmt.Errorf("base.%s: missing or weights sum to 0", p)
		}
	}
	for k, m := range w.Modifiers.Mentality {
		if !k.valid() {
			return fmt.Errorf("modifiers.mentality: unknown mentality %q", k)
		}
		if err := checkPositions("modifiers.mentality."+string(k), m); err != nil {
			return err
		}
	}
	for k, m := range w.Modifiers.Pressing {
		if !k.valid() {
			return fmt.Errorf("modifiers.pressing: unknown pressing %q", k)
		}
		if err := checkPositions("modifiers.pressing."+string(k), m); err != nil {
			return err
		}
	}
	for k, m := range w.Modifiers.Slot {
		if !k.valid() {
			return fmt.Errorf("modifiers.slot: unknown slot kind %q (want L, PP or BP)", k)
		}
		if err := checkPositions("modifiers.slot."+string(k), m); err != nil {
			return err
		}
	}
	for n := MinLines; n <= MaxLines; n++ {
		f, ok := w.Lines.StaminaFactor[strconv.Itoa(n)]
		if !ok || f <= 0 {
			return fmt.Errorf("lines.staminaFactor: missing or non-positive factor for %d lines", n)
		}
	}
	for k := range w.Lines.StaminaFactor {
		if n, err := strconv.Atoi(k); err != nil || n < MinLines || n > MaxLines {
			return fmt.Errorf("lines.staminaFactor: unknown line count %q (want %d..%d)", k, MinLines, MaxLines)
		}
	}
	if err := checkList("lines.positions", w.Lines.Positions); err != nil {
		return err
	}
	if err := checkList("powerFloor.positions", w.PowerFloor.Positions); err != nil {
		return err
	}
	if w.ShapeExponent <= 0 {
		return fmt.Errorf("shapeExponent must be > 0")
	}
	if w.TrainedBonus < 0 || w.TrainedBonus > 1 {
		return fmt.Errorf("trainedBonus must be in 0..1")
	}
	return nil
}
