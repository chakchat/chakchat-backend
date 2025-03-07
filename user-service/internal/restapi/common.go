package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Data: data,
	})
}

func SendUnprocessableJSON(c *gin.Context) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		ErrorType:    ErrTypeInvalidJson,
		ErrorMessage: "Body has invalid JSON",
	})
}

func SendValidationError(c *gin.Context, errors []ErrorDetail) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		ErrorType:    ErrTypeValidationFailed,
		ErrorMessage: "Validation has failed",
		ErrorDetails: errors,
	})
}

func SendUnauthorizedError(c *gin.Context, errors []ErrorDetail) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		ErrorType:    ErrTypeUnautorized,
		ErrorMessage: "Failed JWT token authentication",
		ErrorDetails: errors,
	})
}

func SendInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		ErrorType:    ErrTypeInternal,
		ErrorMessage: "Internal Server Error",
	})
}
