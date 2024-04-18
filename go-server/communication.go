package goserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PostRequest struct {
	Coords      []string `json:"coords"`
	Desc        string   `json:"desc"`
	Nick        string   `json:"nick"`
	Target_nick string   `json:"target_nick"`
	Wpbot       bool     `json:"wpbot"`
}

func InitGame() string {
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

	var url = "https://go-pjatk-server.fly.dev/api/game"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	var jsonStr string
	if resp.StatusCode == 200 {
		//body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	panic(err)
		// }
		jsonStr = resp.Header.Get("x-auth-token")
		// jsonStr = string(body)
		//fmt.Println("Response: ", jsonStr, body)

	} else {
		return "Post failed with error"
		fmt.Println("Get failed with error: ", resp.Status)
	}
	return jsonStr
}
