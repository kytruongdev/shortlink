package auth

// Client-facing error codes returned by the auth controller.
// Input-form codes (INVALID_EMAIL/INVALID_PASSWORD) live in the handler.
const (
	codeEmailTaken         = "EMAIL_TAKEN"
	codeInvalidCredentials = "INVALID_CREDENTIALS"
	codeInvalidRefresh     = "INVALID_REFRESH_TOKEN"
)
