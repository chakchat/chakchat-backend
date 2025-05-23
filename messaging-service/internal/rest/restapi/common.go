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

func SendInternalError(c *gin.Context) {
	errResp := ErrorResponse{
		ErrorType:    ErrTypeInternal,
		ErrorMessage: "Internal Server Error",
	}
	c.JSON(http.StatusInternalServerError, errResp)
}

func SendInvalidChatID(c *gin.Context) {
	SendValidationError(c, []ErrorDetail{
		{
			Field:   "chatId",
			Message: "Invalid chatId route parameter",
		},
	})
}

func SendInvalidUpdateID(c *gin.Context) {
	SendValidationError(c, []ErrorDetail{
		{
			Field:   "updateId",
			Message: "Invalid updateId route parameter",
		},
	})
}

func SendInvalidMemberID(c *gin.Context) {
	SendValidationError(c, []ErrorDetail{
		{
			Field:   "chatId",
			Message: "Invalid chatId route parameter",
		},
	})
}
