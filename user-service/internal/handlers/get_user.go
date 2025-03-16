package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Phone       *string   `json:"phone,omitempty"`
	DateOfBirth *string   `json:"dateOfBirth,omitempty"`
	PhotoURL    *string   `json:"photo"`
	CreatedAt   int64     `json:"createdAt"`
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
	CheckUserByUsername(ctx context.Context, username string) (*bool, error)
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
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		userOwner := claimId.(string)

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		userTarget, err := uuid.Parse(c.Param("userId"))
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "UserId",
					Message: "Invalid UserId parameter",
				},
			})
			return
		}

		user, err := s.service.GetUserByID(c.Request.Context(), ownerId, userTarget)
		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "Not found user with the id",
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
			Phone:       toStrPtr(user.Phone),
			DateOfBirth: toFormatDate(user.DateOfBirth),
			PhotoURL:    toStrPtr(user.PhotoURL),
			CreatedAt:   user.CreatedAt,
		})
	}
}

func (s *GetUserHandler) GetUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c.Request.Context())
		if claims == nil {
			restapi.SendUnauthorizedError(c, nil)
			return
		}
		userOwner, ok := claims[auth.ClaimId].(string)
		if !ok {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "UserId",
					Message: "Invalid user id parametr",
				},
			})
			return
		}

		user, err := s.service.GetUserByUsername(c.Request.Context(), ownerId, c.Param("username"))

		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "Not found user with the username",
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
			Phone:       toStrPtr(user.Phone),
			DateOfBirth: toFormatDate(user.DateOfBirth),
			PhotoURL:    toStrPtr(user.PhotoURL),
			CreatedAt:   user.CreatedAt,
		})
	}
}

func (s *GetUserHandler) GetUsersByCriteria() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Request.URL.Query().Get("name")
		username := c.Request.URL.Query().Get("username")
		offset := c.Request.URL.Query().Get("offset")
		var int_offset *int
		if toStrPtr(offset) == nil {
			int_offset = nil
		} else {
			offset, err := strconv.Atoi(offset)
			if err != nil {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJson,
					ErrorMessage: "Invalid offset parameter",
				})
				return
			}
			int_offset = &offset
		}
		limit := c.Request.URL.Query().Get("limit")
		var int_limit *int
		if toStrPtr(limit) == nil {
			int_limit = nil
		} else {
			limit, err := strconv.Atoi(limit)
			if err != nil {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJson,
					ErrorMessage: "Invalid limit parameter",
				})
				return
			}
			int_limit = &limit
		}
		req := storage.SearchUsersRequest{
			Name:     toStrPtr(name),
			Username: toStrPtr(username),
			Offset:   int_offset,
			Limit:    int_limit,
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
				restapi.SendSuccess(c, SearchUsersResponse{
					Users:  []User{},
					Offset: *int_offset,
					Count:  0,
				})
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
				Phone:       toStrPtr(user.Phone),
				DateOfBirth: toFormatDate(user.DateOfBirth),
				PhotoURL:    toStrPtr(user.PhotoURL),
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
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		userOwner := claimId.(string)

		ownerId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendUnauthorizedError(c, nil)
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
			Phone:       toStrPtr(me.Phone),
			PhotoURL:    toStrPtr(me.PhotoURL),
			DateOfBirth: toFormatDate(me.DateOfBirth),
			CreatedAt:   me.CreatedAt,
		})
	}
}

func (s *GetUserHandler) CheckUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		res, err := s.service.CheckUserByUsername(c.Request.Context(), username)
		if err != nil {
			restapi.SendInternalError(c)
		}
		restapi.SendSuccess(c, res)
	}
}

func toStrPtr(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func toFormatDate(date *time.Time) *string {
	if date == nil {
		return nil
	}
	formatDate := date.Format(time.DateOnly)
	return &formatDate
}
