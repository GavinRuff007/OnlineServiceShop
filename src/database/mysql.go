package database

import (
	"RestGoTest/src/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	var err error
	cnn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Mysql.DB_HOST, cfg.Mysql.DB_PORT, cfg.Mysql.DB_USER, cfg.Mysql.DB_PASSWORD,
		cfg.Mysql.DB_NAME, cfg.Mysql.SSLMode)

	dbClient, err = gorm.Open(mysql.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Mysql.ConnMaxLifetime * time.Minute)

	log.Println("Db connection established")
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	con, _ := dbClient.DB()
	con.Close()
}
