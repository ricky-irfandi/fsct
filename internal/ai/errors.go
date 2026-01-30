package ai

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNoAPIKey         = errors.New("no API key configured")
	ErrInvalidResponse  = errors.New("invalid API response")
	ErrAPITimeout       = errors.New("API request timeout")
	ErrRateLimited      = errors.New("rate limited by API")
	ErrProviderNotFound = errors.New("provider not found")
	ErrRequestFailed    = errors.New("API request failed")
	ErrContextCanceled  = errors.New("request canceled")
)

// APIError represents an error from the AI API
type APIError struct {
	StatusCode int
	Message    string
	Provider   string
	Retryable  bool
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s API error (status %d): %s", e.Provider, e.StatusCode, e.Message)
}

// IsRetryable returns true if the error is retryable
func (e *APIError) IsRetryable() bool {
	return e.Retryable
}

// NewAPIError creates a new API error
func NewAPIError(provider string, statusCode int, message string) *APIError {
	// Determine if error is retryable based on status code
	retryable := statusCode == 429 || statusCode >= 500
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Provider:   provider,
		Retryable:  retryable,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}
	// Also retry on timeout
	return errors.Is(err, ErrAPITimeout)
}
