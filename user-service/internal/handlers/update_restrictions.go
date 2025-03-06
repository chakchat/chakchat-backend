package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRestrictionsServer interface {
	UpdateRestrictions(ctx context.Context, id uuid.UUID, restr storage.UserRestrictions) (*models.UserRestrictions, error)
}

func UpdateRestrictions(restr UpdateRestrictionsServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRestrReq UserRestrictions
		if err := c.ShouldBindBodyWithJSON(&updateRestrReq); err != nil {
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

		phone := storage.FieldRestriction{
			OpenTo:         updateRestrReq.Phone.OpenTo,
			SpecifiedUsers: updateRestrReq.Phone.SpecifiedUsers,
		}

		date := storage.FieldRestriction{
			OpenTo:         updateRestrReq.DateOfBirth.OpenTo,
			SpecifiedUsers: updateRestrReq.DateOfBirth.SpecifiedUsers,
		}

		updatedRestr, err := restr.UpdateRestrictions(c.Request.Context(), ownerId, storage.UserRestrictions{
			Phone:       phone,
			DateOfBirth: date,
		})
		if err != nil {
			if errors.Is(err, services.ErrValidationError) {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeNotFound,
					ErrorMessage: "Restrictions was not found",
				})
				return
			}
			restapi.SendInternalError(c)
		}

		var users_phone []uuid.UUID
		for _, user := range updatedRestr.Phone.SpecifiedUsers {
			users_phone = append(users_phone, user.UserID)
		}

		var users_date []uuid.UUID
		for _, user := range updatedRestr.DateOfBirth.SpecifiedUsers {
			users_date = append(users_date, user.UserID)
		}

		restapi.SendSuccess(c, &UserRestrictions{
			Phone: FieldRestriction{
				OpenTo:         updatedRestr.Phone.OpenTo,
				SpecifiedUsers: users_phone,
			},
			DateOfBirth: FieldRestriction{
				OpenTo:         updatedRestr.DateOfBirth.OpenTo,
				SpecifiedUsers: users_date,
			},
		})
	}
}
