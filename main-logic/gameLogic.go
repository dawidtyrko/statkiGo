package mainlogic

import (
	"bufio"
	"fmt"
	"os"
	goserver "statkiGo/go-server"
	"time"

	"github.com/fatih/color"
	gui "github.com/grupawp/warships-lightgui/v2"
)

type UserPrompt struct {
	Nick string
}

func Logic() {
	user, err := Prompt()
	if err != nil {
		fmt.Print(err)
		return
	}
	DisplayInitialStatus(user.Nick)
	stat := DisplayWaitingStatus()

	var oppShots []string

	if stat == "ready" {
		config := GuiSetup()
		board := gui.New(config)
		boardCoords, err := goserver.Board()
		if err != nil {
			fmt.Println(err)
			return
		} else {
			err := board.Import(boardCoords)
			if err != nil {
				fmt.Println("Problem with the coords import")
				return
			}
			description, err := goserver.GetDescription()
			if err != nil {
				fmt.Println("Problem with importing description")
				return
			}

			//exported := board.Export(gui.Left)
			//fmt.Printf("Exported: %s\n", exported)
			var i = 0
			for {
				board.Display()
				fmt.Printf("Nick: %s\nDescription: %s\nOpponent: %s\nOpponent description: %s\n",
					description.Nick, description.Desc, description.Opponent, description.OpponentDescription)
				//co 60 sekund oddzielic do funkcji, odpalac w tle
				req, err := goserver.GetGameStatus()
				if err != nil {
					fmt.Println(err)
					time.Sleep(1 * time.Second)

				}

				if len(req.OpponentShots) == 0 {
					continue
					//oppShots = append(oppShots, req.OpponentShots[0])
				} else {
					oppShots = append(oppShots, req.OpponentShots[len(req.OpponentShots)-1])
				}
				// if req.GameStatus != "game_in_progress" {
				// 	//wyswietlic zwyciezce i jakies dane
				// 	break
				// }
				if req.ShouldFire {
					output, ok := gui.ReadLineWithTimer("Enter coords: ", time.Minute)

					if !ok {
						fmt.Println("wrong coordinates")
						break
					}
					//fmt.Println(output)
					fireResponse, err := goserver.Fire(output)
					if err != nil {
						fmt.Println(fireResponse)
						fmt.Println(err)
						break
					}

					state, err := board.HitOrMiss(gui.Right, output)
					if err != nil {
						fmt.Printf("error HitOrMissLeft: %v", err)
						break
					}

					err = board.Set(gui.Right, output, state)
					if err != nil {
						fmt.Printf("error with Set: %v", err)
						break
					}

					lastElement := oppShots[len(oppShots)-1]

					state2, err := board.HitOrMiss(gui.Left, lastElement)
					if err != nil {
						fmt.Printf("error HitOrMissRight: %v", err)
						break
					}
					err = board.Set(gui.Left, lastElement, state2)
					if err != nil {
						fmt.Printf("error with Set: %v", err)
						break
					}

				}
				//DO ZROBIENIA
				//board.Export(gui.Left)
				//board.Export(gui.Right)
				//fmt.Print("\033[H\033[2J")
				i++
			}

		}

	} else {
		fmt.Println("Game could not start, waited for too long")
		return
	}

}

func DisplayInitialStatus(name string) {
	response, err := goserver.InitGame(name)
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

func GuiSetup() *gui.Config {

	cfg := gui.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgYellow
	cfg.MissChar = 'X'
	cfg.MissColor = color.BgCyan

	// board.Import(coords)
	// board.Display()
	return cfg
}

func Prompt() (UserPrompt, error) {
	lobbies, err := goserver.GetLobby()
	if err != nil {
		fmt.Print(err)
		return UserPrompt{}, err
	}
	for _, lobby := range lobbies {
		fmt.Printf("User: %s, Status: %s\n", lobby.User, lobby.User)
	}

	fmt.Print("Enter your chosen enemy: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	var prompt UserPrompt
	prompt.Nick = name

	return prompt, nil
}
