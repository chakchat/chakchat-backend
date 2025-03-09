package proto

import (
	"context"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/proto/filestorage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileStorage struct {
	client filestorage.FileStorageServiceClient
}

func NewFileStorage(conn *grpc.ClientConn) *FileStorage {
	return &FileStorage{
		client: filestorage.NewFileStorageServiceClient(conn),
	}
}

func (s *FileStorage) GetById(ctx context.Context, id uuid.UUID) (*external.FileMeta, error) {
	resp, err := s.client.GetFile(ctx, &filestorage.GetFileRequest{
		FileId: &filestorage.UUID{Value: id.String()},
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return nil, external.ErrFileNotFound
		}
		return nil, err
	}

	fileId, err := uuid.Parse(resp.FileId.Value)
	if err != nil {
		return nil, fmt.Errorf("cannot parse fileId from file storage: %s", err)
	}

	return &external.FileMeta{
		FileId:    fileId,
		FileName:  resp.FileName,
		MimeType:  resp.MimeType,
		FileSize:  resp.FileSize,
		FileUrl:   resp.FileUrl,
		CreatedAt: resp.CreatedAtUNIX,
	}, nil
}
