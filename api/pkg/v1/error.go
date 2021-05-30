package v1

import (
	context "context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ServiceError describes a web-service error.
type ServiceError struct {
	Code    int
	Message string
}

// Error returns a string representation of the error.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("status %d: %s", e.Code, e.Message)
}

// StatusCode returns HTTP status code of the error.
func (e *ServiceError) StatusCode() int {
	return e.Code
}

// Encode encodes the error using the given HTTP response writer.
func (e *ServiceError) Encode(w http.ResponseWriter) {
	message := e.Message
	if e.Code == http.StatusInternalServerError {
		message = "internal error"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Decode decodes the error from the given HTTP response.
func (e *ServiceError) Decode(r *http.Response) {
	e.Code = r.StatusCode
	var res struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&res); err == nil && res.Error != "" {
		e.Message = res.Error
	} else {
		e.Message = http.StatusText(r.StatusCode)
	}
}

// ErrBadRequest creates a BadRequest service error.
func ErrBadRequest(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrInternal creates an Internal service error.
func ErrInternal(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, v...),
	}
}

// ****************** Errors *********************

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	defaultErr := &ServiceError{Code: http.StatusInternalServerError, Message: "internal error"}
	if err, ok := err.(*ServiceError); ok {
		defaultErr.Code = err.Code
		defaultErr.Message = err.Message
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultErr.Code)

	_ = json.NewEncoder(w).Encode(defaultErr)
}
