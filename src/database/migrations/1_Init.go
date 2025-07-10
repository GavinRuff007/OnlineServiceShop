package migrations

import (
	"RestGoTest/src/config"
	db "RestGoTest/src/database"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"RestGoTest/src/util"
	"log"

	"gorm.io/gorm"
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

	createDefaultAdmin(database)
	log.Println("✅ Database migration completed successfully.")
	return nil
}

func createDefaultAdmin(database *gorm.DB) {
	const defaultAdminEmail = "admin@onlineshop.local"

	var count int64
	database.Model(&model.User{}).Where("email = ?", defaultAdminEmail).Count(&count)
	if count == 0 {
		admin := model.User{
			FullName:     "Default Admin",
			Email:        defaultAdminEmail,
			PasswordHash: util.HashPassword("admin123"),
			Role:         "admin",
			IsActive:     true,
		}
		if err := database.Create(&admin).Error; err != nil {
			log.Printf("⚠️ Failed to create default admin user: %v", err)
		} else {
			log.Printf("✅ Default admin user created with email: %s", defaultAdminEmail)
		}
	} else {
		log.Println("ℹ️ Default admin user already exists.")
	}
}
