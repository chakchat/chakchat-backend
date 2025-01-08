keys-rsa:
	mkdir keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub

keys-sym:
	openssl rand -hex 64 > keys/sym

gen-grpc:
	protoc --go_out="./internal/userservice" --go-grpc_out="./internal/userservice" --proto_path="../../api" user.proto