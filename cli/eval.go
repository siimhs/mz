package cli

import (
	"fmt"
	"math"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/siimhs/mz/eval"
	"github.com/siimhs/mz/mzhttpclient"
	"github.com/spf13/cobra"
)

var (
	evalPassing     string
	evalPressing    string
	evalOpponent    string
	evalPos         string
	evalNoCondition bool
	evalJSON        bool
	evalWeights     string
	evalDump        bool

	evalMentality string
	evalSlot      string
	evalLines     int
	evalTeam      bool
)

var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Score my players per position for a tactic (weights are unvalidated priors)",
	Long: `Score my players per position for a tactic.

Football (default): positions GK CD FB DM CM W ST, team tactic --passing,
--pressing and optionally --opponent. Weights: docs/player-eval/RESULTS.md.
SP is a set-piece taker rating (low confidence).

Hockey (--sport hockey or MZ_SPORT=hockey): positions GK D C W. Tactics are
per line, so --mentality and --pressing describe the line the player is on,
--slot picks the line (L1-L4, PP1-PP2 powerplay, BP1-BP2 boxplay) and
--lines says how many regular lines the lineup fills (2-4; fewer lines make
Stamina count more). --team suggests lines L1..L<lines> greedily and shows
each line's strength (the sum of its five skater scores); ice time is a
lineup property and is not in the score. Weights: docs/hockey-eval/RESULTS.md.

Weights are unvalidated priors; small score differences mean nothing until
the model is calibrated on match data.

Without --pos, each player's best position and the scores for all positions
are shown. With --pos, the squad is ranked for that position. Scores are about
0..10. --no-condition drops form (and, in football, experience): the
long-term/transfer view.

Weights, modifiers and constants come from an embedded JSON file per sport.
--weights (or $MZ_EVAL_WEIGHTS for football, $MZ_EVAL_HOCKEY_WEIGHTS for
hockey) loads a full replacement; --dump-weights prints the current one as a
starting point.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sp, err := sport()
		if err != nil {
			return err
		}
		if err := evalCheckFlags(cmd, sp); err != nil {
			return err
		}
		if sp == mzhttpclient.Hockey {
			return runHockeyEval()
		}
		m, err := evalModel()
		if err != nil {
			return err
		}
		if evalDump {
			return printJSON(m.W)
		}
		t, err := evalTactic()
		if err != nil {
			return err
		}
		var pos eval.Position
		if evalPos != "" {
			if pos, err = eval.ParsePosition(evalPos); err != nil {
				return err
			}
		}

		c, err := newClient()
		if err != nil {
			return err
		}
		players, err := c.GetMyPlayers()
		if err != nil {
			return err
		}
		rows := make([]evalRow, len(players))
		for i, p := range players {
			s := evalSkills(p)
			best := m.BestPosition(s, t)
			scores := map[string]float64{}
			for _, b := range best {
				scores[string(b.Pos)] = round2(b.Score)
			}
			sp := round2(m.SetPieceScore(s))
			rows[i] = evalRow{ID: p.ID, Name: p.Name, Age: p.Age, Best: string(best[0].Pos),
				Score: round2(best[0].Score), Scores: scores, SetPiece: &sp}
		}
		if pos != "" {
			return printEvalRanked(rows, string(pos))
		}
		var ps []string
		for _, p := range eval.Positions {
			ps = append(ps, string(p))
		}
		return printEvalAll(rows, ps)
	},
}

// Flags that belong to one sport only.
var (
	evalFootballFlags = []string{"passing", "opponent"}
	evalHockeyFlags   = []string{"mentality", "slot", "lines", "team"}
)

func evalCheckFlags(cmd *cobra.Command, sp string) error {
	wrong, name := evalHockeyFlags, "football"
	if sp == mzhttpclient.Hockey {
		wrong, name = evalFootballFlags, "hockey"
	}
	for _, f := range wrong {
		if cmd.Flags().Changed(f) {
			return fmt.Errorf("--%s does not apply to %s (football: --passing, --opponent; hockey: --mentality, --slot, --lines, --team)", f, name)
		}
	}
	return nil
}

type evalRow struct {
	ID       int                `json:"id"`
	Name     string             `json:"name"`
	Age      int                `json:"age"`
	Best     string             `json:"best"`
	Score    float64            `json:"score"`
	Scores   map[string]float64 `json:"scores"`
	SetPiece *float64           `json:"setPiece,omitempty"` // football only
}

func printEvalAll(rows []evalRow, positions []string) error {
	if evalJSON {
		return printJSON(rows)
	}
	sp := len(rows) > 0 && rows[0].SetPiece != nil
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprint(w, "Name\tAge\tBest\t")
	for _, p := range positions {
		fmt.Fprintf(w, "%s\t", p)
	}
	if sp {
		fmt.Fprint(w, "SP\t")
	}
	fmt.Fprintln(w)
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%d\t%s\t", r.Name, r.Age, r.Best)
		for _, p := range positions {
			fmt.Fprintf(w, "%.2f\t", r.Scores[p])
		}
		if sp {
			fmt.Fprintf(w, "%.2f\t", *r.SetPiece)
		}
		fmt.Fprintln(w)
	}
	return w.Flush()
}

type rankedRow struct {
	Rank     int     `json:"rank"`
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Pos      string  `json:"pos"`
	Score    float64 `json:"score"`
	Best     string  `json:"best"`
	BestDiff float64 `json:"bestDiff"` // best score minus this position's
}

func printEvalRanked(rows []evalRow, pos string) error {
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Scores[pos] > rows[j].Scores[pos] })
	out := make([]rankedRow, len(rows))
	for i, r := range rows {
		out[i] = rankedRow{Rank: i + 1, ID: r.ID, Name: r.Name, Age: r.Age, Pos: pos, Score: r.Scores[pos],
			Best: r.Best, BestDiff: round2(r.Score - r.Scores[pos])}
		if out[i].BestDiff == 0 { // a tie at two decimals: pos is as good as any
			out[i].Best = pos
		}
	}
	if evalJSON {
		return printJSON(out)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "#\tName\tAge\t%s\tBest\t\n", pos)
	for _, r := range out {
		best := r.Best
		if r.Best != pos {
			best = fmt.Sprintf("%s +%.2f", r.Best, r.BestDiff)
		}
		fmt.Fprintf(w, "%d\t%s\t%d\t%.2f\t%s\t\n", r.Rank, r.Name, r.Age, r.Score, best)
	}
	return w.Flush()
}

func evalModel() (*eval.Model, error) {
	path := evalWeights
	if path == "" {
		path = os.Getenv("MZ_EVAL_WEIGHTS")
	}
	m := eval.New()
	if path != "" {
		w, err := eval.LoadWeights(path)
		if err != nil {
			return nil, err
		}
		m.W = w
	}
	m.SkipCondition = evalNoCondition
	return m, nil
}

func evalTactic() (eval.Tactic, error) {
	var t eval.Tactic
	var err error
	if t.Passing, err = eval.ParsePassing(evalPassing); err != nil {
		return t, err
	}
	if t.Pressing, err = eval.ParsePressing(evalPressing); err != nil {
		return t, err
	}
	if evalOpponent != "" {
		if t.OpponentPassing, err = eval.ParsePassing(evalOpponent); err != nil {
			return t, fmt.Errorf("--opponent: %w", err)
		}
	}
	return t, nil
}

func evalSkills(p mzhttpclient.FootballPlayer) eval.Skills {
	f := func(v float32) float64 { return float64(v) }
	return eval.Skills{
		Speed: f(p.Speed), Stamina: f(p.Stamina), Intelligence: f(p.Intelligence), Passing: f(p.Passing),
		Shooting: f(p.Shooting), Heading: f(p.Heading), Keeping: f(p.Keeping), BallControl: f(p.BallControl),
		Tackling: f(p.Tackling), AerialPassing: f(p.AerialPassing), SetPlays: f(p.SetPlays),
		Experience: f(p.Experience), Form: f(p.Form),
	}
}

func round2(x float64) float64 { return math.Round(x*100) / 100 }

func init() {
	f := evalCmd.Flags()
	f.StringVar(&evalPassing, "passing", "short", "passing style: short, long or wings")
	f.StringVar(&evalPressing, "pressing", "normal", "pressing: passive, normal or committed")
	f.StringVar(&evalOpponent, "opponent", "", "opponent's passing style, if known: short, long or wings")
	f.StringVar(&evalPos, "pos", "", "rank the squad for one position: GK, CD, FB, DM, CM, W or ST (hockey: GK, D, C or W)")
	f.BoolVar(&evalNoCondition, "no-condition", false, "ignore form and experience (long-term/transfer view)")
	f.StringVar(&evalMentality, "mentality", "balanced", "hockey: the line's mentality: defensive, balanced or attacking")
	f.StringVar(&evalSlot, "slot", "L1", "hockey: the line: L1-L4, PP1-PP2 (powerplay) or BP1-BP2 (boxplay)")
	f.IntVar(&evalLines, "lines", 4, "hockey: regular lines the lineup fills, 2-4")
	f.BoolVar(&evalTeam, "team", false, "hockey: suggest lines L1..L<lines> and show each line's strength")
	f.BoolVar(&evalJSON, "json", false, "output JSON")
	f.StringVar(&evalWeights, "weights", "", "weights JSON file replacing the embedded one (default $MZ_EVAL_WEIGHTS, hockey $MZ_EVAL_HOCKEY_WEIGHTS)")
	f.BoolVar(&evalDump, "dump-weights", false, "print the weights in use as JSON and exit (no login)")
	rootCmd.AddCommand(evalCmd)
}
