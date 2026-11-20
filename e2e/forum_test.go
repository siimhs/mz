package e2e

import (
	"encoding/json"
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

func TestForumList(t *testing.T) {
	mz := Start(t, Login(), ForumIndex())

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

func TestForumTopicsPages(t *testing.T) {
	mz := Start(t, Login(), ForumTopics())

	res := mz.Run("forum", "topics", "10", "--pages", "2", "--delay", "0")

	res.WantExit(t, 0)
	got := decode[struct {
		ForumID, Offset, Pages int
		Forum                  string
		Topics                 []forumTopic
	}](t, res)
	if got.ForumID != 10 || got.Pages != 135 || got.Forum != "English » ManagerZone talk" {
		t.Errorf("header: %+v", got)
	}
	if len(got.Topics) != 100 {
		t.Errorf("%d topics, want 50 per page", len(got.Topics))
	}
	pinned := 0
	for i, tp := range got.Topics {
		if tp.Pinned {
			pinned++
			if i >= 8 {
				t.Errorf("pinned topic %d at position %d", tp.ID, i)
			}
		}
		if tp.ID == 0 || tp.ForumID != 10 || tp.Author.ID == 0 || tp.LastPostID == 0 || tp.LastPostAt.IsZero() {
			t.Errorf("incomplete topic: %+v", tp)
		}
	}
	if pinned != 8 {
		t.Errorf("%d pinned topics, want 8", pinned)
	}
	first := got.Topics[0]
	if first.ID != 13032964 || first.Title != "Stats Xente User Script" || first.Views != 5750 ||
		first.Replies != 257 || first.LastPostOffset != 250 ||
		first.Author.Name != "xente" || first.Author.Country != "es" || first.Author.FederationID != 2155 ||
		len(first.Author.Roles) != 1 || first.Author.Roles[0] != "MA" {
		t.Errorf("first topic: %+v", first)
	}
}

// --since stops paging at the first page holding an older (unpinned) topic,
// so page 2 is never requested (it has no route).
func TestForumTopicsSince(t *testing.T) {
	mz := Start(t, Login(), Route("GET", topicsURL+"0", File("forum-topics-10.html")))

	res := mz.Run("forum", "topics", "10", "--pages", "0", "--since", "2026-10-01", "--delay", "0")

	res.WantExit(t, 0)
	got := decode[struct{ Topics []forumTopic }](t, res)
	if len(got.Topics) == 0 || len(got.Topics) >= 50 {
		t.Fatalf("%d topics", len(got.Topics))
	}
	for _, tp := range got.Topics {
		if tp.LastPostAt.Before(tallinn("2026-10-01 00:00")) {
			t.Errorf("topic %d last post %v is before --since", tp.ID, tp.LastPostAt)
		}
	}
}

func TestForumTopicPages(t *testing.T) {
	mz := Start(t, Login(), ForumTopicPinned())

	res := mz.Run("forum", "topic", "10", "13032964", "--offset", "200", "--delay", "0")

	res.WantExit(t, 0)
	got := decode[struct {
		TopicID, ForumID, Offset, Pages int
		Title                           string
		Posts                           []forumPost
	}](t, res)
	if got.TopicID != 13032964 || got.ForumID != 10 || got.Offset != 200 || got.Pages != 6 ||
		got.Title != "Stats Xente User Script" {
		t.Errorf("header: %+v", got)
	}
	// 258 posts = opening post + 5*50 + 7. Each page repeats the opening post
	// on top; the output has it once.
	if len(got.Posts) != 1+50+7 {
		t.Errorf("%d posts, want 58", len(got.Posts))
	}
	seen := map[int]bool{}
	for i, p := range got.Posts {
		if seen[p.ID] {
			t.Errorf("post %d twice", p.ID)
		}
		seen[p.ID] = true
		if i > 0 && p.ID < got.Posts[i-1].ID {
			t.Errorf("post %d after %d: want oldest first", p.ID, got.Posts[i-1].ID)
		}
		if p.TopicID != 13032964 || p.Author.ID == 0 || p.PostedAt.IsZero() || p.HTML == "" {
			t.Errorf("incomplete post: %+v", p)
		}
	}
}

func TestForumTopicQuotes(t *testing.T) {
	mz := Start(t, Login(), ForumTopicQuotes())

	res := mz.Run("forum", "topic", "https://www.managerzone.com/?p=forum&sub=topic&topic_id=13143067&forum_id=10&sport=soccer")

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
}

// MZ shows no posts for a topic under another forum id, so a bare topic id is refused.
func TestForumTopicNeedsForumID(t *testing.T) {
	mz := Start(t)

	res := mz.Run("forum", "topic", "13143067")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "forum")
	mz.WantRequests()
}
