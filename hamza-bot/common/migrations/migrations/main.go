package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var sqlMigrations embed.FS

var (
	// Migrations provides migration logic for bun
	Migrations = migrate.NewMigrations()
)

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
