package mainlogic

import (
	"fmt"
	goserver "statkiGo/go-server"
	"time"

	"github.com/fatih/color"
	gui "github.com/grupawp/warships-lightgui/v2"
)

func Logic() {
	DisplayStatus()
	stat := DisplayWaitingStatus()

	//fmt.Println(stat)

	if stat == "ready" {
		board, err := goserver.Board()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(board)
			fmt.Println("##########################")
			WarshipsGui()
		}

	} else {
		fmt.Println("Game could not start, waited for too long")
		return
	}

}

func DisplayStatus() {
	response, err := goserver.InitGame()
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(response)
		fmt.Printf("Nick: %s\nGame_status: %s\nLast_game_status: %s\nOpponent: %s\nShould_fire: %t\nTimer: %d\n",
			response.Nick, response.GameStatus, response.LastGameStatus, response.Opponent, response.ShouldFire, response.Timer)
	}

}

func DisplayWaitingStatus() string {
	var initialStatus string

	for i := 1; i <= 50; i++ {
		if i >= 50 {
			initialStatus = "not_ready"
			break
		}
		res, err := goserver.GetGameStatus()
		if err != nil {
			fmt.Println(err)
			initialStatus = "not_ready"
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

func WarshipsGui() {
	GuiSetup()
}

func GuiSetup() {
	coords, errCoords := goserver.Board()

	cfg := gui.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgYellow

	board := gui.New(cfg)
	//coords, err := goserver.Board()
	if errCoords != nil {
		fmt.Println(errCoords)
	}
	board.Import(coords)
	board.Display()

}
