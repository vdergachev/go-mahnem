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
	profile(*User) error
	photos(*User) error
	logout() error
}

// User struct populated bu Mahneclientlient::profile method
type User struct {
	Profile  string
	Name     string
	Location string
	Photos   *[]string
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

	if response.StatusCode != 302 {
		return fmt.Errorf("Login status code is %d", response.StatusCode)
	}

	// TODO Parse file and chech conent - no errors
	fmt.Println(" success [OK]")
	dumpResponse("login.html", response.Body)

	return nil
}

// Mahneclientlient :: login
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
func (wc WebClient) profile(user *User) error {

	profileURL := wc.url(fmt.Sprintf("/web/%s", user.Profile))

	response, err := wc.client.Get(profileURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	LocationSel := "html body table.pagew tbody tr td.pagew table.t tbody tr td div img[src='https://img.zhivem.ru/pic_location.png']"
	doc.Find(LocationSel).Each(func(i int, sel *goquery.Selection) {
		user.Location = strip(sel.Parent().Contents().Text()) // We need parent contents
	})

	UsernameSel := "html body table.pagew tbody tr td.pagew table.t tbody tr td div.header2"
	doc.Find(UsernameSel).Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			user.Name = strip(sel.Contents().Text())
		}
	})

	return nil
}

func (wc WebClient) photos(user *User) error {

	photosURL := wc.url(fmt.Sprintf("/photo/%s", user.Profile))

	response, err := wc.client.Get(photosURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	//dumpResponse("photos.html", response.Body)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	PhotosSel := "html body table.pagew tbody tr td.pagew div.pg-lst"
	doc.Find(PhotosSel).SiblingsFiltered("script").Each(func(i int, sel *goquery.Selection) {
		raw := sel.Contents().Text()
		if strings.Contains(raw, "PS=[") {
			raw = strings.Split(strings.Split(raw, "PS=[")[1], "]")[0]
			data := strings.Split(strings.ReplaceAll(raw, "'", ""), ",")
			user.Photos = &data
		}
	})

	return nil
}

// --------------------------------------------------------------------------------------------------
func main() {

	var nickname = "_760112"

	var client = defaultClient()

	err := client.init()
	if err != nil {
		log.Fatal("Can't init web client", err.Error())
	}

	err = client.login()
	if err != nil {
		log.Fatal("Login failed", err.Error())
	}
	defer client.logout()

	var user = &User{Profile: nickname}

	err = client.profile(user)
	if err != nil {
		log.Fatal("Profile fetch failed, error ", err.Error())
	}

	err = client.photos(user)
	if err != nil {
		log.Fatal("Photos fetch failed, error ", err.Error())
	}

	fmt.Println("profile: " + user.Name)

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

func strip(val string) string {
	re, err := regexp.Compile(`\r?\n`)
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(strings.ReplaceAll(val, "\t", ""), "")
}
