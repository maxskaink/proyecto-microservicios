package dto

type CreateTenantRequest struct {
	TenantID   string `json:"tenant_id" binding:"required"`
	TenantName string `json:"tenant_name" binding:"required"`
}

type TenantResponse struct {
	ID         uint   `json:"id"`
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	CreatedAt  string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
