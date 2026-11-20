# Plan

Personal CLI for [ManagerZone](https://www.managerzone.com) (football). First goal:
automate routine chores from the command line. Later: scrape training, transfer
market and played matches, and store them in a DB.

## Status

Working (tested live 2026-11-20):

- `mz whoami`, `mz players [--json]`, `mz matches [tid] [--played] [--type friendly]`
- `mz friendly candidates`, `mz friendly send <tid> --date …`, `mz friendly auto [--dry-run]`
  (`auto` = old `main.go`: invite all quick-challenge teams for each day without a friendly)
- `mz forum sync` scrapes forums 10 and 12 (by default) into SQLite (`store/`, schema in
  `store/schema.sql`, FTS5 over post text); `mz forum list|topics|topic|search` and the
  hidden `mz sql` read only the DB. Purpose: material for analysing what works in MZ
  (player evaluation, tactics) with Claude. Sync is incremental: forums whose latest post
  is unchanged are skipped; the topic list is walked down to the first unchanged topic;
  changed topics are fetched oldest activity first (so an interrupted sync resumes
  correctly) from the last page read before. Edited posts are only picked up by `--full`;
  deleted posts stay in the DB. Notes on MZ's forum:
  50 topics/posts per page (`offset=`); a topic page is empty unless `forum_id` is the
  topic's own forum; pinned topics repeat the opening post on every page (deduped by id);
  the list's "Posts" column counts replies and its date is the last post's; titles in
  the list are cut with `..` (the full one is on the topic page); recent dates are
  "Today", "Yesterday" or a weekday name. They are read in the local timezone. They are
  not Stockholm time: at 23:00 Tallinn time a post already showed "Today 22:04".
- hidden `mz raw <path>` dumps any logged-in page, for discovering endpoints

Not yet verified after the rewrite: an actual live `send` (dry-run only). Open question:
the hour sent to MZ is converted to Europe/Stockholm; check that a 12:00 local invite
shows up as 12:00 in MZ. (e2e `TestFriendlySend` pins the current behaviour: 12:00 in
Europe/Tallinn is sent as `date_home=<Stockholm midnight>,11`.)

End-to-end tests: `make e2e`, see [e2e/README.md](../e2e/README.md).

## Next steps

1. **Foundation**
   - session: cookiejar, persist to `~/.config/mz/session`, `mz login` / `mz logout`
     (today every command logs in from `.env`)
   - harden the players scraper (it depends on exact whitespace and node layout;
     the old pretty-printed `testdata/players.html` no longer parses)
   - remove hardcoded `country_id=123` from the friendly request
2. **`mz friendly auto-arrange`** — fill friendlies as far ahead as allowed (~5 days).
   Order matters:
   1. accept invites other teams have sent us
   2. accept a quick challenge someone offers for the next day
   3. post our own quick challenge for the next day
   4. send invites to all recent quick-challenge teams

   Needs endpoints for: incoming invites (list/accept), quick challenge (list/accept/post).
3. `mz friendly challenge --team-id <tid>`, `mz matches set-tactic --match-id <mid> --tactic <A|B|C>`
4. `mz team [tid]`, `mz player <pid>`, `mz match <mid>`
5. **Player evaluation** — a numeric score per position: keeper, defender, midfield,
   winger, striker (finer roles like defensive/attacking mid later). Used for
   lineups/tactics and, day to day, for **training plans**.
6. **Scraping + storage** — `mz training`, `mz transfers search …`,
   `mz sync players|matches|training` into SQLite.

## Original command sketch

```
mz login                 # log in, save the session cookie to ~/.config/mz/session
mz logout
mz whoami                # username, team name/id, sport

mz team [tid]
mz players [--json|--csv] [--sort speed]
mz player <pid>

mz matches [tid] [--played|--upcoming] [--type friendly|league]
mz match <mid>
mz matches set-tactic --match-id <mid> --tactic <A|B|C>

mz friendly list                       # my scheduled and pending friendlies
mz friendly candidates                 # team ids from quick-challenge
mz friendly send <tid> --home "2026-10-21 10:00" --away "... 18:00" --tactic a
mz friendly challenge --team-id <tid>
mz friendly auto --days 8 --tactic tue,wed=a,default=b --dry-run
mz friendly auto-arrange

mz training [--week]
mz transfers search ...
mz sync players|matches|training
```
