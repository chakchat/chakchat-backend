package services

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

func GetReceivingMembers(members []domain.UserID, sender domain.UserID) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(members)-1)
	for _, user := range members {
		if user != sender {
			res = append(res, uuid.UUID(user))
		}
	}
	return res
}

func GetReceivingUpdateMembers(members []domain.UserID, sender domain.UserID, update *domain.Update) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(members)-1)
	for _, user := range members {
		if user != sender && !update.DeletedFor(user) {
			res = append(res, uuid.UUID(user))
		}
	}
	return res
}

func NewDomainFileMeta(f *external.FileMeta) domain.FileMeta {
	return domain.FileMeta{
		FileId:    f.FileId,
		FileName:  f.FileName,
		MimeType:  f.MimeType,
		FileSize:  f.FileSize,
		FileURL:   domain.URL(f.FileUrl),
		CreatedAt: domain.Timestamp(f.CreatedAt),
	}
}
