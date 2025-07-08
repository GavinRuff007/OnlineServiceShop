package database

import (
	"RestGoTest/src/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	cfg := config.GetConfig()
	dbUser := cfg.Mysql.DB_USER
	dbPass := cfg.Mysql.DB_PASSWORD
	dbHost := cfg.Mysql.DB_HOST
	dbPort := cfg.Mysql.DB_PORT
	dbName := cfg.Mysql.DB_NAME

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	return DB
}
