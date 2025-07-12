package migrations

import (
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	db "RestGoTest/src/database"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"log"

	"golang.org/x/crypto/bcrypt"
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
		&model.UserRole{},
		&model.Role{},
		&model.Payment{},
		&model.GiftcardCode{},
		&model.DiscountCode{},
		&model.OrderDiscount{},
		&model.AuditLog{},
	)
	createDefaultUserInformation(database)
	if err != nil {
		logger.Fatalf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migration completed successfully.")
	return nil
}

func createDefaultUserInformation(database *gorm.DB) {

	adminRole := model.Role{Name: constant.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)

	defaultRole := model.Role{Name: constant.DefaultRoleName}
	createRoleIfNotExists(database, &defaultRole)

	u := model.User{Username: constant.DefaultUserName, FullName: "Test", LastName: "Test",
		MobileNumber: "09111112222", Email: "admin@admin.com"}
	pass := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &u, adminRole.Id)

}

func createRoleIfNotExists(database *gorm.DB, r *model.Role) {
	exists := 0
	database.
		Model(&model.Role{}).
		Select("1").
		Where("name = ?", r.Name).
		First(&exists)
	if exists == 0 {
		database.Create(r)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, u *model.User, roleId int) {
	exists := 0
	database.
		Model(&model.User{}).
		Select("1").
		Where("username = ?", u.Username).
		First(&exists)

	if exists == 0 {
		database.Create(u)
		ur := model.UserRole{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}
