package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadFileRequest struct {
	FileName string
	MimeType string
	FileSize int64

	File io.ReadSeeker
}

type FileMeta struct {
	FileName  string
	MimeType  string
	FileSize  int64
	FileId    uuid.UUID
	FileUrl   string
	CreatedAt time.Time
}

type FileMetaStorer interface {
	Store(context.Context, *FileMeta) error
}

type UploadService struct {
	storer FileMetaStorer
	client *s3.Client
	conf   *S3Config
}

type S3Config struct {
	Bucket    string
	UrlPrefix string
}

func NewUploadService(storer FileMetaStorer, client *s3.Client, conf *S3Config) *UploadService {
	return &UploadService{
		storer: storer,
		client: client,
		conf:   conf,
	}
}

func (s *UploadService) Upload(ctx context.Context, req *UploadFileRequest) (*FileMeta, error) {
	fileId := uuid.New()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, req.File); err != nil {
		return nil, fmt.Errorf("failed to compute SHA-256 hash: %s", err)
	}
	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	if _, err := req.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:         aws.String(s.conf.Bucket),
		Key:            aws.String(fileId.String()),
		Body:           req.File,
		ContentType:    aws.String(req.MimeType),
		ChecksumSHA256: aws.String(hash),
	})
	if err != nil {
		return nil, fmt.Errorf("file uploading failed: %s", err)
	}

	file := &FileMeta{
		FileName:  req.FileName,
		MimeType:  req.MimeType,
		FileSize:  req.FileSize,
		FileId:    fileId,
		FileUrl:   s.conf.UrlPrefix + fileId.String(),
		CreatedAt: time.Now(),
	}

	if err := s.storer.Store(ctx, file); err != nil {
		return nil, fmt.Errorf("file meta storing failed: %s", err)
	}
	return file, nil
}
