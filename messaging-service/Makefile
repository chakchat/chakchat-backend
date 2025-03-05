gen:
	go generate ./...

test:
	go test -v -race ./...

gen-file-grpc:
	protoc --go_out="./internal/inrastructure/filestorage" \
		--go-grpc_out="./internal/inrastructure/filestorage" \
		--proto_path="../api/file-storage-service" \
		file-storage.proto