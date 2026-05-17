package enums

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeAuth       ErrorType = "AUTH_ERROR"
	ErrorTypeForbidden  ErrorType = "FORBIDDEN_ERROR"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict   ErrorType = "CONFLICT_ERROR"
	ErrorTypeRateLimit  ErrorType = "RATE_LIMIT_ERROR"
	ErrorTypeDatabase   ErrorType = "DATABASE_ERROR"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
)

func (e ErrorType) String() string { return string(e) }
