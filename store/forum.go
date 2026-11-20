package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
)

// ErrNotFound means the forum or topic is not in the database.
var ErrNotFound = errors.New("not in the database")

// SaveForums upserts the forum index of a language section. Sync state is kept.
func (db *DB) SaveForums(lang string, forums []mzhttpclient.Forum) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, f := range forums {
		_, err := tx.Exec(`INSERT INTO forums (id, name, lang, position, description, total_posts, total_topics,
				latest_topic_id, latest_topic, latest_post_id, latest_post_by, latest_post_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO UPDATE SET name = excluded.name, lang = excluded.lang,
				position = excluded.position, description = excluded.description,
				total_posts = excluded.total_posts, total_topics = excluded.total_topics,
				latest_topic_id = excluded.latest_topic_id, latest_topic = excluded.latest_topic,
				latest_post_id = excluded.latest_post_id, latest_post_by = excluded.latest_post_by,
				latest_post_at = excluded.latest_post_at`,
			f.ID, f.Name, nullStr(lang), i, nullStr(f.Description), f.TotalPosts, f.TotalTopics,
			nullInt(f.LatestTopicID), nullStr(f.LatestTopic), nullInt(f.LatestPostID),
			nullStr(f.LatestPostBy), fmtTime(f.LatestPostAt))
		if err != nil {
			return fmt.Errorf("forum %d: %w", f.ID, err)
		}
	}
	return tx.Commit()
}

// ForumSyncedPostID is the forum's latest post id as of its last completed sync, 0 if never.
func (db *DB) ForumSyncedPostID(forumID int) (int, error) {
	var id sql.NullInt64
	err := db.QueryRow(`SELECT synced_post_id FROM forums WHERE id = ?`, forumID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return int(id.Int64), err
}

// MarkForumSynced records a completed sync: every topic active up to latestPostID is stored.
// path is the forum's breadcrumb from the topic list ("English » ManagerZone talk").
func (db *DB) MarkForumSynced(forumID int, path string, latestPostID int, at time.Time) error {
	_, err := db.Exec(`INSERT INTO forums (id, name, path, synced_post_id, synced_at) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET path = coalesce(excluded.path, path),
			synced_post_id = excluded.synced_post_id, synced_at = excluded.synced_at`,
		forumID, lastSegment(path), nullStr(path), nullInt(latestPostID), fmtTime(at))
	return err
}

func lastSegment(path string) string {
	if i := strings.LastIndex(path, "»"); i >= 0 {
		return strings.TrimSpace(path[i+len("»"):])
	}
	return path
}

// TopicState is what a sync needs to know about a stored topic.
type TopicState struct {
	LastPostID   int
	SyncedOffset int
}

// TopicState returns the stored topic's sync state; ok is false if the topic is not stored.
func (db *DB) TopicState(topicID int) (st TopicState, ok bool, err error) {
	var last sql.NullInt64
	err = db.QueryRow(`SELECT last_post_id, synced_offset FROM topics WHERE id = ?`, topicID).
		Scan(&last, &st.SyncedOffset)
	if errors.Is(err, sql.ErrNoRows) {
		return st, false, nil
	}
	st.LastPostID = int(last.Int64)
	return st, err == nil, err
}

// SaveTopic stores a topic as listed, its posts, and the offset of the last page
// fetched, in one transaction: a stored last_post_id means its posts are stored.
// title is the full title from the topic page ("" keeps the list's).
func (db *DB) SaveTopic(t mzhttpclient.ForumTopic, title string, posts []mzhttpclient.ForumPost, syncedOffset int, at time.Time) error {
	if title == "" {
		title = t.Title
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := saveUser(tx, t.Author); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO topics (id, forum_id, title, pinned, views, replies, author_id,
			last_post_id, last_post_offset, last_post_at, synced_offset, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET forum_id = excluded.forum_id, title = excluded.title,
			pinned = excluded.pinned, views = excluded.views, replies = excluded.replies,
			author_id = excluded.author_id, last_post_id = excluded.last_post_id,
			last_post_offset = excluded.last_post_offset, last_post_at = excluded.last_post_at,
			synced_offset = excluded.synced_offset, synced_at = excluded.synced_at`,
		t.ID, t.ForumID, title, t.Pinned, t.Views, t.Replies, nullInt(t.Author.ID),
		nullInt(t.LastPostID), t.LastPostOffset, fmtTime(t.LastPostAt), syncedOffset, fmtTime(at))
	if err != nil {
		return fmt.Errorf("topic %d: %w", t.ID, err)
	}
	for _, p := range posts {
		if err := saveUser(tx, p.Author); err != nil {
			return err
		}
		var quotes any
		if len(p.Quotes) > 0 {
			b, err := json.Marshal(p.Quotes)
			if err != nil {
				return err
			}
			quotes = string(b)
		}
		_, err := tx.Exec(`INSERT INTO posts (id, topic_id, forum_id, title, author_id, team_id, posted_at, text, html, quotes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO UPDATE SET topic_id = excluded.topic_id, forum_id = excluded.forum_id,
				title = excluded.title, author_id = excluded.author_id, team_id = excluded.team_id,
				posted_at = excluded.posted_at, text = excluded.text, html = excluded.html, quotes = excluded.quotes`,
			p.ID, t.ID, t.ForumID, nullStr(p.Title), nullInt(p.Author.ID), nullInt(p.TeamID),
			fmtTime(p.PostedAt), p.Text, p.HTML, quotes)
		if err != nil {
			return fmt.Errorf("post %d: %w", p.ID, err)
		}
	}
	return tx.Commit()
}

func saveUser(tx *sql.Tx, u mzhttpclient.User) error {
	if u.ID == 0 {
		return nil
	}
	var roles any
	if len(u.Roles) > 0 {
		b, _ := json.Marshal(u.Roles)
		roles = string(b)
	}
	_, err := tx.Exec(`INSERT INTO users (id, name, country, federation_id, roles) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET name = excluded.name, country = excluded.country,
			federation_id = excluded.federation_id, roles = excluded.roles`,
		u.ID, u.Name, nullStr(u.Country), nullInt(u.FederationID), roles)
	if err != nil {
		return fmt.Errorf("user %d: %w", u.ID, err)
	}
	return nil
}

// --- reading ---

// Forums returns the stored forum index in page order.
func (db *DB) Forums() ([]mzhttpclient.Forum, error) {
	rows, err := db.Query(`SELECT id, name, coalesce(description, ''), coalesce(total_posts, 0),
			coalesce(total_topics, 0), coalesce(latest_topic_id, 0), coalesce(latest_topic, ''),
			coalesce(latest_post_id, 0), coalesce(latest_post_by, ''), latest_post_at
		FROM forums ORDER BY lang, position, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []mzhttpclient.Forum{}
	for rows.Next() {
		var f mzhttpclient.Forum
		var at sql.NullString
		if err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.TotalPosts, &f.TotalTopics,
			&f.LatestTopicID, &f.LatestTopic, &f.LatestPostID, &f.LatestPostBy, &at); err != nil {
			return nil, err
		}
		f.LatestPostAt = parseTime(at)
		out = append(out, f)
	}
	return out, rows.Err()
}

func (db *DB) forumPath(forumID int) (string, error) {
	var path string
	err := db.QueryRow(`SELECT coalesce(path, name) FROM forums WHERE id = ?`, forumID).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return path, err
}

const userCols = `u.id, u.name, coalesce(u.country, ''), coalesce(u.federation_id, 0), u.roles`

func scanUser(id sql.NullInt64, name, country sql.NullString, fed sql.NullInt64, roles sql.NullString) mzhttpclient.User {
	u := mzhttpclient.User{ID: int(id.Int64), Name: name.String, Country: country.String, FederationID: int(fed.Int64)}
	if roles.Valid {
		json.Unmarshal([]byte(roles.String), &u.Roles)
	}
	return u
}

// TopicsQuery selects a forum's topics: pinned first, then by latest activity.
type TopicsQuery struct {
	ForumID       int
	Offset, Limit int       // Limit 0 = all
	Since         time.Time // only topics with a post since then (pinned ones too)
}

// Topics returns stored topics in the shape of a topic list page. Pages counts
// pages of Limit topics.
func (db *DB) Topics(q TopicsQuery) (mzhttpclient.ForumTopicsPage, error) {
	page := mzhttpclient.ForumTopicsPage{ForumID: q.ForumID, Offset: q.Offset, Topics: []mzhttpclient.ForumTopic{}}
	var err error
	if page.Forum, err = db.forumPath(q.ForumID); err != nil {
		return page, err
	}
	where, args := `t.forum_id = ?`, []any{q.ForumID}
	if !q.Since.IsZero() {
		where += ` AND t.last_post_at >= ?`
		args = append(args, fmtTime(q.Since))
	}
	var total int
	if err := db.QueryRow(`SELECT count(*) FROM topics t WHERE `+where, args...).Scan(&total); err != nil {
		return page, err
	}
	if total == 0 && page.Forum == "" {
		return page, fmt.Errorf("forum %d: %w", q.ForumID, ErrNotFound)
	}
	page.Pages = 1
	if q.Limit > 0 && total > 0 {
		page.Pages = (total + q.Limit - 1) / q.Limit
	}
	limit := q.Limit
	if limit <= 0 {
		limit = -1
	}
	rows, err := db.Query(`SELECT t.id, t.forum_id, t.title, t.pinned, coalesce(t.views, 0), coalesce(t.replies, 0),
			coalesce(t.last_post_id, 0), coalesce(t.last_post_offset, 0), t.last_post_at, `+userCols+`
		FROM topics t LEFT JOIN users u ON u.id = t.author_id
		WHERE `+where+` ORDER BY t.pinned DESC, t.last_post_at DESC, t.id DESC LIMIT ? OFFSET ?`,
		append(args, limit, q.Offset)...)
	if err != nil {
		return page, err
	}
	defer rows.Close()
	for rows.Next() {
		var t mzhttpclient.ForumTopic
		var at, name, country, roles sql.NullString
		var uid, fed sql.NullInt64
		if err := rows.Scan(&t.ID, &t.ForumID, &t.Title, &t.Pinned, &t.Views, &t.Replies,
			&t.LastPostID, &t.LastPostOffset, &at, &uid, &name, &country, &fed, &roles); err != nil {
			return page, err
		}
		t.LastPostAt = parseTime(at)
		t.Author = scanUser(uid, name, country, fed, roles)
		page.Topics = append(page.Topics, t)
	}
	return page, rows.Err()
}

// Topic returns a stored topic's posts, oldest first, from offset on (limit 0 = all).
// Pages counts MZ's pages of 50 posts.
func (db *DB) Topic(topicID, offset, limit int) (mzhttpclient.ForumTopicPage, error) {
	page := mzhttpclient.ForumTopicPage{TopicID: topicID, Offset: offset, Posts: []mzhttpclient.ForumPost{}}
	err := db.QueryRow(`SELECT forum_id, title FROM topics WHERE id = ?`, topicID).Scan(&page.ForumID, &page.Title)
	if errors.Is(err, sql.ErrNoRows) {
		return page, fmt.Errorf("topic %d: %w", topicID, ErrNotFound)
	}
	if err != nil {
		return page, err
	}
	if page.Forum, err = db.forumPath(page.ForumID); err != nil {
		return page, err
	}
	var total int
	if err := db.QueryRow(`SELECT count(*) FROM posts WHERE topic_id = ?`, topicID).Scan(&total); err != nil {
		return page, err
	}
	page.Pages = max(1, (total+mzhttpclient.ForumPageSize-1)/mzhttpclient.ForumPageSize)
	if limit <= 0 {
		limit = -1
	}
	rows, err := db.Query(`SELECT p.id, p.topic_id, p.forum_id, coalesce(p.title, ''), coalesce(p.team_id, 0),
			p.posted_at, p.text, p.html, p.quotes, `+userCols+`
		FROM posts p LEFT JOIN users u ON u.id = p.author_id
		WHERE p.topic_id = ? ORDER BY p.id LIMIT ? OFFSET ?`, topicID, limit, offset)
	if err != nil {
		return page, err
	}
	defer rows.Close()
	for rows.Next() {
		var p mzhttpclient.ForumPost
		var at, quotes, name, country, roles sql.NullString
		var uid, fed sql.NullInt64
		if err := rows.Scan(&p.ID, &p.TopicID, &p.ForumID, &p.Title, &p.TeamID, &at, &p.Text, &p.HTML,
			&quotes, &uid, &name, &country, &fed, &roles); err != nil {
			return page, err
		}
		p.PostedAt = parseTime(at)
		p.Author = scanUser(uid, name, country, fed, roles)
		if quotes.Valid {
			json.Unmarshal([]byte(quotes.String), &p.Quotes)
		}
		page.Posts = append(page.Posts, p)
	}
	return page, rows.Err()
}

// SearchQuery is a full-text search over posts.
type SearchQuery struct {
	Match   string // FTS5 query: words, "a phrase", tactic*, a OR b, NEAR(a b, 5), text:word
	ForumID int
	Author  string // name, case-insensitive
	Since   time.Time
	Limit   int
}

// SearchHit is one matching post, best matches first.
type SearchHit struct {
	PostID   int
	TopicID  int
	ForumID  int
	Topic    string
	Author   string
	PostedAt time.Time
	Snippet  string // matched words in [brackets]
}

func (db *DB) Search(q SearchQuery) ([]SearchHit, error) {
	where, args := `posts_fts MATCH ?`, []any{q.Match}
	if q.ForumID != 0 {
		where += ` AND p.forum_id = ?`
		args = append(args, q.ForumID)
	}
	if q.Author != "" {
		where += ` AND u.name = ? COLLATE NOCASE`
		args = append(args, q.Author)
	}
	if !q.Since.IsZero() {
		where += ` AND p.posted_at >= ?`
		args = append(args, fmtTime(q.Since))
	}
	limit := q.Limit
	if limit <= 0 {
		limit = -1
	}
	rows, err := db.Query(`SELECT p.id, p.topic_id, p.forum_id, coalesce(t.title, ''), coalesce(u.name, ''),
			p.posted_at, snippet(posts_fts, 1, '[', ']', '…', 32)
		FROM posts_fts JOIN posts p ON p.id = posts_fts.rowid
		LEFT JOIN topics t ON t.id = p.topic_id LEFT JOIN users u ON u.id = p.author_id
		WHERE `+where+` ORDER BY bm25(posts_fts, 2.0, 1.0) LIMIT ?`, append(args, limit)...)
	if err != nil {
		return nil, fmt.Errorf("search query %q: %w", q.Match, err)
	}
	defer rows.Close()
	out := []SearchHit{}
	for rows.Next() {
		var h SearchHit
		var at sql.NullString
		if err := rows.Scan(&h.PostID, &h.TopicID, &h.ForumID, &h.Topic, &h.Author, &at, &h.Snippet); err != nil {
			return nil, err
		}
		h.PostedAt = parseTime(at)
		out = append(out, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search query %q: %w", q.Match, err)
	}
	return out, nil
}
