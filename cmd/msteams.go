package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/urfave/cli.v1"

	"github.com/flix-tech/confs.tech.push/confs"
)

func MsteamsCommand() cli.Command {
	return cli.Command{
		Name:   "msteams",
		Usage:  "push to msteams",
		Action: wrapAction(msteamsAction),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "msteams-url",
				Usage:  "Teams Incoming Webhook url",
				EnvVar: "MSTEAMS_URL",
			},
			cli.StringFlag{
				Name:  "state-file, s",
				Value: "state.json",
				Usage: "State file path",
			},
		},
	}
}

func msteamsAction(topic string, conferences []confs.Conference, c *cli.Context) error {
	stateFile := c.String("state-file")
	processedConferences := confs.LoadState(stateFile)

	conferences = confs.FilterConferences(conferences,
		confs.NewTestConferenceIsNotOneOf(processedConferences),
	)

	// Push to slack
	webhookURL := c.String("msteams-url")
	if webhookURL == "" {
		return fmt.Errorf("Please provide Teams Incoming Webhook url")
	}

	for _, c := range conferences {
		err := pushToMsteams(c, webhookURL)
		if err != nil {
			_ = confs.SaveState(stateFile, processedConferences)
			return err
		}

		processedConferences = append(processedConferences, c)
	}

	return confs.SaveState(stateFile, processedConferences)
}

type msteamsMessage struct {
	Text string `json:"text"`
}

func pushToMsteams(c confs.Conference, webhookURL string) error {
	message := msteamsMessage{
		Text: fmt.Sprintf("[%s](%s)\n\n%s・%s", c.Name, c.URL, formatLocation(c), formatDateRange(c)),
	}

	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageString))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Got response code %d when sending message to msteams", resp.StatusCode)
	}

	return nil
}
