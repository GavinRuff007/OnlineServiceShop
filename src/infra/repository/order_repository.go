package repository

import (
	"RestGoTest/src/config"
	"RestGoTest/src/database"
	"RestGoTest/src/dto"
	"RestGoTest/src/model"
	"RestGoTest/src/pkg/logging"
	"RestGoTest/src/pkg/metrics"
	"context"
	"reflect"
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

func (r *MysqlOrderRepository) CreateOrder(ctx context.Context, orderDTO *dto.OrderCreateDTO) (dto.OrderResponseDTO, error) {
	var orderModel model.Order
	var result dto.OrderResponseDTO
	modelName := reflect.TypeOf(orderModel).String()

	if err := validate.Struct(orderDTO); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return result, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	totalPrice := float64(orderDTO.Quantity) * orderDTO.UnitPrice

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

	if orderModel.Status == "" {
		orderModel.Status = "pending"
	}

	tx := r.database.WithContext(ctx).Begin()
	if tx.Error != nil {
		metrics.DbCall.WithLabelValues(modelName, "Create", "Error").Inc()
		return result, tx.Error
	}

	if err := tx.Create(&orderModel).Error; err != nil {
		tx.Rollback()
		metrics.DbCall.WithLabelValues(modelName, "Create", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return result, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Create", "Success").Inc()

	if err := tx.Commit().Error; err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Commit", "Error").Inc()
		return result, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Commit", "Success").Inc()

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

func (r *MysqlOrderRepository) GetOrderByID(ctx context.Context, req dto.OrderGetByIDDTO) (dto.OrderResponseDTO, error) {
	var order model.Order
	var result dto.OrderResponseDTO
	modelName := reflect.TypeOf(order).String()

	if err := validate.Struct(req); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return result, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	err := r.database.WithContext(ctx).
		Where(orderIDFilterExp, req.OrderID).
		Preload("User").
		Preload("Giftcard").
		First(&order).Error
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Select", "Error").Inc()
		return result, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Select", "Success").Inc()

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

func (r *MysqlOrderRepository) GetOrdersByUserID(ctx context.Context, req dto.OrdersByUserDTO) ([]dto.OrderResponseDTO, error) {
	modelName := reflect.TypeOf(model.Order{}).String()

	if err := validate.Struct(req); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return nil, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	var orders []model.Order
	var result []dto.OrderResponseDTO

	err := r.database.WithContext(ctx).
		Where("user_id = ?", req.UserID).
		Preload("Giftcard").
		Find(&orders).Error
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Select", "Error").Inc()
		return nil, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Select", "Success").Inc()

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

func (r *MysqlOrderRepository) ExistsOrder(ctx context.Context, req dto.OrderExistsDTO) (bool, error) {
	modelName := reflect.TypeOf(model.Order{}).String()

	if err := validate.Struct(req); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return false, err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	var exists bool
	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Select("count(*) > 0").
		Where(orderIDFilterExp, req.OrderID).
		Find(&exists).Error
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "ExistsOrder", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return false, err
	}
	metrics.DbCall.WithLabelValues(modelName, "ExistsOrder", "Success").Inc()
	return exists, nil
}

func (r *MysqlOrderRepository) UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateDTO) error {
	modelName := reflect.TypeOf(model.Order{}).String()

	if err := validate.Struct(req); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	err := r.database.WithContext(ctx).
		Model(&model.Order{}).
		Where(orderIDFilterExp, req.OrderID).
		Update("status", req.Status).Error
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "UpdateStatus", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return err
	}
	metrics.DbCall.WithLabelValues(modelName, "UpdateStatus", "Success").Inc()
	return nil
}

func (r *MysqlOrderRepository) DeleteOrder(ctx context.Context, req dto.OrderDeleteDTO) error {
	modelName := reflect.TypeOf(model.Order{}).String()

	if err := validate.Struct(req); err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Validate", "Error").Inc()
		return err
	}
	metrics.DbCall.WithLabelValues(modelName, "Validate", "Success").Inc()

	err := r.database.WithContext(ctx).
		Where(orderIDFilterExp, req.OrderID).
		Delete(&model.Order{}).Error
	if err != nil {
		metrics.DbCall.WithLabelValues(modelName, "Delete", "Error").Inc()
		r.logger.Error(logging.Mysql, logging.Delete, err.Error(), nil)
		return err
	}
	metrics.DbCall.WithLabelValues(modelName, "Delete", "Success").Inc()
	return nil
}
