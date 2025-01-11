package restapi

// Specified by contract.md in /api folder

const (
	ErrTypeInternal              = "internal"
	ErrTypeInvalidJson           = "invalid_json"
	ErrTypeValidationFailed      = "validation_failed"
	ErrTypeNotFound              = "not_found"
	ErrTypeIdempotencyKeyMissing = "idempotency_key_missing"
	ErrTypeUnautorized           = "unauthorized"

	ErrTypeInvalidHeader   = "invalid_header"
	ErrTypeContentTooLarge = "content_too_large"
	ErrTypeInvalidForm     = "invalid_form"

	ErrTypeFileNotFound   = "file_not_found"
	ErrTypeUploadNotFound = "upload_not_found"
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
