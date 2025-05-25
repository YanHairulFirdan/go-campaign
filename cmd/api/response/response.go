package response

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewResponse(status, message string, data any) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(status, message, err string) *Response {
	return NewResponse(status, message, err)
}

func NewValidationErrorResponse(status, message string, errors []map[string]string) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}
