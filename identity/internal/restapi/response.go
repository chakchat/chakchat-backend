package restapi

// Specified by contract.md in /api folder

const (
	ErrTypeInternal         = "internal"
	ErrTypeInvalidJson      = "invalid_json"
	ErrTypeValidationFailed = "validation_failed"
	ErrTypeUserNotFound     = "user_not_found"

	ErrTypeNotFound = "not_found"

	ErrTypeIdempotencyKeyMissing = "idempotency_key_missing"

	ErrTypeSendCodeFreqExceeded = "send_code_freq_exceeded"
	ErrTypeSignInKeyNotFound    = "signin_key_not_found"
	ErrTypeWrongCode            = "wrong_code"

	ErrTypeRefreshTokenExpired     = "refresh_token_expired"
	ErrTypeRefreshTokenInvalidated = "refresh_token_invalidated"
	ErrTypeInvalidJWT              = "invalid_token"
	ErrTypeInvalidTokenType        = "invalid_token_type"

	ErrTypeUnautorized        = "unauthorized"
	ErrTypeAccessTokenExpired = "access_token_expired"

	ErrTypeUserAlreadyExists     = "user_already_exists"
	ErrTypeSignUpKeyNotFound     = "signup_key_not_found"
	ErrTypeUsernameAlreadyExists = "username_already_exists"
	ErrTypePhoneNotVerified      = "phone_not_verified"
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
