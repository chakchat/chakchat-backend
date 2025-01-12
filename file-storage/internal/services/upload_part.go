package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	ErrUploadNotFound = errors.New("upload not found")
)

type UploadPartRequest struct {
	PartNumber int
	UploadId   uuid.UUID
	Part       io.Reader
}

type UploadPartResponse struct {
	ETag string
}

type UploadMetaGetter interface {
	Get(uuid.UUID) (*UploadMeta, bool, error)
}

type UploadPartService struct {
	metaGetter UploadMetaGetter
	client     *s3.Client
	conf       *S3Config
}

func NewUploadPartService(metaGetter UploadMetaGetter, client *s3.Client, conf *S3Config) *UploadPartService {
	return &UploadPartService{
		metaGetter: metaGetter,
		client:     client,
		conf:       conf,
	}
}

func (s *UploadPartService) UploadPart(ctx context.Context, req *UploadPartRequest) (*UploadPartResponse, error) {
	meta, ok, err := s.metaGetter.Get(req.UploadId)
	if err != nil {
		return nil, fmt.Errorf("upload-meta getting failed: %s", err)
	}
	if !ok {
		return nil, ErrUploadNotFound
	}

	res, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(s.conf.Bucket),
		Key:        aws.String(meta.Key),
		PartNumber: aws.Int32(int32(req.PartNumber)),
		UploadId:   aws.String(meta.S3UploadId),
		Body:       req.Part,
	})
	if err != nil {
		return nil, fmt.Errorf("upload part failed: %s", err)
	}

	return &UploadPartResponse{
		ETag: *res.ETag,
	}, nil
}
