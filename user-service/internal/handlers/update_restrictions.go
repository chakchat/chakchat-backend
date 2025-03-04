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
	UpdateRestrictions(ctx context.Context, id uuid.UUID, phone storage.FieldRestriction, date storage.FieldRestriction) (*models.UserRestrictions, error)
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

		updatedRestr, err := restr.UpdateRestrictions(c.Request.Context(), ownerId, phone, date)
		if err != nil {
			if errors.Is(err, services.ErrValidationError) {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeBadRequest,
					ErrorMessage: "Input is invalid",
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
			Phone: struct {
				OpenTo         string      "json:\"open_to\""
				SpecifiedUsers []uuid.UUID "json:\"specified_users\""
			}{
				OpenTo:         updatedRestr.Phone.OpenTo,
				SpecifiedUsers: users_phone,
			},
			DateOfBirth: struct {
				OpenTo         string      "json:\"open_to\""
				SpecifiedUsers []uuid.UUID "json:\"specified_users\""
			}{
				OpenTo:         updatedRestr.DateOfBirth.OpenTo,
				SpecifiedUsers: users_date,
			},
		})
	}
}
