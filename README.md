# mz

Command line tool for [ManagerZone](https://www.managerzone.com) (football).

## Setup

Go runs in Docker; only the built static binary runs on the host.

```sh
cp .env.example .env      # fill in MZ_USERNAME / MZ_PASSWORD
make build                # -> ./mz
make test                 # Go unit tests
make e2e                  # end-to-end tests: the binary against a fake MZ
```

`mz` reads credentials from the environment or `./.env`.

## Usage

```sh
mz whoami
mz players [--json]
mz matches [team-id] [--played] [--type friendly]

mz friendly candidates                       # teams on the quick-challenge page
mz friendly send <team-id> --date 2026-11-22 [--home 12 --away 20 --tactic a]
mz friendly auto --dry-run                   # invite candidates for each free day
mz friendly auto --days 5 --tactics tue=a,wed=a,default=b
```

End-to-end tests are described in [e2e/README.md](e2e/README.md).

See [docs/PLAN.md](docs/PLAN.md) for status and roadmap.
