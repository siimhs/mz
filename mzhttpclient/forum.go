package mzhttpclient

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ForumPageSize is how many topics or posts MZ shows per page; pages are
// addressed with offset=0,50,100,...
const ForumPageSize = 50

// Forum is one entry of the forum index (?p=forum&sub=forums).
type Forum struct {
	ID            int
	Name          string
	Description   string
	TotalPosts    int
	TotalTopics   int
	LatestTopicID int
	LatestTopic   string
	LatestPostID  int
	LatestPostBy  string
	LatestPostAt  time.Time
}

// User is a forum author.
type User struct {
	ID           int
	Name         string
	Country      string   `json:",omitempty"` // flag code, e.g. "ee"
	FederationID int      `json:",omitempty"`
	Roles        []string `json:",omitempty"` // tags after the name, e.g. "MA", "Crew"
}

// ForumTopic is one row of a forum's topic list.
type ForumTopic struct {
	ID             int
	ForumID        int
	Title          string
	TitleTruncated bool // the list cuts long titles with ".."; ForumTopicPage.Title is the full one
	Pinned         bool
	Views          int
	Replies        int // the list's "Posts" column: posts after the first one
	Author         User
	LastPostID     int
	LastPostOffset int // offset of the page holding the last post
	LastPostAt     time.Time
}

// ForumTopicsPage is one page of a forum's topic list.
type ForumTopicsPage struct {
	ForumID int
	Forum   string // e.g. "English » ManagerZone talk"
	Offset  int
	Pages   int // total pages
	Topics  []ForumTopic
}

// ForumPost is one post of a topic.
type ForumPost struct {
	ID       int
	TopicID  int
	ForumID  int
	Title    string
	Author   User
	TeamID   int `json:",omitempty"`
	PostedAt time.Time
	Text     string  // the author's own text as written, without quoted posts
	HTML     string  // rendered body, including quotes
	Quotes   []Quote `json:",omitempty"`
}

// Quote is a quoted post inside a post body.
type Quote struct {
	Author string
	Text   string
}

// ForumTopicPage is one page of a topic.
type ForumTopicPage struct {
	TopicID int
	ForumID int
	Forum   string
	Title   string
	Offset  int
	Pages   int
	Posts   []ForumPost
}

// ForumPath builds a forum URL path for sport (Football or Hockey): sub is
// "forums", "topics" or "topic". A forum's pages are empty unless sport is the
// forum's own.
func ForumPath(sport, sub string, params url.Values) string {
	q := url.Values{"p": {"forum"}, "sub": {sub}, "sport": {sport}}
	for k, v := range params {
		q[k] = v
	}
	return "?" + q.Encode()
}

// GetForums returns the forum index. lang ("en", "se", ...) picks the
// language section; empty uses the site's default for the account.
func (mz *MZHTTPClient) GetForums(lang string, now time.Time) ([]Forum, error) {
	params := url.Values{}
	if lang != "" {
		params.Set("forum_language", lang)
	}
	body, err := mz.Get(ForumPath(mz.sport, "forums", params))
	if err != nil {
		return nil, err
	}
	return ParseForums(body, now)
}

// GetForumTopics returns one page of a forum's topic list.
func (mz *MZHTTPClient) GetForumTopics(forumID, offset int, now time.Time) (ForumTopicsPage, error) {
	body, err := mz.Get(ForumPath(mz.sport, "topics", url.Values{
		"forum_id": {strconv.Itoa(forumID)}, "offset": {strconv.Itoa(offset)}}))
	if err != nil {
		return ForumTopicsPage{}, err
	}
	p, err := ParseForumTopics(body, now)
	if err != nil {
		return p, fmt.Errorf("forum %d offset %d: %w", forumID, offset, err)
	}
	p.ForumID, p.Offset = forumID, offset
	return p, nil
}

// GetForumTopic returns one page of a topic. MZ shows no posts unless
// forumID is the topic's forum.
func (mz *MZHTTPClient) GetForumTopic(forumID, topicID, offset int, now time.Time) (ForumTopicPage, error) {
	body, err := mz.Get(ForumPath(mz.sport, "topic", url.Values{"forum_id": {strconv.Itoa(forumID)},
		"topic_id": {strconv.Itoa(topicID)}, "offset": {strconv.Itoa(offset)}}))
	if err != nil {
		return ForumTopicPage{}, err
	}
	p, err := ParseForumTopic(body, now)
	if err != nil {
		return p, fmt.Errorf("topic %d offset %d: %w", topicID, offset, err)
	}
	p.TopicID, p.Offset = topicID, offset
	for i := range p.Posts {
		p.Posts[i].TopicID = topicID
	}
	return p, nil
}

// ParseForums parses the forum index page.
func ParseForums(body []byte, now time.Time) ([]Forum, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var out []Forum
	for _, dl := range findAll(doc, func(n *html.Node) bool { return n.DataAtom == atom.Dl && hasClass(n, "forums-list") }) {
		var f Forum
		if h3 := find(dl, isTag(atom.H3)); h3 != nil {
			if a := find(h3, isTag(atom.A)); a != nil {
				f.ID = queryInt(attr(a, "href"), "forum_id")
				f.Name = clean(text(a))
			}
		}
		if f.ID == 0 {
			continue
		}
		if p := find(dl, withClass("forum-description")); p != nil {
			f.Description = clean(text(p))
		}
		if p := find(dl, withClass("forum-in-list-stats")); p != nil {
			s := clean(text(p))
			f.TotalPosts = atoiAfter(s, "Total posts:")
			f.TotalTopics = atoiAfter(s, "Total topics:")
		}
		if p := find(dl, withClass("forum-latest-updated-topic")); p != nil {
			if a := find(p, isTag(atom.A)); a != nil {
				f.LatestTopicID = queryInt(attr(a, "href"), "topic_id")
				f.LatestPostID = fragmentInt(attr(a, "href"))
				f.LatestTopic = clean(text(a))
				// " - author, " text node after the link
				if t := a.NextSibling; t != nil && t.Type == html.TextNode {
					f.LatestPostBy = strings.TrimSuffix(strings.TrimPrefix(clean(t.Data), "- "), ",")
				}
			}
			if s := find(p, isTag(atom.Span)); s != nil {
				f.LatestPostAt, _ = ParseForumTime(clean(text(s)), now)
			}
		}
		out = append(out, f)
	}
	if len(out) == 0 {
		return nil, errors.New("forum index: no forums found (layout changed?)")
	}
	return out, nil
}

// ParseForumTopics parses one page of a forum's topic list.
func ParseForumTopics(body []byte, now time.Time) (ForumTopicsPage, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return ForumTopicsPage{}, err
	}
	var p ForumTopicsPage
	p.ForumID, p.Forum = forumTitle(doc)
	list := find(doc, withID("topics-list"))
	if list == nil {
		return p, errors.New("topic list not found (layout changed or no such forum)")
	}
	for row := list.FirstChild; row != nil; row = row.NextSibling {
		if row.DataAtom != atom.Dd {
			continue
		}
		c := find(row, withClass("topics-col-title"))
		if c == nil {
			continue // an ad between rows
		}
		var t ForumTopic
		t.ForumID = p.ForumID
		if a := find(c, isTag(atom.A)); a != nil {
			t.ID = queryInt(attr(a, "href"), "topic_id")
			t.Title = clean(text(a))
			t.TitleTruncated = strings.HasSuffix(t.Title, "..")
		}
		if t.ID == 0 {
			return p, errors.New("topic row without a topic link (layout changed?)")
		}
		if c := find(row, withClass("topics-col-counter")); c != nil {
			s := clean(text(c))
			if _, s, ok := strings.Cut(s, ":"); ok {
				v, r, _ := strings.Cut(s, "/")
				t.Views, _ = strconv.Atoi(strings.TrimSpace(v))
				t.Replies, _ = strconv.Atoi(strings.TrimSpace(r))
			}
		}
		if c := find(row, withClass("topics-col-author")); c != nil {
			t.Author = parseUser(c)
		}
		if c := find(row, withClass("topics-col-date")); c != nil {
			if a := find(c, isTag(atom.A)); a != nil {
				href := attr(a, "href")
				t.LastPostID = fragmentInt(href)
				t.LastPostOffset = queryInt(href, "offset")
				if t.LastPostAt, err = ParseForumTime(clean(text(a)), now); err != nil {
					return p, fmt.Errorf("topic %d: %w", t.ID, err)
				}
			}
		}
		if c := find(row, withClass("topics-col-flags")); c != nil {
			t.Pinned = find(c, withClass("fa-thumb-tack")) != nil
		}
		p.Topics = append(p.Topics, t)
	}
	p.Pages = pageCount(doc)
	return p, nil
}

// ParseForumTopic parses one page of a topic.
func ParseForumTopic(body []byte, now time.Time) (ForumTopicPage, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return ForumTopicPage{}, err
	}
	var p ForumTopicPage
	p.ForumID, p.Forum = forumTitle(doc)
	seen := map[int]bool{}
	for _, dl := range findAll(doc, func(n *html.Node) bool { return n.DataAtom == atom.Dl && hasClass(n, "forum_body") }) {
		post := ForumPost{ForumID: p.ForumID}
		post.ID, _ = strconv.Atoi(attr(dl, "id"))
		if post.ID == 0 {
			return p, errors.New("post without an id (layout changed?)")
		}
		// Pinned topics repeat the opening post at the top of every page.
		if seen[post.ID] {
			continue
		}
		seen[post.ID] = true
		if h := find(dl, withClass("post-title")); h != nil {
			post.Title = clean(text(h))
		}
		if h := find(dl, withClass("post-author")); h != nil {
			post.Author = parseUser(h)
			var date string
			for _, s := range findAll(h, isTag(atom.Span)) {
				if strings.Contains(attr(s, "style"), "nowrap") {
					date = clean(text(s))
				}
			}
			if post.PostedAt, err = ParseForumTime(date, now); err != nil {
				return p, fmt.Errorf("post %d: %w", post.ID, err)
			}
		}
		if b := find(dl, withClass("forum-post-badges")); b != nil {
			for _, img := range findAll(b, isTag(atom.Img)) {
				if id := queryInt(attr(img, "src"), "team_id"); id != 0 {
					post.TeamID = id
				}
			}
			if post.Author.FederationID == 0 {
				if a := find(b, func(n *html.Node) bool { return n.DataAtom == atom.A && queryInt(attr(n, "href"), "fid") != 0 }); a != nil {
					post.Author.FederationID = queryInt(attr(a, "href"), "fid")
				}
			}
		}
		content := find(dl, withClass("forum-post-content"))
		if content == nil {
			return p, fmt.Errorf("post %d: no content (layout changed?)", post.ID)
		}
		var rendered bytes.Buffer
		for c := content.FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom == atom.Textarea {
				post.Text = text(c)
				continue
			}
			if c.DataAtom == atom.Div && hasClass(c, "forum_quoteBox") {
				q := Quote{}
				if s := find(c, isTag(atom.Strong)); s != nil {
					q.Author = strings.TrimSuffix(clean(text(s)), " wrote:")
				}
				if ct := find(c, isTag(atom.Cite)); ct != nil {
					q.Text = strings.TrimSpace(text(ct))
				}
				post.Quotes = append(post.Quotes, q)
			}
			html.Render(&rendered, c)
		}
		post.HTML = strings.TrimSpace(rendered.String())
		p.Posts = append(p.Posts, post)
	}
	if len(p.Posts) == 0 {
		return p, errors.New("no posts found (layout changed, no such topic, or wrong forum id)")
	}
	p.Title = strings.TrimPrefix(p.Posts[0].Title, "Re: ")
	p.Pages = pageCount(doc)
	return p, nil
}

var reForumDate = regexp.MustCompile(`^(\S+) (\d{1,2}):(\d{2})$`)

// ParseForumTime parses forum dates: "12-10-2026 13:52", and for the last
// week "Today 05:45", "Yesterday 18:10", "Monday 09:00" (the latest past
// Monday). Times are in the account's timezone, taken as time.Local.
func ParseForumTime(s string, now time.Time) (time.Time, error) {
	if t, err := time.ParseInLocation("02-01-2006 15:04", s, time.Local); err == nil {
		return t, nil
	}
	m := reForumDate.FindStringSubmatch(s)
	if m == nil {
		return time.Time{}, fmt.Errorf("unknown forum date %q", s)
	}
	h, _ := strconv.Atoi(m[2])
	mi, _ := strconv.Atoi(m[3])
	now = now.In(time.Local)
	day := time.Date(now.Year(), now.Month(), now.Day(), h, mi, 0, 0, time.Local)
	switch m[1] {
	case "Today":
		return day, nil
	case "Yesterday":
		return day.AddDate(0, 0, -1), nil
	}
	for wd := time.Sunday; wd <= time.Saturday; wd++ {
		if m[1] == wd.String() {
			back := (int(now.Weekday()) - int(wd) + 7) % 7
			if back == 0 {
				back = 7
			}
			return day.AddDate(0, 0, -back), nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown forum date %q", s)
}

// forumTitle reads forum id and name from the "Football » English » ManagerZone talk" heading.
func forumTitle(doc *html.Node) (int, string) {
	h := find(doc, withID("forum-title"))
	if h == nil {
		return 0, ""
	}
	for _, a := range findAll(h, isTag(atom.A)) {
		if id := queryInt(attr(a, "href"), "forum_id"); id != 0 {
			return id, clean(text(a))
		}
	}
	return 0, ""
}

// parseUser reads an author from a node holding MZ's user-link markup.
func parseUser(n *html.Node) User {
	var u User
	if w := find(n, withClass("user-link-wrapper")); w != nil {
		n = w.Parent // role tags and federation sit next to the wrapper
		for _, a := range findAll(w, isTag(atom.A)) {
			if id := queryInt(attr(a, "href"), "uid"); id != 0 {
				u.ID, u.Name = id, clean(text(a))
			}
		}
		for _, img := range findAll(w, isTag(atom.Img)) {
			if src := attr(img, "src"); strings.Contains(src, "/flags/") {
				u.Country = strings.TrimSuffix(src[strings.LastIndex(src, "/")+1:], ".png")
			}
		}
	}
	for _, a := range findAll(n, isTag(atom.A)) {
		if strings.Contains(attr(a, "href"), "p=federations") {
			u.FederationID = queryInt(attr(a, "href"), "fid")
		}
	}
	for _, s := range findAll(n, withClass("author-string")) {
		u.Roles = append(u.Roles, strings.Trim(clean(text(s)), "()"))
	}
	return u
}

// pageCount reads the total number of pages from the "Page 1 2 3 ..." navigation.
func pageCount(doc *html.Node) int {
	pages := 1
	for _, nav := range findAll(doc, func(n *html.Node) bool { return n.DataAtom == atom.Dl && hasClass(n, "forum-list-nav") }) {
		if b := find(nav, isTag(atom.B)); b != nil {
			if v, err := strconv.Atoi(clean(text(b))); err == nil && v > pages {
				pages = v
			}
		}
		for _, a := range findAll(nav, isTag(atom.A)) {
			href := attr(a, "href")
			if !strings.Contains(href, "offset=") {
				continue
			}
			if v := queryInt(href, "offset")/ForumPageSize + 1; v > pages {
				pages = v
			}
		}
	}
	return pages
}

// --- small html helpers ---

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func hasClass(n *html.Node, c string) bool {
	if n.Type != html.ElementNode {
		return false
	}
	for _, f := range strings.Fields(attr(n, "class")) {
		if f == c {
			return true
		}
	}
	return false
}

func isTag(a atom.Atom) func(*html.Node) bool {
	return func(n *html.Node) bool { return n.DataAtom == a }
}
func withClass(c string) func(*html.Node) bool {
	return func(n *html.Node) bool { return hasClass(n, c) }
}
func withID(id string) func(*html.Node) bool {
	return func(n *html.Node) bool { return n.Type == html.ElementNode && attr(n, "id") == id }
}

// find returns the first descendant of n matching pred.
func find(n *html.Node, pred func(*html.Node) bool) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if pred(c) {
			return c
		}
		if r := find(c, pred); r != nil {
			return r
		}
	}
	return nil
}

// findAll returns the descendants of n matching pred, not looking inside matches.
func findAll(n *html.Node, pred func(*html.Node) bool) []*html.Node {
	var out []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if pred(c) {
			out = append(out, c)
			continue
		}
		out = append(out, findAll(c, pred)...)
	}
	return out
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(text(c))
	}
	return b.String()
}

// clean collapses whitespace.
func clean(s string) string { return strings.Join(strings.Fields(s), " ") }

// queryInt reads an integer query parameter from a link like "/?p=forum&topic_id=1#2".
func queryInt(href, key string) int {
	u, err := url.Parse(href)
	if err != nil {
		return 0
	}
	v, _ := strconv.Atoi(u.Query().Get(key))
	return v
}

func fragmentInt(href string) int {
	_, f, _ := strings.Cut(href, "#")
	v, _ := strconv.Atoi(f)
	return v
}

func atoiAfter(s, label string) int {
	_, rest, ok := strings.Cut(s, label)
	if !ok {
		return 0
	}
	f := strings.Fields(rest)
	if len(f) == 0 {
		return 0
	}
	v, _ := strconv.Atoi(f[0])
	return v
}
