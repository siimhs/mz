# Handoff brief for research subagents (hockey)

Read `docs/hockey-eval/README.md` first. This is **ice hockey**. Football
research in `docs/player-eval/` is a format template only: do not carry
football facts (positions, passing styles, skills) over.

## Rules

- Use ONLY these sources:
  1. local forum DB through the `./mz` binary in `/home/siimh/src/mz`,
     **always with `--sport hockey`** (the default sport is football):
     - `./mz forum search '<fts5 query>' --sport hockey --limit 20 [--since YYYY-MM-DD] [--forum 138|140]`
       (snippets; FTS5: `"exact phrase"`, `OR`, `NEAR(a b, 10)`, `prefix*`, `title:word`)
     - `./mz forum topic <id> --sport hockey --limit 30 --offset N` (full posts)
     - `./mz sql "SELECT id, topic_id, substr(text,1,600) FROM posts WHERE forum_id IN (138,140) AND …"`
       Tables: `topics(id, forum_id, title, pinned, views, replies, author_id, …)`,
       `posts(id, topic_id, forum_id, title, author_id, team_id, posted_at, text, html, quotes)`,
       `users(id, name, country, federation_id, roles)`.
       **Every SQL query must include `forum_id IN (138,140)`**: football
       forums 10 and 12 are in the same DB.
  2. text files under `docs/hockey-eval/sources/` (grep them).
  No web access, no other files.
- Keep your context under ~120k tokens: always pass `--limit`, read topics in
  pages of 20–30, use `substr(text,1,N)` in SQL, and never dump whole topics
  or whole manual files (use `grep -n` plus `sed -n` ranges).
- Use **your own scratchpad subfolder** (given in your prompt) for any temp
  file. Never use generic shared names: parallel agents overwrote each other's
  files last time.
- Write findings to your output file in `docs/hockey-eval/research/` as you
  go (in case you run out of context), then overwrite it with a clean final
  version at the end.
- Cite every claim: `[post 33886074, 2011-05-05, author]` or
  `[manual A04 Tactics]`. Keep quotes short (≤2 sentences).
- Tag each claim's confidence: **official** (manual, or MZ staff), **consensus**
  (several independent posters agree), **opinion** (single poster),
  **contested** (posters disagree; give both sides). Tag the **engine era**
  too (see "Lessons").
- Prefer numbers (weights, rankings, thresholds, "X matters more than Y") over
  vague text.

## Lessons from the football round

- **`users.roles` is incomplete.** Football's patrick7777 (Power Challenge
  CEO, sim developer) has no role in the DB but is staff; he also posts in the
  hockey forums (2011–12 Crew announcements). Roles seen in hockey: Crew
  (powdersnow, thunderexpress, mikelony), GPA/NC/MA/LA/RA/GA/FA/HMA
  (moderators and admins: floodish, agapitod, sabiutza, …). Moderators are not
  sim developers: tag them official only when they relay Crew or quote release
  notes. Check posters others call "(Crew)" or "(GPA)" the same way.
  **Staff release notes are official.**
- **Engine versions matter.** Find the hockey sim's version history (release
  dates, what changed) and tag each claim by era. Newest-engine consensus beats
  old posts.
- **No source gave numeric weights** in football. Expect the same here, but
  collect every ranking, poll, threshold ("needs X ≥ 8"), poster test with
  numbers, and third-party tool or script formula (football found one: MZ
  Assistant's keeper weights).
- **Check that tactic options exist before modelling them.** Football's tactic
  screen had no counter-attack option, and its mentality setting showed no
  effect on skills. Record which options hockey really has and whether any
  source ties them to skills.
- **Per-player match statistics**: note what the match pages and stats show per
  player (goals, shots, hits, ice time, +/-, …), as data to calibrate on later.

## Output format (research/NN-topic.md)

```
# <topic>
## Summary (≤15 bullets, the answer)
## Findings (grouped, with citations, confidence and era tags)
## Open questions / contradictions
## Queries used (short list, so others don't repeat them)
```

## Background facts (from the manual, don't re-research)

Skills, 0–10 "pucks", **ten trainable** [manual A05 Training]: **Play
Intelligence, Power, Skating, Passing, Quickness, Shooting, Keeping, Puck
Control, Checking, Stamina**. Also Form (and Experience, which the manual barely mentions).
- Power: make and withstand tackles. Skating: balance, turning.
  Quickness: speed with and without the puck. Checking: pressure and 1-vs-1
  puck wins. Puck Control: trap passes, keep possession.

Rink and tactics [manual A04, guide manager_basics]: the rink has 6 zones per
line: **1 goalkeeper, 2 defenders, 2 wingers, 1 centre**. At least 10 skaters + 1 GK
in **3 lines, including powerplay and boxplay lines**. The guide's tactics
table names the shirts Goalkeeper, Defender, Center, Forward. Tactic sections in
the manual: Formation, Mentality (Defensive/Balanced/Attacking), Pressing
(Passive/…/Committed; Committed costs Stamina and causes more penalties). The
A04 intro also names a "Passing style", but **no section describes one** (maybe
copied from the football text). Find out from the forum whether it exists, and
what else the tactic screen has (line ice time / rotation, aggression, …).
