package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/siimhs/mz/store"
)

// fakeForum serves forum 10 from an in-memory model, 50 per page like MZ.
type fakeForum struct {
	topics   map[int]*fakeTopic
	nextPost int
	requests []string
	failOn   string // request that errors once
}

type fakeTopic struct {
	id     int
	pinned bool
	posts  []int
}

const page = mzhttpclient.ForumPageSize

var t0 = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func newFakeForum() *fakeForum { return &fakeForum{topics: map[int]*fakeTopic{}, nextPost: 1000} }

// post adds n posts to a topic, creating it.
func (f *fakeForum) post(topic, n int, pinned bool) {
	t := f.topics[topic]
	if t == nil {
		t = &fakeTopic{id: topic, pinned: pinned}
		f.topics[topic] = t
	}
	for range n {
		f.nextPost++
		t.posts = append(t.posts, f.nextPost)
	}
}

func (f *fakeForum) req(s string) error {
	f.requests = append(f.requests, s)
	if s == f.failOn {
		f.failOn = ""
		return errors.New("connection reset")
	}
	return nil
}

func (f *fakeForum) list() []*fakeTopic {
	var ts []*fakeTopic
	for _, t := range f.topics {
		ts = append(ts, t)
	}
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].pinned != ts[j].pinned {
			return ts[i].pinned
		}
		return last(ts[i].posts) > last(ts[j].posts)
	})
	return ts
}

func last(s []int) int { return s[len(s)-1] }

func postTime(id int) time.Time { return t0.Add(time.Duration(id) * time.Minute) }

func (f *fakeForum) GetForums(lang string, now time.Time) ([]mzhttpclient.Forum, error) {
	if err := f.req("index"); err != nil {
		return nil, err
	}
	latest := 0
	for _, t := range f.topics {
		latest = max(latest, last(t.posts))
	}
	return []mzhttpclient.Forum{{ID: 10, Name: "ManagerZone talk", LatestPostID: latest}}, nil
}

func (f *fakeForum) GetForumTopics(forumID, offset int, now time.Time) (mzhttpclient.ForumTopicsPage, error) {
	if err := f.req(fmt.Sprintf("topics %d", offset)); err != nil {
		return mzhttpclient.ForumTopicsPage{}, err
	}
	ts := f.list()
	p := mzhttpclient.ForumTopicsPage{ForumID: forumID, Forum: "English » ManagerZone talk", Offset: offset,
		Pages: max(1, (len(ts)+page-1)/page)}
	for _, t := range ts[min(offset, len(ts)):min(offset+page, len(ts))] {
		lp := last(t.posts)
		p.Topics = append(p.Topics, mzhttpclient.ForumTopic{ID: t.id, ForumID: forumID, Title: fmt.Sprint("topic ", t.id),
			Pinned: t.pinned, Replies: len(t.posts) - 1, Author: mzhttpclient.User{ID: 1, Name: "op"},
			LastPostID: lp, LastPostOffset: (len(t.posts) - 1) / page * page, LastPostAt: postTime(lp)})
	}
	return p, nil
}

func (f *fakeForum) GetForumTopic(forumID, topicID, offset int, now time.Time) (mzhttpclient.ForumTopicPage, error) {
	if err := f.req(fmt.Sprintf("topic %d %d", topicID, offset)); err != nil {
		return mzhttpclient.ForumTopicPage{}, err
	}
	t := f.topics[topicID]
	p := mzhttpclient.ForumTopicPage{TopicID: topicID, ForumID: forumID, Title: fmt.Sprint("full title ", topicID),
		Offset: offset, Pages: max(1, (len(t.posts)+page-1)/page)}
	ids := t.posts[min(offset, len(t.posts)):min(offset+page, len(t.posts))]
	if t.pinned && offset > 0 && len(ids) > 0 {
		ids = append([]int{t.posts[0]}, ids...)
	}
	for _, id := range ids {
		p.Posts = append(p.Posts, mzhttpclient.ForumPost{ID: id, TopicID: topicID, ForumID: forumID,
			Author: mzhttpclient.User{ID: 2, Name: "poster"}, PostedAt: postTime(id),
			Text: fmt.Sprint("post ", id), HTML: fmt.Sprint("post ", id)})
	}
	return p, nil
}

func testDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "mz.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func runSync(t *testing.T, f *fakeForum, db *store.DB) ([]string, error) {
	t.Helper()
	f.requests = nil
	s := &forumSyncer{src: f, db: db, lang: "en", now: func() time.Time { return t0 }, log: io.Discard}
	_, err := s.sync([]int{10})
	return f.requests, err
}

func wantRequests(t *testing.T, got []string, want ...string) {
	t.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("requests:\n  %q\nwant\n  %q", got, want)
	}
}

// wantStored checks that the database holds exactly the model's posts.
func wantStored(t *testing.T, f *fakeForum, db *store.DB) {
	t.Helper()
	for _, ft := range f.topics {
		p, err := db.Topic(ft.id, 0, 0)
		if err != nil {
			t.Errorf("topic %d: %v", ft.id, err)
			continue
		}
		var got []int
		for _, post := range p.Posts {
			got = append(got, post.ID)
		}
		if fmt.Sprint(got) != fmt.Sprint(ft.posts) {
			t.Errorf("topic %d posts %v, want %v", ft.id, got, ft.posts)
		}
		if p.Title != fmt.Sprint("full title ", ft.id) {
			t.Errorf("topic %d title %q", ft.id, p.Title)
		}
	}
}

func TestForumSyncFirstThenUpToDate(t *testing.T) {
	f := newFakeForum()
	f.post(1, 120, true) // pinned, 3 pages
	for id := 2; id <= 60; id++ {
		f.post(id, 1, false) // 60 topics: a second list page
	}
	db := testDB(t)

	reqs, err := runSync(t, f, db)
	if err != nil {
		t.Fatal(err)
	}
	if len(reqs) != 1+2+59+3 {
		t.Errorf("%d requests, want index + 2 list pages + 59 + 3 topic pages", len(reqs))
	}
	wantStored(t, f, db)

	reqs, err = runSync(t, f, db)
	if err != nil {
		t.Fatal(err)
	}
	wantRequests(t, reqs, "index")
}

func TestForumSyncIncremental(t *testing.T) {
	f := newFakeForum()
	f.post(1, 70, true)
	f.post(2, 60, false)
	for id := 3; id <= 10; id++ {
		f.post(id, 1, false)
	}
	db := testDB(t)
	if _, err := runSync(t, f, db); err != nil {
		t.Fatal(err)
	}

	f.post(2, 45, false) // page 2 fills up and page 3 starts
	f.post(5, 1, false)
	reqs, err := runSync(t, f, db)
	if err != nil {
		t.Fatal(err)
	}
	// List order now: 1 (pinned, unchanged), 5, 2, then 10 (unchanged: stop).
	// Changed topics go oldest activity first; topic 2 resumes on page 2.
	wantRequests(t, reqs, "index", "topics 0", "topic 2 50", "topic 2 100", "topic 5 0")
	wantStored(t, f, db)
}

func TestForumSyncResumesAfterError(t *testing.T) {
	f := newFakeForum()
	for id := 1; id <= 5; id++ {
		f.post(id, 1, false)
	}
	db := testDB(t)
	f.failOn = "topic 3 0"

	if _, err := runSync(t, f, db); err == nil {
		t.Fatal("want the error")
	}
	// 1 and 2 (oldest) were saved; the forum is not marked synced.
	reqs, err := runSync(t, f, db)
	if err != nil {
		t.Fatal(err)
	}
	wantRequests(t, reqs, "index", "topics 0", "topic 3 0", "topic 4 0", "topic 5 0")
	wantStored(t, f, db)
}

func TestForumSyncDeletedPosts(t *testing.T) {
	f := newFakeForum()
	f.post(1, 120, false)
	db := testDB(t)
	if _, err := runSync(t, f, db); err != nil {
		t.Fatal(err)
	}
	// 30 posts deleted, one added: the last post is now on page 2 (offset 50).
	f.topics[1].posts = f.topics[1].posts[30:]
	f.post(1, 1, false)

	reqs, err := runSync(t, f, db)
	if err != nil {
		t.Fatal(err)
	}
	wantRequests(t, reqs, "index", "topics 0", "topic 1 50")
	// Deleted posts stay in the database; the new one is there.
	p, _ := db.Topic(1, 0, 0)
	if len(p.Posts) != 121 || p.Posts[120].ID != last(f.topics[1].posts) {
		t.Errorf("%d posts, last %d", len(p.Posts), p.Posts[len(p.Posts)-1].ID)
	}
}
