package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat/backend/file-storage/internal/restapi"
	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const paramFileId = "fileId"

type GetFileService interface {
	GetFile(context.Context, uuid.UUID) (*services.FileMeta, error)
}

func GetFile(service GetFileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileId, err := uuid.Parse(c.Param(paramFileId))
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{
				{
					Field:   "fileId",
					Message: "Invalid fileId query parameter",
				},
			})
			return
		}

		file, err := service.GetFile(c, fileId)
		if err != nil {
			if err == services.ErrFileNotFound {
				c.JSON(http.StatusNotFound, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeFileNotFound,
					ErrorMessage: "File not found",
				})
				return
			}
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, fileResponse{
			FileName: file.FileName,
			FileSize: file.FileSize,
			MimeType: file.MimeType,
			FileId:   file.FileId,
			FileUrl:  file.FileUrl,
		})
	}
}
