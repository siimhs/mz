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

// New creates a new MZHTTPClient
func New(baseURL string) *MZHTTPClient {
	return &MZHTTPClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// New creates a new MZHTTPClient
func NewFootballClient() *MZHTTPClient {
	return &MZHTTPClient{
		client:  &http.Client{},
		baseURL: "https://www.managerzone.com/",
		sport:   "soccer",
	}
}

// New creates a new MZHTTPClient
func NewHockeyClient() *MZHTTPClient {
	return &MZHTTPClient{
		client:  &http.Client{},
		baseURL: "https://www.managerzone.com/",
		sport:   "hockey",
	}
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
		"logindata[sport]":         {"soccer"},
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

	return nil
}

// https://www.managerzone.com/ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer
func (mz *MZHTTPClient) GetTeamIDsOnQuickChallengePage() ([]string, error) {
	if mz.phpsessid == "" {
		return nil, errors.New("not logged in: PHPSESSID is empty")
	}

	req, err := http.NewRequest("GET", mz.baseURL+"ajax.php?p=challenge&sub=quick-challenge-template&sport=soccer", nil)
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

func (mz *MZHTTPClient) GetMatches(teamId string) []Match {
	if mz.phpsessid == "" {
		return nil
	}

	req, err := http.NewRequest("GET", mz.baseURL+"xml/team_matchlist.php?team_id="+teamId+"&match_status=1", nil)
	if err != nil {
		return nil
	}

	// req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", mz.phpsessid))
	// req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := mz.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	type xmlTeam struct {
		Field  string `xml:"field,attr"`
		TeamID string `xml:"teamId,attr"`
	}

	type xmlMatch struct {
		ID    string    `xml:"id,attr"`
		Date  string    `xml:"date,attr"`
		Type  string    `xml:"type,attr"`
		Teams []xmlTeam `xml:"Team"`
	}

	type xmlMatchList struct {
		Matches []xmlMatch `xml:"Match"`
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	var result xmlMatchList
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil
	}

	var matches []Match
	for _, m := range result.Matches {

		parsedDate, err := time.Parse("2006-01-02 15:04", m.Date)
		if err != nil {
			return nil
		}

		match := Match{
			Id:         m.ID,
			MatchType:  m.Type,
			Date:       parsedDate,
			HomeTeamId: "", // will be filled next
			AwayTeamId: "",
		}

		for _, t := range m.Teams {
			switch t.Field {
			case "home":
				match.HomeTeamId = t.TeamID
			case "away":
				match.AwayTeamId = t.TeamID
			}
		}

		matches = append(matches, match)
	}

	return matches
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

	d := dateTimeHome.Day()
	if d == 12 {
		println("asdf")
	}

	req, err := http.NewRequest("POST", mz.baseURL+"ajax.php?p=challenge&sub=handle-challenge&action=create&tid="+teamId+"&country_id=123&sport=soccer", strings.NewReader(data.Encode()))
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
// Login() must be called before.
func (mz *MZHTTPClient) GetMyPlayers() ([]FootballPlayer, error) {
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
	playersContainerDiv := findByTagAndAttr(doc, atom.Div, "id", "players_container")
	var players []FootballPlayer

	if playersContainerDiv.FirstChild == nil {
		return nil, nil
	}

	players, err = scrapePlayers(playersContainerDiv.FirstChild, players)
	if err != nil {
		return nil, err
	} else {
		return players, nil
	}
}

func scrapePlayers(n *html.Node, players []FootballPlayer) ([]FootballPlayer, error) {

	name := findByTagAndAttr(n, atom.Span, "class", "player_name").FirstChild.Data                          // TODO err handling
	age, _ := strconv.Atoi(findByData(n, "Age: ").NextSibling.FirstChild.Data)                              // TODO err handling
	foot := findByData(n, "Foot: ").NextSibling.FirstChild.Data                                             // TODO err handling
	birthSeason, _ := strconv.Atoi(findByData(n, "Birth Season: ").NextSibling.FirstChild.Data)             // TODO err handling
	height, _ := strconv.Atoi(strings.Split(findByData(n, "Height: ").NextSibling.FirstChild.Data, " ")[0]) // TODO err handling
	weight, _ := strconv.Atoi(strings.Split(findByData(n, "Weight: ").NextSibling.FirstChild.Data, " ")[0]) // TODO err handling
	v := strings.Split(findByData(n, "Value: ").NextSibling.FirstChild.Data, " ")
	value, _ := strconv.Atoi(strings.ReplaceAll(v[0], "\u00a0", ""))
	youthPlayer := findByData(n, "Youth Player")
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
