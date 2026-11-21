package mzhttpclient

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GameDayState is what the friendly schedule says about one game day.
type GameDayState string

const (
	DayOpen    GameDayState = "open"     // no friendly booked, can be booked now
	DayBooked  GameDayState = "booked"   // has at least one friendly
	DayTooFar  GameDayState = "too-far"  // beyond the booking window (7 days)
	DayNoMatch GameDayState = "no-match" // league day: no friendlies at all
	DayPast    GameDayState = "past"     // before today, nothing booked
	DayUnknown GameDayState = "unknown"
)

// GameDay is one day of the friendly schedule. A game day runs from 01:00 to
// 00:59 the next day (local time); Date is its calendar date.
type GameDay struct {
	Number int
	Date   time.Time // local midnight; zero for league days (the page shows no date)
	State  GameDayState
}

// Invite is a friendly invite in the incoming or outgoing list.
type Invite struct {
	Day    int    // game day number
	TeamID string // the other team
	// AcceptURL is the link behind the accept button (incoming only); MZ's
	// page appends the chosen tactics to it. Side is handleChallenge's second
	// argument: 2 = we play away only, 3 = home only, otherwise a double.
	// Not yet seen on a live page: the format is read from challenges.js
	// (handleChallenge(url, side, date1, date2)).
	AcceptURL string
	Side      int
}

// Schedule is the friendly page (Challenges > Friendlies).
type Schedule struct {
	Days     []GameDay
	Incoming []Invite
	Outgoing []Invite
}

// Day returns the game day with the given number.
func (s Schedule) Day(n int) (GameDay, bool) {
	for _, d := range s.Days {
		if d.Number == n {
			return d, true
		}
	}
	return GameDay{}, false
}

// GetSchedule fetches and parses the friendly schedule with the invite lists.
func (mz *MZHTTPClient) GetSchedule(now time.Time) (Schedule, error) {
	body, err := mz.Get("ajax.php?p=challenge&sub=personal-challenge-template&sport=soccer")
	if err != nil {
		return Schedule{}, fmt.Errorf("fetch friendly schedule: %w", err)
	}
	return ParseSchedule(string(body), now)
}

var (
	reDayBlock  = regexp.MustCompile(`<div id="d(\d+)" class="d\d+">`)
	reDayStart  = regexp.MustCompile(`date-range">\s*<strong><span>[^<]*</span>(\d{1,2})/(\d{1,2}) `)
	reInviteRow = regexp.MustCompile(`<tr class="[^"]*\bd(\d+)\b[^"]*"`)
	reRowTeam   = regexp.MustCompile(`href="/\?p=team&(?:amp;)?tid=(\d+)"`)
	reAcceptRef = regexp.MustCompile(`href="([^"]*sub=handle-challenge[^"]*action=accept[^"]*)"`)
	reHandle    = regexp.MustCompile(`handleChallenge\(\s*(?:'([^']*)'|&#0?39;(.*?)&#0?39;|"([^"]*)"|&quot;(.*?)&quot;)\s*,\s*(\d+)`)
)

// ParseSchedule parses the personal-challenge-template page. now decides the
// year of the schedule's day/month dates.
func ParseSchedule(page string, now time.Time) (Schedule, error) {
	var s Schedule
	locs := reDayBlock.FindAllStringSubmatchIndex(page, -1)
	if len(locs) == 0 {
		return s, fmt.Errorf("friendly schedule: no game days found (page layout changed?)")
	}
	for i, l := range locs {
		end := len(page)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		block := page[l[1]:end]
		n, _ := strconv.Atoi(page[l[2]:l[3]])
		d := GameDay{Number: n, State: DayUnknown}
		if m := reDayStart.FindStringSubmatch(block); m != nil {
			day, _ := strconv.Atoi(m[1])
			mon, _ := strconv.Atoi(m[2])
			d.Date = nearestDate(now, time.Month(mon), day)
		}
		switch {
		case strings.Contains(block, "match-wrapper"):
			d.State = DayBooked
		case strings.Contains(block, "fss-icon-day-active"):
			d.State = DayOpen
		case strings.Contains(block, "fss-icon-day-future"):
			d.State = DayTooFar
		case d.Date.IsZero():
			d.State = DayNoMatch
		case d.Date.Before(truncDay(now)):
			d.State = DayPast
		}
		s.Days = append(s.Days, d)
	}

	in, out := section(page, `id="matches_in_wrapper"`), section(page, `id="matches_out_wrapper"`)
	s.Incoming = parseInvites(in, true)
	s.Outgoing = parseInvites(out, false)
	return s, nil
}

// section returns the page from marker up to the next "_wrapper" div, or "".
func section(page, marker string) string {
	i := strings.Index(page, marker)
	if i < 0 {
		return ""
	}
	rest := page[i+len(marker):]
	if j := strings.Index(rest, `_wrapper"`); j >= 0 {
		rest = rest[:j]
	}
	return rest
}

func parseInvites(sec string, incoming bool) []Invite {
	locs := reInviteRow.FindAllStringSubmatchIndex(sec, -1)
	var out []Invite
	for i, l := range locs {
		end := len(sec)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		row := sec[l[1]:end]
		inv := Invite{}
		inv.Day, _ = strconv.Atoi(sec[l[2]:l[3]])
		if m := reRowTeam.FindStringSubmatch(row); m != nil {
			inv.TeamID = m[1]
		}
		if incoming {
			if m := reHandle.FindStringSubmatch(row); m != nil {
				u := m[1] + m[2] + m[3] + m[4]
				inv.AcceptURL = html.UnescapeString(u)
				inv.Side, _ = strconv.Atoi(m[5])
			} else if m := reAcceptRef.FindStringSubmatch(row); m != nil {
				inv.AcceptURL = html.UnescapeString(m[1])
			}
		}
		out = append(out, inv)
	}
	return out
}

// AcceptInvite accepts an incoming invite with the given tactic for our match(es).
func (mz *MZHTTPClient) AcceptInvite(inv Invite, tactic string) error {
	if inv.AcceptURL == "" {
		return fmt.Errorf("no accept link found for the invite from %s", inv.TeamID)
	}
	q := "&tacticset=" + url.QueryEscape(tactic)
	if inv.Side != 2 && inv.Side != 3 {
		q += "&tacticset2=" + url.QueryEscape(tactic)
	}
	return mz.getSuccess(inv.AcceptURL + q)
}

// QuickChallenge is the quick-challenge page: one double (12:00 and 20:00) on
// the next game day, either by opting in or by accepting someone who did.
type QuickChallenge struct {
	Day   int    // game day the page is for
	State string // "opted-in", "opted-out" (can opt in) or "booked"
	Date  string // value of #qc-date, sent back when opting in
	// Offers are teams that opted in, with the link behind their book button.
	Offers []QuickOffer
	// OfferRows counts the rows of the open-challenges table, with or
	// without a recognised book link.
	OfferRows int
	// Teams are all team ids on the page in page order (offers and recently
	// booked), duplicates included.
	Teams []string
}

type QuickOffer struct {
	TeamID  string
	BookURL string
}

const quickChallengePath = "ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer"

// GetQuickChallenge fetches and parses the quick-challenge page.
func (mz *MZHTTPClient) GetQuickChallenge() (QuickChallenge, error) {
	body, err := mz.Get(quickChallengePath)
	if err != nil {
		return QuickChallenge{}, fmt.Errorf("fetch quick challenges: %w", err)
	}
	return ParseQuickChallenge(string(body))
}

var (
	reQCDay    = regexp.MustCompile(`id="qc-opted-out"[^>]*>[^<]*<strong>Day (\d+)</strong>`)
	reQCActive = regexp.MustCompile(`<div id="qc-(opted-in|opted-out|booked)" class="active">`)
	reQCDate   = regexp.MustCompile(`id="qc-date" value="(\d+)"`)
	// The book button is <a href="/ajax.php?p=challenge&sub=quick-challenge-accept&tid=…"
	// onclick="quickChallenge.book(this.href)">.
	reQCBook = regexp.MustCompile(`href="([^"]*sub=quick-challenge-accept[^"]*)"`)
)

// ParseQuickChallenge parses the quick-challenge-template page.
func ParseQuickChallenge(page string) (QuickChallenge, error) {
	var qc QuickChallenge
	m := reQCDay.FindStringSubmatch(page)
	if m == nil {
		return qc, fmt.Errorf("quick challenge: game day not found (page layout changed?)")
	}
	qc.Day, _ = strconv.Atoi(m[1])
	if m := reQCActive.FindStringSubmatch(page); m != nil {
		qc.State = m[1]
	}
	if m := reQCDate.FindStringSubmatch(page); m != nil {
		qc.Date = m[1]
	}
	// Offers: rows of the open-challenges table (not "Recently booked").
	if i := strings.Index(page, `id="quick_challenges"`); i >= 0 {
		tbl := page[i:]
		if j := strings.Index(tbl, "</table>"); j >= 0 {
			tbl = tbl[:j]
		}
		qc.OfferRows = strings.Count(tbl, `<tr class="a `)
		for _, row := range strings.Split(tbl, "<tr")[1:] {
			b := reQCBook.FindStringSubmatch(row)
			t := reRowTeam.FindStringSubmatch(row)
			if b == nil || t == nil {
				continue
			}
			qc.Offers = append(qc.Offers, QuickOffer{TeamID: t[1], BookURL: html.UnescapeString(b[1])})
		}
	}
	for _, m := range reRowTeam.FindAllStringSubmatch(page, -1) {
		qc.Teams = append(qc.Teams, m[1])
	}
	return qc, nil
}

// BookQuickChallenge accepts another team's quick challenge.
func (mz *MZHTTPClient) BookQuickChallenge(o QuickOffer, tacticHome, tacticAway string) error {
	return mz.getSuccess(o.BookURL + "&home_tacticset=" + url.QueryEscape(tacticHome) + "&away_tacticset=" + url.QueryEscape(tacticAway))
}

// OptInQuickChallenge opts our team in for the quick challenge's game day.
func (mz *MZHTTPClient) OptInQuickChallenge(qc QuickChallenge, tacticHome, tacticAway string) error {
	body, err := mz.Get("ajax.php?p=challenge&sub=quick-challenge-opt-in&date=" + url.QueryEscape(qc.Date) +
		"&tactic_home=" + url.QueryEscape(tacticHome) + "&tactic_away=" + url.QueryEscape(tacticAway) + "&sport=soccer")
	if err != nil {
		return err
	}
	var r struct {
		Success      bool   `json:"success"`
		ErrorContent string `json:"error_content"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return fmt.Errorf("opt-in response is not JSON: %w", err)
	}
	if !r.Success {
		return fmt.Errorf("opt-in refused: %s", stripTags(r.ErrorContent))
	}
	return nil
}

// getSuccess GETs a link from an MZ page and expects {"response":"success"}.
func (mz *MZHTTPClient) getSuccess(link string) error {
	body, err := mz.Get(sitePath(link))
	if err != nil {
		return err
	}
	var r struct {
		Response        string `json:"response"`
		ResponseContent string `json:"response_content"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return fmt.Errorf("response is not JSON: %w", err)
	}
	if r.Response != "success" {
		msg := stripTags(r.ResponseContent)
		if msg == "" {
			msg = string(body)
		}
		return fmt.Errorf("refused (response=%q): %s", r.Response, msg)
	}
	return nil
}

// sitePath turns a link from a page (absolute, /-rooted or relative) into a
// path relative to the site root.
func sitePath(link string) string {
	if u, err := url.Parse(link); err == nil && u.IsAbs() {
		link = u.RequestURI()
	}
	return strings.TrimPrefix(link, "/")
}

var reTags = regexp.MustCompile(`<[^>]*>`)

func stripTags(s string) string {
	return strings.Join(strings.Fields(html.UnescapeString(reTags.ReplaceAllString(s, " "))), " ")
}

func truncDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// nearestDate returns the local date with the given month and day closest to now.
func nearestDate(now time.Time, mon time.Month, day int) time.Time {
	best := time.Time{}
	for _, y := range []int{now.Year() - 1, now.Year(), now.Year() + 1} {
		d := time.Date(y, mon, day, 0, 0, 0, 0, time.Local)
		if best.IsZero() || absDur(d.Sub(now)) < absDur(best.Sub(now)) {
			best = d
		}
	}
	return best
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
