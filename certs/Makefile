.PHONY: all ca invalid-ca server client invalid

all: ca invalid-ca server client invalid

RSA_LN := 4096
ORGANIZATION := "Localhost LTD"
COUNTRY := "IT"
STATE := "Sardinia"

ca invalid-ca:
	openssl ecparam -name prime256v1 \
		-genkey -noout -out $@.key

	openssl req -new -x509 -sha256 \
		-key $@.key -out $@.crt \
		-addext "subjectAltName = DNS:localhost" \
		-subj /C=$(COUNTRY)/ST=$(STATE)/O=$(ORGANIZATION)/CN="Localhost Root"

server client:
	openssl ecparam -name prime256v1 \
		-genkey -noout -out $@.key

	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-CA ca.crt -CAkey ca.key \
  		-key $@.key -out $@.crt \
  		-addext "subjectAltName = DNS:localhost" \
  		-subj /C=$(COUNTRY)/ST=$(STATE)/O=$(ORGANIZATION)/CN=localhost

invalid:
	openssl ecparam -name prime256v1 \
		-genkey -noout -out $@.key

	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-CA invalid-ca.crt -CAkey invalid-ca.key \
  		-key $@.key -out $@.crt \
  		-addext "subjectAltName = DNS:localhost" \
  		-subj /C=$(COUNTRY)/ST=$(STATE)/O=$(ORGANIZATION)/CN=localhost
