package mzhttpclient

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	Football string = "soccer"
	Hockey   string = "hockey"
)

// MZHTTPClient struct represents the HTTP client for interacting with ManagerZone
type MZHTTPClient struct {
	client    *http.Client
	phpsessid string
	baseURL   string
	sport     string
}

// New creates a football client; use WithSport to switch to hockey.
func New(baseURL string) *MZHTTPClient {
	return &MZHTTPClient{
		client:  &http.Client{},
		baseURL: baseURL,
		sport:   Football,
	}
}

// WithSport sets the sport (Football or Hockey) used for login and
// sport-specific pages. Call it before Login.
func (mz *MZHTTPClient) WithSport(sport string) *MZHTTPClient {
	mz.sport = sport
	return mz
}

// Sport is the client's sport, Football ("soccer") or Hockey ("hockey").
func (mz *MZHTTPClient) Sport() string { return mz.sport }

// sportID is the numeric sport id some endpoints use (1 football, 2 hockey).
func (mz *MZHTTPClient) sportID() int {
	if mz.sport == Hockey {
		return 2
	}
	return 1
}

// Login authenticates the user and stores the PHPSESSID value
func (mz *MZHTTPClient) Login(username, password string) error {
	// Step 1: Get PHPSESSID from the initial request
	resp, err := mz.client.Get(mz.baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch initial page: %w", err)
	}
	defer resp.Body.Close()

	// Extract PHPSESSID from cookies
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			mz.phpsessid = cookie.Value
			break
		}
	}

	if mz.phpsessid == "" {
		return errors.New("failed to retrieve PHPSESSID")
	}

	// Convert password to lowercase
	password = strings.ToLower(password)

	// This is a MZ's client side hardcoded "secret" value for some reason.
	// logindata[md5] needs to be correctly calculated for successful login.
	secret := "g372=RVx6!yn"
	// Create a new HMAC-MD5 hash using the secret key
	h := hmac.New(md5.New, []byte(secret))
	h.Write([]byte(password))

	form := url.Values{
		"logindata[md5]":           {hex.EncodeToString(h.Sum(nil))}, // Placeholder for now
		"logindata[username]":      {username},
		"logindata[password]":      {password},
		"logindata[markasdefault]": {""},
		"logindata[sport]":         {mz.sport},
		"logindata[remember_me]":   {""},
	}

	req, err := http.NewRequest("POST", mz.baseURL+"?p=login", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	req.Header.Set("User-Agent", "Mozilla/5.0")

	loginResp, err := mz.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform login: %w", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", loginResp.StatusCode)
	}

	// MZ answers 200 even for bad credentials; check that a logged-in page loads.
	if _, err := mz.Whoami(); err != nil {
		mz.phpsessid = ""
		return fmt.Errorf("login failed (wrong credentials?): %w", err)
	}
	return nil
}

// Me is the logged-in manager.
type Me struct {
	Username string
	TeamID   string
}

var (
	reHeaderUsername = regexp.MustCompile(`id="header-username">([^<]+)<`)
	reOwnPlayerTeam  = regexp.MustCompile(`href="/\?p=players&pid=[0-9]+&tid=([0-9]+)"`)
)

// Whoami reads the username and own team id from the players page.
func (mz *MZHTTPClient) Whoami() (Me, error) {
	body, err := mz.Get("?p=players")
	if err != nil {
		return Me{}, err
	}
	u := reHeaderUsername.FindSubmatch(body)
	if u == nil {
		return Me{}, errors.New("not logged in: username not found on players page")
	}
	me := Me{Username: html.UnescapeString(string(u[1]))}
	if t := reOwnPlayerTeam.FindSubmatch(body); t != nil {
		me.TeamID = string(t[1])
	}
	return me, nil
}

// Get fetches a path relative to baseURL with the session cookie and returns the raw body.
func (mz *MZHTTPClient) Get(path string) ([]byte, error) {
	if mz.phpsessid == "" {
		return nil, errors.New("not logged in: PHPSESSID is empty")
	}
	req, err := http.NewRequest("GET", mz.baseURL+strings.TrimPrefix(path, "/"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := mz.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", path, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// GetTeamIDsOnQuickChallengePage returns the team ids on the quick-challenge page.
func (mz *MZHTTPClient) GetTeamIDsOnQuickChallengePage() ([]string, error) {
	if mz.phpsessid == "" {
		return nil, errors.New("not logged in: PHPSESSID is empty")
	}

	req, err := http.NewRequest("GET", mz.baseURL+quickChallengePath+mz.sport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for my team: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := mz.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch my team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch my team: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Define the regex pattern
	pattern := `href="/\?p=team&tid=([0-9]+)"`

	// Compile the regex
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	regex_found_tids := r.FindAllStringSubmatch(string(body), -1)
	var tids []string
	for _, t := range regex_found_tids {
		if len(t) > 1 { // Ensure the capturing group exists
			tids = append(tids, t[1])
		}
	}
	return tids, nil
}

// GetMatches returns the played or upcoming matches of a team (any team id).
func (mz *MZHTTPClient) GetMatches(teamId string, status MatchStatus) ([]Match, error) {
	body, err := mz.Get(fmt.Sprintf("xml/team_matchlist.php?team_id=%s&match_status=%d", url.QueryEscape(teamId), status))
	if err != nil {
		return nil, fmt.Errorf("fetch match list: %w", err)
	}

	type xmlTeam struct {
		Field    string `xml:"field,attr"`
		TeamID   string `xml:"teamId,attr"`
		TeamName string `xml:"teamName,attr"`
		Goals    string `xml:"goals,attr"`
	}
	type xmlMatch struct {
		ID     string    `xml:"id,attr"`
		Date   string    `xml:"date,attr"`
		Type   string    `xml:"type,attr"`
		Status string    `xml:"status,attr"`
		Teams  []xmlTeam `xml:"Team"`
	}
	var result struct {
		Matches []xmlMatch `xml:"Match"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse match list: %w", err)
	}

	matches := make([]Match, 0, len(result.Matches))
	for _, m := range result.Matches {
		d, err := time.ParseInLocation("2006-01-02 15:04", m.Date, time.Local)
		if err != nil {
			return nil, fmt.Errorf("match %s: bad date %q: %w", m.ID, m.Date, err)
		}
		match := Match{Id: m.ID, MatchType: m.Type, Status: m.Status, Date: d}
		for _, t := range m.Teams {
			switch t.Field {
			case "home":
				match.HomeTeamId, match.HomeTeamName, match.HomeGoals = t.TeamID, t.TeamName, t.Goals
			case "away":
				match.AwayTeamId, match.AwayTeamName, match.AwayGoals = t.TeamID, t.TeamName, t.Goals
			}
		}
		matches = append(matches, match)
	}
	return matches, nil
}

// Note: This sent friendly requests to days where I already had a friendly.
func (mz *MZHTTPClient) SendFriendlyRequest(teamId string, dateTimeHome time.Time, tacticHome string, dateTimeAway time.Time, tacticAway string) error {

	if mz.phpsessid == "" {
		return errors.New("not logged in: PHPSESSID is empty")
	}
	data := url.Values{}
	cet, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		return fmt.Errorf("failed to load timezone Europe/Stockholm: %w", err)
	}

	awaysrvtime := dateTimeAway.In(cet)
	homesrvtime := dateTimeHome.In(cet)
	a := time.Date(awaysrvtime.Year(), awaysrvtime.Month(), awaysrvtime.Day(), 0, 0, 0, 0, cet)
	h := time.Date(homesrvtime.Year(), homesrvtime.Month(), homesrvtime.Day(), 0, 0, 0, 0, cet)

	data.Set("date_home", fmt.Sprintf("%d,%d", h.Unix(), homesrvtime.Hour()))
	data.Set("tactic_home", tacticHome)
	data.Set("date_away", fmt.Sprintf("%d,%d", a.Unix(), awaysrvtime.Hour()))
	data.Set("tactic_away", tacticAway)

	req, err := http.NewRequest("POST", mz.baseURL+"ajax.php?p=challenge&sub=handle-challenge&action=create&tid="+teamId+"&country_id=123&sport="+mz.sport, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request for my team: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := mz.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send friendly request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("friendly request returned other than http 200 OK: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	body := string(bodyBytes)

	// Try to parse the JSON and check "response" == "success"
	var jsonResp struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(bodyBytes, &jsonResp); err != nil {
		return fmt.Errorf("response is not valid JSON: %w\nResponse body: %s", err, body)
	}

	if jsonResp.Response != "success" {
		return fmt.Errorf("friendly request failed: response=%q\nResponse body: %s", jsonResp.Response, body)
	}

	return nil
}

// GetMyPlayers scrapes players from "/?p=players" page returns slice of players or error.
// Login() must be called before. For a hockey client use GetMyHockeyPlayers.
func (mz *MZHTTPClient) GetMyPlayers() ([]FootballPlayer, error) {
	doc, err := mz.getPlayersPage()
	if err != nil {
		return nil, err
	}
	playersContainerDiv := findByTagAndAttr(doc, atom.Div, "id", "players_container")
	var players []FootballPlayer

	if playersContainerDiv == nil || playersContainerDiv.FirstChild == nil {
		return nil, nil
	}

	players, err = scrapePlayers(playersContainerDiv.FirstChild, players)
	if err != nil {
		return nil, err
	} else {
		return players, nil
	}
}

// getPlayersPage fetches and parses "/?p=players" for the client's sport.
func (mz *MZHTTPClient) getPlayersPage() (*html.Node, error) {
	if mz.phpsessid == "" {
		return nil, errors.New("not logged in: PHPSESSID is empty")
	}

	req, err := http.NewRequest("GET", mz.baseURL+"?p=players", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for my team: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := mz.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch my team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch my team: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	return doc, nil
}

func scrapePlayers(n *html.Node, players []FootballPlayer) (_ []FootballPlayer, err error) {
	// The scraper assumes an exact page layout and dereferences freely; turn a
	// layout change into an error instead of a crash until it is hardened.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("players page layout not recognised: %v", r)
		}
	}()

	name := findByTagAndAttr(n, atom.Span, "class", "player_name").FirstChild.Data // TODO err handling
	id, _ := strconv.Atoi(strings.TrimSpace(findByTagAndAttr(n, atom.Span, "class", "player_id_span").FirstChild.Data))
	age, _ := strconv.Atoi(findByData(n, "Age: ").NextSibling.FirstChild.Data)                              // TODO err handling
	foot := findByData(n, "Foot: ").NextSibling.FirstChild.Data                                             // TODO err handling
	birthSeason, _ := strconv.Atoi(findByData(n, "Birth Season: ").NextSibling.FirstChild.Data)             // TODO err handling
	height, _ := strconv.Atoi(strings.Split(findByData(n, "Height: ").NextSibling.FirstChild.Data, " ")[0]) // TODO err handling
	weight, _ := strconv.Atoi(strings.Split(findByData(n, "Weight: ").NextSibling.FirstChild.Data, " ")[0]) // TODO err handling
	v := strings.Split(findByData(n, "Value: ").NextSibling.FirstChild.Data, " ")
	value, _ := strconv.Atoi(strings.ReplaceAll(v[0], "\u00a0", ""))
	// Search only this player's box: findByData(n, ...) also walks n's next
	// siblings and would find "Youth Player" in a later player's box.
	youthPlayer := findByData(n.FirstChild, "Youth Player")
	var s []string
	var salary int
	var salaryCurrency string
	if youthPlayer == nil {
		s = strings.Split(findByData(n, "Salary: ").NextSibling.FirstChild.Data, " ")
		salary, _ = strconv.Atoi(strings.ReplaceAll(s[0], "\u00a0", ""))
		salaryCurrency = s[1]
	}

	skills := findAllByTagAndAttr(findByTagAndAttr(n, atom.Div, "class", "pp-3-2"), atom.Td, "class", "skillval", []*html.Node{})

	speed := extractSkillValue(skills[0])
	stamina := extractSkillValue(skills[1])
	intelligence := extractSkillValue(skills[2])
	passing := extractSkillValue(skills[3])
	shooting := extractSkillValue(skills[4])
	heading := extractSkillValue(skills[5])
	keeping := extractSkillValue(skills[6])
	ballControl := extractSkillValue(skills[7])
	tackling := extractSkillValue(skills[8])
	aerialPassing := extractSkillValue(skills[9])
	setPlays := extractSkillValue(skills[10])
	experience := extractSkillValue(skills[11])
	form := extractSkillValue(skills[12])
	p := FootballPlayer{
		ID:             id,
		Name:           name,
		Age:            age,
		Foot:           foot,
		BirthSeason:    birthSeason,
		Height:         height,
		Weight:         weight,
		Value:          value,
		ValueCurrency:  v[1],
		Salary:         salary,
		SalaryCurrency: salaryCurrency,
		Speed:          speed,
		Stamina:        stamina,
		Intelligence:   intelligence,
		Passing:        passing,
		Shooting:       shooting,
		Heading:        heading,
		Keeping:        keeping,
		BallControl:    ballControl,
		Tackling:       tackling,
		AerialPassing:  aerialPassing,
		SetPlays:       setPlays,
		Experience:     experience,
		Form:           form,
	}

	players = append(players, p)

	if n.NextSibling != nil {
		return scrapePlayers(n.NextSibling, players)
	}
	return players, nil
}

func extractSkillValue(n *html.Node) float32 {
	intValue, _ := strconv.Atoi(n.FirstChild.NextSibling.FirstChild.Data)
	floatValue := 0.0
	switch findByTagAndAttr(n, atom.Div, "class", "skill_exact_bar").Attr[1].Val {
	case "width: 0px;":
		floatValue = 0.0
	case "width: 2px;":
		floatValue = 0.25
	case "width: 4px;":
		floatValue = 0.50
	case "width: 6px;":
		floatValue = 0.75
	}
	return float32(intValue) + float32(floatValue)
}

func findByTagAndAttr(n *html.Node, a atom.Atom, attrKey string, attrVal string) *html.Node {

	if n.DataAtom == a {
		for _, a := range n.Attr {
			if a.Key == attrKey && a.Val == attrVal {
				return n
			}
		}
	}

	if n.FirstChild != nil {
		node := findByTagAndAttr(n.FirstChild, a, attrKey, attrVal)
		if node != nil {
			return node
		}
	}

	if n.NextSibling != nil {
		node := findByTagAndAttr(n.NextSibling, a, attrKey, attrVal)
		if node != nil {
			return node
		}
	}

	return nil
}

func findAllByTagAndAttr(n *html.Node, a atom.Atom, attrKey string, attrVal string, results []*html.Node) []*html.Node {

	if n.DataAtom == a {
		for _, a := range n.Attr {
			if a.Key == attrKey && a.Val == attrVal {
				results = append(results, n)
			}
		}
	}

	if n.FirstChild != nil {
		node := findAllByTagAndAttr(n.FirstChild, a, attrKey, attrVal, []*html.Node{})
		results = append(results, node...)
	}

	if n.NextSibling != nil {
		node := findAllByTagAndAttr(n.NextSibling, a, attrKey, attrVal, []*html.Node{})
		results = append(results, node...)
	}

	return results
}

func findByData(n *html.Node, data string) *html.Node {

	if n.Data == data {
		return n
	}

	if n.FirstChild != nil {
		node := findByData(n.FirstChild, data)
		if node != nil {
			return node
		}
	}

	if n.NextSibling != nil {
		node := findByData(n.NextSibling, data)
		if node != nil {
			return node
		}
	}

	return nil
}
