package cli

import (
	"testing"
	"time"
)

func TestParseTactics(t *testing.T) {
	p, err := parseTactics("tue=a, wed=c,default=b")
	if err != nil {
		t.Fatal(err)
	}
	for d, want := range map[time.Weekday]string{time.Tuesday: "a", time.Wednesday: "c", time.Friday: "b"} {
		if got := p.forDay(d); got != want {
			t.Errorf("%v: got %q want %q", d, got, want)
		}
	}
	for _, bad := range []string{"tue=z", "xyz=a"} {
		if _, err := parseTactics(bad); err == nil {
			t.Errorf("%q: expected error", bad)
		}
	}
}

func TestDayArithmeticCrossesMonth(t *testing.T) {
	d := truncateDay(time.Date(2026, 11, 30, 15, 0, 0, 0, time.Local)).AddDate(0, 0, 1)
	if got := at(d, 12); got.Month() != time.December || got.Day() != 1 || got.Hour() != 12 {
		t.Errorf("got %v", got)
	}
}
