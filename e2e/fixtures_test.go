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
func ForumIndex() Routes {
	return Route("GET", "/?p=forum&sub=forums&sport=soccer", File("forum-index.html"))
}

// ForumTopics: page 1 (8 pinned topics on top) and page 2 of forum 10, 135 pages in all.
func ForumTopics() Routes {
	return join(
		Route("GET", topicsURL+"0", File("forum-topics-10.html")),
		Route("GET", topicsURL+"50", File("forum-topics-10-offset-50.html")),
	)
}

// ForumTopicPinned: the last two of six pages of pinned topic 13032964.
// Every page repeats the opening post on top.
func ForumTopicPinned() Routes {
	return join(
		Route("GET", topicURL+"13032964&offset=200", File("forum-topic-13032964-offset-200.html")),
		Route("GET", topicURL+"13032964&offset=250", File("forum-topic-13032964-offset-250.html")),
	)
}

// ForumTopicQuotes: a one-page topic of 5 posts; post 46574097 quotes 46574029.
func ForumTopicQuotes() Routes {
	return Route("GET", topicURL+"13143067&offset=0", File("forum-topic-13143067.html"))
}

const (
	topicsURL = "/?p=forum&sub=topics&forum_id=10&sport=soccer&offset="
	topicURL  = "/?p=forum&sub=topic&forum_id=10&sport=soccer&topic_id="
)
