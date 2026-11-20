package e2e

import (
	"strings"
	"testing"
)

// Teams on the quick-challenge page, in page order, without duplicates.
var candidates = []string{"849080", "946359", "1042969", "998908", "898269", "460800", "387388", "934369", "415203", "446105", "141429"}

// Unix time of local midnight in Europe/Stockholm (MZ server time), as sent
// in date_home/date_away.
const (
	sat21 = "1795215600"
	sun22 = "1795302000"
	mon23 = "1795388400"
)

func TestFriendlyCandidates(t *testing.T) {
	mz := Start(t, Login(), QuickChallenge())

	res := mz.Run("friendly", "candidates")

	res.WantExit(t, 0)
	res.WantStdout(t, strings.Join(candidates, "\n")+"\n")
}

func TestFriendlySend(t *testing.T) {
	mz := Start(t, Login(), SendInvite(inviteOK))

	res := mz.Run("friendly", "send", "849080", "--date", "2026-11-22",
		"--home", "12", "--away", "20", "--tactic", "a", "--tactic-away", "b")

	res.WantExit(t, 0)
	res.WantStdout(t, "sent invite to 849080 for Sun 2026-11-22\n")
	sent := mz.Find("POST", sendInvite)
	if len(sent) != 1 {
		t.Fatalf("%d invites sent, want 1", len(sent))
	}
	if !strings.Contains(sent[0].URL, "tid=849080") {
		t.Errorf("invite to the wrong team: %s", sent[0].URL)
	}
	// Hours are MZ server time: 12:00 and 20:00 in Tallinn are 11 and 19 in Stockholm.
	wantForm(t, sent[0], map[string]string{
		"date_home":   sun22 + ",11",
		"date_away":   sun22 + ",19",
		"tactic_home": "a",
		"tactic_away": "b",
	})
}

func TestFriendlySendAwayTacticDefaultsToHome(t *testing.T) {
	mz := Start(t, Login(), SendInvite(inviteOK))

	mz.Run("friendly", "send", "849080", "--date", "2026-11-22", "--tactic", "c").WantExit(t, 0)

	for _, r := range mz.Find("POST", sendInvite) {
		wantForm(t, r, map[string]string{"tactic_home": "c", "tactic_away": "c"})
	}
}

func TestFriendlySendRejected(t *testing.T) {
	mz := Start(t, Login(), SendInvite(inviteRejected))

	res := mz.Run("friendly", "send", "849080", "--date", "2026-11-22")

	res.WantExit(t, 1)
	res.WantStdout(t, "")
	res.WantStderrContains(t, "Team does not accept challenges") // MZ's reason is shown
}

func TestFriendlySendServerError(t *testing.T) {
	mz := Start(t, Login(), SendInvite(Status(500)))

	res := mz.Run("friendly", "send", "849080", "--date", "2026-11-22")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "500")
}

func TestFriendlySendBadDate(t *testing.T) {
	mz := Start(t)

	res := mz.Run("friendly", "send", "849080", "--date", "22.11.2026")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "--date")
	mz.WantRequests() // rejected before logging in
}

func TestFriendlyAutoDryRun(t *testing.T) {
	mz := Start(t, Login(), UpcomingMatches(), QuickChallenge())

	res := mz.Run("friendly", "auto", "--dry-run")

	res.WantExit(t, 0)
	res.WantStdout(t, `11 candidate teams: `+strings.Join(candidates, ", ")+`
Sat 2026-11-21: already has a friendly, skipping
Sun 2026-11-22: inviting 11 teams (tactic b, 12:00 home / 20:00 away)
Mon 2026-11-23: inviting 11 teams (tactic b, 12:00 home / 20:00 away)
Tue 2026-11-24: inviting 11 teams (tactic a, 12:00 home / 20:00 away)
Wed 2026-11-25: inviting 11 teams (tactic a, 12:00 home / 20:00 away)
dry run: nothing sent
`)
	if sent := mz.Find("POST", sendInvite); len(sent) != 0 {
		t.Errorf("dry run sent %d invites", len(sent))
	}
}

func TestFriendlyAuto(t *testing.T) {
	// The third invite is rejected; the others succeed.
	mz := Start(t, Login(), UpcomingMatches(), QuickChallenge(),
		SendInvite(inviteOK, inviteOK, inviteRejected, inviteOK))

	res := mz.Run("friendly", "auto", "--days", "3", "--tactics", "mon=c,default=b")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "friendly-auto.txt")
	if !strings.HasSuffix(res.Stdout, "sent 21, failed 1\n") {
		t.Errorf("summary: %q", res.Stdout[strings.LastIndex(strings.TrimSuffix(res.Stdout, "\n"), "\n")+1:])
	}

	// Saturday is booked, so every candidate gets one invite for Sunday and
	// one for Monday, with that day's tactic.
	wantTactic := map[string]string{sun22: "b", mon23: "c"}
	invited := map[string]map[string]int{} // day -> team -> invites
	for _, r := range mz.Find("POST", sendInvite) {
		day, _, _ := strings.Cut(r.Form.Get("date_home"), ",")
		if day == sat21 {
			t.Errorf("invite on Saturday, which already has a friendly: %s", r.URL)
		}
		if tac, ok := wantTactic[day]; ok {
			wantForm(t, r, map[string]string{"tactic_home": tac, "tactic_away": tac})
		}
		if invited[day] == nil {
			invited[day] = map[string]int{}
		}
		invited[day][tidOf(r)]++
	}
	for day := range wantTactic {
		for _, tid := range candidates {
			if n := invited[day][tid]; n != 1 {
				t.Errorf("day %s, team %s: %d invites, want 1", day, tid, n)
			}
		}
		if len(invited[day]) != len(candidates) {
			t.Errorf("day %s: invited %d teams, want %d", day, len(invited[day]), len(candidates))
		}
	}
	if len(invited) != 2 {
		t.Errorf("invites on %d days, want 2", len(invited))
	}
}

func wantForm(t *testing.T, r Request, want map[string]string) {
	t.Helper()
	for k, v := range want {
		if got := r.Form.Get(k); got != v {
			t.Errorf("%s: %s = %q, want %q", r.URL, k, got, v)
		}
	}
}

func tidOf(r Request) string {
	_, after, _ := strings.Cut(r.URL, "tid=")
	tid, _, _ := strings.Cut(after, "&")
	return tid
}
