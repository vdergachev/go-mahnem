package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

// ResultDIR is place for downloaded html files
var ResultDIR = "./result"

func main() {

	var LoginURL = "http://mahnem.ru?module=login"
	var LogoutURL = "http://www.mahnem.ru/?module=quit"
	var ProfileURL = "http://www.mahnem.ru/?module=posts&user=_760112"

	initResultStorage()

	client, err := initWebClient()
	if err != nil {
		log.Fatal("Can't init web client", err.Error())
	}

	login(client, LoginURL, os.Args[1], os.Args[2], "login.html")
	visitLink(client, ProfileURL, "_760112.html")
	logout(client, LogoutURL)
}

func initWebClient() (*http.Client, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		CheckRedirect: redirect,
		Jar:           jar,
	}

	return client, nil
}

func redirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func login(client *http.Client, loginURL string, login string, password string, filename string) {

	fmt.Printf("Login url: %s, username: %s, password: %s", loginURL, os.Args[1], os.Args[2])

	loginFormData := url.Values{}
	loginFormData.Set("logon", login)
	loginFormData.Set("pwd", password)

	response, err := client.PostForm(loginURL, loginFormData)

	if err != nil {
		log.Fatalf("Login request failed: %s\n", err.Error())
	}

	fmt.Println(" success [OK]")

	defer response.Body.Close()

	dumpResponse(filename, response.Body)
}

func logout(client *http.Client, link string) {

	fmt.Printf("Logout url: %s", link)

	response, err := client.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fmt.Println(" success [OK]")
}

func visitLink(client *http.Client, link string, filename string) {

	response, err := client.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	dumpResponse(filename, response.Body)
}

func initResultStorage() {

	_, err := os.Stat(ResultDIR)
	if !os.IsExist(err) {
		err = os.RemoveAll(ResultDIR)
		if err != nil {
			log.Fatalf("Can't remove existing fs result storage %s\n", err.Error())
		}
	}

	err = os.MkdirAll(ResultDIR, os.ModeDir)
	if err != nil {
		log.Fatal("Can't init fs result storage")
	}

	err = os.Chdir(ResultDIR)
	if err != nil {
		log.Fatal("Can't change working directory to ", ResultDIR)
	}
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
