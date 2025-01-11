package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadInitRequest struct {
	FileName string
	MimeType string
}

type UploadMeta struct {
	PublicUploadId uuid.UUID
	Key            uuid.UUID
	FileName       string
	MimeType       string
	S3UploadId     string
}

type UploadMetaStorer interface {
	Store(*UploadMeta) error
}

type UploadInitService struct {
	metaStorer UploadMetaStorer
	client     *s3.Client
	conf       *S3Config
}

func NewUploadInitService(metaStorer UploadMetaStorer, client *s3.Client, conf *S3Config) *UploadInitService {
	return &UploadInitService{
		metaStorer: metaStorer,
		client:     client,
		conf:       conf,
	}
}

func (s *UploadInitService) Init(ctx context.Context, req *UploadInitRequest) (uploadId uuid.UUID, err error) {
	fileId := uuid.New()

	res, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(s.conf.Bucket),
		Key:         aws.String(fileId.String()),
		ContentType: aws.String(req.MimeType),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("creating multipart upload failed: %s", err)
	}

	meta := &UploadMeta{
		PublicUploadId: uuid.New(),
		Key:            fileId,
		FileName:       req.FileName,
		MimeType:       req.MimeType,
		S3UploadId:     *res.UploadId,
	}

	if err := s.metaStorer.Store(meta); err != nil {
		return uuid.Nil, fmt.Errorf("upload meta storing failed: %s", err)
	}
	return meta.PublicUploadId, nil
}
