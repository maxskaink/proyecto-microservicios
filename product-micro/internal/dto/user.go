package dto

type UserRequest struct {
	ID          string `json:"id"`
	FirebaseUID string `json:"firabseUID"`
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
	Rol         string `json:"rol"`
}

type UserResponse struct {
	ID          string `json:"id"`
	FirebaseUID string `json:"firabseUID"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Rol         string `json:"rol"`
}
