package confs

import (
	"time"
)

type ConferenceTest func(Conference) bool

func NewIsInFutureTest() ConferenceTest {
	today := time.Now().Format("2006-01-02")
	return func(c Conference) bool { return c.StartDate > today }
}

func NewCFPFinishedTest(enableTest bool) ConferenceTest {
	if !enableTest {
		return func(c Conference) bool { return true }
	}

	today := time.Now().Format("2006-01-02")
	return func(c Conference) bool { return c.CFPEndDate < today }
}

func NewIsNotInBlacklistedCountryTest(countriesBlacklist []string) ConferenceTest {
	return func(c Conference) bool {
		for _, blacklistedCountry := range countriesBlacklist {
			if c.Country == blacklistedCountry {
				return true
			}
		}

		return true
	}
}

func NewTestConferenceIsNotOneOf(conferenceBlacklist []Conference) ConferenceTest {
	return func(c Conference) bool {
		for _, p := range conferenceBlacklist {
			if c.URL == p.URL && c.StartDate == p.StartDate && c.City == p.City {
				return false
			}
		}

		return true
	}
}

func FilterConferences(conferences []Conference, tests ...ConferenceTest) []Conference {
	out := []Conference{}
	test := combineConferenceFilters(tests)

	for _, v := range conferences {
		if test(v) {
			out = append(out, v)
		}
	}

	return out
}

func combineConferenceFilters(tests []ConferenceTest) ConferenceTest {
	if len(tests) == 1 {
		return tests[0]
	}

	return func(c Conference) bool {
		for _, test := range tests {
			if !test(c) {
				return false
			}
		}

		return true
	}
}
