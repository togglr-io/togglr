package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/eriktate/toggle/pg"
)

const maxTries = 5
const waitTime = 3

func connect(tries int) (*sql.DB, error) {
	cfg := pg.ConfigFromEnv("TOGGLE")

	// default configs use the app user
	if cfg.User == "toggle" {
		cfg.User = "root"
		cfg.Password = "toggle"
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Println("Retrying db connection...")
		if tries < maxTries {
			time.Sleep(waitTime * time.Second)
			return connect(tries + 1)
		}
	}

	return db, err
}

func runMigration(db *sql.DB, query string, tries int) error {
	if _, err := db.Exec(query); err != nil {
		if tries < maxTries {
			log.Println("Retrying migration...")
			time.Sleep(waitTime * time.Second)
			return runMigration(db, query, tries+1)
		}

		return err
	}

	return nil
}

func migrate(direction string) error {
	query, err := ioutil.ReadFile(fmt.Sprintf("./migrations/%s.sql", direction))
	if err != nil {
		return fmt.Errorf("failed to load 'up' migration")
	}

	db, err := connect(0)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigration(db, string(query), 0); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("migrate must be called with a direction ('up' or 'down')")
	}

	direction := os.Args[1]
	if err := migrate(direction); err != nil {
		log.Fatalf("migration failed: %s", err)
	}
}
