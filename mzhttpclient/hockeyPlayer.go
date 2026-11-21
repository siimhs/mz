package mzhttpclient

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HockeyPlayer is one player of my hockey squad. Skills are 0..10 in
// quarter pucks (the page shows the fraction as a bar 0/2/4/6 px wide);
// partial skills count in the sim.
type HockeyPlayer struct {
	ID             int
	Name           string
	Age            int
	BirthSeason    int
	Height         int
	Weight         int
	Shoots         string
	Value          int
	ValueCurrency  string
	Salary         int
	SalaryCurrency string

	PlayIntelligence float32
	Power            float32
	Skating          float32
	Passing          float32
	Quickness        float32
	Shooting         float32
	Keeping          float32
	PuckControl      float32
	Checking         float32
	Stamina          float32
	Form             float32
}

// GetMyHockeyPlayers scrapes "/?p=players" of a hockey client (WithSport(Hockey)).
func (mz *MZHTTPClient) GetMyHockeyPlayers() ([]HockeyPlayer, error) {
	doc, err := mz.getPlayersPage()
	if err != nil {
		return nil, err
	}
	return scrapeHockeyPlayers(doc)
}

// hockeySkill maps the page's skill label to the field it fills.
var hockeySkill = map[string]func(*HockeyPlayer) *float32{
	"Play Intelligence": func(p *HockeyPlayer) *float32 { return &p.PlayIntelligence },
	"Power":             func(p *HockeyPlayer) *float32 { return &p.Power },
	"Skating":           func(p *HockeyPlayer) *float32 { return &p.Skating },
	"Passing":           func(p *HockeyPlayer) *float32 { return &p.Passing },
	"Quickness":         func(p *HockeyPlayer) *float32 { return &p.Quickness },
	"Shooting":          func(p *HockeyPlayer) *float32 { return &p.Shooting },
	"Keeping":           func(p *HockeyPlayer) *float32 { return &p.Keeping },
	"Puck Control":      func(p *HockeyPlayer) *float32 { return &p.PuckControl },
	"Checking":          func(p *HockeyPlayer) *float32 { return &p.Checking },
	"Stamina":           func(p *HockeyPlayer) *float32 { return &p.Stamina },
	"Form":              func(p *HockeyPlayer) *float32 { return &p.Form },
}

// scrapeHockeyPlayers reads every div.playerContainer of the players page.
// Skills are found by label, not position, so a reordering cannot shift them.
func scrapeHockeyPlayers(doc *html.Node) ([]HockeyPlayer, error) {
	c := findIn(doc, func(n *html.Node) bool { return n.DataAtom == atom.Div && attr(n, "id") == "players_container" })
	if c == nil {
		return nil, errors.New("players page layout not recognised: no players_container")
	}
	var out []HockeyPlayer
	for n := c.FirstChild; n != nil; n = n.NextSibling {
		if n.DataAtom != atom.Div || !hasClass(n, "playerContainer") {
			continue
		}
		p, err := scrapeHockeyPlayer(n)
		if err != nil {
			return nil, fmt.Errorf("player %d: %w", len(out)+1, err)
		}
		out = append(out, p)
	}
	return out, nil
}

func scrapeHockeyPlayer(n *html.Node) (HockeyPlayer, error) {
	var p HockeyPlayer
	name := findIn(n, func(x *html.Node) bool { return x.DataAtom == atom.Span && hasClass(x, "player_name") })
	if name == nil {
		return p, errors.New("no player_name")
	}
	p.Name = text(name)
	id := findIn(n, func(x *html.Node) bool { return x.DataAtom == atom.Span && hasClass(x, "player_id_span") })
	if id == nil {
		return p, errors.New("no player id")
	}
	p.ID, _ = strconv.Atoi(strings.TrimSpace(text(id)))

	// "Age: <strong>25</strong>": the value is the label's next sibling.
	field := func(label string) string {
		l := findIn(n, func(x *html.Node) bool { return x.Type == html.TextNode && strings.TrimSpace(x.Data) == label })
		if l == nil || l.NextSibling == nil {
			return ""
		}
		return strings.TrimSpace(text(l.NextSibling))
	}
	num := func(s string) int { // "760 324 EUR" -> 760324
		f := strings.Fields(strings.ReplaceAll(s, " ", ""))
		if len(f) == 0 {
			return 0
		}
		v, _ := strconv.Atoi(f[0])
		return v
	}
	cur := func(s string) string {
		f := strings.Fields(strings.ReplaceAll(s, " ", ""))
		if len(f) < 2 {
			return ""
		}
		return f[len(f)-1]
	}
	p.Age = num(field("Age:"))
	p.Shoots = field("Shoots:")
	p.Height = num(field("Height:"))
	p.Weight = num(field("Weight:"))
	p.BirthSeason = num(field("Birth Season:"))
	v := field("Value:")
	p.Value, p.ValueCurrency = num(v), cur(v)
	if s := field("Salary:"); s != "" { // youth players have none
		p.Salary, p.SalaryCurrency = num(s), cur(s)
	}

	// The page holds the skills table twice (desktop and responsive); read the first.
	tbl := findIn(n, func(x *html.Node) bool { return x.DataAtom == atom.Table && hasClass(x, "player_skills") })
	if tbl == nil {
		return p, errors.New("no skills table")
	}
	seen := map[string]bool{}
	var walkRows func(*html.Node)
	walkRows = func(x *html.Node) {
		for r := x.FirstChild; r != nil; r = r.NextSibling {
			if r.DataAtom != atom.Tr {
				walkRows(r) // tbody
				continue
			}
			lbl := findIn(r, func(y *html.Node) bool { return y.DataAtom == atom.Span && hasClass(y, "skill_name") })
			val := findIn(r, func(y *html.Node) bool { return y.DataAtom == atom.Td && hasClass(y, "skillval") })
			if lbl == nil || val == nil || lbl.FirstChild == nil {
				continue
			}
			key := text(lbl.FirstChild) // first span; a scout-report <sup> follows
			f, ok := hockeySkill[key]
			if !ok {
				continue
			}
			whole, err := strconv.Atoi(strings.Trim(strings.TrimSpace(text(val)), "()"))
			if err != nil {
				continue
			}
			*f(&p) = float32(whole) + exactBar(r)
			seen[key] = true
		}
	}
	walkRows(tbl)
	for k := range hockeySkill {
		if !seen[k] {
			return p, fmt.Errorf("%s: skill %q not found", p.Name, k)
		}
	}
	return p, nil
}

// exactBar reads the quarter-puck bar: width 0/2/4/6 px -> 0/.25/.5/.75.
func exactBar(row *html.Node) float32 {
	b := findIn(row, func(x *html.Node) bool { return x.DataAtom == atom.Div && hasClass(x, "skill_exact_bar") })
	if b == nil {
		return 0
	}
	var px int
	if _, err := fmt.Sscanf(strings.TrimSpace(attr(b, "style")), "width: %dpx;", &px); err != nil {
		return 0
	}
	return float32(px) / 8
}

// findIn searches n's subtree only (unlike findByTagAndAttr, not its siblings).
func findIn(n *html.Node, match func(*html.Node) bool) *html.Node {
	if match(n) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if f := findIn(c, match); f != nil {
			return f
		}
	}
	return nil
}
