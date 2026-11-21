package e2e

// Route groups for the recorded pages in testdata/mz (recorded 2026-11-20, scrubbed).

// Login is the flow every command goes through: front page for the session
// cookie, login form (answered with a redirect), then the players page that
// mz reads to confirm the login.
func Login() Routes {
	return join(
		Route("GET", "/", File("index.html")),
		Route("POST", "/?p=login", Redirect("/?p=clubhouse")),
		Route("GET", "/?p=clubhouse", File("clubhouse.html")),
		Route("GET", "/?p=players", File("players.html")),
	)
}

// Upcoming matches of my team. Saturday 2026-11-21 already has friendlies.
func UpcomingMatches() Routes {
	return Route("GET", "/xml/team_matchlist.php?team_id="+myTeam+"&match_status=2", File("matches-upcoming.xml"))
}

func PlayedMatches() Routes {
	return Route("GET", "/xml/team_matchlist.php?team_id="+myTeam+"&match_status=1", File("matches-played.xml"))
}

// QuickChallenge lists 11 teams; 849080 appears twice.
func QuickChallenge() Routes {
	return Route("GET", "/ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer", File("quick-challenge.html"))
}

// Pages recorded on Saturday 2026-11-21 ~12:00 for friendly auto.

// UpcomingMatchesSat: friendlies only on Saturday 2026-11-21.
func UpcomingMatchesSat() Routes {
	return Route("GET", "/xml/team_matchlist.php?team_id="+myTeam+"&match_status=2", File("matches-upcoming-2026-11-21.xml"))
}

// Schedule serves a friendly schedule page (Challenges > Friendlies):
//   - challenges-2026-11-21.html: recorded; days 56, 57, 59, 60 open, no invites
//   - challenges-outgoing.html: recorded earlier that day; 10 outgoing invites
//     each for days 56, 57 and 59
//   - challenges-incoming.html: hand-made from the first, one incoming
//     double from 849080 for day 59 (the accept link format is a guess from
//     MZ's challenges.js; no live sample yet)
func Schedule(file string) Routes {
	return Route("GET", "/ajax.php?p=challenge&sub=personal-challenge-template&sport=soccer", File(file))
}

// QuickChallengeOffer: day 56, we can opt in, 429231 offers a quick challenge.
func QuickChallengeOffer() Routes {
	return Route("GET", "/ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer", File("quick-challenge-offer.html"))
}

// QuickChallengeNone: day 56, we can opt in, nobody else has.
func QuickChallengeNone() Routes {
	return Route("GET", "/ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer", File("quick-challenge-none.html"))
}

const (
	qcAccept     = "/ajax.php?p=challenge&sub=quick-challenge-accept&tid=*&sport=soccer&home_tacticset=*&away_tacticset=*"
	qcOptIn      = "/ajax.php?p=challenge&sub=quick-challenge-opt-in&date=*&tactic_home=*&tactic_away=*&sport=soccer"
	inviteAccept = "/ajax.php?p=challenge&sub=handle-challenge&action=accept&cid=*&sport=soccer&tacticset=*&tacticset2=*"
)

const sendInvite = "/ajax.php?p=challenge&sub=handle-challenge&action=create&tid=*&country_id=123&sport=soccer"

var (
	inviteOK       = JSON(`{"response":"success"}`)
	inviteRejected = JSON(`{"response":"error","message":"Team does not accept challenges"}`)
)

// SendInvite answers friendly invites with the given responses in order.
func SendInvite(responses ...Response) Routes { return Route("POST", sendInvite, responses...) }

func join(groups ...Routes) Routes {
	var out Routes
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

// Forum pages, recorded 2026-11-20 23:07 (a Friday; "Today", "Yesterday" and
// weekday names in them are relative to that day).

// ForumIndex is the English section's index, which forum sync reads.
func ForumIndex() Routes {
	return Route("GET", "/?p=forum&sub=forums&sport=soccer&forum_language=en", File("forum-index.html"))
}

// ForumOneTopic: forum 10's topic list cut down to one page holding only
// topic 13143067 (hand-made from forum-topics-10.html).
func ForumOneTopic() Routes {
	return Route("GET", topicsURL+"0", File("forum-topics-10-one.html"))
}

// ForumTopicQuotes: a one-page topic of 5 posts; post 46574097 quotes 46574029.
func ForumTopicQuotes() Routes {
	return Route("GET", topicURL+"13143067&offset=0", File("forum-topic-13143067.html"))
}

const (
	topicsURL = "/?p=forum&sub=topics&forum_id=10&sport=soccer&offset="
	topicURL  = "/?p=forum&sub=topic&forum_id=10&sport=soccer&topic_id="
)
