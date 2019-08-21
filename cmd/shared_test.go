package cmd

import (
	"testing"

	"github.com/flix-tech/confs.tech.push/confs"
)

func TestFormatLocationAddsFlag(t *testing.T) {
	location := formatLocation(confs.Conference{
		Name:      "Go two",
		URL:       "https://go2.com/",
		StartDate: "2019-08-21",
		EndDate:   "2019-08-21",
		City:      "Mariupol",
		Country:   "Ukraine",
	})
	expected := "Mariupol, Ukraine 🇺🇦"

	if location != expected {
		t.Errorf("Got error when formating location: expected '%s', got '%s'", expected, location)
	}
}

func TestFormatLocationWorksWithUnknownCountries(t *testing.T) {
	location := formatLocation(confs.Conference{
		Name:      "Go two",
		URL:       "https://go2.com/",
		StartDate: "2019-08-21",
		EndDate:   "2019-08-21",
		City:      "Voodoocity",
		Country:   "Voodooland",
	})
	expected := "Voodoocity, Voodooland"

	if location != expected {
		t.Errorf("Got error when formating location: expected '%s', got '%s'", expected, location)
	}
}
