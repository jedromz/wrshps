package api

type ErrorMessage struct {
	Message string `json:"message"`
}

type ApiError struct {
	ErrorMessage
	ErrorType string
}

type UnauthorizedError struct {
	ApiError
}
type ForbiddenError struct {
	ApiError
}
type RateLimitExceededError struct {
	ApiError
}
type BadRequestError struct {
	ApiError
}
type NotFoundError struct {
	ApiError
}

func (e ApiError) Error() string {
	return e.ErrorType + ": " + e.Message
}
