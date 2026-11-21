# Player evaluation research

Goal: a position-based player evaluation that yields one comparable number,
and that can shift with tactics (passing style etc.):

```
evaluate(player, position, tactic) -> double
```

Proposed positions: GK, Def, Mid, Wing, St (open question whether more/fewer).

## Layout

- `sources/manual/` – ManagerZone manual chapters as text (fetched 2026-11-21
  from `?p=manual_faq&cat_id=…`).
- `sources/beginners_guide/` – Beginners Guide chapters as text.
- `research/` – one findings file per research subagent (see HANDOFF.md).
- `RESULTS.md` – synthesis: positions, weights, tactic modifiers, pseudo code.

## Sources and caveats

- Forum: local SQLite (`~/.local/share/mz/mz.db`), forums 10 "ManagerZone
  talk" and 12 "Questions & Answers": 10,928 topics, 139k posts,
  2011-02 … 2026-11. Round 1 (`research/01…05`) saw only ~25k posts up to
  ~2012-03 while the sync was still running. Round 2 (`research/06-sweep-*`,
  brief in `SWEEP.md`) scanned every title. Query with `./mz forum search`,
  `./mz forum topic`, `./mz sql`.
- The match engine changed later ("Sim release confirmed, Season 40", topic
  9506303), so 2011 forum advice may be partly outdated. The manual is current.
