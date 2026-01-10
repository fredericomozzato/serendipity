package main

import (
	"context"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "up or down")
	flag.Parse()

	dsn := "postgres://serendipity:password@serendipity_db:5432/serendipity?sslmode=disable"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m, err := migrate.New(
		"file://db/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("unknown direction: %s", direction)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("All migrations already applied")
			return
		}
		log.Fatal(err)
	}

	log.Println("Migration completed:", direction)
}
