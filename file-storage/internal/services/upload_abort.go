package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadMetaGetRemover interface {
	UploadMetaGetter
	Remove(context.Context, uuid.UUID) error
}

type UploadAbortService struct {
	metaStorage UploadMetaGetRemover
	client      *s3.Client
	conf        *S3Config
}

func NewUploadAbortService(metaStorage UploadMetaGetRemover, client *s3.Client, conf *S3Config) *UploadAbortService {
	return &UploadAbortService{
		metaStorage: metaStorage,
		client:      client,
		conf:        conf,
	}
}

func (s *UploadAbortService) Abort(ctx context.Context, uploadId uuid.UUID) error {
	meta, ok, err := s.metaStorage.Get(ctx, uploadId)
	if err != nil {
		return fmt.Errorf("getting upload meta failed: %s", err)
	}
	if !ok {
		return ErrUploadNotFound
	}

	_, err = s.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(s.conf.Bucket),
		Key:      aws.String(meta.Key),
		UploadId: aws.String(meta.S3UploadId),
	})
	if err != nil {
		return fmt.Errorf("multipart upload aborting failed: %s", err)
	}

	if err := s.metaStorage.Remove(ctx, meta.PublicUploadId); err != nil {
		return fmt.Errorf("upload meta removing failed: %s", err)
	}
	return nil
}
