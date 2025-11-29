package models

import "time"

type Tenant struct {
	ID          uint   `gorm:"primaryKey"`
	TenantID    string `gorm:"unique;index"`
	TenantName  string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Location    string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
