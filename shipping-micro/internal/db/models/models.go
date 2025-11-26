package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserDB representa el modelo de usuario en la base de datos
type UserDB struct {
	ID        string `gorm:"type:string;primaryKey"`
	UUID      string `gorm:"type:string;uniqueIndex;not null"`
	Rol       string `gorm:"not null;default:'user'"`
	Email     string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserDB) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// ProductDB representa el modelo de producto en la base de datos
type ProductDB struct {
	ID          string `gorm:"type:string;primaryKey"`
	Name        string `gorm:"not null"`
	ProducerID  string `gorm:"type:string;not null;index"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null;default:0"`
	PhotoUrl    string  `gorm:"type:text" json:"photo_url"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CartItemDB representa un artículo en el carrito
type CartItemDB struct {
	ID        string     `gorm:"type:uuid;primaryKey"`
	UserID    string     `gorm:"type:string;not null;index"`
	ProductID string     `gorm:"type:string;not null;index"`
	Product   *ProductDB `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity  int        `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *CartItemDB) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// OrderDB representa una orden
type OrderDB struct {
	ID         string  `gorm:"type:uuid;primaryKey"`
	UserID     string  `gorm:"type:string;not null;index"`
	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Items      []OrderItemDB `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (o *OrderDB) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// OrderItemDB representa un producto dentro de una orden
type OrderItemDB struct {
	ID        string     `gorm:"type:uuid;primaryKey"`
	OrderID   string     `gorm:"type:uuid;not null;index"`
	ProductID string     `gorm:"type:string;not null;index"`
	Product   *ProductDB `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity  int        `gorm:"not null"`
	Price     float64    `gorm:"not null"`
	CreatedAt time.Time
}

func (oi *OrderItemDB) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = uuid.New().String()
	}
	return nil
}

// ShippingDB representa un envío
type ShippingDB struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	OrderID         string `gorm:"type:uuid;not null;uniqueIndex"`
	TrackingNumber  string `gorm:"uniqueIndex;not null"`
	ShippingAddress string `gorm:"not null"`
	Status          string `gorm:"not null;default:'pending'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *ShippingDB) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
