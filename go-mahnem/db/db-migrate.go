package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {

	const (
		defaultMode   = "update"
		defaultHost   = "localhost"
		defaultPort   = 5432
		defaultDb     = "goma"
		defaultUser   = "postgres"
		defaultPasswd = "1"
	)

	var mode = flag.String("mode", defaultMode, "migrate mode")
	var hostname = flag.String("hostname", defaultHost, "db hostname")
	var port = flag.Int("port", defaultPort, "db port")
	var dbname = flag.String("db", defaultDb, "db name")
	var username = flag.String("user", defaultUser, "db username")
	var passwd = flag.String("passwd", defaultPasswd, "db username's password")

	flag.Parse()

	if len(*passwd) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var dbURL string
	if strings.ToLower(*mode) == "init" {

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
			*hostname,
			*port,
			*username,
			*passwd,
		)

		err := doInit(psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

	} else {

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			*username,
			*passwd,
			*hostname,
			*port,
			*dbname,
		)

		err := doMigrate(dbURL)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func doInit(psqlInfo string) error {
	const (
		initSQL = `CREATE DATABASE goma`
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	r, err := db.Exec(initSQL)
	if err != nil || r == nil {
		return err
	}

	log.Println("DB Migrate - init db complete")
	return nil
}

func doMigrate(url string) error {

	migration, err := migrate.New(
		"file://migrations",
		url,
	)

	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil {
		return err
	}

	log.Println("DB Migrate - update db complete")
	return nil
}
