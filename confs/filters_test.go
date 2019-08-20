package confs

import (
	"testing"
	"time"
)

func TestFutureConferencePass(t *testing.T) {
	result := NewIsInFutureTest()(Conference{
		Name:      "Future Conference",
		StartDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
	})

	if result == false {
		t.Errorf("Future Conference test must pass testConferenceIsInFuture()")
	}
}

func TestFutureConferenceFail(t *testing.T) {
	result := NewIsInFutureTest()(Conference{
		Name:      "Past Conference",
		StartDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	})

	if result == true {
		t.Errorf("FutPasture Conference test must fail testConferenceIsInFuture()")
	}
}

func TestFilterCFPFinishedConferencesWithNoCFP(t *testing.T) {
	result := NewCFPFinishedTest(true)(Conference{
		Name: "no CFP",
	})

	if result == false {
		t.Errorf("No CFP Conference test must pass testConferenceCFPFinished()")
	}
}

func TestFilterCFPFinishedConferencesWithFinishedCFP(t *testing.T) {
	result := NewCFPFinishedTest(true)(Conference{
		Name:       "CFP finished",
		CFPEndDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	})

	if result == false {
		t.Errorf("Finished CFP Conference test must pass testConferenceCFPFinished()")
	}
}

func TestFilterCFPFinishedConferencesWithNotFinishedCFP(t *testing.T) {
	result := NewCFPFinishedTest(true)(Conference{
		Name:       "CFP not finished",
		CFPEndDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
	})

	if result == true {
		t.Errorf("Not finished CFP Conference test must fail testConferenceCFPFinished()")
	}
}
