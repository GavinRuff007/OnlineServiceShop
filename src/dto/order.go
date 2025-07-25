package dto

import (
	"time"
)

// DTO برای ایجاد سفارش جدید
type OrderCreateDTO struct {
	UserID     int     `json:"user_id" validate:"required,gt=0"`
	GiftcardID int     `json:"giftcard_id" validate:"required,gt=0"`
	Quantity   int     `json:"quantity" validate:"required,gt=0"`
	UnitPrice  float64 `json:"unit_price" validate:"required,gt=0"`
	Status     string  `json:"status" validate:"omitempty,oneof=pending paid failed delivered"`
}

// DTO برای آپدیت وضعیت سفارش
type OrderStatusUpdateDTO struct {
	OrderID int    `json:"order_id" validate:"required,gt=0"`
	Status  string `json:"status" validate:"required,oneof=pending paid failed delivered"`
}

// DTO برای دریافت سفارش با ID
type OrderGetByIDDTO struct {
	OrderID int `json:"order_id" validate:"required,gt=0"`
}

// DTO برای دریافت سفارشات یک کاربر
type OrdersByUserDTO struct {
	UserID int `json:"user_id" validate:"required,gt=0"`
}

// DTO برای چک کردن وجود سفارش
type OrderExistsDTO struct {
	OrderID int `json:"order_id" validate:"required,gt=0"`
}

// DTO برای حذف سفارش
type OrderDeleteDTO struct {
	OrderID int `json:"order_id" validate:"required,gt=0"`
}

// DTO برای بازگشت اطلاعات سفارش
type OrderResponseDTO struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	GiftcardID int       `json:"giftcard_id"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
