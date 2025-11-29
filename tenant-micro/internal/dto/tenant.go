package dto

type CreateTenantRequest struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	TenantName  string `json:"tenant_name" binding:"required"`
	Description string `json:"description" bininding:"required"`
	Location    string `json:"location" binding:"required"`
}

type TenantResponse struct {
	TenantID    string `json:"tenant_id"`
	TenantName  string `json:"tenant_name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
