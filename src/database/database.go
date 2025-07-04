package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// TODO : if want run in local use below function
func InitDatabase() *sql.DB {

	var err error
	DB, err := sql.Open("mysql", "root:P@ssw0rd!2023@tcp(127.0.0.1:3306)/onlineShop?parseTime=true")
	if err != nil {
		panic(err)
	}
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	return DB
}

// func InitDatabase() *sql.DB {
// 	dbUser := os.Getenv("DB_USER")
// 	dbPass := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

// 	DB, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := DB.Ping(); err != nil {
// 		panic(err)
// 	}
// 	return DB
// }
