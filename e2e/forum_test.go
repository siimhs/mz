package e2e

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type forumUser struct {
	ID           int
	Name         string
	Country      string
	FederationID int
	Roles        []string
}

type forumTopic struct {
	ID, ForumID    int
	Title          string
	TitleTruncated bool
	Pinned         bool
	Views, Replies int
	Author         forumUser
	LastPostID     int
	LastPostOffset int
	LastPostAt     time.Time
}

type forumPost struct {
	ID, TopicID, ForumID int
	Title                string
	Author               forumUser
	TeamID               int
	PostedAt             time.Time
	Text, HTML           string
	Quotes               []struct{ Author, Text string }
}

func decode[T any](t *testing.T, res *Result) T {
	t.Helper()
	var v T
	if err := json.Unmarshal([]byte(res.Stdout), &v); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, res.Stdout)
	}
	return v
}

func tallinn(s string) time.Time {
	loc, _ := time.LoadLocation(defaultTZ)
	t, err := time.ParseInLocation("2006-01-02 15:04", s, loc)
	if err != nil {
		panic(err)
	}
	return t
}

// syncedForum syncs forum 10, served as a one-page list holding only topic 13143067.
func syncedForum(t *testing.T) *MZ {
	t.Helper()
	mz := Start(t, Login(), ForumIndex(), ForumOneTopic(), ForumTopicQuotes())
	res := mz.Run("forum", "sync", "--forum", "10", "--delay", "0")
	res.WantExit(t, 0)
	res.WantStderrContains(t, "1 topics", "5 posts")
	return mz
}

func TestForumSync(t *testing.T) {
	mz := syncedForum(t)
	if n := len(mz.Find("GET", topicURL+"*&offset=*")); n != 1 {
		t.Errorf("%d topic page requests, want 1", n)
	}

	// Nothing new on MZ: only the forum index is read.
	before := len(mz.Requests())
	res := mz.Run("forum", "sync", "--forum", "10", "--delay", "0")
	res.WantExit(t, 0)
	res.WantStderrContains(t, "up to date")
	var after []string
	for _, r := range mz.Requests()[before:] {
		if r.Method == "GET" && r.URL != "/" && r.URL != "/?p=clubhouse" && r.URL != "/?p=players" {
			after = append(after, r.URL)
		}
	}
	if len(after) != 1 || after[0] != "/?forum_language=en&p=forum&sport=soccer&sub=forums" {
		t.Errorf("second sync requested %q, want only the forum index", after)
	}
}

func TestForumList(t *testing.T) {
	mz := syncedForum(t)

	res := mz.Run("forum", "list")

	res.WantExit(t, 0)
	forums := decode[[]struct {
		ID                          int
		Name                        string
		TotalPosts, TotalTopics     int
		LatestTopicID, LatestPostID int
		LatestPostBy                string
		LatestPostAt                time.Time
	}](t, res)
	if len(forums) != 9 {
		t.Errorf("%d forums, the page has 9", len(forums))
	}
	for _, f := range forums {
		if f.ID != 10 {
			continue
		}
		if f.Name != "ManagerZone talk" || f.TotalTopics != 6766 || f.LatestTopicID != 13032964 ||
			f.LatestPostID != 46585176 || f.LatestPostBy != "syxiao" ||
			!f.LatestPostAt.Equal(tallinn("2026-11-20 05:45")) { // "Today 05:45"
			t.Errorf("forum 10: %+v", f)
		}
		return
	}
	t.Error("forum 10 missing")
}

func TestForumTopics(t *testing.T) {
	mz := syncedForum(t)

	res := mz.Run("forum", "topics", "10")

	res.WantExit(t, 0)
	got := decode[struct {
		ForumID, Offset, Pages int
		Forum                  string
		Topics                 []forumTopic
	}](t, res)
	if got.ForumID != 10 || got.Pages != 1 || got.Forum != "English » ManagerZone talk" || len(got.Topics) != 1 {
		t.Fatalf("got %+v", got)
	}
	tp := got.Topics[0]
	if tp.ID != 13143067 || tp.Title != "What's happening on the transfer market ?" || tp.Views != 119 ||
		tp.Replies != 4 || tp.LastPostID != 46574111 || !tp.LastPostAt.Equal(tallinn("2026-10-12 13:52")) ||
		tp.Author.Name != "charef" || len(tp.Author.Roles) != 2 {
		t.Errorf("topic: %+v", tp)
	}

	res = mz.Run("forum", "topics", "10", "--since", "2026-11-01")
	res.WantExit(t, 0)
	if got := decode[struct{ Topics []forumTopic }](t, res); len(got.Topics) != 0 {
		t.Errorf("--since 2026-11-01: %d topics, want 0", len(got.Topics))
	}
}

func TestForumTopic(t *testing.T) {
	mz := syncedForum(t)

	res := mz.Run("forum", "topic", "13143067")

	res.WantExit(t, 0)
	res.WantStdoutGolden(t, "forum-topic-13143067.json")
	got := decode[struct{ Posts []forumPost }](t, res)
	if len(got.Posts) != 5 {
		t.Fatalf("%d posts, want 5", len(got.Posts))
	}
	p := got.Posts[2]
	if p.ID != 46574097 || p.Author.Name != "danielchika" || p.TeamID != 960378 ||
		!p.PostedAt.Equal(tallinn("2026-10-12 12:53")) ||
		p.Text != "I thought they fixed this issue years ago?" ||
		len(p.Quotes) != 1 || p.Quotes[0].Author != "nat-go" {
		t.Errorf("post 3: %+v", p)
	}

	// A topic URL, or MZ's forum-id topic-id pair, work too.
	for _, args := range [][]string{
		{"https://www.managerzone.com/?p=forum&sub=topic&topic_id=13143067&forum_id=10&sport=soccer"},
		{"10", "13143067"},
	} {
		res := mz.Run(append([]string{"forum", "topic"}, args...)...)
		res.WantExit(t, 0)
		res.WantStdoutGolden(t, "forum-topic-13143067.json")
	}

	res = mz.Run("forum", "topic", "13143067", "--offset", "3", "--limit", "1")
	res.WantExit(t, 0)
	if one := decode[struct{ Posts []forumPost }](t, res); len(one.Posts) != 1 || one.Posts[0].ID != got.Posts[3].ID {
		t.Errorf("--offset 3 --limit 1: %+v", one.Posts)
	}

	res = mz.Run("forum", "topic", "999")
	res.WantExit(t, 1)
	res.WantStderrContains(t, "999", "not in the database")
}

func TestForumSearch(t *testing.T) {
	mz := syncedForum(t)

	res := mz.Run("forum", "search", "rubbish youth")

	res.WantExit(t, 0)
	hits := decode[[]struct {
		PostID, TopicID int
		Topic, Author   string
		Snippet         string
	}](t, res)
	if len(hits) != 1 || hits[0].PostID != 46573641 || hits[0].TopicID != 13143067 || hits[0].Author != "charef" {
		t.Fatalf("hits: %+v", hits)
	}
	if !strings.Contains(hits[0].Snippet, "[rubbish]") || !strings.Contains(hits[0].Snippet, "[youth]") {
		t.Errorf("snippet: %q", hits[0].Snippet)
	}

	res = mz.Run("forum", "search", `"unclosed`)
	res.WantExit(t, 1)
	res.WantStderrContains(t, "search query")
}

func TestForumSQL(t *testing.T) {
	mz := syncedForum(t)

	res := mz.Run("sql", "SELECT topic_id, count(*) AS posts FROM posts GROUP BY topic_id")

	res.WantExit(t, 0)
	res.WantStdout(t, "[\n  {\"topic_id\": 13143067, \"posts\": 5}\n]\n")

	res = mz.Run("sql", "DELETE FROM posts")
	res.WantExit(t, 1)
}

func TestForumReadBeforeSync(t *testing.T) {
	mz := Start(t)

	res := mz.Run("forum", "topics", "10")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "mz forum sync")
	mz.WantRequests()
}
