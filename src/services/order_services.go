package services

import (
	"RestGoTest/src/config"
	"RestGoTest/src/domain/repository"
	"RestGoTest/src/dto"
	"RestGoTest/src/helper/service_errors"
	"RestGoTest/src/pkg/logging"
	"context"
)

type OrderService struct {
	logger     logging.Logger
	cfg        *config.Config
	repository repository.OrderRepository
}

func NewOrderService(cfg *config.Config, repo repository.OrderRepository) *OrderService {
	logger := logging.NewLogger(cfg)
	return &OrderService{
		cfg:        cfg,
		repository: repo,
		logger:     logger,
	}
}

// ایجاد سفارش جدید
func (s *OrderService) CreateOrder(ctx context.Context, req dto.OrderCreateDTO) (dto.OrderResponseDTO, error) {
	order, err := s.repository.CreateOrder(ctx, &req)
	if err != nil {
		s.logger.Error(logging.General, logging.Api, err.Error(), nil)
		return dto.OrderResponseDTO{}, &service_errors.ServiceError{EndUserMessage: "خطا در ثبت سفارش"}
	}
	return order, nil
}

// دریافت سفارش با ID
func (s *OrderService) GetOrderByID(ctx context.Context, req dto.OrderGetByIDDTO) (dto.OrderResponseDTO, error) {
	order, err := s.repository.GetOrderByID(ctx, req)
	if err != nil {
		s.logger.Error(logging.General, logging.Select, err.Error(), nil)
		return dto.OrderResponseDTO{}, &service_errors.ServiceError{EndUserMessage: "سفارش یافت نشد"}
	}
	return order, nil
}

// دریافت همه سفارشات کاربر
func (s *OrderService) GetOrdersByUserID(ctx context.Context, req dto.OrdersByUserDTO) ([]dto.OrderResponseDTO, error) {
	orders, err := s.repository.GetOrdersByUserID(ctx, req)
	if err != nil {
		s.logger.Error(logging.General, logging.Select, err.Error(), nil)
		return nil, &service_errors.ServiceError{EndUserMessage: "خطا در دریافت سفارشات"}
	}
	return orders, nil
}

// بروزرسانی وضعیت سفارش
func (s *OrderService) UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateDTO) error {
	err := s.repository.UpdateOrderStatus(ctx, req)
	if err != nil {
		s.logger.Error(logging.General, logging.Update, err.Error(), nil)
		return &service_errors.ServiceError{EndUserMessage: "بروزرسانی سفارش ناموفق بود"}
	}
	return nil
}

// حذف سفارش
func (s *OrderService) DeleteOrder(ctx context.Context, req dto.OrderDeleteDTO) error {
	err := s.repository.DeleteOrder(ctx, req)
	if err != nil {
		s.logger.Error(logging.General, logging.Delete, err.Error(), nil)
		return &service_errors.ServiceError{EndUserMessage: "حذف سفارش ناموفق بود"}
	}
	return nil
}

// بررسی وجود سفارش
func (s *OrderService) ExistsOrder(ctx context.Context, req dto.OrderExistsDTO) (bool, error) {
	exists, err := s.repository.ExistsOrder(ctx, req)
	if err != nil {
		s.logger.Error(logging.General, logging.Select, err.Error(), nil)
		return false, &service_errors.ServiceError{EndUserMessage: "خطا در بررسی سفارش"}
	}
	return exists, nil
}
