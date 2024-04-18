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

func InitGame() (string, error) {
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
		Desc:        "My first game",
		Nick:        "John_Doe",
		Target_nick: "",
		Wpbot:       true,
	}

	body, _ := json.Marshal(postRequest)
	var responseToken string
	var url = "https://go-pjatk-server.fly.dev/api/game"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "error: ", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		responseToken = resp.Header.Get("x-auth-token")

	} else {
		responseToken = fmt.Sprint("Post failed with error: " + resp.Status)
		return responseToken, err
	}

	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Set("X-Auth-Token", responseToken)
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)

	if err != nil {
		return "Get request error: ", err
	}
	defer req.Body.Close()

	var getResponseString string

	if resp2.StatusCode == 200 {
		body, err := io.ReadAll(resp2.Body)
		if err != nil {
			return "Get body read error: ", err
		}

		getResponseString = string(body)
		fmt.Println("Response: ", getResponseString, body)
	} else {
		getResponseString = fmt.Sprint("Get failed with error: " + resp2.Status)
		return getResponseString, err
	}
	return getResponseString, nil
}
