package main

import (
	"strings"
	"testing"

	"github.com/h2non/gock"
	. "github.com/stretchr/testify/assert"
)

func TestWebClient_Login(t *testing.T) { // https://github.com/h2non/gock
	defer gock.Off()

	// Reference https://github.com/h2non/gock/blob/master/matchers_test.go

	// given
	var client = client()

	gock.New("http://fake.mahnem.ru").
		Post("/").
		MatchParam("module", "login").
		Reply(200).
		File(file("login_failed.html"))

	// when
	err := client.Login()

	// then
	NotNil(t, err)
	True(t, strings.HasPrefix(err.Error(), "Login failed"))

}

/*
func TestWebClient_Logout(t *testing.T) {}
func TestWebClient_Profile(t *testing.T) {}
func TestWebClient_Photos(t *testing.T) {}
*/

func client() WebClient {

	config := SiteConfig{
		URL:      "http://fake.mahnem.ru",
		Login:    "login",
		Password: "pwd",
	}

	var err error
	if client, err := NewWebClient(config); err == nil {
		return client.(WebClient)
	}
	panic(err)
}

func file(filename string) string {
	return "testdata/client/" + filename
}
