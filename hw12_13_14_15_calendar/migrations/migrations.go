package migrations

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq" // comment for justifying
	"github.com/pressly/goose/v3"
)

//go:embed files/*.sql
//go:embed inserting/*.sql
var embedMigrations embed.FS

func Up(dir string) error {
	dsn := "user=igor dbname=calendardb password=igor"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, dir); err != nil {
		return err
	}

	return db.Close()
}

func Down(dir string) error {
	dsn := "user=igor dbname=calendardb password=igor"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Down(db, dir); err != nil {
		return err
	}
	return db.Close()
}
