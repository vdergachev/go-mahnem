package main

import (
	"fmt"
	"strings"
)

// Location - user location
type Location struct {
	City    string
	Country string
}

// User struct populated bu Mahneclientlient::profile method
type User struct {
	Profile   string
	Name      string
	Location  *Location
	Languages *[]string
	Motto     string
	Photos    *[]string
}

// NewLocation function parses string and returns user location
func newLocation(val string) *Location {
	if len(val) == 0 {
		return &Location{}
	}

	vals := strings.Split(strip(val), ",")

	if len(vals) != 2 {
		return &Location{}
	}

	return &Location{
		City:    strip(vals[0]),
		Country: strip(vals[1]),
	}
}

func (location *Location) toString() string {
	return fmt.Sprintf("country: %s, city: %s", location.Country, location.City)
}

func (user *User) toString() string {
	return fmt.Sprintf("login: %s, name: %s, location: %s",
		user.Profile,
		user.Name,
		user.Location.toString(),
	)
}
