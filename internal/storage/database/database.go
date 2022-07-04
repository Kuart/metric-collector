package database

import (
	"context"
	"database/sql"
	"fmt"
	config "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/metric"
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

func (db DB) Update(ctx context.Context, m metric.Metric) error {
	if m.MType == metric.GaugeTypeName {
		_, err := db.instance.ExecContext(ctx,
			`
				INSERT INTO gauge (name, value)
				VALUES ($1, $2) 
				ON CONFLICT(name) 
				DO UPDATE SET value = $2;
			`,
			m.ID,
			m.Value,
		)

		if err != nil {
			return err
		}
	} else if m.MType == metric.CounterTypeName {
		_, err := db.instance.ExecContext(ctx,
			`
				INSERT INTO counter (name, value)
				VALUES ($1, $2) 
				ON CONFLICT(name) 
				DO UPDATE SET value = counter.value + $2;
			`,
			m.ID,
			m.Delta,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db DB) GetMetric(ctx context.Context, m metric.Metric) (metric.Metric, bool) {
	var err error
	mtr := metric.Metric{
		MType: m.MType,
	}
	row := db.instance.QueryRowContext(ctx, fmt.Sprintf("SELECT name, value FROM %s", m.MType))

	if m.MType == metric.CounterTypeName {
		err = row.Scan(&mtr.ID, &mtr.Delta)
	} else {
		err = row.Scan(&mtr.ID, &mtr.Value)
	}

	if err != nil {
		log.Printf("get metric err: %s", err.Error())
		return mtr, false
	}

	return mtr, true
}

func (db DB) GetAllMetrics(ctx context.Context, MType string) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	rows, err := db.instance.QueryContext(ctx, fmt.Sprintf("SELECT name, value FROM %s", MType))
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		name := ""

		if MType == metric.GaugeTypeName {
			var gauge float64

			if err = rows.Scan(&name, &gauge); err != nil {
				log.Printf("get all scan err: %s", err.Error())
				return nil, err
			}

			metrics[name] = gauge
		} else {
			var counter int64

			if err = rows.Scan(&name, &counter); err != nil {
				log.Printf("get all scan err: %s", err.Error())
				return nil, err
			}

			metrics[name] = counter
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (db DB) Close() {
	db.instance.Close()
}
