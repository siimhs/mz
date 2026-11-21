package cli

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/siimhs/mz/eval/hockey"
	"github.com/siimhs/mz/mzhttpclient"
)

func runHockeyEval() error {
	m, err := hockeyModel()
	if err != nil {
		return err
	}
	if evalDump {
		return printJSON(m.W)
	}
	l, err := hockeyLine()
	if err != nil {
		return err
	}
	var pos hockey.Position
	if evalPos != "" {
		if pos, err = hockey.ParsePosition(evalPos); err != nil {
			return err
		}
	}
	if evalTeam && (evalPos != "" || l.Slot != hockey.L1) {
		return errors.New("--team picks its own positions and lines; drop --pos and --slot")
	}

	c, err := newClient()
	if err != nil {
		return err
	}
	players, err := c.GetMyHockeyPlayers()
	if err != nil {
		return err
	}
	skills := make([]hockey.Skills, len(players))
	for i, p := range players {
		skills[i] = hockeySkills(p)
	}
	if evalTeam {
		return printHockeyTeam(m, players, skills, l)
	}
	rows := make([]evalRow, len(players))
	for i, p := range players {
		best := m.BestPosition(skills[i], l)
		scores := map[string]float64{}
		for _, b := range best {
			scores[string(b.Pos)] = round2(b.Score)
		}
		rows[i] = evalRow{ID: p.ID, Name: p.Name, Age: p.Age, Best: string(best[0].Pos),
			Score: round2(best[0].Score), Scores: scores}
	}
	if pos != "" {
		return printEvalRanked(rows, string(pos))
	}
	var ps []string
	for _, p := range hockey.Positions {
		ps = append(ps, string(p))
	}
	return printEvalAll(rows, ps)
}

type teamMember struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Pos   string  `json:"pos"`
	Score float64 `json:"score"`
}

type teamLine struct {
	Slot     string       `json:"slot"`
	Strength float64      `json:"strength"`
	Members  []teamMember `json:"members"`
}

type teamView struct {
	Goalkeeper *teamMember `json:"goalkeeper"`
	Lines      []teamLine  `json:"lines"`
}

// printHockeyTeam shows the best keeper and greedily picked lines. The
// keeper is picked first and is not used as a skater.
func printHockeyTeam(m *hockey.Model, players []mzhttpclient.HockeyPlayer, skills []hockey.Skills, l hockey.Line) error {
	var tv teamView
	gk, best := -1, 0.0
	for i, s := range skills {
		if sc := m.Evaluate(s, hockey.GK, l); gk < 0 || sc > best {
			gk, best = i, sc
		}
	}
	var skaters []hockey.Skills
	var index []int // skaters[j] is players[index[j]]
	for i, s := range skills {
		if i != gk {
			skaters = append(skaters, s)
			index = append(index, i)
		}
	}
	if gk >= 0 {
		tv.Goalkeeper = &teamMember{players[gk].ID, players[gk].Name, string(hockey.GK), round2(best)}
	}
	for _, lu := range m.PickLines(skaters, l.Lines, l) {
		tl := teamLine{Slot: lu.Slot.String(), Strength: round2(lu.Strength)}
		for _, mb := range lu.Members {
			p := players[index[mb.Player]]
			tl.Members = append(tl.Members, teamMember{p.ID, p.Name, string(mb.Pos), round2(mb.Score)})
		}
		tv.Lines = append(tv.Lines, tl)
	}
	if evalJSON {
		return printJSON(tv)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if g := tv.Goalkeeper; g != nil {
		fmt.Fprintf(w, "GK\t%s\t%.2f\n", g.Name, g.Score)
	}
	for _, tl := range tv.Lines {
		fmt.Fprintf(w, "%s\tstrength\t%.2f\n", tl.Slot, tl.Strength)
		for _, mb := range tl.Members {
			fmt.Fprintf(w, "  %s\t%s\t%.2f\n", mb.Pos, mb.Name, mb.Score)
		}
	}
	if len(tv.Lines) < l.Lines {
		fmt.Fprintf(w, "(only %d of %d lines: not enough skaters)\n", len(tv.Lines), l.Lines)
	}
	return w.Flush()
}

func hockeyModel() (*hockey.Model, error) {
	path := evalWeights
	if path == "" {
		path = os.Getenv("MZ_EVAL_HOCKEY_WEIGHTS")
	}
	m := hockey.New()
	if path != "" {
		w, err := hockey.LoadWeights(path)
		if err != nil {
			return nil, err
		}
		m.W = w
	}
	m.SkipCondition = evalNoCondition
	return m, nil
}

func hockeyLine() (hockey.Line, error) {
	var l hockey.Line
	var err error
	if l.Mentality, err = hockey.ParseMentality(evalMentality); err != nil {
		return l, err
	}
	if l.Pressing, err = hockey.ParsePressing(evalPressing); err != nil {
		return l, err
	}
	if l.Slot, err = hockey.ParseSlot(evalSlot); err != nil {
		return l, err
	}
	if evalLines < hockey.MinLines || evalLines > hockey.MaxLines {
		return l, fmt.Errorf("--lines %d: want %d..%d", evalLines, hockey.MinLines, hockey.MaxLines)
	}
	if l.Slot.Kind == hockey.Regular && l.Slot.N > evalLines {
		return l, fmt.Errorf("--slot %s with --lines %d: that line is not filled", l.Slot, evalLines)
	}
	l.Lines = evalLines
	return l, nil
}

func hockeySkills(p mzhttpclient.HockeyPlayer) hockey.Skills {
	f := func(v float32) float64 { return float64(v) }
	return hockey.Skills{
		Power: f(p.Power), Stamina: f(p.Stamina), Skating: f(p.Skating), Quickness: f(p.Quickness),
		Shooting: f(p.Shooting), Checking: f(p.Checking), PuckControl: f(p.PuckControl),
		PlayIntelligence: f(p.PlayIntelligence), Passing: f(p.Passing), Keeping: f(p.Keeping),
		Form: f(p.Form),
	}
}
