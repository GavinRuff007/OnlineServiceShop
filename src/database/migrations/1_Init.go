package migrations

import (
	"RestGoTest/src/config"
	db "RestGoTest/src/database"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"log"
)

func InitMigrations() error {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)
	database := db.GetDb()
	if database == nil {
		log.Fatal("Database instance is nil")
	}

	err := database.AutoMigrate(
		&model.User{},
		&model.GiftcardProvider{},
		&model.Giftcard{},
		&model.Order{},
		&model.Payment{},
		&model.GiftcardCode{},
		&model.DiscountCode{},
		&model.OrderDiscount{},
		&model.AuditLog{},
	)

	if err != nil {
		logger.Fatalf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migration completed successfully.")
	return nil
}
