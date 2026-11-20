package e2e

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"
)

func TestWhoami(t *testing.T) {
	mz := Start(t, Login())

	res := mz.Run("whoami")

	res.WantExit(t, 0)
	res.WantStdout(t, "user: testmanager\nteam: 100001\n")
}

func TestLoginProtocol(t *testing.T) {
	mz := Start(t, Login())

	mz.Run("whoami").WantExit(t, 0)

	// The front page hands out the session; every later request carries it,
	// including following the login redirect.
	mz.WantRequests(
		"GET / (no session)",
		"POST /?p=login",
		"GET /?p=clubhouse",
		"GET /?p=players",
		"GET /?p=players",
	)

	login := mz.Find("POST", "/?p=login")[0].Form
	// MZ's login form lowercases the password and sends an HMAC-MD5 of it
	// keyed with a constant from the site's JavaScript.
	h := hmac.New(md5.New, []byte("g372=RVx6!yn"))
	h.Write([]byte(strings.ToLower(password)))
	for k, want := range map[string]string{
		"logindata[username]":      username,
		"logindata[password]":      strings.ToLower(password),
		"logindata[md5]":           hex.EncodeToString(h.Sum(nil)),
		"logindata[sport]":         "soccer",
		"logindata[markasdefault]": "",
		"logindata[remember_me]":   "",
	} {
		if got := login.Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
}

func TestLoginWrongPassword(t *testing.T) {
	// MZ answers a failed login with 200 too; the tell is that the players
	// page shows the logged-out front page.
	mz := Start(t, Login(), Route("GET", "/?p=players", File("index.html")))

	res := mz.Run("whoami")

	res.WantExit(t, 1)
	res.WantStdout(t, "")
	res.WantStderrContains(t, "login failed")
}

func TestMissingCredentials(t *testing.T) {
	mz := Start(t)
	mz.Env["MZ_USERNAME"] = ""

	res := mz.Run("whoami")

	res.WantExit(t, 1)
	res.WantStderrContains(t, "MZ_USERNAME")
	mz.WantRequests() // gives up before contacting the site
}
