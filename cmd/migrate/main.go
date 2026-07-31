package main

import (
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"database/sql"

	"github.com/justKody/taskboard-go-api/config"
)

func main() {

	db, err := sql.Open("pgx", config.Envs.DatabaseUrl)

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"pgx",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	argsCommands(m)

}

func argsCommands(m *migrate.Migrate) {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|down|force]")
	}

	switch os.Args[1] {
	case "up":
		if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
