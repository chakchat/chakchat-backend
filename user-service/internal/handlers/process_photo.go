package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updatePhotoRequest struct {
	PhotoId string `json:"photo_id"`
}

type UpdatePhotoServer interface {
	UpdatePhoto(ctx context.Context, id uuid.UUID, photoId string) (*models.User, error)
	DeletePhoto(ctx context.Context, id uuid.UUID) (*models.User, error)
}

func UpdatePhoto(u UpdatePhotoServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updatePhotoRequest
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

		user, err := u.UpdatePhoto(c.Request.Context(), ownerId, req.PhotoId)
		if err != nil {
			if errors.Is(err, services.ErrInvalidPhoto) {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeBadRequest,
					ErrorMessage: "Invalid photo",
				})
				return
			}

			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "Can't find photo",
				})
				return
			}

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

func DeletePhoto(u UpdatePhotoServer) gin.HandlerFunc {
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

		user, err := u.DeletePhoto(c.Request.Context(), ownerId)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "Can't find user with owner id",
				})
				return
			}

			restapi.SendInternalError(c)
			return
		}
		restapi.SendSuccess(c, User{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Phone:       toStrPtr(user.Phone),
			DateOfBirth: toFormatDate(user.DateOfBirth),
			PhotoURL:    nil,
			CreatedAt:   user.CreatedAt,
		})
	}
}
