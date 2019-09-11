package main

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

var _repository *Repository

func init() {
	log.Println("init [db.go]")
}

// Repository represination of user repository
type Repository struct {
	cp *pgx.ConnPool
}

// NewRepository connects with db and allow to interact with it
func NewRepository() (*Repository, error) {

	if _repository != nil {
		return _repository, nil
	}

	cfg := GetAppConfig().Db

	connCfg := pgx.ConnConfig{
		Host:     cfg.URL,
		Port:     cfg.Port,
		User:     cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
	}

	connPoolCfg := pgx.ConnPoolConfig{
		ConnConfig:     connCfg,
		MaxConnections: cfg.MaxConnections,
		AcquireTimeout: time.Millisecond * time.Duration(cfg.AcquireTimeout),
	}

	cp, err := pgx.NewConnPool(connPoolCfg)
	if err != nil {
		return nil, fmt.Errorf("Can't connect with DB, %s", err.Error())
	}

	_repository = &Repository{cp: cp}

	return _repository, nil
}

// Close closes connection with db
func (rep Repository) Close() {
	if rep.cp != nil {
		rep.cp.Close()
	}
}

// Connection retrieve available connection from pool
func (rep Repository) Connection() *pgx.Conn {
	conn, err := rep.cp.Acquire()
	if err != nil {
		log.Printf("Can't acquire db connection from pool: %s\n", err.Error())
	}
	return conn
}

// Release returns connection to the pool
func (rep Repository) Release(conn *pgx.Conn) {
	rep.cp.Release(conn)
}

func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (rep Repository) truncate(sb sq.DeleteBuilder) {
	if connection := rep.Connection(); connection != nil {
		defer rep.Release(connection)
		sql, _, _ := sb.ToSql()
		connection.QueryRow(sql)
	}
}

func (rep Repository) fetch(sb sq.SelectBuilder) uint64 {
	if connection := rep.Connection(); connection != nil {
		defer rep.Release(connection)
		sql, args, _ := sb.ToSql()

		var count uint64
		if row := connection.QueryRow(sql, args...); row != nil {

			row.Scan(&count)
			return count
		}
	}
	return 0
}

func (rep Repository) insert(sb sq.InsertBuilder) uint64 {
	if connection := rep.Connection(); connection != nil {
		defer rep.Release(connection)
		sql, args, _ := sb.ToSql()

		if row := connection.QueryRow(sql, args...); row != nil {
			var count uint64
			row.Scan(&count)
			return count
		}
	}
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
func (rep Repository) StoreUser(login string, name string, locationID uint64, motto string, instagram string) uint64 {

	query := psql().Insert(
		"user_profile",
	).Columns(
		"user_login",
		"user_name",
		"user_location_id",
		"motto",
		"instagram_url",
		"created_date",
	).Values(
		login,
		name,
		locationID,
		motto,
		instagram,
		time.Now(),
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

// PrintStats print number of stored users, langs, photos etc
func (rep Repository) PrintStats() {
	log.Println("###################### STATISTICS ######################")
	log.Println("## users    ", rep.CountUsers())
	log.Println("## languages", rep.CountLanguages())
	log.Println("## locations", rep.CountLocations())
	log.Println("## photos   ", rep.CountUserPhotos())
	log.Println("###################### STATISTICS ######################")
}
