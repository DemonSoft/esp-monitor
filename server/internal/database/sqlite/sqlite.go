package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"remoteesp/internal/database"
	"remoteesp/internal/domain"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed schemes/*.sql
var migrationsFS embed.FS

const (
	defaultDBName = "remoteEsp.db"
	Driver        = "sqlite"
)

type Repository struct {
	db     *sql.DB
	driver string
}

func Open() database.Database {
	db, err := sql.Open(Driver, defaultDBName)
	if err != nil {
		domain.Log.Log(err)
		return &Repository{driver: Driver}
	}

	err = migrations(db)
	if err != nil {
		domain.Log.Log(err)
		return &Repository{db: db, driver: Driver}
	}

	if err := db.Ping(); err != nil {
		domain.Log.Log(err)
		return &Repository{db: db, driver: Driver}
	}

	return &Repository{db: db, driver: Driver}
}

func migrations(db *sql.DB) error {
	// source driver from embed.FS
	sourceDriver, err := iofs.New(migrationsFS, "schemes")
	if err != nil {
		return fmt.Errorf("iofs source: %w", err)
	}

	// sqlite driver for migrate
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("sqlite driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func (r *Repository) Close() {
	if r.db == nil {
		return // It can be, if database wasn't created.
	}
	err := r.db.Close()
	if err != nil {
		domain.Log.Log(err)
	}

	fmt.Printf("Database '%s' closed.\n", r.driver)
}

func (r *Repository) Driver() string {
	return r.driver
}
