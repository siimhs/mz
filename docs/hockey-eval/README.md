# Hockey player evaluation research

Goal: a position-based score for ManagerZone **ice hockey** players that gives
one comparable number per player and can shift with tactic settings:

```
evaluate(player, position, tactic) -> double
```

Same process as `docs/player-eval/` (football). The football files are a
template for process and format only; they are **not** hockey facts.

## Layout

- `sources/manual/` – hockey ManagerZone manual, 40 chapters as text, fetched
  2026-11-21 from `?p=manual_faq&cat_id=N&sport=hockey` (ids 721–769, taken
  from the landing page's `<select>`, stopped at the first language name).
- `sources/beginners_guide/` – hockey Beginners Guide, 4 chapters
  (`get_started`, `manager_basics`, `your_first_season`, `competitions`), from
  `ajax.php?p=rules&sub=beginners_guide&sport=hockey&chapter=<name>`.
- `HANDOFF.md` – rules for research subagents; `SWEEP.md` – title-sweep brief.
- `research/` – one findings file per subagent: `01-manual.md`,
  `02-sweep-<LO>.md` (one per topic-id quartile).
- `RESULTS.md` – synthesis: positions, weights, tactic modifiers, pseudo code.

## Sources and caveats

- **Manual**: hockey chapters have their own cat_ids (721+; football's are
  13xx). `&sport=hockey`, `&share_sport=hockey` and no parameter return the
  same text for a hockey id (checked on A04 Tactics), so the id decides the sport.
  The text was converted with a stdlib HTMLParser script and site navigation was stripped:
  everything before the "Search Manual" box and after "Follow us on:".
  A few chapters (A01, A07, B08, C01, C04, C07, E03) are near-empty on the
  site itself. Some generic sentences still say "pitch".
- **Forum**: local SQLite (`~/.local/share/mz/mz.db`), hockey forums
  **138 "ManagerZone talk"** (742 topics, 9,796 posts) and **140 "Questions &
  Answers"** (449 topics, 1,892 posts): 1,191 topics and 11,688 posts,
  2009-12-11 … 2026-11-19. Counts were stable between 13:21 and 13:28 on
  2026-11-21. The site's own totals are higher (901 / 697 topics), and
  football has a similar gap (forum 10: 6,766 on the site vs 5,110 in the DB),
  so we read it as how the site counts (moved or deleted topics), not an
  unfinished sync.
- The same DB holds football forums 10 and 12. **Every hockey query must be
  restricted**: `--sport hockey` on `./mz forum …`, `forum_id IN (138,140)`
  in SQL.
- The hockey forum is about **10× smaller** than football's (11.7k vs 139k
  posts), so expect less consensus and more single opinions.
