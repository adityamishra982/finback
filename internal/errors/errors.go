package errors

// AppError is the central error structure for the application.
type AppError struct {
	Code    int      `json:"-"` // HTTP Status Code, not exposed directly in JSON response body
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewBadRequest creates a 400 Bad Request error.
func NewBadRequest(msg string, details ...string) *AppError {
	return &AppError{Code: 400, Message: msg, Details: details}
}

// NewUnauthorized creates a 401 Unauthorized error.
func NewUnauthorized(msg string, details ...string) *AppError {
	return &AppError{Code: 401, Message: msg, Details: details}
}

// NewForbidden creates a 403 Forbidden error.
func NewForbidden(msg string, details ...string) *AppError {
	return &AppError{Code: 403, Message: msg, Details: details}
}

// NewNotFound creates a 404 Not Found error.
func NewNotFound(msg string, details ...string) *AppError {
	return &AppError{Code: 404, Message: msg, Details: details}
}

// NewInternalServer creates a 500 Internal Server error.
func NewInternalServer(msg string, details ...string) *AppError {
	return &AppError{Code: 500, Message: msg, Details: details}
}
