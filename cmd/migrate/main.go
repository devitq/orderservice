package main

import (
	"database/sql"
	"embed"
	"errors"
	"flag"
	"log"

	"orderservice/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	var cmd string

	flag.StringVar(&cmd, "cmd", "up", "migration command: up|down|force|version")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.BuildPostgresDSN())
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	drv, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("postgres driver: %v", err)
		return
	}

	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Printf("iofs source: %v", err)
		return
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", drv)
	if err != nil {
		log.Printf("migrate NewWithInstance: %v", err)
		return
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("m.Up failed: %v", err)
			return
		}
		log.Println("migrations applied (up)")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Printf("m.Steps(-1) failed: %v", err)
			return
		}
		log.Println("stepped down 1 migration")
	case "version":
		v, dirty, verr := m.Version()
		if verr != nil {
			log.Printf("version: %v", verr)
			return
		}
		log.Printf("version: %d dirty: %v\n", v, dirty)
	default:
		log.Printf("unknown cmd: %s", cmd)
	}
}
