package e2e

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	mzBinary = flag.String("mz", "../mz", "the mz binary under test")
	update   = flag.Bool("update", false, "rewrite testdata/golden files from actual output")
)

// Defaults for every run; a test changes them through MZ.Env.
const (
	defaultNow = "2026-11-20T10:00:00+02:00" // a Friday
	defaultTZ  = "Europe/Tallinn"
	username   = "testmanager"
	password   = "Secret" // mz must lowercase it
	myTeam     = "100001" // team id in the scrubbed fixtures
)

// MZ is a fake ManagerZone plus the environment the binary runs with.
type MZ struct {
	t    *testing.T
	fake *fakeMZ
	url  string
	// Env is passed to the binary on top of a clean environment.
	Env map[string]string
}

// Start serves the routes for the rest of the test. When the test ends it
// fails if the binary made a request that no route answers.
func Start(t *testing.T, routes ...Routes) *MZ {
	t.Helper()
	var specs []routeSpec
	for _, rs := range routes {
		for _, r := range rs {
			for i, resp := range r.responses {
				if resp.file == "" {
					continue
				}
				body, err := os.ReadFile(filepath.Join("testdata", "mz", resp.file))
				if err != nil {
					t.Fatal(err)
				}
				resp.Body = body
				if resp.Header == nil {
					resp.Header = map[string][]string{"Content-Type": {contentType(resp.file)}}
				}
				r.responses[i] = resp
			}
			specs = append(specs, r)
		}
	}
	fake := newFakeMZ(specs)
	srv := httptest.NewServer(fake)
	t.Cleanup(func() {
		srv.Close()
		for _, r := range fake.unmatched {
			t.Errorf("request without a route: %s", r.describe())
		}
	})
	return &MZ{
		t:    t,
		fake: fake,
		url:  srv.URL + "/",
		Env: map[string]string{
			"TZ":          defaultTZ,
			"MZ_NOW":      defaultNow,
			"MZ_USERNAME": username,
			"MZ_PASSWORD": password,
		},
	}
}

// Result is what one run of the binary did.
type Result struct {
	Args   []string
	Stdout string
	Stderr string
	Code   int
}

// Run runs the binary in an empty directory with a clean environment.
func (mz *MZ) Run(args ...string) *Result {
	t := mz.t
	t.Helper()
	bin, err := filepath.Abs(*mzBinary)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("mz binary: %v (run 'make build' or pass -mz)", err)
	}
	dir := t.TempDir()
	env := []string{"HOME=" + dir, "PATH=/usr/bin:/bin", "MZ_URL=" + mz.url}
	for k, v := range mz.Env {
		env = append(env, k+"="+v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir, cmd.Env = dir, env
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()

	res := &Result{Args: args, Stdout: stdout.String(), Stderr: stderr.String()}
	var exit *exec.ExitError
	switch {
	case errors.As(err, &exit):
		res.Code = exit.ExitCode()
	case err != nil:
		t.Fatalf("mz %s: %v", strings.Join(args, " "), err)
	}
	return res
}

// Requests returns the requests the binary made so far, in order.
func (mz *MZ) Requests() []Request {
	mz.fake.mu.Lock()
	defer mz.fake.mu.Unlock()
	return append([]Request(nil), mz.fake.requests...)
}

// Find returns the requests whose method matches and whose URL matches the
// pattern ("*" is a wildcard).
func (mz *MZ) Find(method, pattern string) []Request {
	p := normalizePattern(pattern)
	var out []Request
	for _, r := range mz.Requests() {
		if r.Method == method && globMatch(p, r.URL) {
			out = append(out, r)
		}
	}
	return out
}

// WantRequests checks the exact request sequence (Request.String form).
func (mz *MZ) WantRequests(want ...string) {
	mz.t.Helper()
	var got []string
	for _, r := range mz.Requests() {
		got = append(got, r.String())
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		mz.t.Errorf("requests:\n  %s\nwant:\n  %s", strings.Join(got, "\n  "), strings.Join(want, "\n  "))
	}
}

func (r *Result) cmd() string { return "mz " + strings.Join(r.Args, " ") }

// WantExit checks the exit status.
func (r *Result) WantExit(t *testing.T, code int) {
	t.Helper()
	if r.Code != code {
		t.Errorf("%s: exit status %d, want %d\nstdout:\n%s\nstderr:\n%s", r.cmd(), r.Code, code, r.Stdout, r.Stderr)
	}
}

// WantStdout checks stdout exactly.
func (r *Result) WantStdout(t *testing.T, want string) {
	t.Helper()
	if r.Stdout != want {
		t.Errorf("%s: stdout differs:\n%s", r.cmd(), diff(want, r.Stdout))
	}
}

// WantStdoutGolden checks stdout against testdata/golden/<name>; -update rewrites it.
func (r *Result) WantStdoutGolden(t *testing.T, name string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if *update {
		if err := os.WriteFile(path, []byte(r.Stdout), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run 'make e2e-update' to create it)", err)
	}
	if string(want) != r.Stdout {
		t.Errorf("%s: stdout differs from %s (if the change is right, run 'make e2e-update'):\n%s",
			r.cmd(), path, diff(string(want), r.Stdout))
	}
}

// WantStderrContains checks that stderr mentions each part. Tests check the
// meaning of an error, not its exact wording, which is implementation detail.
func (r *Result) WantStderrContains(t *testing.T, parts ...string) {
	t.Helper()
	for _, p := range parts {
		if !strings.Contains(r.Stderr, p) {
			t.Errorf("%s: stderr does not mention %q:\n%s", r.cmd(), p, r.Stderr)
		}
	}
}

// WantStderrLacks checks that stderr does not contain any of the parts.
func (r *Result) WantStderrLacks(t *testing.T, parts ...string) {
	t.Helper()
	for _, p := range parts {
		if strings.Contains(r.Stderr, p) {
			t.Errorf("%s: stderr should not contain %q:\n%s", r.cmd(), p, r.Stderr)
		}
	}
}

// diff shows the first differing line with some context.
func diff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	i := 0
	for i < len(w) && i < len(g) && w[i] == g[i] {
		i++
	}
	from, to := max(0, i-3), i+4
	var b strings.Builder
	fmt.Fprintf(&b, "first difference at line %d\n", i+1)
	for j := from; j < to && j < len(w); j++ {
		fmt.Fprintf(&b, "  want %4d | %s\n", j+1, w[j])
	}
	for j := from; j < to && j < len(g); j++ {
		fmt.Fprintf(&b, "  got  %4d | %s\n", j+1, g[j])
	}
	return b.String()
}
