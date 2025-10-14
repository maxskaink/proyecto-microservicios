package dto

type UserRequest struct {
	FirebaseUID string `json:"firabseUID" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
