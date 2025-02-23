test:
	docker compose up --build --abort-on-container-exit --exit-code-from test

gen-grpc:
	protoc --go_out="./service-test/userservice" --go-grpc_out="./service-test/userservice" --proto_path="../../api/user-service" user.proto

keys-rsa:
	mkdir keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub