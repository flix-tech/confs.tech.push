package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"encoding/json"
	"net/http"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "confs.tech.push"
	app.Usage = "push data about tech conferences to somewhere you can read it"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "countries-blacklist, C",
			Usage:  "Countries to be blocked",
			EnvVar: "COUNTRIES_BLACKLIST",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "slack",
			Usage:  "push to slack",
			Action: slackAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "slack-url",
					Usage:  "Slack Incoming Webhook url",
					EnvVar: "SLACK_URL",
				},
				cli.StringFlag{
					Name:   "slack-channel, k",
					Usage:  "Slack channel name",
					EnvVar: "SLACK_CHANNEL",
				},
				cli.StringFlag{
					Name:  "state-file, s",
					Value: "state.json",
					Usage: "State file path",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func slackAction(c *cli.Context) error {
	topic, err := validateTopicArgument(c.Args().Get(0))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Fetch conference data
	conferences, err := getConferences(topic)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Filter out past conferences
	conferences = filterFutureConferences(conferences)

	// Filter by countries
	countriesBlacklist := c.GlobalStringSlice("countries-blacklist")
	conferences = filterConferences(conferences, func(c Conference) bool {
		return !stringInArray(c.Country, countriesBlacklist)
	})

	stateFile := c.String("state-file")
	processedConferences := loadState(stateFile)

	conferences = filterConferences(conferences, func(c Conference) bool {
		for _, p := range processedConferences {
			if c.URL == p.URL && c.StartDate == p.StartDate && c.City == p.City {
				return false
			}
		}

		return true
	})

	// Push to slack
	slackUrl := c.String("slack-url")
	if slackUrl == "" {
		return cli.NewExitError("Please provide slack Incoming Webhook url", 1)
	}
	slackChannel := c.String("slack-channel")

	for _, c := range conferences {
		err = pushToSlack(c, slackUrl, slackChannel)
		if err != nil {
			_ = saveState(stateFile, processedConferences)
			return cli.NewExitError(err, 1)
		}

		processedConferences = append(processedConferences, c)
	}

	err = saveState(stateFile, processedConferences)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func validateTopicArgument(topic string) (string, error) {
	if topic == "" {
		return "", errors.New("Please provide conference topic")
	}

	match, _ := regexp.MatchString("^[a-z\\-]+$", topic)
	if !match {
		return "", errors.New("Invalid conference topic")
	}

	return topic, nil
}

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

func getConferences(topic string) ([]Conference, error) {
	var conferences []Conference

	url := fmt.Sprintf("https://raw.githubusercontent.com/tech-conferences/conference-data/master/conferences/%d/%s.json", time.Now().Year(), topic)
	resp, err := http.Get(url)
	if err != nil {
		return conferences, err
	}
	if resp.StatusCode != 200 {
		return conferences, errors.New(fmt.Sprintf("Got response code %d when calling %s", resp.StatusCode, url))
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&conferences)
	if err != nil {
		return conferences, err
	}

	return conferences, nil
}

func filterFutureConferences(conferences []Conference) []Conference {
	return filterConferences(conferences, func(c Conference) bool {
		return c.StartDate > time.Now().Format("2006-01-02")
	})
}

func filterConferences(conferences []Conference, test func(Conference) bool) []Conference {
	out := make([]Conference, 0, len(conferences))

	for _, v := range conferences {
		if test(v) {
			out = append(out, v)
		}
	}

	return out
}

func stringInArray(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}

	return false
}

type SlackField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

type SlackAttachment struct {
	Fallback  string       `json:"fallback,omitempty"`
	Color     string       `json:"color,omitempty"`
	Title     string       `json:"title,omitempty"`
	TitleLink string       `json:"title_link,omitempty"`
	Text      string       `json:"text,omitempty"`
	Fields    []SlackField `json:"fields,omitempty"`
}

type SlackMessage struct {
	Channel     string            `json:"channel,omitempty"`
	Text        string            `json:"text,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
	UnfurlLinks bool              `json:"unfurl_links,omitempty"`
	Markdown    bool              `json:"mrkdwn,omitempty"`
}

func pushToSlack(c Conference, slackUrl string, slackChannel string) error {
	dateRange := c.StartDate
	if c.StartDate != c.EndDate {
		dateRange = fmt.Sprintf("%s — %s", c.StartDate, c.EndDate)
	}

	message := SlackMessage{
		Channel: slackChannel,
		Text:    fmt.Sprintf("*%s*\n<%s>", c.Name, c.URL),
		Attachments: []SlackAttachment{
			SlackAttachment{
				Fields: []SlackField{
					SlackField{
						Title: "Location",
						Value: fmt.Sprintf("%s, %s", c.City, c.Country),
						Short: true,
					},
					SlackField{
						Title: "Dates",
						Value: dateRange,
						Short: true,
					},
				},
			},
		},
		UnfurlLinks: true,
		Markdown:    true,
	}
	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(slackUrl, "application/json", bytes.NewBuffer(messageString))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Got response code %d when sending message to slack", resp.StatusCode))
	}

	return nil
}

func loadState(finename string) []Conference {
	state, err := ioutil.ReadFile(finename)
	if err != nil {
		return []Conference{}
	}

	var conferences = []Conference{}
	json.Unmarshal(state, &conferences)

	return filterFutureConferences(conferences)
}

func saveState(filename string, conferences []Conference) error {
	stateString, err := json.Marshal(conferences)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, stateString, 0644)
}
