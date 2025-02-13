gen-ssl-cert:
	mkdir -p ssl
	openssl genrsa -out ssl/self-signed.key 2048
	openssl req -new -x509 -key ssl/self-signed.key -out ssl/self-signed.crt -days 365
	
keys-rsa:
	mkdir -p keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub

keys-sym:
	openssl rand -hex 64 | tr -d '\n' > keys/sym

gen: gen-ssl-cert keys-rsa keys-sym

run: 
	docker-compose up -d --build

down:
	docker-compose down

clean:
	docker-compose down --volumes

unit-test:
	cd identity-service && go test ./...
	cd file-storage-service && go test ./...
	cd shared/go && go test ./...

identity-service-test:
	cd tests/identity-service && make test

.PHONY: test
test: unit-test identity-service-test
	echo "All tests passed"