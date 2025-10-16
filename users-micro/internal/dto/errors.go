package dto

import "time"

// Error resposponse for exceptions
type ErrorDTO struct {
	TimeStamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
}

// Return a new error, also the status for better reading
func NewErrorDTO(status int, msg_error string) (int, *ErrorDTO) {
	return status, &ErrorDTO{
		TimeStamp: time.Now(),
		Status:    status,
		Error:     msg_error,
	}
}
