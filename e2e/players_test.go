package e2e

import (
	"encoding/json"
	"testing"
)

type player struct {
	Name           string
	Age            int
	Value          int
	ValueCurrency  string
	Salary         int
	SalaryCurrency string

	Speed, Stamina, Intelligence, Passing, Shooting, Heading, Keeping,
	BallControl, Tackling, AerialPassing, SetPlays, Experience, Form float64
}

func (p player) skills() map[string]float64 {
	return map[string]float64{
		"Speed": p.Speed, "Stamina": p.Stamina, "Intelligence": p.Intelligence, "Passing": p.Passing,
		"Shooting": p.Shooting, "Heading": p.Heading, "Keeping": p.Keeping, "BallControl": p.BallControl,
		"Tackling": p.Tackling, "AerialPassing": p.AerialPassing, "SetPlays": p.SetPlays,
		"Experience": p.Experience, "Form": p.Form,
	}
}

func TestPlayersTable(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("players")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "players.txt")
}

func TestPlayersJSON(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("players", "--json")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "players.json")

	var players []player
	if err := json.Unmarshal([]byte(res.Stdout), &players); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if len(players) != 46 {
		t.Errorf("%d players, the page has 46", len(players))
	}
	for _, p := range players {
		for name, v := range p.skills() {
			if v < 0 || v > 10 {
				t.Errorf("%s: %s = %v, want 0..10", p.Name, name, v)
			}
		}
		// Youth players (16-18 in this squad) have no salary; everyone else does.
		if youth := p.Age <= 18; youth != (p.Salary == 0) {
			t.Errorf("%s (age %d): salary %d", p.Name, p.Age, p.Salary)
		}
	}
}

// A youth player later on the page must not erase the salary of the players
// before them (regression: all 46 salaries were 0).
func TestPlayersSalaryOfFirstPlayer(t *testing.T) {
	mz := Start(t, Login())

	var players []player
	if err := json.Unmarshal([]byte(mz.Run("players", "--json").Stdout), &players); err != nil || len(players) == 0 {
		t.Fatalf("no players: %v", err)
	}
	p := players[0]
	if p.Name != "Aurelio Theodore" || p.Salary != 8728 || p.SalaryCurrency != "EUR" || p.Value != 621210 {
		t.Errorf("first player: %+v", p)
	}
}

func TestPlayersLayoutChanged(t *testing.T) {
	// The players container is there, but the player boxes look different.
	mz := Start(t, Login(), Route("GET", "/?p=players", File("players-new-layout.html")))

	res := mz.Run("players")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "players page")
	res.WantStderrLacks(t, "goroutine", "panic:") // an error, not a crash
}
