package handlers

var (
	errTypeUserNotFound            = "user_not_found"
	errTypeSendCodeFreqExceeded    = "send_code_freq_exceeded"
	errTypeSignInKeyNotFound       = "signin_key_not_found"
	errTypeWrongCode               = "wrong_code"
	errTypeRefreshTokenExpired     = "refresh_token_expired"
	errTypeRefreshTokenInvalidated = "refresh_token_invalidated"
	errTypeInvalidJWT              = "invalid_token"
	errTypeInvalidTokenType        = "invalid_token_type"
	errTypeAccessTokenExpired      = "access_token_expired"
	errTypeUserAlreadyExists       = "user_already_exists"
	errTypeSignUpKeyNotFound       = "signup_key_not_found"
	errTypeUsernameAlreadyExists   = "username_already_exists"
	errTypePhoneNotVerified        = "phone_not_verified"
)
