package restapi

// Specified by contract.md in /api folder

const (
	ErrTypeInternal         = "internal"
	ErrTypeInvalidJson      = "invalid_json"
	ErrTypeValidationFailed = "validation_failed"
	ErrTypeUserNotFound     = "user_not_found"

	ErrTypeIdempotencyKeyMissing = "idempotency_key_missing"

	ErrTypeSendCodeFreqExceeded = "send_code_freq_exceeded"
	ErrTypeSignInKeyNotFound    = "signin_key_not_found"
	ErrTypeWrongCode            = "wrong_code"
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
