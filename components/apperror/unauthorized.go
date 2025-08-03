package apperror

import "net/http"

type UnauthorizedError struct {
	AppError
	Type string `json:"type"`
}

func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusUnauthorized,
		},
		Type: "INVALID_TOKEN",
	}
}
