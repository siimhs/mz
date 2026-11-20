package cli

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/siimhs/mz/store"
	"github.com/spf13/cobra"
)

var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "Sync the ManagerZone forums into a local database and read them (JSON)",
	Long: `Sync the ManagerZone forums into a local SQLite database and read them as JSON.

'mz forum sync' fetches from MZ; every other forum command reads only the
database, so it is fast and works offline. The database is $MZ_DB, or
~/.local/share/mz/mz.db. For queries of your own see 'mz sql'.

Dates are read in the local timezone, which must match the account's MZ setting.`,
}

// Forums synced by default: English » ManagerZone talk, English » Questions & Answers.
var defaultSyncForums = []int{10, 12}

var (
	syncForums []int
	syncFull   bool
	syncDelay  time.Duration
)

var forumSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch new posts from MZ into the database",
	Long: `Fetch new posts from MZ into the database.

The first sync of a forum fetches all of it, which takes a while (one request
per page of 50 topics or posts, --delay apart). Later syncs fetch only topics
with new posts, from the last page read before. An interrupted sync can simply
be run again. Progress goes to stderr.`,
	Example: `  mz forum sync                    # ManagerZone talk and Questions & Answers
  mz forum sync --forum 10 --forum 13
  mz forum sync --full             # re-read everything, e.g. to pick up edited posts`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDB(false)
		if err != nil {
			return err
		}
		defer db.Close()
		c, err := newClient()
		if err != nil {
			return err
		}
		s := &forumSyncer{src: c, db: db, lang: "en", full: syncFull, delay: syncDelay, now: now, log: os.Stderr}
		start := time.Now()
		st, err := s.sync(syncForums)
		fmt.Fprintf(os.Stderr, "synced %d topics, %d posts fetched in %s\n",
			st.topics, st.posts, time.Since(start).Round(time.Second))
		return err
	},
}

var forumListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the forums of the English section",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDB(true)
		if err != nil {
			return err
		}
		defer db.Close()
		fs, err := db.Forums()
		if err != nil {
			return err
		}
		return printJSON(fs)
	},
}

var (
	readOffset  int
	topicsLimit int
	topicLimit  int
	searchLimit int
	readSince   string
)

var forumTopicsCmd = &cobra.Command{
	Use:   "topics <forum-id|url>",
	Short: "List a forum's topics: pinned first, then most recently active",
	Example: `  mz forum topics 10
  mz forum topics 10 --limit 0 --since 2026-11-01   # every topic active since Nov 1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := idArg(args[0], "forum_id")
		if err != nil {
			return err
		}
		q := store.TopicsQuery{ForumID: id, Offset: readOffset, Limit: topicsLimit}
		if q.Since, err = sinceFlag(); err != nil {
			return err
		}
		db, err := openDB(true)
		if err != nil {
			return err
		}
		defer db.Close()
		p, err := db.Topics(q)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w (sync it with 'mz forum sync --forum %d')", err, id)
		}
		if err != nil {
			return err
		}
		return printJSON(p)
	},
}

var forumTopicCmd = &cobra.Command{
	Use:   "topic <topic-id|url>",
	Short: "Get the posts of a topic, oldest first",
	Example: `  mz forum topic 13143067
  mz forum topic 13032964 --offset 250 --limit 50
  mz forum topic 'https://www.managerzone.com/?p=forum&sub=topic&topic_id=13143067&forum_id=10&sport=soccer'`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// "<forum-id> <topic-id>" is accepted too: MZ's URLs need both.
		id, err := idArg(args[len(args)-1], "topic_id")
		if err != nil {
			return err
		}
		db, err := openDB(true)
		if err != nil {
			return err
		}
		defer db.Close()
		p, err := db.Topic(id, readOffset, topicLimit)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w (its forum is not synced, or the topic is newer than the last sync)", err)
		}
		if err != nil {
			return err
		}
		return printJSON(p)
	},
}

var (
	searchForum  int
	searchAuthor string
)

var forumSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Full-text search over posts, best matches first",
	Long: `Full-text search over the posts' own text (quotes excluded) and titles.

The query is SQLite FTS5 syntax: words (all must match; "tactics" also
matches "tactic"), "an exact phrase", prefix*, a OR b, NOT, NEAR(a b, 5),
title:word. Snippets mark matches with [brackets].`,
	Example: `  mz forum search 'training camp'
  mz forum search '"youth academy" scout*' --since 2025-01-01
  mz forum search 'offside OR "high line"' --forum 12 --limit 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := store.SearchQuery{Match: args[0], ForumID: searchForum, Author: searchAuthor, Limit: searchLimit}
		var err error
		if q.Since, err = sinceFlag(); err != nil {
			return err
		}
		db, err := openDB(true)
		if err != nil {
			return err
		}
		defer db.Close()
		hits, err := db.Search(q)
		if err != nil {
			return err
		}
		return printJSON(hits)
	},
}

func openDB(readOnly bool) (*store.DB, error) {
	path, err := store.DefaultPath()
	if err != nil {
		return nil, err
	}
	if readOnly {
		return store.OpenReadOnly(path)
	}
	return store.Open(path)
}

// idArg accepts a plain id or a forum URL holding it as the `key` query parameter.
func idArg(s, key string) (int, error) {
	if id, err := strconv.Atoi(s); err == nil && id > 0 {
		return id, nil
	}
	if u, err := url.Parse(s); err == nil {
		if id, err := strconv.Atoi(u.Query().Get(key)); err == nil && id > 0 {
			return id, nil
		}
	}
	return 0, fmt.Errorf("%q: want an id or a forum URL with %s", s, key)
}

func sinceFlag() (time.Time, error) {
	if readSince == "" {
		return time.Time{}, nil
	}
	return parseSince(readSince)
}

func parseSince(s string) (time.Time, error) {
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("--since %q: want YYYY-MM-DD or RFC3339", s)
}

func init() {
	forumSyncCmd.Flags().IntSliceVar(&syncForums, "forum", defaultSyncForums, "forum ids to sync (repeat or comma-separate)")
	forumSyncCmd.Flags().BoolVar(&syncFull, "full", false, "re-fetch every topic from its first page")
	forumSyncCmd.Flags().DurationVar(&syncDelay, "delay", 500*time.Millisecond, "pause between requests to MZ")

	forumTopicsCmd.Flags().IntVar(&topicsLimit, "limit", 50, "topics to show, 0 = all")
	forumTopicCmd.Flags().IntVar(&topicLimit, "limit", 0, "posts to show, 0 = all")
	forumSearchCmd.Flags().IntVar(&searchLimit, "limit", 20, "hits to show, 0 = all")
	for _, c := range []*cobra.Command{forumTopicsCmd, forumTopicCmd} {
		c.Flags().IntVar(&readOffset, "offset", 0, "skip this many")
	}
	forumTopicsCmd.Flags().StringVar(&readSince, "since", "", "only topics with a post since this date (YYYY-MM-DD or RFC3339)")
	forumSearchCmd.Flags().StringVar(&readSince, "since", "", "only posts since this date (YYYY-MM-DD or RFC3339)")
	forumSearchCmd.Flags().IntVar(&searchForum, "forum", 0, "only this forum id")
	forumSearchCmd.Flags().StringVar(&searchAuthor, "author", "", "only posts by this user name")

	forumCmd.AddCommand(forumSyncCmd, forumListCmd, forumTopicsCmd, forumTopicCmd, forumSearchCmd)
	rootCmd.AddCommand(forumCmd)
}
