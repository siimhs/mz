# End-to-end tests

Black-box tests for the `mz` binary. Each test is a Go function that:
1. starts a fake ManagerZone (`httptest`) serving saved real responses,
2. runs the built binary against it,
3. checks the binary's output, exit status and the requests it sent.

This module never imports mz. It only talks to the binary through arguments,
environment, stdout/stderr and HTTP, so mz can be rewritten in another language
without changing these tests.

```sh
make e2e                     # build mz, run every test
make e2e-update              # rewrite testdata/golden from actual output; review with git diff
```

## Layout

```
server.go            fake ManagerZone: routes, request log
harness_test.go      Start, Run, assertions, golden files
fixtures_test.go     route groups for the recorded pages: Login(), UpcomingMatches(), ...
*_test.go            the tests
testdata/mz/         recorded responses (scrubbed)
testdata/golden/     expected output that is too long to write inline
cmd/fakemz/          fixture tool: record, scrub
record.sh            records responses for a new test
```

## Writing a test

```go
func TestFriendlySend(t *testing.T) {
	mz := Start(t, Login(), SendInvite(inviteOK))

	res := mz.Run("friendly", "send", "849080", "--date", "2026-11-22", "--tactic", "a")

	res.WantExit(t, 0)
	res.WantStdout(t, "sent invite to 849080 for Sun 2026-11-22\n")
	sent := mz.Find("POST", sendInvite)
	wantForm(t, sent[0], map[string]string{"date_home": sun22 + ",11"})
}
```

- **Routes.** `Route(method, url, responses...)` answers a request.
  - Query parameter order doesn't matter, and `*` matches anything.
  - With several responses, the n-th request gets the n-th response and the last one repeats.
  - A later route for the same method and URL replaces an earlier one, so
    `Start(t, Login(), Route("GET", "/?p=players", File("index.html")))` overrides one page of the login.
  - Responses: `File(name)`, `JSON(body)`, `Redirect(location)`, `Status(code)`.
- **Unrouted requests fail.** Any request without a route fails the test when it ends.
- **Running.** `Run` executes the binary in an empty directory with a clean environment:
  `MZ_URL`, `MZ_USERNAME=testmanager`, `MZ_PASSWORD=Secret`, `TZ=Europe/Tallinn`,
  `MZ_NOW=2026-11-20T10:00:00+02:00`. Change these through `mz.Env` before `Run`.
- **Checks.** `WantExit`, `WantStdout`, `WantStdoutGolden`, `WantStderrContains`,
  `WantStderrLacks`, `mz.WantRequests`, `mz.Find`.
  Check what an error means, not Go's exact wording. A rewrite should not have
  to reproduce `parsing time "x" as "2006-01-02"`.

`MZ_NOW` is the binary's only test hook. Date-dependent commands such as
`friendly auto` use it instead of the clock, and any reimplementation must
support it.

## Fixtures

`testdata/mz/` holds real pages recorded on 2026-11-20 and scrubbed: forum pages (`forum-*`), the
logged-out front page, the clubhouse, the players page, the upcoming and played
match lists, and the quick-challenge page. `players-new-layout.html` and
`forum-topics-10-one.html` (forum 10's list cut to one topic) are hand-made. The friendly-invite JSON is written inline in the tests; only
`{"response": ...}` is known to matter.

### Recording

```sh
make build fakemz
e2e/record.sh team-page raw '?p=team&tid=849080'   # real credentials from .env
```

`record.sh`:
1. Runs mz through `fakemz record`, a proxy to the real site.
2. Saves each response into `testdata/mz/<name>/` and prints `Route(...)` lines to paste into a test.
3. Scrubs the saved files.

The proxy forwards only GETs and the login, so recording never changes
anything on the site (`-allow-post` lifts that). Recordings include the login
pages again; delete them and use `Login()`.

### Scrubbing personal data

`e2e/scrub.local` (gitignored) lists `real => fake` pairs: username, team name,
team id. `record.sh` applies it to the files and their names. Before committing,
check:

```sh
e2e/bin/fakemz scrub -map e2e/scrub.local -check e2e/testdata
```

Other teams' names and ids and player ids are public game data and stay as
recorded.
