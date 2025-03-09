package proto

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/file-storage-service/internal/proto/filestorage"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/services"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCService struct {
	service *services.GetFileService
	filestorage.UnimplementedFileStorageServiceServer
}

func NewGRPCServer(service *services.GetFileService) *GRPCService {
	return &GRPCService{
		service: service,
	}
}

func (s *GRPCService) GetFile(ctx context.Context, req *filestorage.GetFileRequest) (*filestorage.GetFileResponse, error) {
	fileId, err := uuid.Parse(req.GetFileId().Value)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "Cannot parse UUID").Err()
	}

	file, err := s.service.GetFile(ctx, fileId)
	if err != nil {
		if errors.Is(err, services.ErrFileNotFound) {
			return nil, status.New(codes.NotFound, "File not found").Err()
		}
		return nil, err
	}

	return &filestorage.GetFileResponse{
		FileId:        &filestorage.UUID{Value: file.FileId.String()},
		FileName:      file.FileName,
		FileSize:      file.FileSize,
		MimeType:      file.MimeType,
		FileUrl:       file.FileUrl,
		CreatedAtUNIX: file.CreatedAt.Unix(),
	}, nil
}
