package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/siimhs/mz/store"
)

// forumSource is the part of the MZ client a forum sync uses.
type forumSource interface {
	GetForums(lang string, now time.Time) ([]mzhttpclient.Forum, error)
	GetForumTopics(forumID, offset int, now time.Time) (mzhttpclient.ForumTopicsPage, error)
	GetForumTopic(forumID, topicID, offset int, now time.Time) (mzhttpclient.ForumTopicPage, error)
}

type forumSyncer struct {
	src   forumSource
	db    *store.DB
	lang  string // forum index section whose forums are recorded
	full  bool   // re-fetch every topic from its first page
	delay time.Duration
	now   func() time.Time
	log   io.Writer

	requests int
}

type syncStats struct{ topics, posts int }

// wait spaces out requests to MZ.
func (s *forumSyncer) wait() {
	if s.requests > 0 && s.delay > 0 {
		time.Sleep(s.delay)
	}
	s.requests++
}

// sync brings the given forums up to date.
//
// A forum is skipped when its latest post (from the index) is the one recorded
// at its last completed sync. Otherwise the topic list is read newest activity
// first, down to the first non-pinned topic that is stored with the same last
// post; the changed topics are then fetched oldest activity first, each saved
// in its own transaction. If a sync stops halfway, the topics it saved are the
// oldest of the changed ones, so the next sync's walk still reaches the rest.
// A first sync finds nothing stored and so fetches the whole forum.
func (s *forumSyncer) sync(forumIDs []int) (syncStats, error) {
	var total syncStats
	s.wait()
	index, err := s.src.GetForums(s.lang, s.now())
	if err != nil {
		return total, fmt.Errorf("forum index: %w", err)
	}
	if err := s.db.SaveForums(s.lang, index); err != nil {
		return total, err
	}
	latest := map[int]int{}
	for _, f := range index {
		latest[f.ID] = f.LatestPostID
	}
	for _, fid := range forumIDs {
		st, err := s.syncForum(fid, latest[fid])
		total.topics += st.topics
		total.posts += st.posts
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (s *forumSyncer) syncForum(fid, latestPostID int) (syncStats, error) {
	var st syncStats
	if latestPostID != 0 && !s.full {
		synced, err := s.db.ForumSyncedPostID(fid)
		if err != nil {
			return st, err
		}
		if synced == latestPostID {
			fmt.Fprintf(s.log, "forum %d: up to date\n", fid)
			return st, nil
		}
	}
	changed, path, err := s.changedTopics(fid)
	if err != nil {
		return st, err
	}
	fmt.Fprintf(s.log, "forum %d (%s): %d topics to fetch\n", fid, path, len(changed))
	for i := len(changed) - 1; i >= 0; i-- {
		t := changed[i]
		n, err := s.syncTopic(t)
		if err != nil {
			return st, err
		}
		st.topics++
		st.posts += n
		fmt.Fprintf(s.log, "  [%d/%d] topic %d: %d posts fetched  %s\n", len(changed)-i, len(changed), t.ID, n, t.Title)
	}
	if err := s.db.MarkForumSynced(fid, path, latestPostID, s.now()); err != nil {
		return st, err
	}
	return st, nil
}

// changedTopics walks the topic list and returns the topics whose last post
// is not the stored one, newest activity first.
func (s *forumSyncer) changedTopics(fid int) ([]mzhttpclient.ForumTopic, string, error) {
	var out []mzhttpclient.ForumTopic
	var path string
	seen := map[int]bool{}
	for offset := 0; ; offset += mzhttpclient.ForumPageSize {
		s.wait()
		p, err := s.src.GetForumTopics(fid, offset, s.now())
		if err != nil {
			return nil, "", err
		}
		if offset == 0 {
			path = p.Forum
		}
		done := false
		for _, t := range p.Topics {
			if seen[t.ID] {
				continue // shifted down a page by new activity while we read
			}
			seen[t.ID] = true
			if !s.full {
				stored, ok, err := s.db.TopicState(t.ID)
				if err != nil {
					return nil, "", err
				}
				if ok && stored.LastPostID == t.LastPostID {
					// Pinned topics sit on top whatever their date; the rest are ordered.
					done = done || !t.Pinned
					continue
				}
			}
			out = append(out, t)
		}
		if done || len(p.Topics) == 0 || offset/mzhttpclient.ForumPageSize+1 >= p.Pages {
			return out, path, nil
		}
		if offset%(10*mzhttpclient.ForumPageSize) == 0 && offset > 0 {
			fmt.Fprintf(s.log, "forum %d: topic list page %d/%d\n", fid, offset/mzhttpclient.ForumPageSize+1, p.Pages)
		}
	}
}

// syncTopic fetches a topic's new pages and saves them; it returns the number
// of posts fetched. A stored topic is read again from the last page fetched
// before (it may have filled up since), or from the page now holding the last
// post if that is earlier (posts were deleted).
func (s *forumSyncer) syncTopic(t mzhttpclient.ForumTopic) (int, error) {
	start := 0
	if !s.full {
		stored, ok, err := s.db.TopicState(t.ID)
		if err != nil {
			return 0, err
		}
		if ok {
			start = min(stored.SyncedOffset, t.LastPostOffset)
		}
	}
	var posts []mzhttpclient.ForumPost
	var title string
	seen := map[int]bool{}
	last := start
	for offset := start; ; offset += mzhttpclient.ForumPageSize {
		s.wait()
		p, err := s.src.GetForumTopic(t.ForumID, t.ID, offset, s.now())
		if err != nil {
			return 0, err
		}
		if title == "" {
			title = p.Title
		}
		if len(p.Posts) == 0 {
			break
		}
		last = offset
		for _, post := range p.Posts {
			// Pinned topics repeat the opening post on every page.
			if !seen[post.ID] {
				seen[post.ID] = true
				posts = append(posts, post)
			}
		}
		if offset/mzhttpclient.ForumPageSize+1 >= p.Pages {
			break
		}
	}
	if err := s.db.SaveTopic(t, title, posts, last, s.now()); err != nil {
		return 0, err
	}
	return len(posts), nil
}
