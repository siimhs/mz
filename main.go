package main

import (
	"log"
	"os"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
)

func main() {
	username := os.Getenv("MZ_USERNAME")
	password := os.Getenv("MZ_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Environment variables MZ_USERNAME and MZ_PASSWORD must be set")
	}
	mz_url := os.Getenv("MZ_URL")
	if mz_url == "" {
		mz_url = "https://www.managerzone.com/"
	}

	client := mzhttpclient.New(mz_url)

	// Login
	err := client.Login(username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// matches := client.GetMatches("364704")
	// if matches != nil {
	// 	for _, m := range matches {
	// 		fmt.Printf("Match %s\n", m.Date)
	// 	}
	// }
	// homeDate := time.Date(2025, 7, 12, 11, 0, 0, 0, time.Local)
	// awayDate := time.Date(2025, 7, 12, 19, 0, 0, 0, time.Local)
	// err = client.SendFriendlyRequest("121197", homeDate, "a", awayDate, "a")
	// if err != nil {
	// 	println("Insuccessful friendly! A")
	// } else {
	// 	println("Successful friendly! A")
	// }

	tids, err := client.GetTeamIDsOnQuickChallengePage()
	if err != nil {
		log.Fatalf("Failed to team Ids from quick challenge page: %v", err)
	}
	for _, t := range tids {
		println(t)
	}
	cur := time.Now()
	for i := 0; i < 8; i++ {
		for _, t := range tids {
			homeDate := time.Date(cur.Year(), cur.Month(), cur.AddDate(0, 0, i).Day(), 10, 0, 0, 0, time.Local)
			awayDate := time.Date(cur.Year(), cur.Month(), cur.AddDate(0, 0, i).Day(), 18, 0, 0, 0, time.Local)
			if cur.AddDate(0, 0, i).Weekday() == time.Tuesday || cur.AddDate(0, 0, i).Weekday() == time.Wednesday {
				err := client.SendFriendlyRequest(t, homeDate, "a", awayDate, "a")
				if err != nil {
					println("Insuccessful friendly! A")
				} else {
					println("Successful friendly! A")
				}
			} else {
				err := client.SendFriendlyRequest(t, homeDate, "b", awayDate, "b")
				if err != nil {
					println("Insuccessful friendly! B")
				} else {
					println("Successful friendly! B")
				}
			}

		}
	}

	// client.SendFriendlyRequest("829622", time.Date(2025, time.January, 24, 12, 0, 0, 0, time.Local), "b", time.Date(2025, time.January, 24, 20, 0, 0, 0, time.Local), "b")
	// // Get My Team
	// players, err := client.GetMyPlayers()
	// if err != nil {
	// 	log.Fatalf("Failed to get my team: %v", err)
	// }

	// println(len(players))
	// println(players[0].Name)
	// println(players[0].Age)
	// println(players[0].Foot)
	// println(players[0].BirthSeason)
	// println(players[0].Height)
	// println(players[0].Weight)
	// println(players[0].Value)
	// println(players[0].ValueCurrency)
	// println(players[0].Salary)
	// println(players[0].SalaryCurrency)
	// println(fmt.Sprintf("%.2f", players[0].Speed))
	// println(fmt.Sprintf("%.2f", players[0].Stamina))
	// println(fmt.Sprintf("%.2f", players[0].Intelligence))
	// println(fmt.Sprintf("%.2f", players[0].Passing))
	// println(fmt.Sprintf("%.2f", players[0].Shooting))
	// println(fmt.Sprintf("%.2f", players[0].Heading))
	// println(fmt.Sprintf("%.2f", players[0].Keeping))
	// println(fmt.Sprintf("%.2f", players[0].BallControl))
	// println(fmt.Sprintf("%.2f", players[0].Tackling))
	// println(fmt.Sprintf("%.2f", players[0].AerialPassing))
	// println(fmt.Sprintf("%.2f", players[0].SetPlays))
	// println(fmt.Sprintf("%.2f", players[0].Experience))
	// println(fmt.Sprintf("%.2f", players[0].Form))
}
