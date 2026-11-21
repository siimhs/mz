# mz

Command line tool for [ManagerZone](https://www.managerzone.com) (football and hockey).

## Setup

Go runs in Docker; only the built static binary runs on the host.

```sh
cp .env.example .env      # fill in MZ_USERNAME / MZ_PASSWORD
make build                # -> ./mz
make test                 # Go unit tests
make e2e                  # end-to-end tests: the binary against a fake MZ
```

`mz` reads credentials from the environment or `./.env`.

The sport is football unless `--sport hockey` (any command) or `MZ_SPORT=hockey`
says otherwise; the flag wins over the variable. So far `whoami`, `mz friendly …`
and `mz forum …` are tested with hockey; `players` and `eval` assume football.

## Usage

```sh
mz whoami
mz players [--json]
mz matches [team-id] [--played] [--type friendly]

mz eval [--passing short|long|wings] [--pressing passive|normal|committed]
        [--opponent short|long|wings] [--pos GK|CD|FB|DM|CM|W|ST]
        [--no-condition] [--json] [--weights file] [--dump-weights]

mz friendly candidates                       # teams on the quick-challenge page
mz friendly send <team-id> --date 2026-11-22 [--home 12 --away 20 --tactic a]
mz friendly auto --dry-run                   # accept invites, quick challenge, invite for free days
mz friendly auto --tactics tue=a,wed=a,default=b

mz forum sync [--forum 10 --forum 12] [--full]   # fetch new posts from MZ into SQLite
mz forum list                                     # the rest read only the database (JSON)
mz forum topics <forum-id> [--limit N] [--offset N] [--since 2026-11-01]
mz forum topic <topic-id|url> [--limit N] [--offset N]
mz forum search '"youth academy" scout*' [--forum 12] [--author name] [--since …]
mz sql 'SELECT …'                                 # read-only SQL, JSON rows (hidden; '.schema')
```

The database is `$MZ_DB`, else `~/.local/share/mz/mz.db`. `sync` defaults to
the sport's English » ManagerZone talk and Questions & Answers (football 10 and
12, hockey 138 and 140); the first sync fetches all of a forum, later ones only
topics with new posts. `forum list` and `forum search` show the sport's forums
only; `topics` and `topic` take ids, which are unique across sports.

`mz eval` scores the squad per position for a tactic: best position and all
scores, or with `--pos` the squad ranked for one position. The model is
[docs/player-eval/RESULTS.md](docs/player-eval/RESULTS.md); **its weights are
unvalidated priors**. They live in `eval/weights.json` (embedded). To tune
without rebuilding, `mz eval --dump-weights > w.json`, edit it, and pass
`--weights w.json` or set `$MZ_EVAL_WEIGHTS`.

End-to-end tests are described in [e2e/README.md](e2e/README.md).

See [docs/PLAN.md](docs/PLAN.md) for status and roadmap.
