package main

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var _repository *Repository

func init() {
	log.Println("init [db.go]")
}

// Repository represination of user repository
type Repository struct {
	db *sql.DB
}

// NewRepository connects with db and allow to interact with it
func NewRepository() (*Repository, error) {

	if _repository != nil {
		return _repository, nil
	}

	cfg := GetAppConfig().Db

	dataSource := fmt.Sprintf(""+
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.URL,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	d, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("Can't connect with DB, %s", err.Error())
	}

	d.SetMaxOpenConns(100)
	d.SetMaxIdleConns(5)

	_repository = &Repository{db: d}

	return _repository, nil
}

// Close closes connection with db
func (rep Repository) Close() {
	if rep.db != nil {
		rep.db.Close()
	}
}

func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (rep Repository) truncate(sb sq.DeleteBuilder) {
	rows, err := sb.RunWith(rep.db).Query()
	if err == nil {
		rows.Close()
	} else {
		sql, args, _ := sb.ToSql()
		log.Printf("Query '%s' with args '%+v' failed, %s\n", sql, args, err.Error())
	}
}

func (rep Repository) fetch(sb sq.SelectBuilder) uint64 {

	rows, err := sb.RunWith(rep.db).Query()

	if err == nil {

		defer rows.Close()

		if rows.Next() {
			var count uint64
			rows.Scan(&count)
			return count
		}
		return 0
	}

	sql, args, _ := sb.ToSql()
	log.Printf("Query '%s' with args '%+v' failed, %s\n", sql, args, err.Error())
	return 0
}

func (rep Repository) insert(sb sq.InsertBuilder) uint64 {

	rows, err := sb.RunWith(rep.db).Query()

	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var id uint64
			rows.Scan(&id)
			return id
		}
		return 0
	}

	sql, args, _ := sb.ToSql()
	log.Printf("Query '%s' with args '%+v' failed, %s\n", sql, args, err.Error())
	return 0
}

// TODO 		WE DEFENETLY NEED A CASCADE SETUP IN TABLES
// DeleteAll...
func (rep Repository) deleteAllUsers() {
	rep.truncate(sq.Delete("user_profile"))
}

func (rep Repository) deleteAllLanguages() {
	rep.truncate(sq.Delete("user_language"))
}

func (rep Repository) deleteAllLocations() {
	rep.truncate(sq.Delete("user_location"))
}

func (rep Repository) deleteAllUserPhotos() {
	rep.truncate(sq.Delete("user_photo"))
}

func (rep Repository) deleteAllUserLanguages() {
	rep.truncate(sq.Delete("user_to_language"))
}

// Count...

// CountUsers counts existing user profiles
func (rep Repository) CountUsers() uint64 {
	return rep.fetch(sq.Select("count(*)").From("user_profile"))
}

// CountLanguages counts user languages
func (rep Repository) CountLanguages() uint64 {
	return rep.fetch(sq.Select("count(*)").From("user_language"))
}

// CountLocations counts user locations
func (rep Repository) CountLocations() uint64 {
	return rep.fetch(sq.Select("count(*)").From("user_location"))
}

// CountUserPhotos counts user photos
func (rep Repository) CountUserPhotos() uint64 {
	return rep.fetch(sq.Select("count(*)").From("user_photo"))
}

// Store...

// StoreUser save user profile to database
func (rep Repository) StoreUser(login string, name string, locationID uint64, motto string) uint64 {

	query := psql().Insert(
		"user_profile",
	).Columns(
		"user_login",
		"user_name",
		"user_location_id",
		"motto",
	).Values(
		login,
		name,
		locationID,
		motto,
	).Suffix("RETURNING user_profile_id")

	return rep.insert(query)
}

// StoreLanguage save language to database
func (rep Repository) StoreLanguage(language string) uint64 {

	query := psql().Insert(
		"user_language",
	).Columns(
		"language_name",
	).Values(
		language,
	).Suffix("RETURNING user_language_id")

	return rep.insert(query)
}

// StoreUserLanguage add language to the user
func (rep Repository) StoreUserLanguage(userID uint64, languageID uint64) uint64 {

	query := psql().Insert(
		"user_to_language",
	).Columns(
		"user_profile_id",
		"user_language_id",
	).Values(
		userID,
		languageID,
	).Suffix("ON CONFLICT (user_profile_id, user_language_id) DO NOTHING")

	return rep.insert(query)
}

// StoreLocation save user location to db
func (rep Repository) StoreLocation(country string, city string) uint64 {
	query := psql().Insert(
		"user_location",
	).Columns(
		"country",
		"city",
	).Values(
		country,
		city,
	).Suffix("RETURNING user_location_id")

	return rep.insert(query)
}

// StoreUserPhoto save user photo link to db
func (rep Repository) StoreUserPhoto(userID uint64, url string) uint64 {
	query := psql().Insert(
		"user_photo",
	).Columns(
		"user_profile_id",
		"url",
	).Values(
		userID,
		url,
	).Suffix("ON CONFLICT (user_profile_id, url) DO NOTHING RETURNING user_photo_id")

	return rep.insert(query)
}

// Find...

// FindUserByLogin finds user profile id by login
func (rep Repository) FindUserByLogin(login string) uint64 {
	query := psql().
		Select("user_profile_id").
		From("user_profile").
		Where(sq.Eq{"user_login": login})

	return rep.fetch(query)
}

// FindLanguageByName finds language id by language name
func (rep Repository) FindLanguageByName(language string) uint64 {
	query := psql().
		Select("user_language_id").
		From("user_language").
		Where(sq.Eq{"language_name": language})

	return rep.fetch(query)
}

// FindLocation finds location id with given country and city
func (rep Repository) FindLocation(country string, city string) uint64 {
	query := psql().
		Select("user_location_id").
		From("user_location").
		Where(sq.Eq{"country": country, "city": city})

	return rep.fetch(query)
}

// FindUserPhoto finds user photo by user id and link
func (rep Repository) FindUserPhoto(userID uint64, url string) uint64 {
	query := psql().
		Select("user_photo_id").
		From("user_photo").
		Where(sq.Eq{"user_profile_id": userID, "url": url})

	return rep.fetch(query)
}
