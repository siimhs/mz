# Title sweep brief (agents B–E)

Read `README.md` and `HANDOFF.md` first; their rules (sources, citations,
confidence and era tags, context limits, own scratchpad folder) apply.
A separate agent covers the manual and guide (`research/01-manual.md`); you
may grep `sources/` to check a term, but your job is the forum.

## Task

You get a topic id range `id >= LO AND id < HI` (one quartile of the 1,191
hockey topics; ids grow with time) and the date span of its first posts.

1. **Scan every title** in your range, in pages of 400:
   `./mz sql "SELECT id, forum_id, replies, substr(title,1,80) t FROM topics WHERE forum_id IN (138,140) AND id >= LO AND id < HI ORDER BY id LIMIT 400 OFFSET N"`
   Note candidates about: skills (each of the ten), positions (G, D, C, W/F),
   lines and ice time, powerplay/boxplay, player quality or comparison, which
   player to buy, sell or play, tactics, mentality, pressing, passing,
   aggression, formations, match engine and sim updates, form, experience,
   stamina, training priorities and camps, "is X skill useful", match stats,
   third-party tools or scripts.
2. **Write the candidate list to your file before digging.**
3. **Dig**: read candidates, most promising first: staff answers, Q&A
   (forum 140), many replies, recent. Use `--limit 30` pages and `substr` in
   SQL. When a topic points to another topic, follow it, even outside your range.
   - Run full-text searches limited to your date span for what you learn
     (new terms, staff names, version numbers). `forum search` has `--since`
     only; for an upper bound use SQL on `posts` with `posted_at < 'YYYY-MM-DD'`.
   - Staff: check `users.roles`, plus posters called Crew/GPA/dev by others
     (see HANDOFF "Lessons").
4. Write `research/02-sweep-<LO>.md` incrementally in HANDOFF format, plus:
   - `## Engine versions`: every sim release or change you find (date, what changed).
   - `## Changes to the model`: concrete proposals (position X exists or not,
     skill Y matters most for Z, tactic T shifts weight from A to B, line/ice-time
     effects), each with citations and a confidence tag.
   - `## Candidate topics`: id and title of each candidate, marked read or unread.

Budget: stay under ~120k context. A quartile is ~300 titles (~6k tokens), so
you have room to dig deep, but keep notes short and write to the file often.
