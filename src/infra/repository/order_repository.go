package repository

import (
	"RestGoTest/src/config"
	"RestGoTest/src/database"
	"RestGoTest/src/dto"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

const orderIDFilterExp string = "id = ?"

var validate = validator.New()

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
func (r *MysqlOrderRepository) CreateOrder(ctx context.Context, orderDTO *dto.OrderCreateDTO) (dto.OrderResponseDTO, error) {
	var orderModel model.Order
	var result dto.OrderResponseDTO

	// --- اعتبارسنجی ورودی ---
	if err := validate.Struct(orderDTO); err != nil {
		return result, err
	}

	// --- محاسبه TotalPrice ---
	totalPrice := float64(orderDTO.Quantity) * orderDTO.UnitPrice

	// --- پر کردن مدل ---
	orderModel = model.Order{
		UserID:     orderDTO.UserID,
		GiftcardID: orderDTO.GiftcardID,
		Quantity:   orderDTO.Quantity,
		UnitPrice:  orderDTO.UnitPrice,
		TotalPrice: totalPrice,
		Status:     orderDTO.Status,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	// پیش‌فرض pending
	if orderModel.Status == "" {
		orderModel.Status = "pending"
	}

	tx := r.database.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result, tx.Error
	}

	if err := tx.Create(&orderModel).Error; err != nil {
		tx.Rollback()
		r.logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return result, err
	}
	if err := tx.Commit().Error; err != nil {
		return result, err
	}

	// مپ به DTO برای خروجی
	result = dto.OrderResponseDTO{
		ID:         orderModel.ID,
		UserID:     orderModel.UserID,
		GiftcardID: orderModel.GiftcardID,
		Quantity:   orderModel.Quantity,
		UnitPrice:  orderModel.UnitPrice,
		TotalPrice: orderModel.TotalPrice,
		Status:     orderModel.Status,
		CreatedAt:  orderModel.CreatedAt,
		UpdatedAt:  orderModel.UpdatedAt,
	}
	return result, nil
}

// دریافت سفارش با شناسه
func (r *MysqlOrderRepository) GetOrderByID(ctx context.Context, req dto.OrderGetByIDDTO) (dto.OrderResponseDTO, error) {
	var order model.Order
	var result dto.OrderResponseDTO

	if err := validate.Struct(req); err != nil {
		return result, err
	}

	err := r.database.WithContext(ctx).
		Where(orderIDFilterExp, req.OrderID).
		Preload("User").
		Preload("Giftcard").
		First(&order).Error
	if err != nil {
		return result, err
	}

	result = dto.OrderResponseDTO{
		ID:         order.ID,
		UserID:     order.UserID,
		GiftcardID: order.GiftcardID,
		Quantity:   order.Quantity,
		UnitPrice:  order.UnitPrice,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
	return result, nil
}

// دریافت همه سفارشات یک کاربر
func (r *MysqlOrderRepository) GetOrdersByUserID(ctx context.Context, req dto.OrdersByUserDTO) ([]dto.OrderResponseDTO, error) {
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	var orders []model.Order
	var result []dto.OrderResponseDTO

	err := r.database.WithContext(ctx).
		Where("user_id = ?", req.UserID).
		Preload("Giftcard").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	for _, o := range orders {
		result = append(result, dto.OrderResponseDTO{
			ID:         o.ID,
			UserID:     o.UserID,
			GiftcardID: o.GiftcardID,
			Quantity:   o.Quantity,
			UnitPrice:  o.UnitPrice,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		})
	}

	return result, nil
}

// چک کردن وجود سفارش
func (r *MysqlOrderRepository) ExistsOrder(ctx context.Context, req dto.OrderExistsDTO) (bool, error) {
	if err := validate.Struct(req); err != nil {
		return false, err
	}

	var exists bool
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Select("count(*) > 0").
		Where(orderIDFilterExp, req.OrderID).
		Find(&exists).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

// بروزرسانی وضعیت سفارش
func (r *MysqlOrderRepository) UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateDTO) error {
	if err := validate.Struct(req); err != nil {
		return err
	}

	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Where(orderIDFilterExp, req.OrderID).
		Update("status", req.Status).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}

// حذف سفارش
func (r *MysqlOrderRepository) DeleteOrder(ctx context.Context, req dto.OrderDeleteDTO) error {
	if err := validate.Struct(req); err != nil {
		return err
	}

	err := r.database.WithContext(ctx).
		Where(orderIDFilterExp, req.OrderID).
		Delete(&model.Order{}).Error
	if err != nil {
		r.logger.Error(logging.Mysql, logging.Delete, err.Error(), nil)
		return err
	}
	return nil
}
