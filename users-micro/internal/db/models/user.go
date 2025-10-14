package db_models

import "time"

type UserDB struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`                                // Clave primaria con UUID
	FirebaseUID string     `gorm:"uniqueIndex;size:255;not null"`                                                 // UID único de Firebase
	Email       string     `gorm:"uniqueIndex;size:255;not null"`                                                 // Email único
	Name        string     `gorm:"size:255;not null"`                                                             // Nombre del usuario
	CreatedAt   time.Time  `gorm:"autoCreateTime"`                                                                // Fecha de creación
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`                                                                // Fecha de actualización
	DeletedAt   *time.Time `gorm:"index"`                                                                         // Fecha de eliminación (soft delete)
	Profile     ProfileDB  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Relación uno a uno con ProfileDB
}

type ProfileDB struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // Clave primaria con UUID
	UserID    string    `gorm:"uniqueIndex;not null"`                           // Relación uno a uno con UserDB
	Address   string    `gorm:"size:255"`                                       // Dirección
	Phone     string    `gorm:"size:20"`                                        // Teléfono
	AvatarURL string    `gorm:"size:255"`                                       // URL del avatar
	CreatedAt time.Time `gorm:"autoCreateTime"`                                 // Fecha de creación
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                                 // Fecha de actualización
}
