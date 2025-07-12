package model

import (
	"time"

	"gorm.io/datatypes"
)

// User represents the users table
type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:20;not null;unique"`
	FullName     string `gorm:"type:string;size:100;not null"`
	LastName     string `gorm:"type:string;size:25;null"`
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`
	Email        string `gorm:"type:string;size:64;null;unique;default:null"`
	Password     string `gorm:"type:string;size:64;not null"`
	Enabled      bool   `gorm:"default:true"`
	UserRoles    *[]UserRole
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:10;not null,unique"`
	UserRoles *[]UserRole
}

type UserRole struct {
	BaseModel
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserId int
	RoleId int
}

// GiftcardProvider represents the giftcard_providers table
type GiftcardProvider struct {
	ID        int     `gorm:"primaryKey"`
	Name      string  `gorm:"size:50;not null"`
	LogoURL   *string `gorm:"type:text"`
	IsActive  bool    `gorm:"default:true"`
	SortOrder int     `gorm:"default:0"`
}

// Giftcard represents the giftcards table
type Giftcard struct {
	ID          int `gorm:"primaryKey"`
	ProviderID  *int
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
	ID         int `gorm:"primaryKey"`
	UserID     int
	User       User
	GiftcardID int
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
	ID        int     `gorm:"primaryKey"`
	OrderID   int     `gorm:"uniqueIndex"`
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
	ID              int `gorm:"primaryKey"`
	GiftcardID      *int
	Giftcard        *Giftcard
	CodeEncrypted   string `gorm:"type:text;not null"`
	IsUsed          bool   `gorm:"default:false"`
	UsedAt          *time.Time
	AssignedOrderID *int
	AssignedOrder   *Order `gorm:"foreignKey:AssignedOrderID"`
	CreatedAt       time.Time
}

// DiscountCode represents the discount_codes table
type DiscountCode struct {
	ID            int     `gorm:"primaryKey"`
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
	ID             int `gorm:"primaryKey"`
	OrderID        *int
	Order          *Order
	DiscountCodeID *int
	DiscountCode   *DiscountCode
	AppliedValue   *float64 `gorm:"type:decimal(10,2)"`
}

// AuditLog represents the audit_logs table
type AuditLog struct {
	ID        int `gorm:"primaryKey"`
	UserID    *int
	User      *User
	Action    string `gorm:"size:100"`
	Details   datatypes.JSON
	CreatedAt time.Time
}
