package main

import (
	"testing"
	"time"
)

func TestFilterFutureConferences(t *testing.T) {
	conferences := []Conference{
		Conference{
			Name:      "Past Conference",
			StartDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
		Conference{
			Name:      "Future Conference",
			StartDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		},
	}

	filtered := filterFutureConferences(conferences)

	if len(filtered) != 1 {
		t.Errorf("FutureConferences length was incorrect, got: %d, want: %d.", len(filtered), 1)
	}
	if filtered[0].Name != "Future Conference" {
		t.Errorf("FutureConference name was incorrect, got: %s, want: %s.", filtered[0].Name, "Future Conference")
	}
}

func TestFilterCFPFinishedConferences(t *testing.T) {
	conferences := []Conference{
		Conference{
			Name: "no CFP",
		},
		Conference{
			Name:       "CFP not finished",
			CFPEndDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		},
		Conference{
			Name:       "CFP finished",
			CFPEndDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
	}

	filtered := filterCFPFinishedConferences(conferences)

	if len(filtered) != 2 {
		t.Errorf("FutureConferences length was incorrect, got: %d, want: %d.", len(filtered), 1)
	}
	if filtered[0].Name != "no CFP" {
		t.Errorf("FutureConference name was incorrect, got: %s, want: %s.", filtered[0].Name, "no CFP")
	}
	if filtered[1].Name != "CFP finished" {
		t.Errorf("FutureConference name was incorrect, got: %s, want: %s.", filtered[0].Name, "CFP finished")
	}
}
