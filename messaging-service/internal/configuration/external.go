package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/proto"
	"google.golang.org/grpc"
)

type External struct {
	Publisher   publish.Publisher
	FileStorage external.FileStorage
}

func NewExternal(fileStConn *grpc.ClientConn, mq external.MqPublisher) *External {
	return &External{
		Publisher:   publish.NewUserEventPublisher(mq),
		FileStorage: proto.NewFileStorage(fileStConn),
	}
}
