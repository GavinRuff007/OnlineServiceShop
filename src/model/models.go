package model

import (
	"time"

	"gorm.io/datatypes"
)

// User represents the users table
type User struct {
	ID           uint    `gorm:"primaryKey"`
	FullName     string  `gorm:"size:100;not null"`
	Email        string  `gorm:"size:120;not null;unique"`
	PasswordHash string  `gorm:"type:text;not null"`
	PhoneNumber  *string `gorm:"size:20"`
	Role         string  `gorm:"type:enum('user','admin');default:'user'"`
	IsActive     bool    `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GiftcardProvider represents the giftcard_providers table
type GiftcardProvider struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"size:50;not null"`
	LogoURL   *string `gorm:"type:text"`
	IsActive  bool    `gorm:"default:true"`
	SortOrder int     `gorm:"default:0"`
}

// Giftcard represents the giftcards table
type Giftcard struct {
	ID          uint `gorm:"primaryKey"`
	ProviderID  *uint
	Provider    *GiftcardProvider `gorm:"foreignKey:ProviderID"`
	Region      string            `gorm:"size:20;not null"`
	AmountUSD   float64           `gorm:"type:decimal(10,2);not null"`
	PriceLocal  float64           `gorm:"type:decimal(10,2);not null"`
	Currency    string            `gorm:"size:10;default:'IRT'"`
	Description *string           `gorm:"type:text"`
	IsActive    bool              `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Order represents the orders table
type Order struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	User       User
	GiftcardID uint
	Giftcard   Giftcard
	Quantity   int
	UnitPrice  float64 `gorm:"type:decimal(10,2)"`
	TotalPrice float64 `gorm:"type:decimal(12,2)"`
	Status     string  `gorm:"type:enum('pending','paid','failed','delivered');default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Payment represents the payments table
type Payment struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"uniqueIndex"`
	Order     Order   `gorm:"foreignKey:OrderID;references:ID"`
	Provider  string  `gorm:"size:50;not null"`
	Status    string  `gorm:"type:enum('unpaid','paid','failed');default:'unpaid'"`
	Amount    float64 `gorm:"type:decimal(12,2);not null"`
	Authority *string `gorm:"size:100"`
	PaidAt    *time.Time
	CreatedAt time.Time
}

// GiftcardCode represents the giftcard_codes table
type GiftcardCode struct {
	ID              uint `gorm:"primaryKey"`
	GiftcardID      *uint
	Giftcard        *Giftcard
	CodeEncrypted   string `gorm:"type:text;not null"`
	IsUsed          bool   `gorm:"default:false"`
	UsedAt          *time.Time
	AssignedOrderID *uint
	AssignedOrder   *Order `gorm:"foreignKey:AssignedOrderID"`
	CreatedAt       time.Time
}

// DiscountCode represents the discount_codes table
type DiscountCode struct {
	ID            uint    `gorm:"primaryKey"`
	Code          string  `gorm:"size:20;unique;not null"`
	Description   *string `gorm:"type:text"`
	DiscountType  string  `gorm:"type:enum('fixed','percent');not null"`
	DiscountValue float64 `gorm:"type:decimal(10,2);not null"`
	MaxUsage      *int
	ExpiresAt     *time.Time
	IsActive      bool `gorm:"default:true"`
	CreatedAt     time.Time
}

// OrderDiscount represents the order_discounts table
type OrderDiscount struct {
	ID             uint `gorm:"primaryKey"`
	OrderID        *uint
	Order          *Order
	DiscountCodeID *uint
	DiscountCode   *DiscountCode
	AppliedValue   *float64 `gorm:"type:decimal(10,2)"`
}

// AuditLog represents the audit_logs table
type AuditLog struct {
	ID        uint `gorm:"primaryKey"`
	UserID    *uint
	User      *User
	Action    string `gorm:"size:100"`
	Details   datatypes.JSON
	CreatedAt time.Time
}
