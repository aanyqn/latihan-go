package helper

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	CodeValidation         = "VALIDATION_ERROR"
	CodeBadRequest         = "BAD_REQUEST"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeUnsupportedMedia   = "UNSUPPORTED_MEDIA_TYPE"
	CodeNotAcceptable      = "NOT_ACCEPTABLE"
	CodeTooManyRequests    = "TOO_MANY_REQUESTS"
	CodeInternal           = "INTERNAL_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Errors  map[string]string
	cause   error
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Errors    map[string]string `json:"errors,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.cause }

func BadRequest(message string) *AppError {
	return &AppError{Status: fiber.StatusBadRequest, Code: CodeBadRequest, Message: message}
}
func Unauthorized(message string) *AppError {
	return &AppError{Status: fiber.StatusUnauthorized, Code: CodeUnauthorized, Message: message}
}
func Forbidden(message string) *AppError {
	return &AppError{Status: fiber.StatusForbidden, Code: CodeForbidden, Message: message}
}
func NotFound(message string) *AppError {
	return &AppError{Status: fiber.StatusNotFound, Code: CodeNotFound, Message: message}
}
func Conflict(message string) *AppError {
	return &AppError{Status: fiber.StatusConflict, Code: CodeConflict, Message: message}
}
func UnprocessableEntity(message string) *AppError {
	return &AppError{Status: fiber.StatusUnprocessableEntity, Code: "UNPROCESSABLE_ENTITY", Message: message}
}
func Validation(fields map[string]string) *AppError {
	return &AppError{
		Status: fiber.StatusUnprocessableEntity, Code: CodeValidation,
		Message: "failed validation", Errors: fields,
	}
}
func UnsupportedMedia(message string) *AppError {
	return &AppError{Status: fiber.StatusUnsupportedMediaType, Code: CodeUnsupportedMedia, Message: message}
}
func TooManyRequests(message string) *AppError {
	return &AppError{Status: fiber.StatusTooManyRequests, Code: CodeTooManyRequests, Message: message}
}
func ServiceUnavailable(message string) *AppError {
	return &AppError{Status: fiber.StatusServiceUnavailable, Code: CodeServiceUnavailable, Message: message}
}
func NotAcceptable(message string) *AppError {
	return &AppError{
		Status: fiber.StatusNotAcceptable, Code: CodeNotAcceptable, Message: message,
	}
}
func Internal(cause error, message string) *AppError {
	return &AppError{
		Status: fiber.StatusInternalServerError, Code: CodeInternal,
		Message: message, cause: cause,
	}
}

var (
	ErrNotFound       = errors.New("Not found")
	ErrDuplicate      = errors.New("Already used")
	ErrInvalidInput   = errors.New("Input isn't valid")
	ErrNoFieldsChange = errors.New("No changes input")
	ErrQuotaFull      = errors.New("Kuota penuh")
)

func TranslateError(c *fiber.Ctx, err error, generalMessage string) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return NotFound(err.Error())
	case errors.Is(err, ErrDuplicate):
		return Conflict(err.Error())
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrNoFieldsChange):
		return BadRequest(err.Error())
	case errors.Is(err, ErrQuotaFull):
		return UnprocessableEntity(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return Fail(c, fiber.StatusGatewayTimeout, "Request timeout")
	default:
		return Internal(err, generalMessage)
	}
}
