keys-rsa:
	mkdir keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub

keys-sym:
	openssl rand -hex 64 > keys/sym

test:
	docker compose up --build --abort-on-container-exit --exit-code-from test

gen-grpc:
	protoc --go_out="./user-service-stub/userservice" --go-grpc_out="./user-service-stub/userservice" --proto_path="../../../api" user.proto