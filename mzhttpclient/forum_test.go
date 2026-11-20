package mzhttpclient

import (
	"testing"
	"time"
)

func TestParseForumTime(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Tallinn")
	old := time.Local
	time.Local = loc
	defer func() { time.Local = old }()

	now := time.Date(2026, 11, 20, 10, 0, 0, 0, loc) // a Friday
	for in, want := range map[string]string{
		"12-10-2026 13:52": "2026-10-12 13:52",
		"Today 05:45":      "2026-11-20 05:45",
		"Yesterday 18:10":  "2026-11-19 18:10",
		"Wednesday 09:00":  "2026-11-18 09:00",
		"Saturday 23:59":   "2026-11-14 23:59",
		"Friday 08:00":     "2026-11-13 08:00", // a week ago; today would be "Today"
	} {
		got, err := ParseForumTime(in, now)
		if err != nil {
			t.Errorf("%q: %v", in, err)
			continue
		}
		if s := got.Format("2006-01-02 15:04"); s != want {
			t.Errorf("%q = %s, want %s", in, s, want)
		}
	}
	for _, in := range []string{"", "Tomorrow 10:00", "2026-10-12 13:52"} {
		if _, err := ParseForumTime(in, now); err == nil {
			t.Errorf("%q: no error", in)
		}
	}
}
