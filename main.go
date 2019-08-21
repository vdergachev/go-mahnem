package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Configuration configuration
type Configuration struct {
	BaseURL    string
	Login      string
	Password   string
	DumpFolder string
}

// WebClient is basic web client struct
type WebClient struct {
	Config *Configuration
	client *http.Client
}

// Mahneclientlient defenition of web client
type Mahneclientlient interface {
	init() error
	login() error
	profile(string) error
	logout() error
}

func defaultClient() *WebClient {

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	webClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	return &WebClient{
		Config: &Configuration{
			Login:      os.Args[1],
			Password:   os.Args[2],
			BaseURL:    "http://mahnem.ru",
			DumpFolder: "./result",
		},
		client: webClient,
	}
}

// Mahneclientlient :: url
func (wc WebClient) url(path string) string {
	return wc.Config.BaseURL + path // TODO use something sereous than stupid concat
}

// Mahneclientlient :: initStorage
func (wc WebClient) initStorage() error {

	dir := wc.Config.DumpFolder

	_, err := os.Stat(dir)
	if !os.IsExist(err) {
		err = os.RemoveAll(dir)
		if err != nil {
			return fmt.Errorf("Can't remove existing fs result storage %s", err.Error())
		}
	}

	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return fmt.Errorf("Can't init fs result storage")
	}

	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("Can't change working directory to %s", dir)
	}

	return nil
}

// Mahneclientlient :: init
func (wc WebClient) init() error {

	err := wc.initStorage()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Mahneclientlient :: login
// TODO :: We have to check that 302 been recieved and profile is available
func (wc WebClient) login() error {

	loginURL := wc.url("?module=login")
	login := wc.Config.Login
	passwd := wc.Config.Password

	fmt.Printf("Login url: %s, username: %s, password: %s", loginURL, login, passwd)

	form := url.Values{}
	form.Set("logon", login)
	form.Set("pwd", passwd)

	response, err := wc.client.PostForm(loginURL, form)
	if err != nil {
		return fmt.Errorf("Login request failed: %s", err.Error())
	}

	defer response.Body.Close()
	fmt.Println(" success [OK]")

	dumpResponse("login.html", response.Body)

	if response.StatusCode != 302 {
		return fmt.Errorf("Login status code is %d", response.StatusCode)
	}

	return nil
}

// Mahneclientlient :: login
// TODO :: We have to check that 302 been recieved and profile is anavailable
func (wc WebClient) logout() error {

	logoutURL := wc.url("/?module=quit")

	fmt.Printf("Logout url: %s", logoutURL)

	response, err := wc.client.Get(logoutURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 302 {
		return fmt.Errorf("Logout status code is %d", response.StatusCode)
	}

	fmt.Println(" success [OK]")

	return nil
}

// Mahneclientlient :: profile
// Returns true - own profile is available, auth is successful
func (wc WebClient) profile(username string) (bool, error) {

	if len(username) == 0 {
		username = wc.Config.Login
	}

	profileURL := wc.url(fmt.Sprintf("/web/%s", username))

	fmt.Println("Profile url: %s\n", profileURL)

	response, err := wc.client.Get(profileURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	re, err := regexp.Compile(`\r?\n`)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("html body table.pagew tbody tr td.pagew table.t tbody tr td div.header2").Each(func(i int, sel *goquery.Selection) {
		val := strings.ReplaceAll(strings.TrimSpace(sel.Contents().Text()), "\t", "")
		val = re.ReplaceAllString(val, " ")
		if len(username) == 0 {
			return
		}

		fmt.Println(username)
	})

	//fmt.Println(" success [OK]")
	//dumpResponse("profile.html", response.Body)
	return true, nil
}

// --------------------------------------------------------------------------------------------------
func main() {

	//var ProfileURL = "_760112"

	var client = defaultClient()

	err := client.init()
	if err != nil {
		log.Fatal("Can't init web client", err.Error())
	}

	err = client.login()
	if err != nil {
		log.Fatal("Login failed", err.Error())
	}

	v, err := client.profile("_760112")
	if err != nil {
		log.Fatal("Profile fetch failed, error ", err.Error())
	} else if !v {
		log.Fatal("Profile fetch failed")
	}

	// Do something realy usefull ^_^
	defer client.logout()

}

func dumpResponse(filename string, val io.ReadCloser) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Can't create file %s, error: %s\n", filename, err.Error())
	}
	defer file.Close()

	written, err := io.Copy(file, val)
	if err != nil {
		log.Fatalf("Can't save file %s, error: %s\n", filename, err.Error())
	}

	fmt.Printf("\tfile %s (%d bytes) created\n",
		filename,
		written)
}
