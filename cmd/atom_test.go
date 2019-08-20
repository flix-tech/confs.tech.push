package cmd

import (
	"testing"

	"github.com/flix-tech/confs.tech.push/confs"
)

func TestAtomGeneration(t *testing.T) {
	_, err := generateAtomFeed("golang", []confs.Conference{
		confs.Conference{
			Name:      "Go one",
			URL:       "https://go1.com/",
			StartDate: "2019-08-20",
			EndDate:   "2019-08-20",
			City:      "Berlin",
			Country:   "Germany",
		},
		confs.Conference{
			Name:      "Go two",
			URL:       "https://go2.com/",
			StartDate: "2019-08-21",
			EndDate:   "2019-08-21",
			City:      "Mariupol",
			Country:   "Ukraine",
		},
	})

	if err != nil {
		t.Errorf("Got error when generating conferences atom: %s", err)
	}
}
