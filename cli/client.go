package cli

import (
	"bufio"
	"errors"
	"fmt"
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

// sport returns the MZ sport name from --sport, else MZ_SPORT, else football.
func sport() (string, error) {
	s := sportFlag
	if s == "" {
		s = os.Getenv("MZ_SPORT")
	}
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "football", "soccer":
		return mzhttpclient.Football, nil
	case "hockey":
		return mzhttpclient.Hockey, nil
	}
	return "", fmt.Errorf("unknown sport %q (use football or hockey)", s)
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
	sp, err := sport()
	if err != nil {
		return nil, err
	}
	c := mzhttpclient.New(baseURL).WithSport(sp)
	if err := c.Login(username, password); err != nil {
		return nil, err
	}
	return c, nil
}
