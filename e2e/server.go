// Package e2e holds black-box tests for the mz binary. Tests run the built
// binary against a fake ManagerZone (this file) that answers from saved real
// responses in testdata/mz, then check its output and the requests it sent.
// Nothing here imports mz, so mz can be rewritten in any language.
package e2e

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
)

// Response is what the fake server answers to one request.
type Response struct {
	Status int
	Header http.Header
	Body   []byte
	file   string // under testdata/mz; loaded by Start
}

// File answers with a saved response from testdata/mz.
func File(name string) Response { return Response{Status: http.StatusOK, file: name} }

// JSON answers with a JSON body.
func JSON(body string) Response {
	return Response{Status: http.StatusOK, Body: []byte(body), Header: http.Header{"Content-Type": {"application/json"}}}
}

// Redirect answers with 302 to loc.
func Redirect(loc string) Response {
	return Response{Status: http.StatusFound, Header: http.Header{"Location": {loc}}}
}

// Status answers with an empty body and the given status.
func Status(code int) Response { return Response{Status: code} }

// Routes are request patterns with their responses. A later route with the
// same method and URL replaces an earlier one, so a test can override part of
// a group such as Login().
type Routes []routeSpec

type routeSpec struct {
	method, url string
	responses   []Response
}

// Route answers METHOD URL. URL is path plus query; query order does not
// matter and "*" matches anything. With several responses the n-th matching
// request gets the n-th response and the last one repeats.
func Route(method, url string, responses ...Response) Routes {
	if len(responses) == 0 {
		panic("Route " + method + " " + url + ": no response")
	}
	return Routes{{method: method, url: url, responses: responses}}
}

// Request is one request the binary made.
type Request struct {
	Method  string
	URL     string     // normalized: query sorted by key
	Form    url.Values // parsed body of a form POST, else nil
	Body    string
	Session bool // had a PHPSESSID cookie
}

func (r Request) String() string {
	s := r.Method + " " + r.URL
	if !r.Session {
		s += " (no session)"
	}
	return s
}

type route struct {
	method, url string
	glob        bool
	responses   []Response
	served      int
}

type fakeMZ struct {
	mu        sync.Mutex
	routes    []*route
	requests  []Request
	unmatched []Request
}

func newFakeMZ(specs []routeSpec) *fakeMZ {
	f := &fakeMZ{}
	index := map[string]int{}
	for _, s := range specs {
		rt := &route{method: strings.ToUpper(s.method), url: normalizePattern(s.url), responses: s.responses}
		rt.glob = strings.Contains(rt.url, "*")
		key := rt.method + " " + rt.url
		if i, ok := index[key]; ok {
			f.routes[i] = rt
			continue
		}
		index[key] = len(f.routes)
		f.routes = append(f.routes, rt)
	}
	return f
}

func (f *fakeMZ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	req := Request{Method: r.Method, URL: normalizeURL(r.URL.RequestURI()), Body: string(body)}
	if c, err := r.Cookie("PHPSESSID"); err == nil && c.Value != "" {
		req.Session = true
	}
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mt == "application/x-www-form-urlencoded" {
		req.Form, _ = url.ParseQuery(req.Body)
	}

	f.mu.Lock()
	rt := f.match(req.Method, req.URL)
	if rt == nil {
		f.unmatched = append(f.unmatched, req)
		f.mu.Unlock()
		http.Error(w, "fake mz: no route for "+req.Method+" "+req.URL, http.StatusNotFound)
		return
	}
	f.requests = append(f.requests, req)
	resp := rt.responses[min(rt.served, len(rt.responses)-1)]
	rt.served++
	f.mu.Unlock()

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	http.SetCookie(w, &http.Cookie{Name: "PHPSESSID", Value: "test", Path: "/"})
	w.WriteHeader(resp.Status)
	w.Write(resp.Body)
}

// match returns the route for a request; exact routes win over globs.
func (f *fakeMZ) match(method, nurl string) *route {
	for _, rt := range f.routes {
		if !rt.glob && rt.method == method && rt.url == nurl {
			return rt
		}
	}
	for _, rt := range f.routes {
		if rt.glob && rt.method == method && globMatch(rt.url, nurl) {
			return rt
		}
	}
	return nil
}

func contentType(file string) string {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".json":
		return "application/json"
	case ".xml":
		return "text/xml; charset=utf-8"
	}
	return "text/plain; charset=utf-8"
}

// normalizeURL turns "path?query" into a canonical form with sorted query keys.
func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	p := u.EscapedPath()
	if p == "" {
		p = "/"
	}
	q, _ := url.ParseQuery(u.RawQuery) // keep whatever parses
	if len(q) > 0 {
		p += "?" + q.Encode()
	}
	return p
}

func normalizePattern(raw string) string {
	return strings.ReplaceAll(normalizeURL(raw), "%2A", "*")
}

// globMatch matches s against pattern where "*" is any run of characters.
func globMatch(pattern, s string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == s
	}
	first, last := parts[0], parts[len(parts)-1]
	if len(s) < len(first)+len(last) || !strings.HasPrefix(s, first) || !strings.HasSuffix(s, last) {
		return false
	}
	s = s[len(first) : len(s)-len(last)]
	for _, p := range parts[1 : len(parts)-1] {
		i := strings.Index(s, p)
		if i < 0 {
			return false
		}
		s = s[i+len(p):]
	}
	return true
}

func (r Request) describe() string {
	var b strings.Builder
	b.WriteString(r.String())
	if r.Body != "" {
		fmt.Fprintf(&b, "\n    body: %s", r.Body)
	}
	return b.String()
}
