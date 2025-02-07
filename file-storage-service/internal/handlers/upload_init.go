package handlers

import (
	"context"

	"github.com/chakchat/chakchat-backend/file-storage-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MultipartUploadConfig struct {
	MinFileSize int64
	MaxPartSize int64
}

type UploadInitService interface {
	Init(context.Context, *services.UploadInitRequest) (uploadId uuid.UUID, err error)
}

func UploadInit(conf *MultipartUploadConfig, service UploadInitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req uploadInitRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		uploadId, err := service.Init(c.Request.Context(), &services.UploadInitRequest{
			FileName: req.FileName,
			MimeType: req.MimeType,
		})
		if err != nil {
			// TODO: for now I don't know what may occur here
			// But please handle errors properly.
			c.Error(err)
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, uploadInitResponse{
			UploadId: uploadId,
		})
	}
}

type uploadInitRequest struct {
	FileName string `json:"file_name" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

type uploadInitResponse struct {
	UploadId uuid.UUID `json:"upload_id"`
}
