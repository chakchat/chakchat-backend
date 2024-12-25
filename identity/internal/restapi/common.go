package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendUnprocessableJSON(c *gin.Context) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		ErrorType:    ErrorTypeInvalidJson,
		ErrorMessage: "Body has invalid JSON",
	})
}

func SendValidationError(c *gin.Context, errors []ErrorDetail) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		ErrorType:    ErrorTypeValidationFailed,
		ErrorMessage: "Validation has failed",
		ErrorDetails: errors,
	})
}

func SendInternalError(c *gin.Context) {
	errResp := ErrorResponse{
		ErrorType:    ErrorTypeInternal,
		ErrorMessage: "Internal Server Error",
	}
	c.JSON(http.StatusInternalServerError, errResp)
}
