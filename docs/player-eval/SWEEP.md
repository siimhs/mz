# Title sweep brief (round 2)

Read `README.md` and `HANDOFF.md` first; their rules (sources, citations,
confidence tags, context limits) still apply. Also read `RESULTS.md`: it is
the current model, and your job is to confirm, refute or refine it.

## What changed since round 1

The forum sync finished. The DB now has 10,928 topics and 139k posts,
2011-02 … 2026-11, across forum 10 (ManagerZone talk) and forum 12
(Questions & Answers). Round 1 saw only ~25k posts from 2011-02 … 2012-03.
Posts after ~2012 describe the current engine and **outweigh** 2011 posts.

## Task

You get a topic id range: `id >= LO AND id < HI`.

1. **Scan every title** in your range, in pages of 400:
   `./mz sql "SELECT id, forum_id, replies, title FROM topics WHERE id >= LO AND id < HI ORDER BY id LIMIT 400 OFFSET N"`
   Note candidate topics about: skills, positions, player quality or comparison,
   which player to buy, sell or play, tactics, passing styles, formations, mentality,
   pressing, match engine, sim updates, form, experience, stamina, training
   priorities, and "is X skill useful".
2. Skip topics already listed in `research/covered-topics.json`.
3. **Dig**: read candidate topics, most promising first (Q&A answers, Crew
   answers, many replies, recent). Use `--limit 30` pages and `substr` in SQL.
   When a topic points to another topic or thread, follow it.
   - Also run full-text searches limited to your date span for
     what you learn (e.g. new terms, Crew names). Staff/Crew posters: check
     `users.roles` and known staff such as patrick7777.
4. Write `research/06-sweep-<LO>.md` incrementally, in HANDOFF format, plus:
   - `## Changes to RESULTS.md`: concrete proposals (weight X up/down,
     position added or removed, tactic modifier changed), each with citations.
   - `## Candidate topics`: the id and title of each candidate, marked read or
     unread.

Budget: stay under ~120k context. Titles cost about 17k tokens for a quarter, so
keep notes short while scanning, and write candidates to the file before
digging.
