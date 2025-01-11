package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/chakchat/chakchat/backend/file-storage/internal/restapi"
	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	fieldPartNumber = "part_number"
	fieldUploadId   = "upload_id"
)

type UploadPartService interface {
	UploadPart(ctx context.Context, partNumber int, uploadId uuid.UUID, part io.Reader) error
}

func UploadPart(conf *MultipartUploadConfig, service UploadPartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(conf.MaxPartSize); err != nil {
			if err == multipart.ErrMessageTooLarge {
				c.JSON(http.StatusRequestEntityTooLarge, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeContentTooLarge,
					ErrorMessage: fmt.Sprintf("Conent is too large. It must be not greater than %d bytes", conf.MaxPartSize),
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

		partNumber, err := strconv.Atoi(c.Request.FormValue(fieldPartNumber))
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{{
				Field:   fieldPartNumber,
				Message: "Part number is missing",
			}})
			return
		}
		uploadId, err := uuid.Parse(c.Request.FormValue(fieldUploadId))
		if err != nil {
			restapi.SendValidationError(c, []restapi.ErrorDetail{{
				Field:   fieldPartNumber,
				Message: "",
			}})
			return
		}

		filePart, _, err := c.Request.FormFile(fieldFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeInvalidForm,
				ErrorMessage: "Can't parse form file",
			})
			return
		}

		err = service.UploadPart(c, partNumber, uploadId, filePart)

		if err != nil {
			if err == services.ErrUploadNotFound {
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUploadNotFound,
					ErrorMessage: "Upload not found",
				})
				return
			}
			restapi.SendInternalError(c)
			return
		}

		restapi.SendSuccess(c, struct{}{})
	}
}
