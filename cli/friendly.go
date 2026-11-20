package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/spf13/cobra"
)

var friendlyCmd = &cobra.Command{
	Use:     "friendly",
	Aliases: []string{"friendlies"},
	Short:   "Arrange friendly matches",
}

var friendlyCandidatesCmd = &cobra.Command{
	Use:   "candidates",
	Short: "List team ids found on the quick-challenge page",
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
		tids, err := candidates(c, me.TeamID)
		if err != nil {
			return err
		}
		for _, t := range tids {
			fmt.Println(t)
		}
		return nil
	},
}

var (
	sendDate       string
	sendHomeHour   int
	sendAwayHour   int
	sendTactic     string
	sendTacticAway string
)

var friendlySendCmd = &cobra.Command{
	Use:   "send <team-id>",
	Short: "Send one friendly invite (a home and an away match on the same day)",
	Example: `  mz friendly send 121197 --date 2026-11-22
  mz friendly send 121197 --date 2026-11-22 --home 10 --away 18 --tactic a --tactic-away b`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		day, err := time.ParseInLocation("2006-01-02", sendDate, time.Local)
		if err != nil {
			return fmt.Errorf("--date: %w", err)
		}
		ta := sendTacticAway
		if ta == "" {
			ta = sendTactic
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		home, away := at(day, sendHomeHour), at(day, sendAwayHour)
		if err := c.SendFriendlyRequest(args[0], home, sendTactic, away, ta); err != nil {
			return err
		}
		fmt.Printf("sent invite to %s for %s\n", args[0], day.Format("Mon 2006-01-02"))
		return nil
	},
}

var (
	autoDays     int
	autoHomeHour int
	autoAwayHour int
	autoTactics  string
	autoDryRun   bool
)

var friendlyAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Invite all quick-challenge teams for every upcoming day without a friendly",
	Long: `For each of the next --days days (starting tomorrow) that has no friendly yet,
send an invite to every team on the quick-challenge page.

This is the old main.go behaviour, made safer. It is a stepping stone towards
'auto-arrange' (accept invites, accept/post quick challenges, then invite).`,
	Example: `  mz friendly auto --dry-run
  mz friendly auto --days 5 --tactics tue=a,wed=a,default=b`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tactics, err := parseTactics(autoTactics)
		if err != nil {
			return err
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		me, err := c.Whoami()
		if err != nil {
			return err
		}
		booked, err := friendlyDays(c, me.TeamID)
		if err != nil {
			return err
		}
		tids, err := candidates(c, me.TeamID)
		if err != nil {
			return err
		}
		fmt.Printf("%d candidate teams: %s\n", len(tids), strings.Join(tids, ", "))

		today := truncateDay(now())
		sent, failed := 0, 0
		for i := 1; i <= autoDays; i++ {
			day := today.AddDate(0, 0, i)
			label := day.Format("Mon 2006-01-02")
			if booked[day] {
				fmt.Printf("%s: already has a friendly, skipping\n", label)
				continue
			}
			tactic := tactics.forDay(day.Weekday())
			fmt.Printf("%s: inviting %d teams (tactic %s, %02d:00 home / %02d:00 away)\n",
				label, len(tids), tactic, autoHomeHour, autoAwayHour)
			if autoDryRun {
				continue
			}
			for _, t := range tids {
				err := c.SendFriendlyRequest(t, at(day, autoHomeHour), tactic, at(day, autoAwayHour), tactic)
				if err != nil {
					failed++
					fmt.Printf("  %s: %v\n", t, firstLine(err))
				} else {
					sent++
					fmt.Printf("  %s: sent\n", t)
				}
			}
		}
		if autoDryRun {
			fmt.Println("dry run: nothing sent")
		} else {
			fmt.Printf("sent %d, failed %d\n", sent, failed)
		}
		return nil
	},
}

// candidates returns unique team ids from the quick-challenge page, excluding my own.
func candidates(c *mzhttpclient.MZHTTPClient, myTID string) ([]string, error) {
	raw, err := c.GetTeamIDsOnQuickChallengePage()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{myTID: true}
	var out []string
	for _, t := range raw {
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out, nil
}

// friendlyDays returns the local calendar days that already have an upcoming friendly.
func friendlyDays(c *mzhttpclient.MZHTTPClient, tid string) (map[time.Time]bool, error) {
	ms, err := c.GetMatches(tid, mzhttpclient.MatchesUpcoming)
	if err != nil {
		return nil, err
	}
	days := map[time.Time]bool{}
	for _, m := range ms {
		if m.MatchType == mzhttpclient.Friendly {
			days[truncateDay(m.Date)] = true
		}
	}
	return days, nil
}

type tacticPlan struct {
	byDay       map[time.Weekday]string
	defaultTact string
}

func (p tacticPlan) forDay(d time.Weekday) string {
	if t, ok := p.byDay[d]; ok {
		return t
	}
	return p.defaultTact
}

// parseTactics parses "tue=a,wed=a,default=b" (a bare "b" means default=b).
func parseTactics(s string) (tacticPlan, error) {
	p := tacticPlan{byDay: map[time.Weekday]string{}, defaultTact: "a"}
	days := map[string]time.Weekday{}
	for d := time.Sunday; d <= time.Saturday; d++ {
		days[strings.ToLower(d.String()[:3])] = d
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			k, v = "default", part
		}
		if !validTactic(v) {
			return p, fmt.Errorf("--tactics: %q is not a tactic (a-h)", v)
		}
		if k == "default" {
			p.defaultTact = v
			continue
		}
		d, ok := days[k]
		if !ok {
			keys := make([]string, 0, len(days))
			for k := range days {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return p, fmt.Errorf("--tactics: unknown day %q (use %s or default)", k, strings.Join(keys, ","))
		}
		p.byDay[d] = v
	}
	return p, nil
}

func validTactic(t string) bool { return len(t) == 1 && t[0] >= 'a' && t[0] <= 'h' }

func truncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func at(day time.Time, hour int) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), hour, 0, 0, 0, time.Local)
}

func firstLine(err error) string {
	s, _, _ := strings.Cut(err.Error(), "\n")
	return s
}

func init() {
	friendlySendCmd.Flags().StringVar(&sendDate, "date", "", "day of the friendly, YYYY-MM-DD (local time)")
	friendlySendCmd.Flags().IntVar(&sendHomeHour, "home", 12, "hour of the home match (local time)")
	friendlySendCmd.Flags().IntVar(&sendAwayHour, "away", 20, "hour of the away match (local time)")
	friendlySendCmd.Flags().StringVar(&sendTactic, "tactic", "a", "tactic for the home match (a-h)")
	friendlySendCmd.Flags().StringVar(&sendTacticAway, "tactic-away", "", "tactic for the away match (default: same as --tactic)")
	friendlySendCmd.MarkFlagRequired("date")

	friendlyAutoCmd.Flags().IntVar(&autoDays, "days", 5, "how many days ahead to fill, starting tomorrow")
	friendlyAutoCmd.Flags().IntVar(&autoHomeHour, "home", 12, "hour of the home match (local time)")
	friendlyAutoCmd.Flags().IntVar(&autoAwayHour, "away", 20, "hour of the away match (local time)")
	friendlyAutoCmd.Flags().StringVar(&autoTactics, "tactics", "tue=a,wed=a,default=b", "tactic per weekday")
	friendlyAutoCmd.Flags().BoolVar(&autoDryRun, "dry-run", false, "show what would be sent without sending")

	friendlyCmd.AddCommand(friendlyCandidatesCmd, friendlySendCmd, friendlyAutoCmd)
	rootCmd.AddCommand(friendlyCmd)
}
