package e2e

import (
	"net/url"
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

// friendly auto runs on Saturday 2026-11-21 (game day 54, booked). Bookable:
// Mon 23 (day 56, the quick-challenge day), Tue 24, Thu 26, Fri 27; Sun and
// Wed are league days, Sat 28 is past the 7-day limit.
const autoNow = "2026-11-21T12:00:00+02:00"

const (
	tue24 = "1795474800"
	thu26 = "1795647600"
	fri27 = "1795734000"
)

func autoMZ(t *testing.T, routes ...Routes) *MZ {
	mz := Start(t, append([]Routes{Login(), UpcomingMatchesSat()}, routes...)...)
	mz.Env["MZ_NOW"] = autoNow
	return mz
}

func TestFriendlyAutoDryRun(t *testing.T) {
	mz := autoMZ(t, Schedule("challenges-2026-11-21.html"), QuickChallengeOffer())

	res := mz.Run("friendly", "auto", "--dry-run")

	res.WantExit(t, 0)
	res.WantStdout(t, `incoming invites: 0
quick challenge for Mon 2026-11-23 (day 56): 1 offers
  429231: would accept (tactic b)
11 candidate teams: 429231, 343098, 143217, 378645, 868864, 114651, 91096, 254695, 1010152, 361140, 378314
Mon 2026-11-23 (day 56): booked
Tue 2026-11-24 (day 57): would invite 11 teams (tactic a, 12:00 home / 20:00 away)
Thu 2026-11-26 (day 59): would invite 11 teams (tactic b, 12:00 home / 20:00 away)
Fri 2026-11-27 (day 60): would invite 11 teams (tactic b, 12:00 home / 20:00 away)
Sat 2026-11-28 (day 61): too far ahead to book yet
dry run: nothing sent
`)
	if n := len(mz.Find("GET", qcAccept)) + len(mz.Find("POST", sendInvite)); n != 0 {
		t.Errorf("dry run sent %d requests that change something", n)
	}
}

func TestFriendlyAutoAcceptsQuickChallenge(t *testing.T) {
	mz := autoMZ(t, Schedule("challenges-2026-11-21.html"), QuickChallengeOffer(),
		Route("GET", qcAccept, JSON(`{"response":"success","day":56}`)), SendInvite(inviteOK))

	res := mz.Run("friendly", "auto", "--tactics", "default=a")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "friendly-auto-qc.txt")
	acc := mz.Find("GET", qcAccept)
	if len(acc) != 1 || !strings.Contains(acc[0].URL, "tid=429231") {
		t.Fatalf("quick-challenge accepts: %v", acc)
	}
	wantQuery(t, acc[0], map[string]string{"home_tacticset": "a", "away_tacticset": "a"})
	if n := len(mz.Find("GET", qcOptIn)); n != 0 {
		t.Errorf("opted in although a quick challenge was accepted")
	}
	wantInviteDays(t, mz, tue24, thu26, fri27)
}

func TestFriendlyAutoOptsInWithoutOffers(t *testing.T) {
	mz := autoMZ(t, Schedule("challenges-2026-11-21.html"), QuickChallengeNone(),
		Route("GET", qcOptIn, JSON(`{"success":true}`)), SendInvite(inviteOK))

	res := mz.Run("friendly", "auto")

	res.WantExit(t, 0)
	res.WantStdoutContains(t, "quick challenge for Mon 2026-11-23 (day 56): 0 offers\n  opted in (tactic b)\n")
	res.WantStdoutContains(t, "Mon 2026-11-23 (day 56): opted in to the quick challenge, not inviting\n")
	opt := mz.Find("GET", qcOptIn)
	if len(opt) != 1 {
		t.Fatalf("%d opt-ins, want 1", len(opt))
	}
	wantQuery(t, opt[0], map[string]string{"date": "1795428376", "tactic_home": "b", "tactic_away": "b"})
	wantInviteDays(t, mz, tue24, thu26, fri27)
}

func TestFriendlyAutoOptInFailsStillInvitesOtherDays(t *testing.T) {
	mz := autoMZ(t, Schedule("challenges-2026-11-21.html"), QuickChallengeOffer(),
		Route("GET", qcAccept, JSON(`{"response":"error","response_content":"<p>Already booked</p>"}`)),
		Route("GET", qcOptIn, JSON(`{"success":false,"error_content":"<div>Too late</div>"}`)),
		SendInvite(inviteOK))

	res := mz.Run("friendly", "auto")

	res.WantExit(t, 0)
	res.WantStdoutContains(t, "  429231: accept failed: ", "Already booked", "  opt in failed: ", "Too late")
	// The quick-challenge day fell through, so it gets invites like the others.
	wantInviteDays(t, mz, mon23, tue24, thu26, fri27)
}

func TestFriendlyAutoAcceptsIncomingInvite(t *testing.T) {
	mz := autoMZ(t, Schedule("challenges-incoming.html"), QuickChallengeNone(),
		Route("GET", inviteAccept, JSON(`{"response":"success"}`)),
		Route("GET", qcOptIn, JSON(`{"success":true}`)), SendInvite(inviteOK))

	res := mz.Run("friendly", "auto", "--tactics", "thu=c,default=b")

	res.WantExit(t, 0)
	res.WantStdoutContains(t, "incoming invites: 1\n  849080 for Thu 2026-11-26 (day 59): accepted (tactic c)\n",
		"Thu 2026-11-26 (day 59): booked\n")
	acc := mz.Find("GET", inviteAccept)
	if len(acc) != 1 {
		t.Fatalf("%d invite accepts, want 1", len(acc))
	}
	wantQuery(t, acc[0], map[string]string{"cid": "228600001", "tacticset": "c", "tacticset2": "c"})
	wantInviteDays(t, mz, tue24, fri27)
}

func TestFriendlyAutoSkipsTeamsAlreadyInvited(t *testing.T) {
	// Outgoing: 10 teams each for Mon, Tue and Thu, 6 of them on today's
	// quick-challenge page; the other 4 there are new.
	mz := autoMZ(t, Schedule("challenges-outgoing.html"), QuickChallengeNone(),
		Route("GET", qcOptIn, JSON(`{"success":true}`)), SendInvite(inviteOK))

	res := mz.Run("friendly", "auto")

	res.WantExit(t, 0)
	res.WantStdoutContains(t,
		"Tue 2026-11-24 (day 57): inviting 4 teams (tactic a, 12:00 home / 20:00 away, 6 already invited)\n",
		"Fri 2026-11-27 (day 60): inviting 10 teams (tactic b, 12:00 home / 20:00 away)\n")
	wantInviteDays(t, mz, tue24, thu26, fri27)
	tue := 0
	for _, r := range mz.Find("POST", sendInvite) {
		if strings.HasPrefix(r.Form.Get("date_home"), tue24+",") {
			tue++
			if tid := tidOf(r); tid == "361140" || tid == "8895" {
				t.Errorf("Tuesday invite to %s, which was already invited", tid)
			}
		}
	}
	if tue != 4 {
		t.Errorf("%d Tuesday invites, want 4", tue)
	}
}

// wantInviteDays checks that invites went out on exactly these days (Stockholm
// midnights) with the server hours for 12:00 and 20:00 Tallinn time.
func wantInviteDays(t *testing.T, mz *MZ, days ...string) {
	t.Helper()
	got := map[string]bool{}
	for _, r := range mz.Find("POST", sendInvite) {
		day, _, _ := strings.Cut(r.Form.Get("date_home"), ",")
		got[day] = true
		wantForm(t, r, map[string]string{"date_home": day + ",11", "date_away": day + ",19"})
	}
	want := map[string]bool{}
	for _, d := range days {
		want[d] = true
		if !got[d] {
			t.Errorf("no invites for day %s", d)
		}
	}
	for d := range got {
		if !want[d] {
			t.Errorf("unexpected invites for day %s", d)
		}
	}
}

func wantQuery(t *testing.T, r Request, want map[string]string) {
	t.Helper()
	u, err := url.Parse(r.URL)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range want {
		if got := u.Query().Get(k); got != v {
			t.Errorf("%s: %s = %q, want %q", r.URL, k, got, v)
		}
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
