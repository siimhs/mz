package cli

import (
	"fmt"
	"strings"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/spf13/cobra"
)

var (
	autoDays     int
	autoHomeHour int
	autoAwayHour int
	autoTactics  string
	autoDryRun   bool
)

var friendlyAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Fill every bookable day with a friendly: accept invites, quick challenge, invite",
	Long: `Fills the game days MZ lets you book (up to 7 days ahead, not league days)
that have no friendly yet, in this order:

  1. accept invites other teams have sent us (one per day)
  2. accept a quick challenge another team offers for the next game day
  3. otherwise opt in to the quick challenge for the next game day
  4. invite every team seen on the quick-challenge page for each other
     free day (teams already invited for that day are skipped)

Bookable days come from MZ's friendly schedule page.`,
	Example: `  mz friendly auto --dry-run
  mz friendly auto --tactics tue=a,wed=a,default=b`,
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
		t := now()
		sched, err := c.GetSchedule(t)
		if err != nil {
			return err
		}
		withMatch, err := friendlyDays(c, me.TeamID)
		if err != nil {
			return err
		}

		today := truncateDay(t)
		last := today.AddDate(0, 0, autoDays)
		booked := map[int]bool{}
		var window []mzhttpclient.GameDay
		for _, d := range sched.Days {
			if d.State == mzhttpclient.DayBooked || (!d.Date.IsZero() && withMatch[d.Date]) {
				booked[d.Number] = true
			}
			if !d.Date.IsZero() && d.Date.After(today) && !d.Date.After(last) {
				window = append(window, d)
			}
		}
		label := func(n int) string {
			if d, ok := sched.Day(n); ok && !d.Date.IsZero() {
				return fmt.Sprintf("%s (day %d)", d.Date.Format("Mon 2006-01-02"), n)
			}
			return fmt.Sprintf("day %d", n)
		}
		tacticFor := func(n int) string {
			d, _ := sched.Day(n)
			return tactics.forDay(d.Date.Weekday())
		}
		// did("accepted", "accept") reads "would accept" in a dry run.
		did := func(past, base string) string {
			if autoDryRun {
				return "would " + base
			}
			return past
		}

		// 1. Incoming invites.
		fmt.Printf("incoming invites: %d\n", len(sched.Incoming))
		for _, inv := range sched.Incoming {
			who := fmt.Sprintf("  %s for %s", inv.TeamID, label(inv.Day))
			if booked[inv.Day] {
				fmt.Printf("%s: day already booked, skipping\n", who)
				continue
			}
			tac := tacticFor(inv.Day)
			if !autoDryRun {
				if err := c.AcceptInvite(inv, tac); err != nil {
					fmt.Printf("%s: accept failed: %v\n", who, firstLine(err))
					continue
				}
			}
			fmt.Printf("%s: %s (tactic %s)\n", who, did("accepted", "accept"), tac)
			booked[inv.Day] = true
		}

		// 2 + 3. Quick challenge for the next game day.
		optedIn := map[int]bool{}
		qc, err := c.GetQuickChallenge()
		if err != nil {
			return err
		}
		fmt.Printf("quick challenge for %s: ", label(qc.Day))
		switch {
		case booked[qc.Day] || qc.State == "booked":
			fmt.Println("day already booked")
			booked[qc.Day] = true
		case qc.State == "opted-in":
			fmt.Println("already opted in")
			optedIn[qc.Day] = true
		default:
			tac := tacticFor(qc.Day)
			fmt.Printf("%d offers\n", len(qc.Offers))
			if qc.OfferRows > len(qc.Offers) {
				fmt.Printf("  warning: %d teams listed but only %d book links recognised\n", qc.OfferRows, len(qc.Offers))
			}
			for _, o := range qc.Offers {
				if !autoDryRun {
					if err := c.BookQuickChallenge(o, tac, tac); err != nil {
						fmt.Printf("  %s: accept failed: %v\n", o.TeamID, firstLine(err))
						continue
					}
				}
				fmt.Printf("  %s: %s (tactic %s)\n", o.TeamID, did("accepted", "accept"), tac)
				booked[qc.Day] = true
				break
			}
			if booked[qc.Day] {
				break
			}
			if !autoDryRun {
				if err := c.OptInQuickChallenge(qc, tac, tac); err != nil {
					fmt.Printf("  opt in failed: %v\n", firstLine(err))
					break
				}
			}
			fmt.Printf("  %s (tactic %s)\n", did("opted in", "opt in"), tac)
			optedIn[qc.Day] = true
		}

		// 4. Invites for the remaining free days.
		tids := uniqueTeams(qc.Teams, me.TeamID)
		fmt.Printf("%d candidate teams: %s\n", len(tids), strings.Join(tids, ", "))
		invited := map[int]map[string]bool{}
		for _, o := range sched.Outgoing {
			if invited[o.Day] == nil {
				invited[o.Day] = map[string]bool{}
			}
			invited[o.Day][o.TeamID] = true
		}
		sent, failed := 0, 0
		for _, d := range window {
			l := label(d.Number)
			switch {
			case booked[d.Number]:
				fmt.Printf("%s: booked\n", l)
				continue
			case optedIn[d.Number]:
				fmt.Printf("%s: opted in to the quick challenge, not inviting\n", l)
				continue
			case d.State == mzhttpclient.DayNoMatch:
				continue // league day
			case d.State == mzhttpclient.DayTooFar:
				fmt.Printf("%s: too far ahead to book yet\n", l)
				continue
			case d.State != mzhttpclient.DayOpen:
				fmt.Printf("%s: %s, skipping\n", l, d.State)
				continue
			}
			var todo []string
			for _, tid := range tids {
				if !invited[d.Number][tid] {
					todo = append(todo, tid)
				}
			}
			tactic := tactics.forDay(d.Date.Weekday())
			already := ""
			if n := len(tids) - len(todo); n > 0 {
				already = fmt.Sprintf(", %d already invited", n)
			}
			fmt.Printf("%s: %s %d teams (tactic %s, %02d:00 home / %02d:00 away%s)\n",
				l, did("inviting", "invite"), len(todo), tactic, autoHomeHour, autoAwayHour, already)
			if autoDryRun {
				continue
			}
			for _, tid := range todo {
				err := c.SendFriendlyRequest(tid, at(d.Date, autoHomeHour), tactic, at(d.Date, autoAwayHour), tactic)
				if err != nil {
					failed++
					fmt.Printf("  %s: %v\n", tid, firstLine(err))
				} else {
					sent++
					fmt.Printf("  %s: sent\n", tid)
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
