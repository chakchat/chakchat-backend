package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignUpSendCodeService interface {
	SendCode(ctx context.Context, phone string) (signUpKey uuid.UUID, err error)
}

func SignUpSendCode(service SignUpSendCodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signUpSendCodeRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		if errors := validateSignUpSendCode(&req); len(errors) != 0 {
			restapi.SendValidationError(c, errors)
			return
		}

		signUpKey, err := service.SendCode(c, req.Phone)

		if err != nil {
			switch err {
			case services.ErrUserAlreadyExists:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUserAlreadyExists,
					ErrorMessage: "Such user already exists",
				})
			case services.ErrSendCodeFreqExceeded:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeSendCodeFreqExceeded,
					ErrorMessage: "Send code operation frequency exceeded",
				})
			default:
				log.Printf("send code endpoint failed: %s", err)
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, signUpSendCodeResponse{
			SignUpKey: signUpKey,
		})
	}
}

type signUpSendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type signUpSendCodeResponse struct {
	SignUpKey uuid.UUID `json:"signup_key"`
}

func validateSignUpSendCode(req *signUpSendCodeRequest) []restapi.ErrorDetail {
	var errors []restapi.ErrorDetail
	if !phoneRegex.MatchString(req.Phone) {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "phone",
			Message: "phone number must match a regex " + phoneRegex.String(),
		})
	}
	return errors
}
