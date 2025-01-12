package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/chakchat/chakchat/backend/file-storage/internal/restapi"
	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	fieldFile         = "file"
	headerContentType = "Content-Type"
)

type UploadService interface {
	Upload(context.Context, *services.UploadFileRequest) (*services.FileMeta, error)
}

type UploadConfig struct {
	// File size limit in bytes
	FileSizeLimit int64
}

func Upload(conf *UploadConfig, service UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(conf.FileSizeLimit); err != nil {
			if err == multipart.ErrMessageTooLarge {
				c.JSON(http.StatusRequestEntityTooLarge, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeContentTooLarge,
					ErrorMessage: fmt.Sprintf("Conent is too large. It must be not greater than %d bytes", conf.FileSizeLimit),
				})
				return
			}
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeInvalidForm,
				ErrorMessage: restapi.ErrTypeInvalidForm,
				ErrorDetails: []restapi.ErrorDetail{},
			})
			return
		}

		file, fileHeader, err := c.Request.FormFile(fieldFile)
		if err != nil {
			restapi.SendInternalError(c)
			return
		}

		defer file.Close()
		resp, err := service.Upload(c, &services.UploadFileRequest{
			FileName: fileHeader.Filename,
			MimeType: fileHeader.Header.Get(headerContentType),
			FileSize: fileHeader.Size,
			File:     file,
		})
		if err != nil {
			// TODO: for now I don't know what may occur here
			// But please handle errors properly.
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, fileResponse{
			FileName:  resp.FileName,
			FileSize:  resp.FileSize,
			MimeType:  resp.MimeType,
			FileId:    resp.FileId,
			FileUrl:   resp.FileUrl,
			CreatedAt: resp.CreatedAt,
		})
	}
}

type fileResponse struct {
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	FileId    uuid.UUID `json:"file_id"`
	FileUrl   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}
