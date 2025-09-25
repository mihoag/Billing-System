package common

// StandardResponse represents a standard API response format
type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewResponse creates a new standardized response
func NewResponse(code int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// SuccessResponse creates a standard success response
func SuccessResponse(data interface{}) StandardResponse {
	return NewResponse(200, "success", data)
}

// ErrorResponse creates a standard error response
func ErrorResponse(code int, message string) StandardResponse {
	return NewResponse(code, message, nil)
}
