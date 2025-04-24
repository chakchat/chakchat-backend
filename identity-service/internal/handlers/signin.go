package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat-backend/identity-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/identity-service/internal/services"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignInService interface {
	SignIn(ctx context.Context, signInKey uuid.UUID, code string, device *services.DeviceInfo) (jwt.Pair, error)
}

func SignIn(service SignInService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signInRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		if errors := validateSignIn(&req); len(errors) != 0 {
			restapi.SendValidationError(c, errors)
			return
		}

		var deviceInfo *services.DeviceInfo
		if req.Device != nil {
			deviceInfo = &services.DeviceInfo{
				DeviceToken: req.Device.DeviceToken,
				Type:        req.Device.Type,
			}
		}

		tokens, err := service.SignIn(c.Request.Context(), req.SignInKey, req.Code, deviceInfo)

		if err != nil {
			switch err {
			case services.ErrSignInKeyNotFound:
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeSignInKeyNotFound,
					ErrorMessage: "Sign in key not found",
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

		restapi.SendSuccess(c, signInResponse{
			AccessToken:  string(tokens.Access),
			RefreshToken: string(tokens.Refresh),
		})
	}
}

type DeviceInfo struct {
	Type        string `json:"type"`
	DeviceToken string `json:"device_token"`
}

type signInRequest struct {
	SignInKey uuid.UUID   `json:"signin_key" binding:"required"`
	Code      string      `json:"code" binding:"required"`
	Device    *DeviceInfo `json:"device"`
}

type signInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func validateSignIn(req *signInRequest) []restapi.ErrorDetail {
	var errors []restapi.ErrorDetail
	if req.SignInKey == uuid.Nil {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "signin_key",
			Message: "it shouldn't be nil",
		})
	}
	return errors
}
