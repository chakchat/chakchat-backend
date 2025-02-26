package restapi

const (
	ErrTypeNotFound         = "invalid_input"
	ErrTypeUnautorized      = "unauthorized"
	ErrTypeInternal         = "internal"
	ErrTypeInvalidJson      = "invalid_json"
	ErrTypeValidationFailed = "validation_failed"
)

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	ErrorType    string        `json:"error_type,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	ErrorDetails []ErrorDetail `json:"error_details,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data,omitempty"`
}
