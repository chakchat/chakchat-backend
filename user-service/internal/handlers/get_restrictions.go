package handlers

import (
	"context"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRestrictions struct {
	Phone       FieldRestriction `json:"phone"`
	DateOfBirth FieldRestriction `json:"dateOfBirth"`
}

type FieldRestriction struct {
	OpenTo         string      `json:"open_to"`
	SpecifiedUsers []uuid.UUID `json:"specified_users"`
}

type GetRestrictionsServer interface {
	GetRestrictions(ctx context.Context, id uuid.UUID) (*models.UserRestrictions, error)
}

func GetRestrictions(service GetRestrictionsServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimId, ok := auth.GetClaims(c.Request.Context())[auth.ClaimId]
		if !ok {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		userOwner := claimId.(string)

		meId, err := uuid.Parse(userOwner)
		if err != nil {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		restr, err := service.GetRestrictions(c.Request.Context(), meId)
		if err != nil {
			if err == services.ErrNotFound {
				restapi.SendUnauthorizedError(c, nil)
				return
			}
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		var users_phone []uuid.UUID
		for _, user := range restr.Phone.SpecifiedUsers {
			users_phone = append(users_phone, user.UserID)
		}

		var users_date []uuid.UUID
		for _, user := range restr.DateOfBirth.SpecifiedUsers {
			users_date = append(users_date, user.UserID)
		}
		restapi.SendSuccess(c, &UserRestrictions{
			Phone: FieldRestriction{
				OpenTo:         restr.Phone.OpenTo,
				SpecifiedUsers: users_phone,
			},
			DateOfBirth: FieldRestriction{
				OpenTo:         restr.DateOfBirth.OpenTo,
				SpecifiedUsers: users_date,
			},
		})
	}
}
