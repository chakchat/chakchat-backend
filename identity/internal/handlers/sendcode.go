package handlers

import (
	"net/http"
	"regexp"

	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Only russian phone numbers for now. Just hard-coded
var phoneRegex = regexp.MustCompile(`^[78]9\d{9}$`)

type SendCodeService interface {
	SendCode(phone string, signInKey uuid.UUID) error
}

func SendCode(service SendCodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req sendCodeRequest
		err := c.ShouldBindBodyWithJSON(&req)
		if err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		if errors := validateSendCode(&req); len(errors) != 0 {
			restapi.SendValidationError(c, errors)
			return
		}

		err = service.SendCode(req.Phone, req.SignInKey)
		if err == services.ErrUserNotFound {
			resp := restapi.ErrorResponse{
				ErrorType:    restapi.ErrorTypeUserNotFound,
				ErrorMessage: "Such user doesn't exist",
			}
			c.JSON(http.StatusNotFound, resp)
			return
		} else if err != nil {
			restapi.SendInternalError(c)
		}

		c.JSON(http.StatusOK, restapi.SuccessResponse{
			Data: sendCodeResponse{},
		})
	}
}

type sendCodeRequest struct {
	Phone     string    `json:"phone"`
	SignInKey uuid.UUID `json:"signin_key"`
}

type sendCodeResponse struct {
}

func validateSendCode(req *sendCodeRequest) []restapi.ErrorDetail {
	var errors []restapi.ErrorDetail
	if !phoneRegex.MatchString(req.Phone) {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "phone",
			Message: "phone number must match a regex " + phoneRegex.String(),
		})
	}
	if req.SignInKey == uuid.Nil {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "signin_key",
			Message: "it shouldn't be nil",
		})
	}
	return errors
}
