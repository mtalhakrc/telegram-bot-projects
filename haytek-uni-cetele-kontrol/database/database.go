package database

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/logx"
	"time"
)

var db *bun.DB

func New(config config.DbConfig) *bun.DB {
	var dsn string
	if config.Socket != "" {
		dsn = fmt.Sprintf("unix://%s:%s@%s/%s?sslmode=disable", config.Username, config.Password, config.Name, config.Socket)
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timeout=30s", config.Username, config.Password, config.Host, config.Port, config.Name)
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn) /*, pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})*/))
	sqldb.SetMaxIdleConns(config.MaxIdleConn)
	sqldb.SetMaxOpenConns(config.MaxPoolSize)
	sqldb.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	db = bun.NewDB(sqldb, pgdialect.New())

	// tablo isimlerini plural yapmasın.
	schema.SetTableNameInflector(func(s string) string {
		return s
	})

	if err := db.Ping(); err != nil {
		logx.SendLog("db nil")
		panic(err)
	}

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(config.Debug)))

	return db
}

func Get() *bun.DB {
	return db
}
