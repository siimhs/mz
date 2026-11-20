package cli

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/siimhs/mz/mzhttpclient"
)

// loadDotEnv sets variables from a KEY=VALUE file without overriding ones already set.
// A missing file is not an error.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		if _, set := os.LookupEnv(k); !set {
			os.Setenv(k, v)
		}
	}
}

// newClient returns a client logged in with MZ_USERNAME / MZ_PASSWORD (env or ./.env).
func newClient() (*mzhttpclient.MZHTTPClient, error) {
	loadDotEnv(".env")
	username := os.Getenv("MZ_USERNAME")
	password := os.Getenv("MZ_PASSWORD")
	if username == "" || password == "" {
		return nil, errors.New("MZ_USERNAME and MZ_PASSWORD must be set (environment or .env)")
	}
	baseURL := os.Getenv("MZ_URL")
	if baseURL == "" {
		baseURL = "https://www.managerzone.com/"
	}
	c := mzhttpclient.New(baseURL)
	if err := c.Login(username, password); err != nil {
		return nil, err
	}
	return c, nil
}
