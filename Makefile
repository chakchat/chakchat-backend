gen-ssl-cert:
	sudo mkdir -p ssl
	sudo openssl genrsa -out ssl/self-signed.key 2048
	sudo openssl req -new -x509 -key ssl/self-signed.key -out ssl/self-signed.crt -days 365
	
keys-rsa:
	mkdir keys
	openssl genrsa -out keys/rsa 2048
	openssl rsa -in keys/rsa -pubout -out keys/rsa.pub

keys-sym:
	openssl rand -hex 64 | tr -d '\n' > keys/sym