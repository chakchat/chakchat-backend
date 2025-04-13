package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateUserRequest struct {
	Name        string  `json:"name"`
	Username    string  `json:"username"`
	DateOfBirth *string `json:"date_of_birth"`
}

type UpdateUserServer interface {
	UpdateUser(ctx context.Context, user *models.User, req *storage.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

func UpdateUser(service UpdateUserServer, getter GetUserServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}
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

		user, err := getter.GetUserByID(c.Request.Context(), ownerId, ownerId)
		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "User not found",
				})
				return
			}
			restapi.SendInternalError(c)
			return
		}
		var date *time.Time
		if req.DateOfBirth != nil {
			cpDate, err := time.Parse(time.DateOnly, *req.DateOfBirth)
			if err != nil {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJson,
					ErrorMessage: "Wrong data format",
				})
				return
			}
			date = &cpDate
		} else {
			date = nil
		}

		updatedUser, err := service.UpdateUser(c.Request.Context(), user, &storage.UpdateUserRequest{
			Name:        req.Name,
			Username:    req.Username,
			DateOfBirth: date,
		})

		if err != nil {
			if errors.Is(err, services.ErrValidationError) {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJson,
					ErrorMessage: "Wrong username",
				})
				return
			}
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "UserId",
					Message: "Invalid UserId query parameter",
				},
			})
			return
		}

		restapi.SendSuccess(c, User{
			ID:          updatedUser.ID,
			Username:    updatedUser.Username,
			Name:        updatedUser.Name,
			Phone:       toStrPtr(updatedUser.Phone),
			DateOfBirth: toFormatDate(updatedUser.DateOfBirth),
			PhotoURL:    toStrPtr(updatedUser.PhotoURL),
		})
	}
}

func DeleteMe(service UpdateUserServer) gin.HandlerFunc {
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

		err = service.DeleteUser(c.Request.Context(), ownerId)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				restapi.SendUnauthorizedError(c, nil)
			}
			restapi.SendInternalError(c)
		}

		restapi.SendSuccess(c, struct{}{})
	}
}
