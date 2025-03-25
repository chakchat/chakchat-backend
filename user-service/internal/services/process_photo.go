package services

import (
	"context"
	"errors"
	"log"

	"github.com/chakchat/chakchat-backend/user-service/internal/filestorage"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidPhoto = errors.New("Invalid photo")

var groupPhotoMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
	"image/heif": true,
	"image/heic": true,
}

type ProcessPhotoRepo interface {
	UpdatePhoto(ctx context.Context, id uuid.UUID, photoURL string) (*models.User, error)
	DeletePhoto(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type ProcessPhotoService struct {
	repo        ProcessPhotoRepo
	fileStorage filestorage.FileStorageServiceClient
}

func NewProcessPhotoService(repo ProcessPhotoRepo, fileStorage filestorage.FileStorageServiceClient) *ProcessPhotoService {
	return &ProcessPhotoService{
		repo:        repo,
		fileStorage: fileStorage,
	}
}

func (u *ProcessPhotoService) UpdatePhoto(ctx context.Context, id uuid.UUID, photoId string) (*models.User, error) {
	photo, err := u.fetchPhotoURL(ctx, photoId)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.UpdatePhoto(ctx, id, photo.FileUrl)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *ProcessPhotoService) DeletePhoto(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.repo.DeletePhoto(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *ProcessPhotoService) fetchPhotoURL(ctx context.Context, photo string) (*filestorage.GetFileResponse, error) {
	photoURL, err := u.fileStorage.GetFile(ctx, &filestorage.GetFileRequest{
		FileId: &filestorage.UUID{Value: photo},
	})

	if err != nil {
		log.Println("Error fetching photo")
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := validatePhoto(photoURL); err != nil {
		return nil, err
	}

	return photoURL, nil
}

func validatePhoto(photo *filestorage.GetFileResponse) error {
	if photo.FileSize > 2<<20 {
		return ErrInvalidPhoto
	}

	if !groupPhotoMimes[photo.MimeType] {
		return ErrInvalidPhoto
	}
	return nil
}
