package goserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PostRequest struct {
	Coords      []string `json:"coords"`
	Desc        string   `json:"desc"`
	Nick        string   `json:"nick"`
	Target_nick string   `json:"target_nick"`
	Wpbot       bool     `json:"wpbot"`
}

type Response struct {
	Nick           string `json:"nick"`
	GameStatus     string `json:"game_status"`
	LastGameStatus string `json:"last_game_status"`
	Opponent       string `json:"opponent"`
	ShouldFire     bool   `json:"should_fire"`
	Timer          int    `json:"timer"`
}

type BoardResponse struct {
	Board []string `json:"board"`
}

var url = "https://go-pjatk-server.fly.dev/api/game"
var responseToken string
var boardUrl = "https://go-pjatk-server.fly.dev/api/game/board"

func InitGame() (Response, error) {

	postResponse, err := gameInitialization()
	if err != nil {
		return Response{}, err
	} else {
		responseToken = postResponse
	}

	response, err := GetGameStatus()
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func GetGameStatus() (Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Set("X-Auth-Token", responseToken)
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)

	if err != nil {
		return Response{}, err
	}
	defer req.Body.Close()

	var getResponseString string

	if resp2.StatusCode == 200 {
		body, err := io.ReadAll(resp2.Body)
		if err != nil {
			return Response{}, err
		}

		getResponseString = string(body)
	} else {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal([]byte(getResponseString), &response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func gameInitialization() (string, error) {

	postRequest := PostRequest{
		Coords: []string{
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
			"J8"},
		Desc:        "USS Missouri",
		Nick:        "William_M_Callaghan",
		Target_nick: "",
		Wpbot:       true,
	}

	body, _ := json.Marshal(postRequest)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "Error occured: ", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		responseToken = resp.Header.Get("x-auth-token")

	} else {
		return fmt.Sprintf("Error occured, status code: %d", resp.StatusCode), err
	}
	return responseToken, nil
}

func Board() ([]string, error) {
	req, _ := http.NewRequest(http.MethodGet, boardUrl, http.NoBody)
	req.Header.Set("X-Auth-Token", responseToken)
	req.Header.Set("Content-Type", "application/json")

	req2, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer req2.Body.Close()

	if req2.StatusCode != 200 {
		return nil, err
	}

	var boardRes BoardResponse

	err = json.NewDecoder(req2.Body).Decode(&boardRes)
	if err != nil {
		return nil, err
	}

	return boardRes.Board, nil
}
