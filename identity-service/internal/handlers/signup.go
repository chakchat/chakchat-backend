package handlers

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/chakchat/chakchat-backend/identity-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/identity-service/internal/services"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-z][_a-z0-9]{2,19}$`)
	nameRegex     = regexp.MustCompile(`^.{1,50}$`)
)

type SignUpService interface {
	SignUp(ctx context.Context, signUpKey uuid.UUID, user services.CreateUserData, deviceInfo *services.DeviceInfo) (jwt.Pair, error)
}

func SignUp(service SignUpService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signUpRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		if errors := validateSignUp(&req); len(errors) != 0 {
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

		tokens, err := service.SignUp(c.Request.Context(), req.SignUpKey, services.CreateUserData{
			Username: req.Username,
			Name:     req.Name,
		}, deviceInfo)

		if err != nil {
			switch err {
			case services.ErrUsernameAlreadyExists:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUsernameAlreadyExists,
					ErrorMessage: "Username already exists",
				})
			case services.ErrCreateUserValidationError:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeValidationFailed,
					ErrorMessage: "Validation failed while user creation",
				})
			case services.ErrSignUpKeyNotFound:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeSignUpKeyNotFound,
					ErrorMessage: "Sign up key not found",
				})
			case services.ErrPhoneNotVerified:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypePhoneNotVerified,
					ErrorMessage: "Phone is not verified",
				})
			default:
				c.Error(err)
				restapi.SendInternalError(c)
			}
			return
		}

		log.Println("sending success response")
		restapi.SendSuccess(c, signUpResponse{
			AccessToken:  string(tokens.Access),
			RefreshToken: string(tokens.Refresh),
		})
	}
}

type signUpRequest struct {
	SignUpKey uuid.UUID   `json:"signup_key" binding:"required"`
	Username  string      `json:"username" binding:"required"`
	Name      string      `json:"name" binding:"required"`
	Device    *DeviceInfo `json:"device"`
}

type signUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func validateSignUp(req *signUpRequest) []restapi.ErrorDetail {
	var errors []restapi.ErrorDetail

	if req.SignUpKey == uuid.Nil {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "signup_key",
			Message: "Sign Up key shouldn't be zero",
		})
	}

	if !usernameRegex.MatchString(req.Username) {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "username",
			Message: "Username should match regex: " + usernameRegex.String(),
		})
	}

	if !nameRegex.MatchString(req.Name) {
		errors = append(errors, restapi.ErrorDetail{
			Field:   "name",
			Message: "Name should match regex: " + nameRegex.String(),
		})
	}

	return errors
}
