package sphelper

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var version string

// URLs
const (
	urlRoot  = "https://www.spieleplanet.eu/"
	urlLogin = "https://www.spieleplanet.eu/login.php?do=login"
)

// regex
const (
	securityTokenRegex = "<input type=\"hidden\" name=\"securitytoken\" value=\"(?P<token>.*)\" \\/>"
)

// errors
const (
	loginFailed           = Error("login failed, please check username and password")
	logoutFailed          = Error("logout failed")
	securityTokenNotFound = Error("failed to find security token")
)

func GetVersion() string {
	return version
}

type session struct {
	SecurityToken string
	logoutUrl     string
}

type Client struct {
	client  *http.Client
	Session session
}

func (r *Client) setSession(session session) {
	r.Session = session
}

func GetClient() (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}
	return &Client{client: &http.Client{Jar: jar}}, nil
}

func (c *Client) Login(credentials Credentials) error {
	values := url.Values{
		"vb_login_username":        {credentials.Username},
		"vb_login_password":        {credentials.Password},
		"s":                        {""},
		"securitytoken":            {"guest"},
		"do":                       {"login"},
		"vb_login_md5password":     {credentials.PasswordMd5},
		"vb_login_md5password_utf": {credentials.PasswordMd5},
	}
	body, err := c.Post(urlLogin, values)
	if err != nil {
		return err
	}
	if !strings.Contains(strings.ToUpper(string(body)), fmt.Sprintf("DANKE FÜR DEINE ANMELDUNG, %s.", strings.ToUpper(credentials.Username))) {
		return loginFailed
	}
	session, err := c.getSessionInformation()
	if err != nil {
		return err
	}
	c.setSession(session)
	return nil
}

func (c Client) getSessionInformation() (session, error) {
	body, err := c.Get(urlRoot)
	if err != nil {
		return session{}, err
	}
	re := regexp.MustCompile(securityTokenRegex)
	l := re.FindStringSubmatch(body)
	if len(l) == 0 {
		return session{}, securityTokenNotFound
	}
	token := l[1]
	return session{
		SecurityToken: token,
		logoutUrl:     fmt.Sprintf("%slogin.php?do=logout&logouthash=%s", urlRoot, token),
	}, nil
}

func (c *Client) Logout() error {
	body, err := c.Get(c.Session.logoutUrl)
	if err != nil {
		return err
	}
	if !strings.Contains(body, "Du hast dich erfolgreich vom Forum abgemeldet.") {
		return logoutFailed
	}
	c.Session = session{}
	return nil
}

func (c Client) Get(url string) (string, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c Client) Post(url string, values url.Values) (string, error) {
	resp, err := c.client.PostForm(url, values)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func WriteBody(filename string, body string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(body)
	if err != nil {
		return err
	}
	return nil
}
