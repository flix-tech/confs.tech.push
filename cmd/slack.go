package cmd

import (
	"bytes"
	"fmt"

	"encoding/json"
	"net/http"

	"gopkg.in/urfave/cli.v1"

	"github.com/flix-tech/confs.tech.push/confs"
)

func SlackCommand() cli.Command {
	return cli.Command{
		Name:   "slack",
		Usage:  "push to slack",
		Action: wrapAction(slackAction),
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
	}
}

func slackAction(topic string, conferences []confs.Conference, c *cli.Context) error {
	stateFile := c.String("state-file")
	processedConferences := confs.LoadState(stateFile)

	conferences = confs.FilterConferences(conferences,
		confs.NewTestConferenceIsNotOneOf(processedConferences),
	)

	// Push to slack
	slackURL := c.String("slack-url")
	if slackURL == "" {
		return fmt.Errorf("Please provide slack Incoming Webhook url")
	}
	slackChannel := c.String("slack-channel")

	for _, c := range conferences {
		err := pushToSlack(c, slackURL, slackChannel)
		if err != nil {
			_ = confs.SaveState(stateFile, processedConferences)
			return err
		}

		processedConferences = append(processedConferences, c)
	}

	return confs.SaveState(stateFile, processedConferences)
}

type slackField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

type slackAttachment struct {
	Fallback  string       `json:"fallback,omitempty"`
	Color     string       `json:"color,omitempty"`
	Title     string       `json:"title,omitempty"`
	TitleLink string       `json:"title_link,omitempty"`
	Text      string       `json:"text,omitempty"`
	Fields    []slackField `json:"fields,omitempty"`
}

type slackMessage struct {
	Channel     string            `json:"channel,omitempty"`
	Text        string            `json:"text,omitempty"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
	UnfurlLinks bool              `json:"unfurl_links,omitempty"`
	Markdown    bool              `json:"mrkdwn,omitempty"`
}

func pushToSlack(c confs.Conference, slackURL string, slackChannel string) error {
	message := slackMessage{
		Channel: slackChannel,
		Text:    fmt.Sprintf("*%s*\n<%s>", c.Name, c.URL),
		Attachments: []slackAttachment{
			slackAttachment{
				Fields: []slackField{
					slackField{
						Title: "Location",
						Value: formatLocation(c),
						Short: true,
					},
					slackField{
						Title: "Dates",
						Value: formatDateRange(c),
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

	resp, err := http.Post(slackURL, "application/json", bytes.NewBuffer(messageString))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Got response code %d when sending message to slack", resp.StatusCode)
	}

	return nil
}
