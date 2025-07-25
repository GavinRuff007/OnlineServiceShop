package repository

import (
	"RestGoTest/src/config"
	"RestGoTest/src/database"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"context"
)

const orderStatusFilterExp string = "status = ?"
const orderIDFilterExp string = "id = ?"

type MysqlOrderRepository struct {
	*BaseRepository[model.Order]
}

func NewOrderRepository(cfg *config.Config) *MysqlOrderRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{
		{Entity: "User"},
		{Entity: "Giftcard"},
	}
	return &MysqlOrderRepository{
		BaseRepository: NewBaseRepository[model.Order](cfg, preloads),
	}
}

// ثبت سفارش جدید
func (r *MysqlOrderRepository) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return order, err
	}
	tx.Commit()
	return order, nil
}

// دریافت سفارش با شناسه
func (r *MysqlOrderRepository) GetOrderByID(ctx context.Context, orderID int) (model.Order, error) {
	var order model.Order
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Where(orderIDFilterExp, orderID).
		Preload("User").
		Preload("Giftcard").
		First(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}

// دریافت همه سفارشات یک کاربر
func (r *MysqlOrderRepository) GetOrdersByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	var orders []model.Order
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Where("user_id = ?", userID).
		Preload("Giftcard").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// چک کردن وجود سفارش بر اساس شناسه
func (r *MysqlOrderRepository) ExistsOrder(ctx context.Context, orderID int) (bool, error) {
	var exists bool
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Select("count(*) > 0").
		Where(orderIDFilterExp, orderID).
		Find(&exists).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

// بروزرسانی وضعیت سفارش (مثلا از pending به paid)
func (r *MysqlOrderRepository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Where(orderIDFilterExp, orderID).
		Update("status", status).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}

// حذف سفارش
func (r *MysqlOrderRepository) DeleteOrder(ctx context.Context, orderID int) error {
	err := r.database.WithContext(ctx).
		Where(orderIDFilterExp, orderID).
		Delete(&model.Order{}).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Delete, err.Error(), nil)
		return err
	}
	return nil
}
