package convert

import (
	"errors"
	"fmt"
	"time"
)

var layouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"02/Jan/2006:15:04:05 -0700",
	"2006/01/02 15:04:05",
	`20060102T150405-07`,
	"2006-01-02 15:04:05 -0700 MST",
	"2006-01-02 15:04:05 -0700",
	`2006-01-02T15:04:05-07:00`,
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05 -0700 MST",
	"2006/01/02 15:04:05 -0700",
	"2006-01-02 -0700 MST",
	"2006-01-02 -0700",
	"2006-01-02",
	"2006/01/02 -0700 MST",
	"2006/01/02 -0700",
	"2006/01/02",
	"02/01/2006--15:04:05",
	"02 Jan 06 15:04",
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func StrToTimeLocation(value string, loc *time.Location) (time.Time, error) {
	if value == "" {
		return time.Now(), errors.New("empty time string")
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, value, loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("can not find any layout to parse %v", value)
}
