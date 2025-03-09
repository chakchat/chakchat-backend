gen:
	go generate ./...

test:
	go test -v -race ./...

gen-file-grpc:
	protoc --go_out="./internal/infrastructure/proto" \
		--go-grpc_out="./internal/inrastructure/proto" \
		--proto_path="../api/file-storage-service" \
		file-storage.proto