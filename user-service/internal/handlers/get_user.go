package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"

	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"name"`
	Name        string     `json:"username"`
	Phone       *string    `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	PhotoURL    string     `json:"photo"`
	CreatedAt   int64      `json:"createdAt"`
}

type SearchUsersResponse struct {
	Users  []User `json:"users"`
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}

type GetUserServer interface {
	GetUserByID(ctx context.Context, ownerId uuid.UUID, targetId uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, ownerId uuid.UUID, username string) (*models.User, error)
	GetUsersByCriteria(ctx context.Context, req storage.SearchUsersRequest) (*storage.SearchUsersResponse, error)
}

type GetUserHandler struct {
	service GetUserServer
}

func NewGetUserHandler(service GetUserServer) *GetUserHandler {
	return &GetUserHandler{
		service: service,
	}
}

func (s *GetUserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimId, ok := auth.GetClaims(c.Request.Context())[auth.ClaimId]
		if !ok {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
			return
		}

		userOwner := claimId.(string)

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
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

		restapi.SendSuccess(c, User{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Phone:       &user.Phone,
			DateOfBirth: user.DateOfBirth,
			PhotoURL:    user.PhotoURL,
			CreatedAt:   user.CreatedAt,
		})
	}
}

func (s *GetUserHandler) GetUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c.Request.Context())
		if claims == nil {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
			return
		}
		userOwner, ok := claims[auth.ClaimUsername].(string)
		if !ok {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
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

		user, err := s.service.GetUserByUsername(c.Request.Context(), ownerId, c.Param("username"))

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

		restapi.SendSuccess(c, User{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Phone:       &user.Phone,
			DateOfBirth: user.DateOfBirth,
			PhotoURL:    user.PhotoURL,
			CreatedAt:   user.CreatedAt,
		})
	}
}

func (s *GetUserHandler) GetUsersByCriteria() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req storage.SearchUsersRequest
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
					ErrorMessage: "Criteria wasn't specified correctly",
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
				ErrorMessage: "Internal error",
			})
		}
		var users []User
		for _, user := range response.Users {
			users = append(users, User{
				ID:          user.ID,
				Username:    user.Username,
				Name:        user.Name,
				Phone:       &user.Phone,
				DateOfBirth: user.DateOfBirth,
				PhotoURL:    user.PhotoURL,
				CreatedAt:   user.CreatedAt,
			})
		}
		restapi.SendSuccess(c, SearchUsersResponse{
			Users:  users,
			Offset: response.Offset,
			Count:  response.Count,
		})
	}
}

func (s *GetUserHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimId, ok := auth.GetClaims(c.Request.Context())[auth.ClaimId]
		if !ok {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
			return
		}

		userOwner := claimId.(string)

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendUnauthorizedError(c, []restapi.ErrorDetail{})
			return
		}

		me, err := s.service.GetUserByID(c.Request.Context(), ownerId, ownerId)
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

		restapi.SendSuccess(c, User{
			ID:          me.ID,
			Name:        me.Name,
			Username:    me.Username,
			Phone:       &me.Phone,
			PhotoURL:    me.PhotoURL,
			DateOfBirth: me.DateOfBirth,
			CreatedAt:   me.CreatedAt,
		})
	}
}
