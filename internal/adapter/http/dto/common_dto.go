package dto

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Success bool       `json:"success" example:"false"`
	Error   ErrorInfo  `json:"error"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string            `json:"code" example:"VALIDATION_ERROR"`
	Message string            `json:"message" example:"Validation failed"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents a success response wrapper
// @Description Success response structure
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// MetaInfo contains pagination information
type MetaInfo struct {
	Page       int   `json:"page,omitempty" example:"1"`
	PerPage    int   `json:"per_page,omitempty" example:"20"`
	Total      int64 `json:"total,omitempty" example:"100"`
	TotalPages int   `json:"total_pages,omitempty" example:"5"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
