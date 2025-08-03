package apperror

import "net/http"

type DatabaseError struct {
	AppError
	Type string `json:"type"`
}

func NewDatabaseError(msg string) *DatabaseError {
	return &DatabaseError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusInternalServerError,
		},
		Type: "DATABASE_ERROR",
	}
}

func NewNotFoundError(msg string) *DatabaseError {
	return &DatabaseError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusNotFound,
		},
		Type: "NOT_FOUND_ERROR",
	}
}

func NewDuplicateEntryError(msg string) *DatabaseError {
	return &DatabaseError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusNotFound,
		},
		Type: "DUPLICATE_ENTRY",
	}
}
