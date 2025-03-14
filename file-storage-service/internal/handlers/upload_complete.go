package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat-backend/file-storage-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadCompleteService interface {
	Complete(context.Context, *services.UploadCompleteRequest) (*services.FileMeta, error)
}

func UploadComplete(service UploadCompleteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req uploadCompleteRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		parts := make([]services.UploadPart, 0, len(req.Parts))
		for _, part := range req.Parts {
			parts = append(parts, services.UploadPart{
				PartNumber: part.PartNumber,
				ETag:       part.ETag,
			})
		}
		file, err := service.Complete(c.Request.Context(), &services.UploadCompleteRequest{
			UploadId: req.UploadId,
			Parts:    parts,
		})

		if err != nil {
			switch err {
			case services.ErrUploadNotFound:
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUploadNotFound,
					ErrorMessage: "Upload not found",
				})
			// TODO: handle occured errors
			default:
				c.Error(err)
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, fileResponse{
			FileName:  file.FileName,
			FileSize:  file.FileSize,
			MimeType:  file.MimeType,
			FileId:    file.FileId,
			FileUrl:   file.FileUrl,
			CreatedAt: file.CreatedAt,
		})
	}
}

type uploadCompleteRequest struct {
	UploadId uuid.UUID `json:"upload_id" binding:"required"`
	Parts    []struct {
		PartNumber int    `json:"part_number" binding:"required"`
		ETag       string `json:"e_tag" binding:"required"`
	} `json:"parts" binding:"required"`
}
