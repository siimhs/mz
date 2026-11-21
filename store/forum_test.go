package store

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
)

func TestTopicRoundTripAndSearch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mz.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 10, 12, 13, 52, 0, 0, time.UTC)
	op := mzhttpclient.User{ID: 7, Name: "Charef", Country: "tn", FederationID: 4491, Roles: []string{"FA", "NC"}}
	topic := mzhttpclient.ForumTopic{ID: 100, ForumID: 10, Title: "Training..", TitleTruncated: true,
		Replies: 1, Author: op, LastPostID: 2, LastPostAt: at}
	posts := []mzhttpclient.ForumPost{
		{ID: 1, Title: "Training camps", Author: op, TeamID: 5, PostedAt: at.Add(-time.Hour),
			Text: "Are training camps worth the money?", HTML: "Are training camps worth the money?"},
		{ID: 2, Title: "Re: Training camps", Author: mzhttpclient.User{ID: 8, Name: "nat"}, PostedAt: at,
			Text: "Only for youngsters.", HTML: "<quote>..</quote>Only for youngsters.",
			Quotes: []mzhttpclient.Quote{{Author: "Charef", Text: "Are training camps worth the money?"}}},
	}
	if err := db.SaveForums(mzhttpclient.Football, "en", []mzhttpclient.Forum{{ID: 10, Name: "ManagerZone talk", LatestPostID: 2}}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveTopic(topic, "Training camps", posts, 0, at); err != nil {
		t.Fatal(err)
	}
	if err := db.MarkForumSynced(mzhttpclient.Football, 10, "English » ManagerZone talk", 2, at); err != nil {
		t.Fatal(err)
	}
	// Saving again (a re-sync) updates in place, including the FTS index.
	posts[1].Text = "Only for youngsters, and only before 21."
	if err := db.SaveTopic(topic, "Training camps", posts, 0, at); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db, err = OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	p, err := db.Topic(100, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if p.Forum != "English » ManagerZone talk" || p.Title != "Training camps" || len(p.Posts) != 2 {
		t.Fatalf("topic: %+v", p)
	}
	if got := p.Posts[0]; got.Author.Name != "Charef" || len(got.Author.Roles) != 2 || got.TeamID != 5 ||
		!got.PostedAt.Equal(at.Add(-time.Hour)) {
		t.Errorf("post 1: %+v", got)
	}
	if got := p.Posts[1]; len(got.Quotes) != 1 || got.Quotes[0].Author != "Charef" {
		t.Errorf("post 2 quotes: %+v", got.Quotes)
	}

	tp, err := db.Topics(TopicsQuery{ForumID: 10, Limit: 50})
	if err != nil || len(tp.Topics) != 1 || tp.Topics[0].Title != "Training camps" || tp.Topics[0].TitleTruncated {
		t.Errorf("topics: %+v %v", tp, err)
	}
	if _, err := db.Topics(TopicsQuery{ForumID: 99}); err == nil {
		t.Error("unknown forum: want an error")
	}

	// "camp" matches "camps" (porter); quotes are not indexed, so post 2 only
	// matches through its title.
	hits, err := db.Search(SearchQuery{Match: "camp money", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].PostID != 1 || hits[0].Author != "Charef" || hits[0].Topic != "Training camps" {
		t.Errorf("hits: %+v", hits)
	}
	if hits, _ := db.Search(SearchQuery{Match: "youngsters 21"}); len(hits) != 1 || hits[0].PostID != 2 {
		t.Errorf("updated text not indexed: %+v", hits)
	}
	if hits, _ := db.Search(SearchQuery{Match: "training", Author: "NAT"}); len(hits) != 1 || hits[0].PostID != 2 {
		t.Errorf("author filter: %+v", hits)
	}
	if _, err := db.Search(SearchQuery{Match: `"unclosed`}); err == nil {
		t.Error("bad query: want an error")
	}
	if _, err := db.Exec(`DELETE FROM posts`); err == nil {
		t.Error("read-only database accepted a write")
	}
}

func TestForumsBySport(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "mz.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	at := time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)
	for _, f := range []struct {
		sport string
		id    int
	}{{mzhttpclient.Football, 10}, {mzhttpclient.Hockey, 138}} {
		if err := db.SaveForums(f.sport, "en", []mzhttpclient.Forum{{ID: f.id, Name: "ManagerZone talk"}}); err != nil {
			t.Fatal(err)
		}
		topic := mzhttpclient.ForumTopic{ID: f.id * 100, ForumID: f.id, Title: "Tactics"}
		post := mzhttpclient.ForumPost{ID: f.id * 1000, Text: "which tactic?", HTML: "which tactic?", PostedAt: at}
		if err := db.SaveTopic(topic, "", []mzhttpclient.ForumPost{post}, 0, at); err != nil {
			t.Fatal(err)
		}
	}
	// A forum synced by id without being in the index gets its sport too
	// (it has no lang or position, so it lists first).
	if err := db.MarkForumSynced(mzhttpclient.Hockey, 999, "English » Other", 0, at); err != nil {
		t.Fatal(err)
	}

	for sport, want := range map[string]string{mzhttpclient.Football: "[10]", mzhttpclient.Hockey: "[999 138]"} {
		fs, err := db.Forums(sport)
		if err != nil {
			t.Fatal(err)
		}
		var ids []int
		for _, f := range fs {
			ids = append(ids, f.ID)
		}
		if fmt.Sprint(ids) != want {
			t.Errorf("Forums(%s) = %v, want %s", sport, ids, want)
		}
	}
	for sport, want := range map[string]string{mzhttpclient.Football: "[10000]", mzhttpclient.Hockey: "[138000]", "": "[10000 138000]"} {
		hits, err := db.Search(SearchQuery{Match: "tactic", Sport: sport})
		if err != nil {
			t.Fatal(err)
		}
		var ids []int
		for _, h := range hits {
			ids = append(ids, h.PostID)
		}
		sort.Ints(ids)
		if fmt.Sprint(ids) != want {
			t.Errorf("Search sport %q = %v, want %s", sport, ids, want)
		}
	}
}

// A version 1 database (before forums.sport) held only football forums.
func TestMigrateV1(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mz.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`ALTER TABLE forums DROP COLUMN sport;
		INSERT INTO forums (id, name) VALUES (10, 'ManagerZone talk');
		PRAGMA user_version = 1`); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var v int
	var sport string
	if err := db.QueryRow(`PRAGMA user_version`).Scan(&v); err != nil || v != schemaVersion {
		t.Errorf("user_version = %d, %v; want %d", v, err, schemaVersion)
	}
	if err := db.QueryRow(`SELECT sport FROM forums WHERE id = 10`).Scan(&sport); err != nil || sport != mzhttpclient.Football {
		t.Errorf("forum 10 sport = %q, %v; want %q", sport, err, mzhttpclient.Football)
	}
}
