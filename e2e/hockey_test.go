package e2e

import (
	"slices"
	"strings"
	"testing"
)

// Teams on the hockey quick-challenge page, in page order, without duplicates.
var hockeyCandidates = []string{"52998", "41639", "161069", "44180", "56078", "38304", "4469", "161897", "110608", "109975", "198753", "112299"}

func hockeyMZ(t *testing.T, routes ...Routes) *MZ {
	mz := Start(t, append([]Routes{HockeyAuto()}, routes...)...)
	mz.Env["MZ_NOW"] = autoNow
	return mz
}

func TestHockeyFriendlyAutoDryRun(t *testing.T) {
	mz := hockeyMZ(t)

	res := mz.Run("--sport", "hockey", "friendly", "auto", "--dry-run")

	res.WantExit(t, 0)
	res.WantStdout(t, `incoming invites: 0
quick challenge for Mon 2026-11-23 (day 56): 2 offers
  52998: would accept (tactic b)
12 candidate teams: `+strings.Join(hockeyCandidates, ", ")+`
Mon 2026-11-23 (day 56): booked
Tue 2026-11-24 (day 57): would invite 12 teams (tactic a, 12:00 home / 20:00 away)
Thu 2026-11-26 (day 59): would invite 12 teams (tactic b, 12:00 home / 20:00 away)
Fri 2026-11-27 (day 60): would invite 12 teams (tactic b, 12:00 home / 20:00 away)
Sat 2026-11-28 (day 61): too far ahead to book yet
dry run: nothing sent
`)
	login := mz.Find("POST", "/?p=login")[0].Form
	if got := login.Get("logindata[sport]"); got != "hockey" {
		t.Errorf("logindata[sport] = %q, want hockey", got)
	}
}

func TestHockeyFriendlyAuto(t *testing.T) {
	mz := hockeyMZ(t,
		Route("GET", hockeyQCAccept, inviteOK),
		Route("POST", hockeySendInvite, inviteOK),
	)
	mz.Env["MZ_SPORT"] = "hockey"

	res := mz.Run("friendly", "auto")

	res.WantExit(t, 0)
	res.WantStdoutContains(t, "  52998: accepted (tactic b)\n")
	res.WantStdoutContains(t, "sent 36, failed 0\n")
	acc := mz.Find("GET", hockeyQCAccept)
	if len(acc) != 1 || !strings.Contains(acc[0].URL, "tid=52998") {
		t.Fatalf("quick challenge accepts: %v, want one for 52998", acc)
	}
	sent := mz.Find("POST", hockeySendInvite)
	if len(sent) != 36 {
		t.Fatalf("%d invites sent, want 36 (12 teams x 3 days)", len(sent))
	}
	// Every request that names a sport names hockey.
	for _, r := range mz.Requests() {
		if strings.Contains(r.URL, "sport=") && !strings.Contains(r.URL, "sport=hockey") {
			t.Errorf("football request in hockey mode: %s", r)
		}
	}
}

func TestHockeyFriendlyCandidates(t *testing.T) {
	mz := hockeyMZ(t)
	mz.Env["MZ_SPORT"] = "hockey"

	res := mz.Run("friendly", "candidates")

	res.WantExit(t, 0)
	res.WantStdout(t, strings.Join(hockeyCandidates, "\n")+"\n")
}

func TestSportFlagOverridesEnv(t *testing.T) {
	mz := Start(t, Login(), QuickChallenge())
	mz.Env["MZ_SPORT"] = "hockey"

	res := mz.Run("--sport", "football", "friendly", "candidates")

	res.WantExit(t, 0)
	res.WantStdout(t, strings.Join(candidates, "\n")+"\n")
}

func TestUnknownSport(t *testing.T) {
	mz := Start(t)

	res := mz.Run("--sport", "curling", "whoami")

	res.WantExit(t, 1)
	res.WantStderrContains(t, `unknown sport "curling"`)
	mz.WantRequests()
}

// hockeyForumSynced syncs hockey forum 140 into a fresh database, next to
// football forum 10 synced as in syncedForum.
func hockeyForumSynced(t *testing.T) *MZ {
	t.Helper()
	mz := Start(t, HockeyForum(), ForumIndex(), ForumOneTopic(), ForumTopicQuotes())
	res := mz.Run("forum", "sync", "--forum", "10", "--delay", "0")
	res.WantExit(t, 0)
	mz.Env["MZ_NOW"] = hockeyForumNow
	res = mz.Run("--sport", "hockey", "forum", "sync", "--forum", "140", "--delay", "0")
	res.WantExit(t, 0)
	res.WantStderrContains(t, "forum 140 (English » Questions & Answers): 1 topics to fetch", "1 topics, 1 posts")
	return mz
}

func TestHockeyForumSync(t *testing.T) {
	mz := hockeyForumSynced(t)

	res := mz.Run("forum", "topic", "13147751")
	res.WantExit(t, 0)
	got := decode[struct {
		ForumID int
		Forum   string
		Posts   []forumPost
	}](t, res)
	if got.ForumID != 140 || got.Forum != "English » Questions & Answers" || len(got.Posts) != 1 {
		t.Fatalf("topic 13147751: %+v", got)
	}
	if p := got.Posts[0]; p.ID != 46585059 || p.Author.Name != "barca87" ||
		!p.PostedAt.Equal(tallinn("2026-11-19 16:23")) || // "Thursday 16:23"
		!strings.HasPrefix(p.Text, "When will the points for the 2026 World Cup prediction event") {
		t.Errorf("post: %+v", p)
	}
}

// list and search show only the sport's forums.
func TestHockeyForumListAndSearch(t *testing.T) {
	mz := hockeyForumSynced(t)

	forumIDs := func(args ...string) []int {
		t.Helper()
		res := mz.Run(append(args, "forum", "list")...)
		res.WantExit(t, 0)
		var ids []int
		for _, f := range decode[[]struct{ ID int }](t, res) {
			ids = append(ids, f.ID)
		}
		return ids
	}
	if ids := forumIDs("--sport", "hockey"); len(ids) != 7 || ids[0] != 136 || ids[6] != 142 {
		t.Errorf("hockey forums: %v, want 136..142", ids)
	}
	if ids := forumIDs(); len(ids) != 9 || slices.Contains(ids, 140) {
		t.Errorf("football forums: %v, want the 9 football ones", ids)
	}

	hits := func(sport, query string) int {
		t.Helper()
		res := mz.Run("--sport", sport, "forum", "search", query)
		res.WantExit(t, 0)
		return len(decode[[]struct{ PostID int }](t, res))
	}
	if n := hits("hockey", "prediction"); n != 1 {
		t.Errorf("hockey search 'prediction': %d hits, want 1", n)
	}
	if n := hits("football", "prediction"); n != 0 {
		t.Errorf("football search 'prediction': %d hits, want 0", n)
	}
	if n := hits("hockey", "rubbish youth"); n != 0 {
		t.Errorf("hockey search found a football post: %d hits", n)
	}
}

// Without --forum, a hockey sync reads hockey's talk (138) and Q&A (140).
// Forum 138 is served the Q&A page; only the requests matter here.
func TestHockeyForumSyncDefaultForums(t *testing.T) {
	mz := Start(t, HockeyForum(),
		Route("GET", "/?p=forum&sub=topics&forum_id=138&sport=hockey&offset=0", File("hockey/forum-topics-140-one.html")))
	mz.Env["MZ_NOW"] = hockeyForumNow
	mz.Env["MZ_SPORT"] = "hockey"

	res := mz.Run("forum", "sync", "--delay", "0")

	res.WantExit(t, 0)
	var lists []string
	for _, r := range mz.Requests() {
		if strings.Contains(r.URL, "sub=topics") {
			lists = append(lists, r.URL)
		}
		if strings.Contains(r.URL, "sport=soccer") {
			t.Errorf("football request in hockey mode: %s", r)
		}
	}
	if len(lists) != 2 || !strings.Contains(lists[0], "forum_id=138") || !strings.Contains(lists[1], "forum_id=140") {
		t.Errorf("topic list requests: %q, want forum 138 then 140", lists)
	}
}
