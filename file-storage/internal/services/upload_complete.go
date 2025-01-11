package services

import "github.com/google/uuid"

type UploadCompleteRequest struct {
	UploadId uuid.UUID
	Parts    []UploadPart
}

type UploadPart struct {
	PartNumber int
}
