package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/flix-tech/confs.tech.push/confs"
)

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

func formatDateRange(c confs.Conference) string {
	dateRange := c.StartDate
	if c.StartDate != c.EndDate {
		dateRange = fmt.Sprintf("%s — %s", c.StartDate, c.EndDate)
	}
	return dateRange
}

func formatLocation(c confs.Conference) string {
	return fmt.Sprintf("%s, %s", c.City, c.Country)
}
