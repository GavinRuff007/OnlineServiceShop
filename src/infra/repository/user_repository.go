package repository

import (
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	"RestGoTest/src/database"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"RestGoTest/src/pkg/metrics"
	"context"
	"reflect"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const userFilterExp string = "username = ?"
const countFilterExp string = "count(*) > 0"

type MysqlUserRepository struct {
	*BaseRepository[model.User]
}

func NewUserRepository(cfg *config.Config) *MysqlUserRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{}
	return &MysqlUserRepository{BaseRepository: NewBaseRepository[model.User](cfg, preloads)}
}

func (r *MysqlUserRepository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	modelName := reflect.TypeOf(u).String()

	roleId, err := r.GetDefaultRole(ctx)
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "GetDefaultRole", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.DefaultRoleNotFound, err.Error(), nil)
		return u, err
	}
	metrics.DbCall.WithLabelValues(modelName, "GetDefaultRole", "Success").Inc()

	tx := r.database.WithContext(ctx).Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCall.WithLabelValues(modelName, "Create", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return u, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Create", "Success").Inc()

	err = tx.Create(&model.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		metrics.DbCall.WithLabelValues("model.UserRole", "Create", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return u, err
	}
	metrics.DbCall.WithLabelValues("model.UserRole", "Create", "Success").Inc()

	tx.Commit()
	return u, nil
}

func (r *MysqlUserRepository) FetchUserInfo(ctx context.Context, username string, password string) (model.User, error) {
	var user model.User
	modelName := reflect.TypeOf(user).String()

	err := r.database.WithContext(ctx).
		Model(&model.User{}).
		Where(userFilterExp, username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "FetchUserInfo", "Error").Inc()
		return user, err
	}
	metrics.DbCall.WithLabelValues(modelName, "FetchUserInfo", "Success").Inc()

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "PasswordCheck", "Error").Inc()
		return user, err
	}
	metrics.DbCall.WithLabelValues(modelName, "PasswordCheck", "Success").Inc()

	return user, nil
}

func (r *MysqlUserRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	modelName := reflect.TypeOf(model.User{}).String()

	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		metrics.DbCall.WithLabelValues(modelName, "ExistsEmail", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	metrics.DbCall.WithLabelValues(modelName, "ExistsEmail", "Success").Inc()
	return exists, nil
}

func (r *MysqlUserRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	modelName := reflect.TypeOf(model.User{}).String()

	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where(userFilterExp, username).
		Find(&exists).
		Error; err != nil {
		metrics.DbCall.WithLabelValues(modelName, "ExistsUsername", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	metrics.DbCall.WithLabelValues(modelName, "ExistsUsername", "Success").Inc()
	return exists, nil
}

func (r *MysqlUserRepository) ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error) {
	var exists bool
	modelName := reflect.TypeOf(model.User{}).String()

	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		metrics.DbCall.WithLabelValues(modelName, "ExistsMobileNumber", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	metrics.DbCall.WithLabelValues(modelName, "ExistsMobileNumber", "Success").Inc()
	return exists, nil
}

func (r *MysqlUserRepository) GetDefaultRole(ctx context.Context) (roleId int, err error) {
	modelName := reflect.TypeOf(model.Role{}).String()

	if err = r.database.WithContext(ctx).Model(&model.Role{}).
		Select("id").
		Where("name = ?", constant.DefaultRoleName).
		First(&roleId).Error; err != nil {
		metrics.DbCall.WithLabelValues(modelName, "GetDefaultRole", "Error").Inc()
		return 0, err
	}
	metrics.DbCall.WithLabelValues(modelName, "GetDefaultRole", "Success").Inc()
	return roleId, nil
}
