package handlers

import (
	"errors"
	"net/http"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"

	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// type GetUserServer interface {
// 	GetUserById(ctx context.Context, ownerId uuid.UUID, targetId uuid.UUID) (*storage.User, error)
// 	GetUserByUsername(ctx context.Context, ownerId uuid.UUID, username string) (*storage.User, error)
// 	GetUsersByCriteria(ctx context.Context, req storage.SearchUsersRequest) (*storage.SearchUsersResponse, error)
// }

type GetUser struct {
	service services.GetUserService
}

func NewGetUserService(service services.GetUserService) *GetUser {
	return &GetUser{
		service: service,
	}
}

func (s *GetUser) GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c.Request.Context())
		if claims == nil {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Invalid JWT token",
			})
			return
		}
		if claims[auth.ClaimId] == nil {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Input is invalid",
			})
			return
		}
		userOwner, ok := claims[auth.ClaimId].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Input is invalid",
			})
			return
		}

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Yout JWT token is invalid",
			})
			return
		}

		userTarget, err := uuid.Parse(c.Param("userId"))
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "UserId",
					Message: "Invalid UserId query parameter",
				},
			})
			return
		}

		user, err := s.service.GetUserByID(c.Request.Context(), ownerId, userTarget)
		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "User not found",
				})
				return
			}
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Phone:       user.Phone,
			DateOfBirth: user.DateOfBirth,
			PhotoURL:    user.PhotoURL,
		})
	}
}

func (s *GetUser) GetUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c.Request.Context())
		if claims == nil {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Input is invalid",
			})
			return
		}
		userOwner, ok := claims[auth.ClaimUsername].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Input is invalid",
			})
			return
		}

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "UserId",
					Message: "Invalid UserId query parameter",
				},
			})
			return
		}

		var req models.GetUserByUsernameRequest
		user, err := s.service.GetUserByUsername(c.Request.Context(), ownerId, req.Username)

		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "User not found",
				})
				return
			}
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Phone:       user.Phone,
			DateOfBirth: user.DateOfBirth,
			PhotoURL:    user.PhotoURL,
		})
	}
}

func (s *GetUser) GetUsersByCriteria() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SearchUsersRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeInvalidJson,
				ErrorMessage: "Invalid query parameters",
			})
			return
		}

		response, err := s.service.GetUsersByCriteria(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, services.ErrNoCriteriaCpecified) {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeBadRequest,
					ErrorMessage: "Input is invalid",
				})
				return
			}
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "User not found",
				})
				return
			}
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeInternal,
				ErrorMessage: "Failed",
			})
		}
		restapi.SendSuccess(c, models.SearchUsersResponse{
			Users: response.Users,
		}.Users)
	}
}
