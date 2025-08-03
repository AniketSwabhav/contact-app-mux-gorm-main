package apperror

type AppError struct {
	HTTPStatus int    `json:"status"`
	Message    string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}
