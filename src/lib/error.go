package lib

type ApiError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewApiErrror(message string, status int) *ApiError {
	return &ApiError{
		Message: message,
		Status:  status,
	}
}

func (s *ApiError) Error() string {
	return s.Message
}

var (
	ErrBadRequest       = NewApiErrror("Bad request", 400)
	ErrUnauthorized     = NewApiErrror("Unauthorized", 401)
	ErrForbidden        = NewApiErrror("This resource is not found", 403)
	ErrNotFound         = NewApiErrror("This resource is not found", 404)
	ErrMethodNotAllowed = NewApiErrror("This resource is not found", 405)
	ErrInternal         = NewApiErrror("Internal server error", 500)
	ErrBadGateway       = NewApiErrror("Internal server error", 502)
	ErrGatewayTimeout   = NewApiErrror("Internal server error", 504)
)
