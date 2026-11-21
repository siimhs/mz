package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Log in and show username and team id",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		me, err := c.Whoami()
		if err != nil {
			return err
		}
		fmt.Printf("user: %s\nteam: %s\n", me.Username, me.TeamID)
		return nil
	},
}

var playersJSON bool

var playersCmd = &cobra.Command{
	Use:   "players",
	Short: "List my squad with skills",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if c.Sport() == mzhttpclient.Hockey {
			return printHockeyPlayers(c)
		}
		players, err := c.GetMyPlayers()
		if err != nil {
			return err
		}
		if playersJSON {
			return printJSON(players)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintln(w, "Name\tAge\tSpd\tSta\tInt\tPas\tSho\tHea\tKee\tBal\tTac\tAer\tSet\tExp\tFrm\t")
		for _, p := range players {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				p.Name, p.Age, sk(p.Speed), sk(p.Stamina), sk(p.Intelligence), sk(p.Passing),
				sk(p.Shooting), sk(p.Heading), sk(p.Keeping), sk(p.BallControl), sk(p.Tackling),
				sk(p.AerialPassing), sk(p.SetPlays), sk(p.Experience), sk(p.Form))
		}
		return w.Flush()
	},
}

var (
	matchesPlayed bool
	matchesType   string
	matchesJSON   bool
)

var matchesCmd = &cobra.Command{
	Use:   "matches [team-id]",
	Short: "List upcoming (default) or played matches; defaults to my team",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		tid := ""
		if len(args) == 1 {
			tid = args[0]
		} else {
			me, err := c.Whoami()
			if err != nil {
				return err
			}
			tid = me.TeamID
		}
		status := mzhttpclient.MatchesUpcoming
		if matchesPlayed {
			status = mzhttpclient.MatchesPlayed
		}
		all, err := c.GetMatches(tid, status)
		if err != nil {
			return err
		}
		var ms []mzhttpclient.Match
		for _, m := range all {
			if matchesType == "" || m.MatchType == matchesType {
				ms = append(ms, m)
			}
		}
		if matchesJSON {
			return printJSON(ms)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, m := range ms {
			score := ""
			if m.Status != "scheduled" {
				score = m.HomeGoals + "-" + m.AwayGoals
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", m.Date.Format("Mon 2006-01-02 15:04"), m.MatchType,
				m.HomeTeamName, score, m.AwayTeamName, m.Id)
		}
		return w.Flush()
	},
}

func printHockeyPlayers(c *mzhttpclient.MZHTTPClient) error {
	players, err := c.GetMyHockeyPlayers()
	if err != nil {
		return err
	}
	if playersJSON {
		return printJSON(players)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Name\tAge\tPI\tPow\tSka\tPas\tQui\tSho\tKee\tPC\tChk\tSta\tFrm\t")
	for _, p := range players {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
			p.Name, p.Age, sk(p.PlayIntelligence), sk(p.Power), sk(p.Skating), sk(p.Passing),
			sk(p.Quickness), sk(p.Shooting), sk(p.Keeping), sk(p.PuckControl), sk(p.Checking),
			sk(p.Stamina), sk(p.Form))
	}
	return w.Flush()
}

// sk formats a skill like 7, 7.25, 7.5.
func sk(v float32) string { return fmt.Sprintf("%g", v) }

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func init() {
	playersCmd.Flags().BoolVar(&playersJSON, "json", false, "output JSON")
	matchesCmd.Flags().BoolVar(&matchesPlayed, "played", false, "show played matches instead of upcoming")
	matchesCmd.Flags().StringVar(&matchesType, "type", "", "filter by type: friendly, league, cup_group, ...")
	matchesCmd.Flags().BoolVar(&matchesJSON, "json", false, "output JSON")
	rootCmd.AddCommand(whoamiCmd, playersCmd, matchesCmd)
}
