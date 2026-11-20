package mzhttpclient

import "time"

type Match struct {
	Id           string
	HomeTeamId   string
	HomeTeamName string
	HomeGoals    string
	AwayTeamId   string
	AwayTeamName string
	AwayGoals    string
	Date         time.Time // MZ's displayed time, interpreted in time.Local
	MatchType    string    // friendly, league, cup_group, ...
	Status       string    // scheduled, playing, played
}

const (
	Friendly string = "friendly"
	League   string = "league"
)

// MatchStatus selects which list team_matchlist.php returns.
type MatchStatus int

const (
	MatchesPlayed   MatchStatus = 1
	MatchesUpcoming MatchStatus = 2
)
