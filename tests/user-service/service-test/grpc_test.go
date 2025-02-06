package main

import (
	"log"
	"os"
	"test/userservice"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpc(t *testing.T) {
	addr := os.Getenv("USER_SERVICE_ADDR")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connecting to UserService failed: %s", err)
	}
	client := userservice.NewUserServiceClient(conn)
	_ = client
}
