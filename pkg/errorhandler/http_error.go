package errorhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGraphQLErrors maps GraphQL errors to HTTP status codes and sends the response
func HandleGraphQLErrors(c *gin.Context, gqlErrors []*GraphQLError) {
	for _, err := range gqlErrors {
		switch err.Extensions["code"] {
		case INVALID_INPUT: // Replace INVALID_INPUT with your actual error code constant
			c.JSON(http.StatusBadRequest, err)
		case CONFLICT: // Replace INVALID_INPUT with your actual error code constant
			c.JSON(http.StatusBadRequest, err)
		case NOT_FOUND: // Replace INVALID_INPUT with your actual error code constant
			c.JSON(http.StatusNotFound, err)
		// Add more cases for other error codes as needed
		default:
			c.JSON(http.StatusInternalServerError, &GraphQLError{
				Message:    "Internal Server Error",
				Extensions: nil,
			}) // Generic error for client
		}
	}
}
