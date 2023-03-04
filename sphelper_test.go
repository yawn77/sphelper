package sphelper

import (
	"flag"
	"strings"
	"testing"
)

// credentials for command-line tests
var username = flag.String("username", "", "Username")
var password = flag.String("password", "", "User password")

// get credentials either from env or from cmdline if env does not exists
func getCredentials() Credentials {
	cred, err := GetCredentials()
	if err != nil {
		cred = Credentials{
			Username:    *username,
			Password:    *password,
			PasswordMd5: GetMD5Hash(*password),
		}
	}
	return cred
}

func TestLoginAndLogout(t *testing.T) {
	// arrange
	cred := getCredentials()
	c, err := GetClient()
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	// act
	err = c.Login(cred)

	// assert
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	s := session{}
	if c.Session == s {
		t.Fatalf("session could not be set %v", c.Session)
	}

	// act
	err = c.Logout()

	// assert
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}
	if c.Session != s {
		t.Fatalf("session was not cleared after logout: %v", c.Session)
	}
}

func TestLoginFails(t *testing.T) {
	// arrange
	cred := Credentials{Username: "", Password: "", PasswordMd5: ""}
	c, err := GetClient()
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	// act
	err = c.Login(cred)

	// assert
	if err != loginFailed {
		t.Fatalf("%v != %v", err, loginFailed)
	}
}

func TestGet(t *testing.T) {
	// arrange
	c, err := GetClient()
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	// act
	body, err := c.Get(urlRoot)

	// assert
	if err != nil {
		t.Fatalf("could not get %s: %v", urlRoot, err)
	}
	if !strings.Contains(body, "<span class=\"smallfont\"><strong>&raquo; <a href=\"/forum.php\">Aktuelle Threads</a></strong></span></td>") {
		t.Fatalf("unexpected page content for %s", urlRoot)
	}
}
