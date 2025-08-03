package apperror

import "net/http"

type UnauthorizedError struct {
	AppError
	Type string `json:"type"`
}

func NewInValidTokenError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusUnauthorized,
		},
		Type: "INVALID_TOKEN",
	}
}

func NewInValidPasswordError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusUnauthorized,
		},
		Type: "INVALID_PASSWORD",
	}
}

func NewUnauthorizedUserError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusUnauthorized,
		},
		Type: "ACCESS_DENIED",
	}
}

func NewInActiveUserError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Message:    msg,
			HTTPStatus: http.StatusUnauthorized,
		},
		Type: "INACTIVE_USER",
	}
}
