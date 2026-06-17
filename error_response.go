package apperror

import "net/http"

var (
	InvalidBody = ErrorResponse{
		Message:  "invalid request body",
		AppCode:  101,
		HttpCode: http.StatusBadRequest,
		Success:  false,
	}
	InvalidUuid = ErrorResponse{
		Message:  "invalid uuid 7 format",
		AppCode:  101,
		HttpCode: http.StatusBadRequest,
		Success:  false,
	}
	InvalidTimestamp = ErrorResponse{
		Message:  "invalid timestamp",
		AppCode:  101,
		HttpCode: http.StatusBadRequest,
		Success:  false,
	}
	InvalidParam = ErrorResponse{
		Message:  "invalid request params",
		AppCode:  102,
		HttpCode: http.StatusBadRequest,
		Success:  false,
	}
	InvalidCredential = ErrorResponse{
		Message:  "invalid email or password",
		AppCode:  103,
		HttpCode: http.StatusUnauthorized,
		Success:  false,
	}
	InvalidToken = ErrorResponse{
		Message:  "invalid token",
		AppCode:  104,
		HttpCode: http.StatusUnauthorized,
		Success:  false,
	}
	TryLoginWithGoogle = ErrorResponse{
		Message:  "try login with Google",
		AppCode:  105,
		HttpCode: http.StatusUnauthorized,
		Success:  false,
	}

	UserAlreadyExist = ErrorResponse{
		Message:  "user already exists",
		AppCode:  202,
		HttpCode: http.StatusConflict,
		Success:  false,
	}

	NoSuchRecord = ErrorResponse{
		Message:  "no such record",
		AppCode:  203,
		HttpCode: http.StatusNotFound,
		Success:  false,
	}

	UnverifiedUser = ErrorResponse{
		Message:  "account is not verified, please verify your email",
		AppCode:  205,
		HttpCode: http.StatusForbidden,
		Success:  false,
	}

	RestrictedUser = ErrorResponse{
		Message:  "account is restricted, please contact support",
		AppCode:  204,
		HttpCode: http.StatusUnauthorized,
		Success:  false,
	}

	InternalError = ErrorResponse{
		Message:  "internal server error",
		AppCode:  301,
		HttpCode: http.StatusInternalServerError,
		Success:  false,
	}
)

type ErrorResponse struct {
	Message  string `json:"message"`
	AppCode  int    `json:"app_code"`
	HttpCode int    `json:"-"`
	Success  bool   `json:"success"`
}

// Ref returns a pointer to a fresh copy, preventing callers from mutating the package-level global.
func (e ErrorResponse) Ref() *ErrorResponse {
	return &e
}

func GenerateValidationError(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:  message,
		AppCode:  101,
		HttpCode: http.StatusBadRequest,
		Success:  false,
	}
}
