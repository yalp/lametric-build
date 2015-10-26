package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)
import "encoding/json"

type TravisPayload struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	Branch        string `json:"branch"`
	/*Repository    struct {
		Name      string `json:"name"`
		OwnerName string `json:"ower_name"`
	} `json:"repository"`*/
}

type CoverallsPayload struct {
	Branch   string
	Coverage float64
	Delta    float64
	//Repository string
}

var URL = os.Getenv("LAMETRIC_URL")
var TOKEN = os.Getenv("LAMETRIC_TOKEN")
var APP_ADDR = os.Getenv("APP_ADDR")

var buildChan chan *TravisPayload
var coverChan chan *CoverallsPayload

func coverageIcon(coverage float64) string {
	if coverage > 90 {
		return "i667" // green
	} else if coverage > 80 {
		return "i665" // orange
	}
	return "i660" // red
}

func progressLabel(delta float64) string {
	if delta > 0 {
		return "+" + strconv.FormatFloat(delta, 'f', 2, 64)
	}
	return strconv.FormatFloat(delta, 'f', 2, 64)
}

func progressIcon(delta float64) string {
	if delta > 0.1 {
		return "i120" // up
	} else if delta < -0.1 {
		return "i124" // down
	}
	return "i401" // equal
}

func statusIcon(status int) string {
	if status == 0 {
		return "i606"
	}
	return "i605"
}

func updateLametricApp(t *TravisPayload, c *CoverallsPayload) {
	app := App{
		Frames: []Frame{
			Frame{
				Index: 0,
				Text:  t.Branch + ": " + t.StatusMessage,
				Icon:  statusIcon(t.Status),
			},
			Frame{
				Index: 1,
				Icon:  coverageIcon(c.Coverage),
				GoalData: &Goal{
					End:     100,
					Current: int(c.Coverage + 0.5), // round
					Unit:    "%",
				},
			},
			Frame{
				Index: 2,
				Text:  progressLabel(c.Delta) + "%",
				Icon:  progressIcon(c.Delta),
			},
		},
	}
	if err := Push(URL, TOKEN, &app); err != nil {
		log.Println("Failed to update LaMetric app:", err)
	}
}

func handleTravis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm()
	payload := r.Form["payload"][0]
	log.Println("DEBUG: ", payload)
	var p TravisPayload
	err := json.Unmarshal([]byte(payload), &p)
	if err != nil {
		log.Println("Failed to parse Travis payload: ", err)
		return
	}

	buildChan <- &p
}

func handleCoveralls(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm()
	log.Println("DEBUG: ", r.Form)
	coverage, err := strconv.ParseFloat(r.Form["covered_percent"][0], 64)
	if err != nil {
		log.Println("Failed to get coverage", err)
		return
	}
	delta, err := strconv.ParseFloat(r.Form["coverage_change"][0], 64)
	if err != nil {
		log.Println("Failed to get coverage delta", err)
		return
	}
	//repo_name := r.Form["repo_name"][0]
	branch := r.Form["branch"][0]

	coverChan <- &CoverallsPayload{branch, coverage, delta}
}

func loop() {
	var lastCover *CoverallsPayload
	var lastBuild *TravisPayload
	var timeout <-chan time.Time
	for {
		select {
		case lastCover = <-coverChan:
			timeout = time.After(4 * time.Second)
		case lastBuild = <-buildChan:
			timeout = time.After(4 * time.Second)
		case <-timeout:
			if lastBuild != nil && lastCover != nil {
				updateLametricApp(lastBuild, lastCover)
			}
		}
	}
}

func main() {
	buildChan = make(chan *TravisPayload)
	coverChan = make(chan *CoverallsPayload)

	go loop()

	http.HandleFunc("/travis", handleTravis)
	http.HandleFunc("/coveralls", handleCoveralls)
	log.Fatal(http.ListenAndServe(APP_ADDR, nil))
}
