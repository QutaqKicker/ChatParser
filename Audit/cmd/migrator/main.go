package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path of queryBuilder storage")
	flag.StringVar(&migrationPath, "migrations-path", "", "path of migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "path of migrations table")
	flag.Parse()

	if storagePath == "" || migrationPath == "" {
		panic("invalid arguments")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		storagePath)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

}
