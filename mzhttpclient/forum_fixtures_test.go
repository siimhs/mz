package mzhttpclient

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Recorded forum pages shared with the e2e tests (recorded 2026-11-20 23:07, a Friday).
func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "e2e", "testdata", "mz", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func fixtureNow(t *testing.T) time.Time {
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Skip(err)
	}
	old := time.Local
	time.Local = loc
	t.Cleanup(func() { time.Local = old })
	return time.Date(2026, 11, 20, 23, 7, 0, 0, loc)
}

func TestParseForumTopicsFixture(t *testing.T) {
	now := fixtureNow(t)
	p, err := ParseForumTopics(fixture(t, "forum-topics-10.html"), now)
	if err != nil {
		t.Fatal(err)
	}
	if p.ForumID != 10 || p.Pages != 135 || p.Forum != "English » ManagerZone talk" || len(p.Topics) != 50 {
		t.Errorf("header: id %d pages %d %q, %d topics", p.ForumID, p.Pages, p.Forum, len(p.Topics))
	}
	pinned := 0
	for i, tp := range p.Topics {
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
	first := p.Topics[0]
	if first.ID != 13032964 || first.Title != "Stats Xente User Script" || first.Views != 5750 ||
		first.Replies != 257 || first.LastPostOffset != 250 ||
		first.Author.Name != "xente" || first.Author.Country != "es" || first.Author.FederationID != 2155 ||
		len(first.Author.Roles) != 1 || first.Author.Roles[0] != "MA" {
		t.Errorf("first topic: %+v", first)
	}

	p2, err := ParseForumTopics(fixture(t, "forum-topics-10-offset-50.html"), now)
	if err != nil || len(p2.Topics) != 50 || p2.Topics[0].Pinned {
		t.Errorf("page 2: %d topics, err %v", len(p2.Topics), err)
	}
}

// Pinned topic 13032964 has 258 posts on 6 pages; each page repeats the opening post on top.
func TestParseForumTopicPinnedFixture(t *testing.T) {
	now := fixtureNow(t)
	for name, want := range map[string]int{
		"forum-topic-13032964-offset-200.html": 1 + 50,
		"forum-topic-13032964-offset-250.html": 1 + 7,
	} {
		p, err := ParseForumTopic(fixture(t, name), now)
		if err != nil {
			t.Fatal(err)
		}
		if p.ForumID != 10 || p.Pages != 6 || p.Title != "Stats Xente User Script" || len(p.Posts) != want {
			t.Errorf("%s: forum %d pages %d %q, %d posts (want %d)", name, p.ForumID, p.Pages, p.Title, len(p.Posts), want)
		}
		for i, post := range p.Posts {
			if i > 0 && post.ID < p.Posts[i-1].ID {
				t.Errorf("%s: post %d after %d: want oldest first", name, post.ID, p.Posts[i-1].ID)
			}
			if post.Author.ID == 0 || post.PostedAt.IsZero() || post.HTML == "" {
				t.Errorf("%s: incomplete post: %+v", name, post)
			}
		}
	}
}
