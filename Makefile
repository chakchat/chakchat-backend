gen-ssl-cert:
	sudo mkdir -p ssl
	sudo openssl genrsa -out ssl/self-signed.key 2048
	sudo openssl req -new -x509 -key ssl/self-signed.key -out ssl/self-signed.crt -days 365
