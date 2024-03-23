package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"matheuslanduci.com/api-fiber/config"
)

func New(config *config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", config.Database.Url)

	if err != nil {
		log.Fatalf("An error occurred during the connection to the database. Error: %v", err)
	}
	return db
}
