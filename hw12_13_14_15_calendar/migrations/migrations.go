package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed files/*.sql
var embedMigrations embed.FS

func Up() error {
	dsn := "user=igor dbname=calendardb password=igor"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "files"); err != nil {
		return err
	}

	return db.Close()
}

func Down() error {
	dsn := "user=igor dbname=calendardb password=igor"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Down(db, "files"); err != nil {
		return err
	}
	return db.Close()
}
