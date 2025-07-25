package repository

import (
	"RestGoTest/src/dto"
	"RestGoTest/src/model"
	"context"
)

type UserRepository interface {
	ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error)
	ExistsUsername(ctx context.Context, username string) (bool, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	FetchUserInfo(ctx context.Context, username string, password string) (model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderDTO *dto.OrderCreateDTO) (dto.OrderResponseDTO, error)
	GetOrderByID(ctx context.Context, req dto.OrderGetByIDDTO) (dto.OrderResponseDTO, error)
	GetOrdersByUserID(ctx context.Context, req dto.OrdersByUserDTO) ([]dto.OrderResponseDTO, error)
	ExistsOrder(ctx context.Context, req dto.OrderExistsDTO) (bool, error)
	UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateDTO) error
	DeleteOrder(ctx context.Context, req dto.OrderDeleteDTO) error
}
