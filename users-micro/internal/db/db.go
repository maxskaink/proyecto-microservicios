package db

import (
	"context"

	"gorm.io/gorm"
)

// DBProvider expone un método para obtener *gorm.DB.
type DBProvider interface {
	DB(ctx context.Context) (*gorm.DB, error)
}
