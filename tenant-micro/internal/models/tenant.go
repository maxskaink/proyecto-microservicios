package models

import "time"

type Tenant struct {
	ID         uint   `gorm:"primaryKey"`
	TenantID   string `gorm:"unique;index"`
	TenantName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
