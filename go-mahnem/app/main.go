package main

import (
	"fmt"
	"log"
)

func main() {

	fmt.Printf("App config is %+v\n", *GetAppConfig())

	const (
		nickname = "_760112"
	)

	client, err := newClient()
	if err != nil {
		log.Fatal("Web client init failed", err.Error())
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

	log.Println("profile: " + user.toString())
}
