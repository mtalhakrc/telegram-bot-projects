package database

import (
	"database/sql"
	"fmt"
	"github.com/haytek-uni-bot-yeniden/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
)

var db *bun.DB

func New(config config.DbConfig) {
	sqlite, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("%s?cache=shared", config.Path))
	if err != nil {
		panic(err)
	}

	db = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if db.Ping() != nil {
		panic(err)
	}
	log.Println("db initialized")
}

func Get() *bun.DB {
	return db
}
