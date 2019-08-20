package cmd

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
	"gopkg.in/urfave/cli.v1"

	"github.com/flix-tech/confs.tech.push/confs"
)

func AtomCommand() cli.Command {
	return cli.Command{
		Name:   "atom",
		Usage:  "generate atom feed",
		Action: atomAction,
		Flags:  []cli.Flag{},
	}
}

func atomAction(c *cli.Context) error {
	topic, err := validateTopicArgument(c.Args().Get(0))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Fetch conference data
	conferences, err := confs.GetConferences(topic)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	conferences = confs.FilterConferences(conferences,
		confs.NewIsInFutureTest(),
		confs.NewCFPFinishedTest(c.GlobalBool("cfp-finished")),
		confs.NewIsNotInBlacklistedCountryTest(c.GlobalStringSlice("countries-blacklist")),
	)

	atom, err := generateAtomFeed(topic, conferences)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(atom)

	return nil
}

func generateAtomFeed(topic string, conferences []confs.Conference) (string, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:   topic + " tech conferences",
		Link:    &feeds.Link{Href: fmt.Sprintf("https://confs.tech/%s", topic)},
		Author:  &feeds.Author{Name: "https://confs.tech/"},
		Created: now,
	}

	items := []*feeds.Item{}
	for _, c := range conferences {
		items = append(items, &feeds.Item{
			Title:       c.Name,
			Link:        &feeds.Link{Href: c.URL},
			Id:          c.URL,
			Description: fmt.Sprintf("%s・%s", formatLocation(c), formatDateRange(c)),
			Created:     now,
		})
	}

	feed.Items = items

	return feed.ToAtom()
}
