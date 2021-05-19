package response

type ErrorResponse struct {
	Internal    error
	IsInternal  bool
	Message     string
	StatusCode  int
}

func (e ErrorResponse) Error() string {
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

func (e *ErrorResponse) Response() map[string]interface{} {
	var message = "Internal server error"
	if e.Message != "" {
		message = e.Message
	}
	return map[string]interface{}{
		"error": message,
	}
}
