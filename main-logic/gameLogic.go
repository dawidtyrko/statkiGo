package mainlogic

import (
	"fmt"
	goserver "statkiGo/go-server"
	"time"
)

func Logic() {
	DisplayStatus()
	stat := DisplayWaitingStatus()

	if stat == "not_ready" {
		fmt.Println("Game could not start, waited for too long")
		return
	}

}

func DisplayStatus() {
	response, err := goserver.InitGame()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Nick: %s\nGame_status: %s\nLast_game_status: %s\nOpponent: %s\nShould_fire: %t\nTimer: %d\n",
			response.Nick, response.GameStatus, response.LastGameStatus, response.Opponent, response.ShouldFire, response.Timer)
	}

}

func DisplayWaitingStatus() string {
	var initialStatus string

	for i := 1; i <= 50; i++ {
		if i == 50 {
			initialStatus = "not_ready"
		}
		res, err := goserver.GetGameStatus()
		if err != nil {
			fmt.Println(err)
			break
		}

		if res.GameStatus == "game_in_progress" {
			fmt.Printf("game_status: %s\n", res.GameStatus)
			initialStatus = "ready"
			break
		} else if res.GameStatus == "waiting_wpbot" {
			fmt.Printf("game_status: %s\n", res.GameStatus)
			time.Sleep(1 * time.Second)
			continue
		}

	}
	return initialStatus
}
