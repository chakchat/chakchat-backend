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

const (
	phoneField string = "Phone"
	dateField  string = "DateOfBirth"
)

type UpdateRestrictionsServer interface {
	UpdateRestrictions(ctx context.Context, id uuid.UUID, restr storage.FieldRestrictions) (*storage.FieldRestrictions, error)
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

		var phoneRestriction FieldRestriction
		var dateRestriction FieldRestriction

		if updateRestrReq.Phone.OpenTo != models.RestrictionSpecified {
			phoneRestriction = FieldRestriction{
				OpenTo:         updateRestrReq.Phone.OpenTo,
				SpecifiedUsers: nil,
			}
		} else {
			if updateRestrReq.Phone.SpecifiedUsers == nil {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeBadRequest,
					ErrorMessage: "Specified users for phone restrictions were not specified",
				})
				return
			}
			phone := storage.FieldRestrictions{
				Field:          phoneField,
				OpenTo:         models.RestrictionSpecified,
				SpecifiedUsers: updateRestrReq.Phone.SpecifiedUsers,
			}
			phoneRestriction = FieldRestriction{
				OpenTo:         updateRestrReq.Phone.OpenTo,
				SpecifiedUsers: phone.SpecifiedUsers,
			}
			_, err := restr.UpdateRestrictions(c.Request.Context(), ownerId, phone)
			if err != nil {
				if errors.Is(err, services.ErrValidationError) {
					c.JSON(http.StatusNotFound, restapi.ErrorResponse{
						ErrorType:    restapi.ErrTypeNotFound,
						ErrorMessage: "Phone restrictions was not found",
					})
					return
				}
				c.Error(err)
				restapi.SendInternalError(c)
				return
			}
		}

		if updateRestrReq.DateOfBirth.OpenTo != models.RestrictionSpecified {
			dateRestriction = FieldRestriction{
				OpenTo:         updateRestrReq.DateOfBirth.OpenTo,
				SpecifiedUsers: nil,
			}
		} else {
			if updateRestrReq.DateOfBirth.SpecifiedUsers == nil {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeBadRequest,
					ErrorMessage: "Specified users for dateOfBirth restrictions were not specified",
				})
				return
			}
			date := storage.FieldRestrictions{
				Field:          dateField,
				OpenTo:         models.RestrictionSpecified,
				SpecifiedUsers: updateRestrReq.DateOfBirth.SpecifiedUsers,
			}

			dateRestriction = FieldRestriction{
				OpenTo:         updateRestrReq.DateOfBirth.OpenTo,
				SpecifiedUsers: date.SpecifiedUsers,
			}

			_, err := restr.UpdateRestrictions(c.Request.Context(), ownerId, date)
			if err != nil {
				if errors.Is(err, services.ErrValidationError) {
					c.JSON(http.StatusNotFound, restapi.ErrorResponse{
						ErrorType:    restapi.ErrTypeNotFound,
						ErrorMessage: "Date of birth restrictions was not found",
					})
					return
				}
				c.Error(err)
				restapi.SendInternalError(c)
				return
			}
		}

		restapi.SendSuccess(c, &UserRestrictions{
			Phone:       phoneRestriction,
			DateOfBirth: dateRestriction,
		})
	}
}
