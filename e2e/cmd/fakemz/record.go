package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// record is a reverse proxy to the real ManagerZone that saves every response
// into a directory under e2e/testdata/mz and appends a Route(...) line for it
// to routes.txt there, ready to paste into a test. Point mz at it with
// MZ_URL=http://ADDR/.
//
// POSTs other than the login are refused unless -allow-post is given: they
// change things on the real site (e.g. send friendly invites).
func recordMain(args []string) {
	fs := flag.NewFlagSet("record", flag.ExitOnError)
	upstream := fs.String("upstream", "https://www.managerzone.com/", "real site")
	dir := fs.String("dir", "", "directory to write responses and routes.txt into (required)")
	ref := fs.String("ref", "", "how tests name the directory in File(...) (default: its base name)")
	addr := fs.String("addr", "127.0.0.1:8090", "listen address")
	allowPost := fs.Bool("allow-post", false, "forward POSTs other than the login")
	fs.Parse(args)
	if *dir == "" {
		fs.Usage()
		os.Exit(2)
	}
	up, err := url.Parse(*upstream)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		log.Fatal(err)
	}
	if *ref == "" {
		*ref = filepath.Base(*dir)
	}
	rec := &recorder{dir: *dir, ref: *ref, next: nextFileNumber(*dir)}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(up)
			pr.Out.Host = up.Host
			// Let the transport negotiate gzip and decode it, so files are plain.
			pr.Out.Header.Del("Accept-Encoding")
		},
		ModifyResponse: func(resp *http.Response) error {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}
			resp.Body = io.NopCloser(bytes.NewReader(body))
			// Keep redirects inside the proxy.
			if loc := resp.Header.Get("Location"); loc != "" {
				if lu, err := url.Parse(loc); err == nil && lu.Host == up.Host {
					resp.Header.Set("Location", lu.RequestURI())
				}
			}
			return rec.save(resp.Request.Method, requestURI(resp.Request), resp, body)
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/_/health", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok\n") })
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && !*allowPost && !(r.Method == http.MethodPost && r.URL.RawQuery == "p=login") {
			log.Printf("record: refused %s %s (use -allow-post to forward)", r.Method, r.URL.RequestURI())
			http.Error(w, "fakemz record: POST refused, use -allow-post", http.StatusForbidden)
			return
		}
		proxy.ServeHTTP(w, r)
	})
	log.Printf("fakemz: recording %s into %s on %s", up, *dir, *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

type recorder struct {
	dir  string
	ref  string
	mu   sync.Mutex
	next int
}

func (rec *recorder) save(method, uri string, resp *http.Response, body []byte) error {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	name := fmt.Sprintf("%03d-%s%s", rec.next, slug(uri), extFor(resp.Header.Get("Content-Type")))
	rec.next++
	if err := os.WriteFile(filepath.Join(rec.dir, name), body, 0o644); err != nil {
		return err
	}
	answer := fmt.Sprintf("File(%q)", rec.ref+"/"+name)
	switch loc := resp.Header.Get("Location"); {
	case resp.StatusCode >= 300 && resp.StatusCode < 400 && loc != "":
		answer = fmt.Sprintf("Redirect(%q)", loc)
	case resp.StatusCode != http.StatusOK:
		answer += " /* status " + strconv.Itoa(resp.StatusCode) + " */"
	}
	line := fmt.Sprintf("Route(%q, %q, %s),", method, uri, answer)
	f, err := os.OpenFile(filepath.Join(rec.dir, "routes.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, line)
	log.Printf("record: %s (%d bytes)", line, len(body))
	return err
}

func requestURI(r *http.Request) string {
	u := r.URL.EscapedPath()
	if r.URL.RawQuery != "" {
		u += "?" + r.URL.RawQuery
	}
	return u
}

var reNumbered = regexp.MustCompile(`^(\d{3})-`)

func nextFileNumber(dir string) int {
	n := 1
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if m := reNumbered.FindStringSubmatch(e.Name()); m != nil {
			if v, _ := strconv.Atoi(m[1]); v >= n {
				n = v + 1
			}
		}
	}
	return n
}

var reNonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slug(uri string) string {
	s := strings.Trim(reNonSlug.ReplaceAllString(strings.ToLower(uri), "-"), "-")
	if s == "" {
		s = "index"
	}
	if len(s) > 50 {
		s = strings.TrimRight(s[:50], "-")
	}
	return s
}

func extFor(contentType string) string {
	mt, _, _ := mime.ParseMediaType(contentType)
	switch {
	case mt == "text/html":
		return ".html"
	case strings.HasSuffix(mt, "json"):
		return ".json"
	case strings.HasSuffix(mt, "xml"):
		return ".xml"
	case strings.HasPrefix(mt, "text/"):
		return ".txt"
	}
	return ".bin"
}
