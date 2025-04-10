package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat-backend/file-storage-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadAbortService interface {
	Abort(ctx context.Context, uploadId uuid.UUID) error
}

func UploadAbort(service UploadAbortService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req uploadAbortRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		err := service.Abort(c.Request.Context(), req.UploadId)

		if err != nil {
			if err == services.ErrUploadNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUploadNotFound,
					ErrorMessage: "Upload not found",
				})
				return
			}
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, struct{}{})
	}
}

type uploadAbortRequest struct {
	UploadId uuid.UUID `json:"upload_id" binding:"required"`
}
