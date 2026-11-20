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
