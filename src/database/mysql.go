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
	cnn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.DB_USER, cfg.Mysql.DB_PASSWORD,
		cfg.Mysql.DB_HOST, cfg.Mysql.DB_PORT,
		cfg.Mysql.DB_NAME,
	)

	dbClient, err = gorm.Open(mysql.Open(cnn), &gorm.Config{}) // ← اصلاح‌شده
	if err != nil {
		println(err.Error())
		return err
	}

	sqlDb, err := dbClient.DB()
	if err != nil {
		return err
	}

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
