package domain

import "time"

// User representa un usuario del sistema.
type User struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid"`
	FirebaseUID string     `json:"firebase_uid" gorm:"uniqueIndex;size:128;not null"`
	Email       string     `json:"email" gorm:"index;size:255;not null"`
	Name        string     `json:"name" gorm:"size:255"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// Profile representa el perfil de un usuario.
type Profile struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    string    `json:"user_id" gorm:"index;type:uuid;not null"`
	Address   string    `json:"address" gorm:"size:255"`
	Phone     string    `json:"phone" gorm:"size:50"`
	AvatarURL string    `json:"avatar_url" gorm:"size:512"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order representa un pedido realizado por un usuario.
type Order struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    string     `json:"user_id" gorm:"index;type:uuid;not null"`
	Status    string     `json:"status" gorm:"size:50;index"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
