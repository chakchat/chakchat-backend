package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type UploadCompleteRequest struct {
	UploadId uuid.UUID
	Parts    []UploadPart
}

type UploadPart struct {
	PartNumber int
	ETag       string
}

type UploadCompleteService struct {
	fileMetaStorage FileMetaStorer
	uploadStorage   UploadMetaGetRemover
	client          *s3.Client
	conf            *S3Config
}

func NewUploadCompleteService(fileMetaStorage FileMetaStorer, uploadStorage UploadMetaGetRemover,
	client *s3.Client, conf *S3Config) *UploadCompleteService {
	return &UploadCompleteService{
		fileMetaStorage: fileMetaStorage,
		uploadStorage:   uploadStorage,
		client:          client,
		conf:            conf,
	}
}

func (s *UploadCompleteService) Complete(ctx context.Context, req *UploadCompleteRequest) (*FileMeta, error) {
	upload, ok, err := s.uploadStorage.Get(ctx, req.UploadId)
	if err != nil {
		return nil, fmt.Errorf("upload-meta getting failed: %s", err)
	}
	if !ok {
		return nil, ErrUploadNotFound
	}

	res, err := s.completeS3Upload(ctx, upload, req)
	if err != nil {
		return nil, fmt.Errorf("complete multipart upload failed: %s", err)
	}

	if err := s.uploadStorage.Remove(ctx, upload.PublicUploadId); err != nil {
		return nil, fmt.Errorf("upload-meta removing failed: %s", err)
	}

	head, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: res.Bucket,
		Key:    res.Key,
	})
	if err != nil {
		return nil, fmt.Errorf("getting head object failed: %s", err)
	}

	file := &FileMeta{
		FileName:  upload.FileName,
		MimeType:  upload.MimeType,
		FileSize:  *head.ContentLength,
		FileId:    upload.FileId,
		FileUrl:   *res.Location,
		CreatedAt: time.Now(),
	}
	if err := s.fileMetaStorage.Store(ctx, file); err != nil {
		return nil, fmt.Errorf("storing file metadata failed: %s", err)
	}
	return file, nil
}

func (s *UploadCompleteService) completeS3Upload(ctx context.Context, upload *UploadMeta, req *UploadCompleteRequest) (*s3.CompleteMultipartUploadOutput, error) {
	parts := make([]types.CompletedPart, 0, len(req.Parts))
	for _, p := range req.Parts {
		parts = append(parts, types.CompletedPart{
			ETag:       aws.String(p.ETag),
			PartNumber: aws.Int32(int32(p.PartNumber)),
		})
	}
	return s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.conf.Bucket),
		Key:      aws.String(upload.Key),
		UploadId: aws.String(upload.S3UploadId),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
}
