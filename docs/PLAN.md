# Plan

Personal CLI for [ManagerZone](https://www.managerzone.com) (football). First goal:
automate routine chores from the command line. Later: scrape training, transfer
market and played matches, and store them in a DB.

## Status

Working (tested live 2026-11-20):

- `mz whoami`, `mz players [--json]`, `mz matches [tid] [--played] [--type friendly]`
- `mz friendly candidates`, `mz friendly send <tid> --date …`, `mz friendly auto [--dry-run]`
  (`auto` fills every bookable game day, read from the friendly schedule page
  `ajax.php?p=challenge&sub=personal-challenge-template`: 1. accept incoming invites,
  2. accept a quick challenge for the next game day (`sub=quick-challenge-accept&tid=`),
  3. else opt in (`sub=quick-challenge-opt-in&date=<#qc-date>&tactic_home=&tactic_away=`),
  4. invite the quick-challenge page's teams for the other free days, skipping teams
  already invited for that day. Game days run 01:00–00:59; friendlies can be booked
  up to 7 days ahead, never on league days (Wed, Sun). **Accepting an incoming invite is
  untested live**: the link format is read from MZ's challenges.js `handleChallenge(url,
  side, …)` + `&tacticset=&tacticset2=`; no incoming invite has been seen on the page yet.)
- `mz forum sync` scrapes forums 10 and 12 (by default) into SQLite (`store/`, schema in
  `store/schema.sql`, FTS5 over post text); `mz forum list|topics|topic|search` and the
  hidden `mz sql` read only the DB. Purpose: material for analysing what works in MZ
  (player evaluation, tactics) with Claude. Sync is incremental: forums whose latest post
  is unchanged are skipped; the topic list is walked down to the first unchanged topic;
  changed topics are fetched oldest activity first (so an interrupted sync resumes
  correctly) from the last page read before. Edited posts are only picked up by `--full`;
  deleted posts stay in the DB. Hockey (`--sport hockey`) syncs its own talk and Q&A
  (138, 140) into the same DB; `forums.sport` (schema v2) says which sport a forum is,
  and `list`/`search` are scoped to the current sport. Notes on MZ's forum:
  a forum's pages are empty unless `sport=` is the forum's own sport;
  50 topics/posts per page (`offset=`); a topic page is empty unless `forum_id` is the
  topic's own forum; pinned topics repeat the opening post on every page (deduped by id);
  the list's "Posts" column counts replies and its date is the last post's; titles in
  the list are cut with `..` (the full one is on the topic page); recent dates are
  "Today", "Yesterday" or a weekday name. They are read in the local timezone. They are
  not Stockholm time: at 23:00 Tallinn time a post already showed "Today 22:04".
- `mz eval [--passing …] [--pressing …] [--opponent …] [--pos X] [--no-condition] [--json]`:
  per-position scores (GK CD FB DM CM W ST) from the model in
  [player-eval/RESULTS.md](player-eval/RESULTS.md), implemented in the pure `eval` package.
  Weights, tactic modifiers and constants are in `eval/weights.json` (embedded; replaced
  by `--weights file` / `$MZ_EVAL_WEIGHTS`, starting from `--dump-weights`). **All
  weights are unvalidated priors.** First look at our squad: outfield scores for one
  player are close together (often < 0.2 apart), so the best position is decided by
  small, unmeasured differences. `mz players` now also reports the player `ID`.
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
2. `mz friendly challenge --team-id <tid>`, `mz matches set-tactic --match-id <mid> --tactic <A|B|C>`
3. `mz team [tid]`, `mz player <pid>`, `mz match <mid>`
4. **Player evaluation: calibrate** (`mz eval` exists, its weights are priors). Next session:
   scrape played match pages (since 2024-11-25 they show per-player passes, failed
   passes, tackles, assists, saves, distance and avg speed) for our and opponents'
   matches into SQLite, then fit the weights in `eval/weights.json`: tackles won vs
   TAC/PI/SPD, pass success vs PAS/AP per passing style, saves vs KEE/PI, goals per shot
   vs SHO against the opposing KEE; test the core term, the size of Experience and
   whether Speed flattens above 8 (RESULTS.md "Next steps" / "Still open"). Then use the
   scores for lineups and **training plans**.
5. **Scraping + storage** — `mz training`, `mz transfers search …`,
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
