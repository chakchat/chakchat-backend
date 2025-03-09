package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
)

type External struct {
	Publisher   publish.Publisher
	FileStorage external.FileStorage
}

func NewExternal() *External {
	return &External{
		Publisher:   publish.PublisherStub{},
		FileStorage: nil,
	}
}
