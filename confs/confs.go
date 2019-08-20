package confs

import (
	"fmt"
	"io/ioutil"
	"time"

	"encoding/json"
	"net/http"
)

type Conference struct {
	Name       string
	URL        string
	StartDate  string
	EndDate    string
	City       string
	Country    string
	CFPUrl     string
	CFPEndDate string
	Twitter    string
}

func GetConferences(topic string) ([]Conference, error) {
	var conferences []Conference

	url := fmt.Sprintf("https://raw.githubusercontent.com/tech-conferences/conference-data/master/conferences/%d/%s.json", time.Now().Year(), topic)
	resp, err := http.Get(url)
	if err != nil {
		return conferences, err
	}
	if resp.StatusCode != 200 {
		return conferences, fmt.Errorf("Got response code %d when calling %s", resp.StatusCode, url)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&conferences)
	if err != nil {
		return conferences, err
	}

	return conferences, nil
}

func LoadState(finename string) []Conference {
	state, err := ioutil.ReadFile(finename)
	if err != nil {
		return []Conference{}
	}

	var conferences = []Conference{}
	json.Unmarshal(state, &conferences)

	return FilterConferences(conferences, NewIsInFutureTest())
}

func SaveState(filename string, conferences []Conference) error {
	stateString, err := json.Marshal(conferences)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, stateString, 0644)
}
