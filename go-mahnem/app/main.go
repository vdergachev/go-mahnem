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
	rep.deleteAllUserPhotos()
	rep.deleteAllLanguages()
	rep.deleteAllUsers()
	rep.deleteAllLocations()
	rep.deleteAllUserLanguages()

	const (
		nickname = "_760112"
	)

	client, err := NewWebClient()
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

	// TODO Add locId to Location struct
	var locationID uint64
	if locationID = rep.FindLocation(user.Location.Country, user.Location.City); locationID == 0 {
		locationID = rep.StoreLocation(user.Location.Country, user.Location.City)
	}

	// TODO Add userId to User struct
	userID := rep.StoreUser(user.Profile, user.Name, locationID, user.Motto)

	// LOOP IT -->

	// TODO Add langId to Languages struct
	var ulang = (*user.Languages)[0]
	var languageID uint64
	if languageID = rep.FindLanguageByName(ulang); languageID == 0 {
		languageID = rep.StoreLanguage(ulang)
	}
	rep.StoreUserLanguage(userID, languageID)

	// LOOP IT -->

	// TODO Define UserPhoto struct
	rep.StoreUserPhoto(userID, (*user.Photos)[0])

	log.Println("###################### STATISTICS ######################")
	log.Println("## users    ", rep.CountUsers())
	log.Println("## languages", rep.CountLanguages())
	log.Println("## locations", rep.CountLocations())
	log.Println("## photos   ", rep.CountUserPhotos())
	log.Println("###################### STATISTICS ######################")
}
