package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
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
	GetRestrictions(ctx context.Context, id uuid.UUID, field string) (*storage.FieldRestrictions, error)
}

func GetRestrictions(service GetRestrictionsServer, getUser GetUserServer) gin.HandlerFunc {
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

		user, err := getUser.GetUserByID(c.Request.Context(), meId, meId)
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

		var phoneRestriction FieldRestriction
		var dateRestrictions FieldRestriction

		log.Println(user.PhoneVisibility)

		if user.PhoneVisibility == models.RestrictionAll {
			phoneRestriction = FieldRestriction{
				OpenTo:         "everyone",
				SpecifiedUsers: nil,
			}
		} else if user.PhoneVisibility == models.RestrictionNone {
			phoneRestriction = FieldRestriction{
				OpenTo:         "only_me",
				SpecifiedUsers: nil,
			}
		} else {
			restrPhone, err := service.GetRestrictions(c.Request.Context(), meId, "Phone")
			if err != nil {
				if err == services.ErrNotFound {
					c.JSON(http.StatusNotFound, restapi.ErrorResponse{
						ErrorType:    restapi.ErrTypeNotFound,
						ErrorMessage: "Phone restrictions were not found",
					})
					return
				}
				c.Error(err)
				restapi.SendInternalError(c)
				return
			}
			phoneRestriction = FieldRestriction{
				OpenTo:         "specified",
				SpecifiedUsers: restrPhone.SpecifiedUsers,
			}
		}

		if user.DateOfBirthVisibility == models.RestrictionAll {
			dateRestrictions = FieldRestriction{
				OpenTo:         "everyone",
				SpecifiedUsers: nil,
			}
		} else if user.DateOfBirthVisibility == models.RestrictionNone {
			dateRestrictions = FieldRestriction{
				OpenTo:         "only_me",
				SpecifiedUsers: nil,
			}
		} else {
			restrDate, err := service.GetRestrictions(c.Request.Context(), meId, "DateOfBirth")
			if err != nil {
				if err == services.ErrNotFound {
					c.JSON(http.StatusNotFound, restapi.ErrorResponse{
						ErrorType:    restapi.ErrTypeNotFound,
						ErrorMessage: "Date of birth restrictions were not found",
					})
					return
				}
				c.Error(err)
				restapi.SendInternalError(c)
				return
			}
			dateRestrictions = FieldRestriction{
				OpenTo:         "specified",
				SpecifiedUsers: restrDate.SpecifiedUsers,
			}
		}

		restapi.SendSuccess(c, &UserRestrictions{
			Phone:       phoneRestriction,
			DateOfBirth: dateRestrictions,
		})
	}
}
