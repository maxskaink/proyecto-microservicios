package db_models

import "time"

type UserDB struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // Clave primaria con UUID
	FirebaseUID string     `gorm:"uniqueIndex;size:255;not null"`                  // UID único de Firebase
	Email       string     `gorm:"uniqueIndex;size:255;not null"`                  // Email único
	Name        string     `gorm:"size:255;not null"`                              // Nombre del usuario
	Rol         string     `gorm:"size:50"`
	DeletedAt   *time.Time `gorm:"index"` // Fecha de eliminación (soft delete)
}
