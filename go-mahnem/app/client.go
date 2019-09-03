package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Configuration configuration
type Configuration struct {
	BaseURL  string
	Login    string
	Password string
}

// WebClient is basic web client struct
type WebClient struct { // TODO We have duplicate struct in config.go
	Config *Configuration
	client *http.Client
}

// MahnemClient defenition of web client
type MahnemClient interface { // TODO Check documentation how to implement interfaces
	Login() error
	Profile(*User) error
	Photos(*User) error
	Logout() error
}

func mapOf(key, value string) map[string]string {
	return map[string]string{key: value}
}

func (wc WebClient) url(path string, params map[string]string) string {

	reqURL, err := url.ParseRequestURI(wc.Config.BaseURL)
	if err != nil {
		log.Printf("Can't parse base site url '%s', error '%s'\n", wc.Config.BaseURL, err.Error())
	}

	reqURL.Path = path

	if params != nil {
		reqParams := url.Values{}
		for k, v := range params {
			if len(k) > 0 && len(v) > 0 {
				reqParams.Add(k, v)
			}
		}
		if len(reqParams) > 0 {
			reqURL.RawQuery = reqParams.Encode()
		}
	}

	return reqURL.String()
}

// NewWebClient creates new web client
func NewWebClient(config SiteConfig) (MahnemClient, error) { // TODO Fix to *Mahneclientlient

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	webClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	return WebClient{
		Config: &Configuration{
			Login:    config.Login,
			Password: config.Password,
			BaseURL:  config.URL,
		},
		client: webClient,
	}, nil
}

// Login log in
func (wc WebClient) Login() error {

	login := wc.Config.Login
	passwd := wc.Config.Password

	loginURL := wc.url("/", mapOf("module", "login"))

	log.Printf("Login url: %s, username: %s, password: %s\n", loginURL, login, passwd)

	form := url.Values{}
	form.Set("logon", login)
	form.Set("pwd", passwd)

	response, err := wc.client.PostForm(loginURL, form)
	if err != nil {
		return fmt.Errorf("Login failed: %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 302 {
		return fmt.Errorf("Login failed, status code: %d", response.StatusCode)
	}

	// TODO check body

	return nil
}

// Logout log out
func (wc WebClient) Logout() error {

	logoutURL := wc.url("/", mapOf("module", "quit"))

	log.Printf("Logout url: %s\n", logoutURL)

	response, err := wc.client.Get(logoutURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 302 {
		return fmt.Errorf("Logout status code is %d", response.StatusCode)
	}

	return nil
}

// Profile fetch user profile
func (wc WebClient) Profile(user *User) error {

	const (
		selLocation = "html body table.pagew tbody tr td.pagew table.t tbody tr td div img[src='https://img.zhivem.ru/pic_location.png']"

		selUsername = "html body table.pagew tbody tr td.pagew table.t tbody tr td div.header2"

		selLanguages = "html body table.pagew tbody tr td.pagew table.t tbody tr td a.black"

		selMotto = "html body table.pagew tbody tr td.pagew table.t tbody tr td table.t tbody tr td"
	)

	profileURL := wc.url(fmt.Sprintf("/web/%s", user.Profile), nil)

	response, err := wc.client.Get(profileURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	user.Location = newLocation(doc.Find(selLocation).First().Parent().Contents().Text()) // We need parent contents

	user.Name = strip(doc.Find(selUsername).First().Contents().Text())

	var langs []string
	doc.Find(selLanguages).Each(func(i int, sel *goquery.Selection) {
		langs = append(langs, sel.Contents().Text())
	})
	user.Languages = &langs

	user.Motto = doc.Find(selMotto).Last().Contents().Text()

	return nil
}

// Photos fetch user photos
func (wc WebClient) Photos(user *User) error {

	const (
		selPhotos = "html body table.pagew tbody tr td.pagew div.pg-lst"
	)

	photosURL := wc.url(fmt.Sprintf("/photo/%s", user.Profile), nil)

	response, err := wc.client.Get(photosURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	doc.Find(selPhotos).SiblingsFiltered("script").Each(func(i int, sel *goquery.Selection) {
		raw := sel.Contents().Text()
		if strings.Contains(raw, "PS=[") {
			raw = strings.Split(strings.Split(raw, "PS=[")[1], "]")[0]
			data := strings.Split(strings.ReplaceAll(raw, "'", ""), ",")
			user.Photos = &data
		}
	})

	return nil
}
