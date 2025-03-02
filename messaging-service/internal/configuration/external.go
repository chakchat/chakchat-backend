package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"

type External struct {
	Publisher publish.Publisher
}

func NewExternal() *External {
	return &External{
		Publisher: publish.PublisherStub{},
	}
}
