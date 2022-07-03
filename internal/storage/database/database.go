package database

import (
	"database/sql"
	config "github.com/Kuart/metric-collector/config/server"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	instance    *sql.DB
	isConnected bool
}

func New(cfg config.Config) (DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseDSN)

	if err != nil {
		log.Printf("database didn't connect: %s", err)
		return DB{}, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("database ping error: %s", err)
		return DB{}, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Printf("database driver error: %s", err)
		return DB{}, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DatabaseDSN, driver)

	if err != nil {
		log.Printf("database migration error: %s", err)
		return DB{}, err
	}

	m.Up()

	log.Println("Database connection successful ended")

	return DB{
		instance:    db,
		isConnected: true,
	}, nil
}

func (db DB) Ping() bool {
	if !db.isConnected {
		return false
	}

	err := db.instance.Ping()

	return err == nil
}

func (db DB) Close() {
	db.instance.Close()
}
