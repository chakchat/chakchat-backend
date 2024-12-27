.PHONY: cert
cert:
	mkdir cert
	openssl genrsa -out cert/id_rsa 2048
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub