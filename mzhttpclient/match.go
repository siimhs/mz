package mzhttpclient

import "time"

type Match struct {
	Id         string
	HomeTeamId string
	AwayTeamId string
	Date       time.Time
	MatchType  string
}

const (
	Friendly string = "friendly"
	League   string = "league"
)
