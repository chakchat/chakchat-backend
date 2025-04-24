keys-rsa:
	mkdir keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub

keys-sym:
	openssl rand -hex 64 > keys/sym

gen: keys-rsa keys-sym

gen-grpc:
	protoc --go_out="./internal/grpcservice" --go-grpc_out="./internal/grpcservice" --proto_path="../api/user-service" user.proto