package e2e

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests of the fake server itself; they do not need the mz binary.

func TestFakeNormalizeSortsQuery(t *testing.T) {
	if a, b := normalizeURL("/x.php?b=2&a=1"), normalizeURL("/x.php?a=1&b=2"); a != b {
		t.Errorf("%q != %q", a, b)
	}
	if got := normalizeURL(""); got != "/" {
		t.Errorf("empty: %q", got)
	}
}

func TestFakeGlob(t *testing.T) {
	for _, c := range []struct {
		p, s string
		want bool
	}{
		{"/a?tid=*", "/a?tid=12", true},
		{"/a?tid=*&x=1", "/a?tid=12&x=1", true},
		{"/a?tid=*&x=1", "/a?tid=12&x=2", false},
		{"/a", "/a", true},
		{"/a", "/ab", false},
		{"*", "", true},
	} {
		if got := globMatch(c.p, c.s); got != c.want {
			t.Errorf("globMatch(%q, %q) = %v", c.p, c.s, got)
		}
	}
}

func TestFakeServer(t *testing.T) {
	f := newFakeMZ(join(
		Route("GET", "/?p=players", JSON("old")),
		Route("GET", "/?p=players", JSON("override")),
		Route("POST", "/ajax.php?p=x&tid=*", JSON("1"), Status(500)),
	))
	ts := httptest.NewServer(f)
	defer ts.Close()

	do := func(method, path, body string) (int, string) {
		t.Helper()
		req, _ := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
		if body != "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(b)
	}

	if _, b := do("GET", "/?p=players", ""); b != "override" {
		t.Errorf("override: %q", b)
	}
	if s, _ := do("POST", "/ajax.php?tid=5&p=x", "a=1"); s != 200 {
		t.Errorf("sequence 1: %d", s)
	}
	for i := 0; i < 2; i++ {
		if s, _ := do("POST", "/ajax.php?p=x&tid=6", ""); s != 500 {
			t.Errorf("sequence 2+ (last repeats): %d", s)
		}
	}
	if s, _ := do("GET", "/nope", ""); s != 404 {
		t.Errorf("unmatched: %d", s)
	}
	if len(f.unmatched) != 1 || f.unmatched[0].URL != "/nope" {
		t.Errorf("unmatched: %v", f.unmatched)
	}
	if r := f.requests[1]; r.URL != "/ajax.php?p=x&tid=5" || r.Form.Get("a") != "1" {
		t.Errorf("request log: %+v", r)
	}
}
