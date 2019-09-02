package main

import (
	"log"
)

func main() {

	//log.Printf("App config is %+v\n", *GetAppConfig())

	rep, err := NewRepository()
	if err != nil {
		log.Fatal("Db connection init failed", err.Error())
	}
	defer rep.Close()

	// TODO Remove later
	/*
		rep.deleteAllUserPhotos()
		rep.deleteAllLanguages()
		rep.deleteAllUsers()
		rep.deleteAllLocations()
		rep.deleteAllUserLanguages()
	*/

	const (
		nickname = "_760112"
		//nickname = "evilcat777"
		//nickname = "_760110"
	)

	client, err := NewWebClient(GetAppConfig().Site)
	if err != nil {
		log.Fatal("Web client init failed", err.Error())
	}

	err = client.Login()
	if err != nil {
		log.Fatal("Login failed", err.Error())
	}
	defer client.Logout()

	var user = &User{Profile: nickname}

	err = client.Profile(user)
	if err != nil {
		log.Fatal("Profile fetch failed, error ", err.Error())
	}

	err = client.Photos(user)
	if err != nil {
		log.Fatal("Photos fetch failed, error ", err.Error())
	}

	// TODO Add locId to Location struct
	var locationID uint64
	if locationID = rep.FindLocation(user.Location.Country, user.Location.City); locationID == 0 {
		locationID = rep.StoreLocation(user.Location.Country, user.Location.City)
	}

	// TODO Add userId to User struct
	var userID uint64
	if userID = rep.FindUserByLogin(user.Profile); userID == 0 {
		userID = rep.StoreUser(user.Profile, user.Name, locationID, user.Motto)
	}

	// TODO Add langId to Languages struct
	if user.Languages != nil {
		for _, lang := range *user.Languages {
			var languageID uint64
			if languageID = rep.FindLanguageByName(lang); languageID == 0 {
				languageID = rep.StoreLanguage(lang)
			}
			rep.StoreUserLanguage(userID, languageID)
		}
	}

	// TODO Define UserPhoto struct
	if user.Photos != nil {
		for _, photo := range *user.Photos {
			rep.StoreUserPhoto(userID, photo)
		}
	}

	rep.PrintStats()
}
