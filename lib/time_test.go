package mybot

import (
	"testing"
	"time"
)

var timeExamples = []string{
	"Mon Jan 2 15:04:05 2006",
	"Mon Jan 2 15:04:05 MST 2006",
	"Mon Jan 02 15:04:05 -0700 2006",
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	// TODO: modify to make parsing possible
	// "2006-01-02T15:04:05Z07:00",
	// "2006-01-02T15:04:05.999999999Z07:00",
	time.Kitchen,
	"Jan 2 15:04:05",
	"Jan 2 15:04:05.000",
	"Jan 2 15:04:05.000000",
	"Jan 2 15:04:05.000000000",
}

func TestParseTime(t *testing.T) {
	for _, time := range timeExamples {
		_, err := ParseTime(time)
		if err != nil {
			t.Error(err)
		}
	}
}
