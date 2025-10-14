package dto

import "time"

type ProfileRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	AvatarURL string `json:"avatar_url" binding:"required"`
}

type ProfileResonse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
