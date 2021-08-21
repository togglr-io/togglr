package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/eriktate/toggle/env"
)

const maxTries = 5
const waitTime = 3

type config struct {
	Host     string
	User     string
	Password string
	Database string
	Port     uint
}

func fromEnv(prefix string) config {
	return config{
		Host:     env.GetString(fmt.Sprintf("%s_DB_HOST", prefix), "localhost"),
		User:     env.GetString(fmt.Sprintf("%s_DB_USER", prefix), "root"),
		Password: env.GetString(fmt.Sprintf("%s_DB_PASSWORD", prefix), "toggle"),
		Database: env.GetString(fmt.Sprintf("%s_DB_NAME", prefix), "toggle"),
		Port:     env.GetUint(fmt.Sprintf("%s_DB_PORT", prefix), 5432),
	}
}

func connect(tries int) (*sql.DB, error) {
	cfg := fromEnv("TOGGLE")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	db, err := sql.Open("postgres", dsn)
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
