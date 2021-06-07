package response

type ErrorResponseParameters struct {
	Internal   error
	IsInternal bool
	Message    string
	StatusCode int
}

func (e ErrorResponseParameters) Error() string {
	var message string
	if e.Message != "" {
		message += e.Message
		if e.Internal != nil {
			message += ": "
		}
	}
	if e.Internal != nil {
		message += e.Internal.Error()
	}
	return message
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func (e *ErrorResponseParameters) Response() interface{} {
	var message = "Internal server error"
	if e.Message != "" {
		message = e.Message
	}
	return ErrorResponse{
		Message: message,
	}
}
