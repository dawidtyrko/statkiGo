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
	Gamemode string
	Token string
}


func Logic() {
	user, err := Prompt()
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(user.Nick)
	
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

			
			var i = 0
			for {
				board.Display()
				fmt.Printf("Nick: %s\nDescription: %s\nOpponent: %s\nOpponent description: %s\n",
					description.Nick, description.Desc, description.Opponent, description.OpponentDescription)
				
				time.Sleep(1 * time.Second)

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


func DisplayWaitingStatus() string {
	var initialStatus string

	for i := 1; i <= 100; i++ {
		if i >= 100 {
			initialStatus = "not_ready"
			break
		}
		res, err := goserver.GetGameStatus()
		//time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println(err)
			initialStatus = "not_ready"
			break
		}
		//fmt.Println(res.GameStatus)
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


	scanner := bufio.NewScanner(os.Stdin)
	var prompt UserPrompt
	var name string
	
	fmt.Print("Choose mode: multi | single\n")
	scanner.Scan()
	mode := scanner.Text()

	if mode == "single"{
		resToken,err := goserver.GameInitialization("","single")
		if err != nil{
			fmt.Printf("Error game initialization: %v",err)
			return UserPrompt{}, err
		}
		prompt.Token = resToken
		prompt.Gamemode = mode
		prompt.Nick = ""

	}else if mode == "multi"{

		resToken,err := goserver.GameInitialization("","multi")
		if err != nil{
			fmt.Printf("Error game initialization: %v",err)
			return UserPrompt{}, err
		}
		prompt.Token = resToken
		prompt.Gamemode = mode

			for {
				lobbies, err := goserver.GetLobby()
				if err != nil {
					fmt.Print(err)
					break
				}
				if len(lobbies) == 0{
					fmt.Println("Lobby is empty")


				}else{
					for _, lobby := range lobbies {
						fmt.Printf("User: %s, Status: %s\n", lobby.User, lobby.Status)
					}
					fmt.Print("Enter your chosen enemy: ")
					scanner.Scan()
					name = scanner.Text()

					found := false
					for _, lobby := range lobbies {
						if lobby.User == name {
							found = true
							break
						}
					}
					if !found {
						fmt.Printf("%s is not in the lobby.\n", name)
						continue 
					}else{
						fmt.Printf("%s is your chosen enemy.\n",name)
						break
					}
				}
				
				time.Sleep(2 * time.Second) 
				fmt.Print("\033[H\033[2J")
			}

		prompt.Nick = name
		

		res, err := goserver.GameInitialization(name,mode)
		if err != nil{
			fmt.Println(err)
		}
		//fmt.Println(res)
		prompt.Token = res

	}else{
		return UserPrompt{},nil
	}
	return prompt, nil
	
}

func ImportExport() ([]string, error){
	config := GuiSetup()
	board := gui.New(config)
	coords := []string{
		"A1",
    "A3",
    "B9",
    "C7",
    "D1",
    "D2",
    "D3",
    "D4",
    "D7",
    "E7",
    "F1",
    "F2",
    "F3",
    "F5",
    "G5",
    "G8",
    "G9",
    "I4",
    "J4",
    "J8"}
	err := board.Import(coords)
	if err != nil{
		fmt.Println(err)
		return nil,err
	}
	exported := board.Export(gui.Left)
	fmt.Println(exported)
	return exported,nil
}
