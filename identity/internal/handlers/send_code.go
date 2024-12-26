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
	SendCode(phone string) (signInKey uuid.UUID, err error)
}

func SendCode(service SendCodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req sendCodeRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		if errors := validateSendCode(&req); len(errors) != 0 {
			restapi.SendValidationError(c, errors)
			return
		}

		signInKey, err := service.SendCode(req.Phone)

		if err != nil {
			switch err {
			case services.ErrUserNotFound:
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUserNotFound,
					ErrorMessage: "Such user doesn't exist",
				})
			case services.ErrSendCodeFreqExceeded:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeSendCodeFreqExceeded,
					ErrorMessage: "Send code operation frequency exceeded",
				})
			default:
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, sendCodeResponse{
			SignInKey: signInKey,
		})
	}
}

type sendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type sendCodeResponse struct {
	SignInKey uuid.UUID `json:"signin_key"`
}

func validateSendCode(req *sendCodeRequest) []restapi.ErrorDetail {
	var errors []restapi.ErrorDetail
	if !phoneRegex.MatchString(req.Phone) {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "phone",
			Message: "phone number must match a regex " + phoneRegex.String(),
		})
	}
	return errors
}
