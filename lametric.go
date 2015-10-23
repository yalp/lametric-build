package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var client = &http.Client{}

type App struct {
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Index    int    `json:"index"`
	Text     string `json:"text,omitempty"`
	Icon     string `json:"icon,omitempty"`
	GoalData *Goal  `json:"goalData,omitempty"`
}

type Goal struct {
	Start   int    `json:"start"`
	Current int    `json:"current"`
	End     int    `json:"end"`
	Unit    string `json:"unit"`
}

func Push(url, accessToken string, app *App) error {
	payload, err := json.Marshal(app)
	if err != nil {
		return err
	}
	//log.Println("DEBUG:", string(payload))
	r, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	r.Header.Set("Accept", "application/json")
	r.Header.Set("X-Access-Token", accessToken)
	r.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
