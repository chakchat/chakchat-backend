package handlers

import (
	"context"
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
	Name        string
	Username    string
	DateOfBirth *time.Time
}

type UpdateUserserver interface {
	UpdateUser(ctx context.Context, user *models.User, req *storage.UpdateUserRequest) (*models.User, error)
}

func UpdateUser(service UpdateUserserver, getter GetUserServer) gin.HandlerFunc {
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
					ErrorType:    restapi.ErrTypeValidationFailed,
					ErrorMessage: "Input is invalid",
				})
				return
			}
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		updatedUser, err := service.UpdateUser(c.Request.Context(), user, &storage.UpdateUserRequest{
			Name:        req.Name,
			Username:    req.Username,
			DateOfBirth: req.DateOfBirth,
		})
		if err != nil {
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
			PhotoURL:    updatedUser.PhotoURL,
		})
	}
}
