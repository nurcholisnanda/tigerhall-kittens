package errorhandler

type GraphQLError struct {
	Message    string                 `json:"message"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Error method to satisfy the error interface
func (e *GraphQLError) Error() string {
	return e.Message
}

// NewGraphQLError creates a new GraphQLError with the given message and extensions
func NewGraphQLError(message string, extensions map[string]interface{}) *GraphQLError {
	return &GraphQLError{Message: message, Extensions: extensions}
}

func GetErrorCode(err error) ErrorCode {
	switch err.(type) {
	case *NotFoundError:
		return NOT_FOUND
	case *InvalidInputError:
		return INVALID_INPUT
	case *SightingTooCloseError:
		return CONFLICT
	case *InternalServerError:
		return INTERNAL_SERVER_ERROR
	default:
		return "UNKNOWN_ERROR" // Default code for unexpected errors
	}
}
