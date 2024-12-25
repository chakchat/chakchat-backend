package restapi

// Specified by contract.md in /api folder

const (
	ErrorTypeInternal         = "internal"
	ErrorTypeInvalidJson      = "invalid_json"
	ErrorTypeValidationFailed = "validation_failed"
	ErrorTypeUserNotFound     = "user_not_found"

	ErrorTypeIdempotencyKeyMissing = "idempotency_key_missing"
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
