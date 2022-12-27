package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func GetConnection() *bun.DB {
	sqldb, err := sql.Open("mysql", "root:fourkeys@/fourkeys")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())
	return db
}
