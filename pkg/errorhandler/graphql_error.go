package errorhandler

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// NewGraphQLError creates a new GraphQLError with the given message and extensions
func NewGraphQLError(message string, extensions map[string]interface{}) *gqlerror.Error {
	return &gqlerror.Error{Message: message, Extensions: extensions}
}

func GetErrorCode(err error) ErrorCode {
	switch err.(type) {
	case *NotFoundError:
		return NOT_FOUND
	case *InvalidInputError:
		return INVALID_INPUT
	case *SightingTooCloseError:
		return CONFLICT
	case *CustomError:
		return err.(*CustomError).Code
	case *InternalServerError:
		return INTERNAL_SERVER_ERROR
	default:
		return "UNKNOWN_ERROR" // Default code for unexpected errors
	}
}
