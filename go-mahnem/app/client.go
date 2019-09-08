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

	loginURL := NewEndpointBuilder(wc.Config.BaseURL).
		WithPath("/").
		WithQueryParam("module", "login").
		String()

	login := wc.Config.Login
	passwd := wc.Config.Password

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

	logoutURL := NewEndpointBuilder(wc.Config.BaseURL).
		WithPath("/").
		WithQueryParam("module", "quit").
		String()

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
		selTable    = "html body table.pagew tbody tr td.pagew table.t tbody tr td table.t tbody tr"
	)

	profileURL := NewEndpointBuilder(wc.Config.BaseURL).
		WithPath(fmt.Sprintf("/web/%s", user.Profile)).
		String()

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
	/*
		var langs []string
		doc.Find(selLanguages).Each(func(i int, sel *goquery.Selection) {
			langs = append(langs, sel.Contents().Text())
		})
		user.Languages = &langs

		user.Motto = doc.Find(selMotto).Last().Contents().Text()
	*/

	table := make(map[string]string)
	doc.Find(selTable).Each(func(i int, sel *goquery.Selection) {
		k := sel.Find("td.grey")
		if key := k.Contents().Text(); len(key) > 0 {
			table[key] = k.Siblings().Contents().Text()
		}
	})

	// TODO Process map

	return nil
}

// Photos fetch user photos
func (wc WebClient) Photos(user *User) error {

	const (
		selPhotos = "html body table.pagew tbody tr td.pagew div.pg-lst"
	)

	photosURL := NewEndpointBuilder(wc.Config.BaseURL).
		WithPath(fmt.Sprintf("/photo/%s", user.Profile)).
		String()

	response, err := wc.client.Get(photosURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if doc, err := goquery.NewDocumentFromReader(response.Body); err == nil {
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

	return err
}
