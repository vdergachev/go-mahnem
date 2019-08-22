package main

import (
	"log"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {

	migration, err := migrate.New(
		"file://migrations",
		"postgres://postgres:admin@localhost:5432/testo?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := migration.Up(); err != nil {
		log.Fatal(err)
	}
}
