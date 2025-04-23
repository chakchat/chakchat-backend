package grpc

type GRPCClients struct {
}

func NewGrpcClients() *GRPCClients {
	return &GRPCClients{}
}

func (c *GRPCClients) GetChatType() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetGroupName() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetUsername() (string, error) {
	return "", nil
}
