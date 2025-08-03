package apperror

import "net/http"

type ValidationError struct {
	AppError
	Type string `json:"type"`
}

func NewInvalidJSONError(msg string) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusBadRequest,
		},
		Type: "INVALID_JSON",
	}
}

func NewMissingFieldsError(msg string) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusBadRequest,
		},
		Type: "MISSING_FIELDS",
	}
}

func NewValidationError(errType, msg string) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusBadRequest,
		},
		Type: errType,
	}
}
