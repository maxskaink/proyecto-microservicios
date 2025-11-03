package dto

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
