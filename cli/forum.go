package cli

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/siimhs/mz/mzhttpclient"
	"github.com/spf13/cobra"
)

var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "Scrape forums, topic lists and posts (JSON)",
	Long: `Scrape the ManagerZone forums as JSON.

Pages hold 50 topics or posts and are addressed by --offset (0, 50, 100, ...).
Dates are read in the local timezone, which must match the account's MZ setting.`,
}

var forumLang string

var forumListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the forums of a language section",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		fs, err := c.GetForums(forumLang, now())
		if err != nil {
			return err
		}
		return printJSON(fs)
	},
}

var (
	forumOffset int
	topicsPages int
	topicPages  int
	forumSince  string
	forumDelay  time.Duration
)

var forumTopicsCmd = &cobra.Command{
	Use:   "topics <forum-id|url>",
	Short: "List topics of a forum, most recently active first",
	Example: `  mz forum topics 10
  mz forum topics 10 --pages 0 --since 2026-11-01   # every topic active since Nov 1
  mz forum topics 'https://www.managerzone.com/?p=forum&sub=topics&forum_id=10&sport=soccer'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := idArg(args[0], "forum_id")
		if err != nil {
			return err
		}
		var since time.Time
		if forumSince != "" {
			if since, err = parseSince(forumSince); err != nil {
				return err
			}
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		var out mzhttpclient.ForumTopicsPage
		seen := map[int]bool{}
		err = eachPage(forumOffset, topicsPages, func(offset int) (int, bool, error) {
			p, err := c.GetForumTopics(id, offset, now())
			if err != nil {
				return 0, false, err
			}
			if offset == forumOffset {
				out = p
				out.Topics = nil
			}
			older := false
			for _, t := range p.Topics {
				if !since.IsZero() && t.LastPostAt.Before(since) {
					// Pinned topics sit on top regardless of date; only the rest are ordered.
					older = older || !t.Pinned
					continue
				}
				if !seen[t.ID] {
					seen[t.ID] = true
					out.Topics = append(out.Topics, t)
				}
			}
			return p.Pages, !older, nil
		})
		if err != nil {
			return err
		}
		return printJSON(out)
	},
}

var forumTopicCmd = &cobra.Command{
	Use:   "topic <forum-id> <topic-id> | topic <url>",
	Short: "Get the posts of a topic, oldest first (all pages by default)",
	Long: `Get the posts of a topic, oldest first (all pages by default).

MZ only shows a topic under its own forum, so the forum id is required.`,
	Example: `  mz forum topic 10 13143067
  mz forum topic 10 13032964 --offset 250          # from page 6 on
  mz forum topic 'https://www.managerzone.com/?p=forum&sub=topic&topic_id=13143067&forum_id=10&sport=soccer'`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fid, id, err := topicArgs(args)
		if err != nil {
			return err
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		var out mzhttpclient.ForumTopicPage
		seen := map[int]bool{}
		err = eachPage(forumOffset, topicPages, func(offset int) (int, bool, error) {
			p, err := c.GetForumTopic(fid, id, offset, now())
			if err != nil {
				return 0, false, err
			}
			if offset == forumOffset {
				out = p
				out.Posts = nil
			}
			for _, post := range p.Posts {
				// Pinned topics repeat the opening post on every page.
				if !seen[post.ID] {
					seen[post.ID] = true
					out.Posts = append(out.Posts, post)
				}
			}
			return p.Pages, true, nil
		})
		if err != nil {
			return err
		}
		return printJSON(out)
	},
}

// eachPage calls fetch for offset, offset+50, ... until it has fetched `pages`
// pages (0 = all), reached the last page, or fetch says stop.
func eachPage(offset, pages int, fetch func(offset int) (total int, more bool, err error)) error {
	if offset < 0 || offset%mzhttpclient.ForumPageSize != 0 {
		return fmt.Errorf("--offset must be a multiple of %d", mzhttpclient.ForumPageSize)
	}
	for i := 0; pages == 0 || i < pages; i++ {
		if i > 0 && forumDelay > 0 {
			time.Sleep(forumDelay)
		}
		total, more, err := fetch(offset)
		if err != nil {
			return err
		}
		offset += mzhttpclient.ForumPageSize
		if !more || offset/mzhttpclient.ForumPageSize >= total {
			return nil
		}
	}
	return nil
}

// topicArgs reads "<forum-id> <topic-id>" or a topic URL.
func topicArgs(args []string) (forumID, topicID int, err error) {
	if len(args) == 2 {
		if forumID, err = idArg(args[0], "forum_id"); err != nil {
			return
		}
		topicID, err = idArg(args[1], "topic_id")
		return
	}
	if _, err := strconv.Atoi(args[0]); err == nil {
		return 0, 0, fmt.Errorf("want <forum-id> <topic-id> or a topic URL")
	}
	if forumID, err = idArg(args[0], "forum_id"); err != nil {
		return
	}
	topicID, err = idArg(args[0], "topic_id")
	return
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
	forumListCmd.Flags().StringVar(&forumLang, "lang", "", "language section, e.g. en, se, ar (default: the account's)")
	for _, c := range []*cobra.Command{forumTopicsCmd, forumTopicCmd} {
		c.Flags().IntVar(&forumOffset, "offset", 0, "start offset (a multiple of 50)")
		c.Flags().DurationVar(&forumDelay, "delay", 500*time.Millisecond, "pause between page requests")
	}
	forumTopicsCmd.Flags().IntVar(&topicsPages, "pages", 1, "pages to fetch, 0 = all")
	forumTopicsCmd.Flags().StringVar(&forumSince, "since", "", "stop at topics with no post since this date (YYYY-MM-DD or RFC3339)")
	forumTopicCmd.Flags().IntVar(&topicPages, "pages", 0, "pages to fetch, 0 = all")
	forumCmd.AddCommand(forumListCmd, forumTopicsCmd, forumTopicCmd)
	rootCmd.AddCommand(forumCmd)
}
