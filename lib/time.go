package mybot

import (
	"fmt"
	"time"
)

var timeLayouts = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func ParseTime(val string) (*time.Time, error) {
	for _, l := range timeLayouts {
		t, err := time.Parse(l, val)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("Failed to parse time: %s", val)
}
