# Handoff brief for research subagents

Read `docs/player-eval/README.md` first.

## Rules

- Use ONLY these sources:
  1. local forum DB through the `./mz` binary in `/home/siimh/src/mz`:
     - `./mz forum search '<fts5 query>' --limit 20` (snippets; FTS5 syntax:
       `"exact phrase"`, `OR`, `NEAR(a b, 10)`, `prefix*`)
     - `./mz forum topic <id> --limit 30 --offset N` (full posts)
     - `./mz sql "SELECT id, topic_id, substr(text,1,600) FROM posts WHERE id IN (…)"`
  2. text files under `docs/player-eval/sources/` (grep them).
  No web access, no other files.
- Keep your context under ~120k tokens: always pass `--limit`, read topics in
  pages of 20–30, use `substr(text,1,N)` in SQL, and do not dump whole topics
  or whole manual files. Use grep -n plus `sed -n` ranges for the manual.
- Write findings as you go (in case you run out of context) to your output file
  in `docs/player-eval/research/`. Overwrite it with a clean final version
  at the end.
- Cite every claim: `[post 33886074, 2011-05-05, author]` or
  `[manual A04 Tactics]`. Keep quotes short (≤2 sentences).
- Tag each claim's confidence: **official** (manual or MZ staff/crew), **consensus**
  (several independent posters agree), **opinion** (single poster), **contested**
  (posters disagree, give both sides). Note whether a post is before or after the
  Season 40 sim release if it matters.
- Prefer numbers (weights, rankings, "X matters more than Y") over vague text.

## Output format (research/NN-topic.md)

```
# <topic>
## Summary (≤15 bullets, the answer)
## Findings (grouped, with citations and confidence tags)
## Open questions / contradictions
## Queries used (short list, so others don't repeat them)
```

## Background facts (don't re-research)

MZ player skills (scale 0–10, shown as balls): Speed, Stamina, Play Intelligence,
Passing, Shooting, Heading, Keeping, Ball Control, Tackling, Aerial Passing,
Set Plays; plus Experience and Form. Verify names against the manual if unsure.
Tactic elements per manual A04: formation, substitutions, passing style,
mentality, pressing, player roles, rules.
