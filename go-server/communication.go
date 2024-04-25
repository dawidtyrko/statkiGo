package goserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)
type FireResponse struct{
	Result string `json:"result"`
	Message string `json:"message"`
}
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
	OpponentShots []string `json:"opp_shots"`
}

type BoardResponse struct {
	Board []string `json:"board"`
}

type PostFire struct {
	Coord string `json:"coord"`
}

var url = "https://go-pjatk-server.fly.dev/api/game"
var responseToken string
var boardUrl = "https://go-pjatk-server.fly.dev/api/game/board"
var fireUrl = "https://go-pjatk-server.fly.dev/api/game/fire"

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

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	var boardRes BoardResponse

	err = json.NewDecoder(resp.Body).Decode(&boardRes)
	if err != nil {
		return nil, err
	}

	return boardRes.Board, nil
}

func Fire(input string) (string,error){
	postFire := PostFire{
		Coord: input,
	}
	requestBody,err := json.Marshal(postFire)
	if err != nil{
		return "",err
	}

	req,err := http.NewRequest(http.MethodPost,fireUrl,bytes.NewBuffer(requestBody))
	if err != nil{
		return "",err
	}
	//fmt.Print(responseToken)
	req.Header.Set("X-Auth-Token",responseToken)
	req.Header.Set("Content-type","application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		return "",fmt.Errorf("client do error: %w",err)
	}
	defer resp.Body.Close()

	


	
	body,err := io.ReadAll(resp.Body)
	if err != nil{
		return "",fmt.Errorf("Reading body error: %w",err)
	}

	//fmt.Printf("body: %s\n",string(body))

	var responseFire FireResponse
	
	err = json.Unmarshal(body, &responseFire)
	if err != nil {
		return "", fmt.Errorf("Decoder error: %w",err)
	}
	if resp.StatusCode != 200{
		return responseFire.Message,fmt.Errorf("Response status: %s",resp.Status)
	}
	return responseFire.Result, nil
}
