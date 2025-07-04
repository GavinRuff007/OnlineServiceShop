package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	var err error
	DB, err := sql.Open("mysql", "root:P@ssw0rd!2023@tcp(127.0.0.1:3306)/onlineShop?parseTime=true")
	if err != nil {
		panic(err)
	}
	// بهتره health check بزنی:
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	return DB
}
