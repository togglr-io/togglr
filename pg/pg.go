package pg

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/togglr-io/togglr/env"

	// importing postgres driver implementation for database/sql and goqu
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// A Config captures information required to make a postgres connection
type Config struct {
	Host     string
	Port     uint
	User     string
	Password string
	Database string
}

// ConfigFromEnv builds a Config from environment variables using a given prefix and falling back to
// local defaults
func ConfigFromEnv(prefix string) Config {
	return Config{
		Host:     env.GetString(fmt.Sprintf("%s_DB_HOST", prefix), "localhost"),
		Port:     env.GetUint(fmt.Sprintf("%s_DB_PORT", prefix), 5432),
		User:     env.GetString(fmt.Sprintf("%s_DB_USER", prefix), "toggle"),
		Password: env.GetString(fmt.Sprintf("%s_DB_PASSWORD", prefix), "toggle"),
		Database: env.GetString(fmt.Sprintf("%s_DB_NAME", prefix), "toggle"),
	}
}

// DSN generates a connection string from a Config
func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", c.Host, c.Port, c.Database, c.User, c.Password)
}

// A Client is a container for things required to talk to postgres. It also implements any postgres service
// interfaces.
type Client struct {
	rawDB *sql.DB
	db    *goqu.Database
	log   *zap.Logger
}

// NewClient returns a new Client given a Config.
func NewClient(cfg Config) (Client, error) {
	dialect := goqu.Dialect("postgres")
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return Client{}, err
	}

	return Client{
		rawDB: db,
		db:    dialect.DB(db),
	}, nil
}
