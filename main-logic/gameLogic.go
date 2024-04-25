package mainlogic

import (
	"fmt"
	goserver "statkiGo/go-server"
	"time"

	"github.com/fatih/color"
	gui "github.com/grupawp/warships-lightgui/v2"
)

func Logic() {
	DisplayInitialStatus()
	stat := DisplayWaitingStatus()


	if stat == "ready" {
		config := GuiSetup()
		board := gui.New(config)
		boardCoords, err := goserver.Board()
		if err != nil {
			fmt.Println(err)
		} else {
			for {
				board.Import(boardCoords)
				board.Display()

				//co 60 sekund oddzielic do funkcji, odpalac w tle
				req, err := goserver.GetGameStatus()
				if err != nil {
					fmt.Println(err)
				}
				if req.GameStatus != "game_in_progress" {
					//wyswietlic zwyciezce i jakies dane
					break
				}
				if req.ShouldFire {
					output, ok := gui.ReadLineWithTimer("Enter coords: ",time.Minute)
					
					if !ok {
						break
					}
					fmt.Println(output)
					fireResponse, err := goserver.Fire(output)
					if err != nil {
						fmt.Println(fireResponse)
						fmt.Println(err)
						break
					}
					
					if fireResponse == "miss"{
						fmt.Println(fireResponse)
						break
					}
				}
			}

		}

	} else {
		fmt.Println("Game could not start, waited for too long")
		return
	}

}

func UserInput() (string, error) {
	var input string
	fmt.Print("Enter the field (A8): ")

	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", err
	}
	return input, nil
}

func DisplayInitialStatus() {
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

// func WarshipsGui() {
// 	GuiSetup()
// }

func GuiSetup() *gui.Config{
	
	

	cfg := gui.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgYellow

	
	
	// board.Import(coords)
	// board.Display()
	return cfg
}
