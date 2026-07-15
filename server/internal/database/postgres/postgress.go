package postgres

import (
	"fmt"
	"remoteesp/internal/database"
	"remoteesp/internal/domain"

	"github.com/jmoiron/sqlx"
)

const (
	Driver = "pgx"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	SSL  string
}

type Repository struct {
	db     *sqlx.DB
	driver string
}

func Open() database.Database {
	cfg := Config{}
	domain.Log.Log("Config", cfg)

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSL))

	if err != nil {
		domain.Log.Log("Connection error: %v", err)
		return &Repository{db: db, driver: Driver}
	}

	if err := db.Ping(); err != nil {
		domain.Log.Log("DB ping error: %v", err)
		return &Repository{db: db, driver: Driver}
	}

	return &Repository{db: db, driver: Driver}
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
