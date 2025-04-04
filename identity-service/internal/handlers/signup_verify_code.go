package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat-backend/identity-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/identity-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignUpVerifyCodeService interface {
	VerifyCode(ctx context.Context, signUpKey uuid.UUID, code string) error
}

func SignUpVerifyCode(service SignUpVerifyCodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signUpVerifyCodeRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		err := service.VerifyCode(c.Request.Context(), req.SignUpKey, req.Code)
		if err != nil {
			switch err {
			case services.ErrSignUpKeyNotFound:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeSignUpKeyNotFound,
					ErrorMessage: "Sign-up key doesn't exist",
				})
			case services.ErrWrongCode:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeWrongCode,
					ErrorMessage: "Wrong phone verification code",
				})
			default:
				c.Error(err)
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, signUpVerifyCodeResponse{})
	}
}

type signUpVerifyCodeRequest struct {
	SignUpKey uuid.UUID `json:"signup_key" binding:"required"`
	Code      string    `json:"code" binding:"required"`
}

type signUpVerifyCodeResponse struct{}
