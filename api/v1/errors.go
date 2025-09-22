package v1

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewErrorResponse(message string, err error) *ErrorResponse {
	res := &ErrorResponse{Message: message}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}
