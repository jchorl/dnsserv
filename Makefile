UID=$(shell id -u)
GID=$(shell id -g)
GOVERSION=1.13

build:
	docker container run --rm -it \
		-v $(PWD):/dnsserv \
		-w /dnsserv \
		golang:$(GOVERSION) \
		go build -o dnsserv main.go

build-pi:
	docker container run --rm -it \
		-v $(PWD):/dnsserv \
		-w /dnsserv \
		-e GOOS=linux \
		-e GOARCH=arm \
		-e GOARM=5 \
		golang:$(GOVERSION) \
		go build -o dnsserv main.go

serve:
	sudo ./dnsserv serve \
		--ca-path $(PWD)/certs/ca.pem \
		--cert-path $(PWD)/certs/client.pem \
		--key-path $(PWD)/certs/client-key.pem \
		--dns-port 53 \
		--https-port 443

test-serve:
	docker container run --rm -it \
		-v $(PWD):/dnsserv \
		-w /dnsserv \
		--net=host \
		golang:$(GOVERSION) \
		./dnsserv serve \
			--ca-path /dnsserv/certs/ca.pem \
			--cert-path /dnsserv/certs/devserver.pem \
			--key-path /dnsserv/certs/devserver-key.pem \
			--dns-port 2012 \
			--https-port 3242 \
			--logtostderr=true

test-client:
	docker container run --rm -it \
		-v $(PWD):/dnsserv \
		-w /dnsserv \
		--net=host \
		golang:$(GOVERSION) \
		./dnsserv update \
			 --ca-path /dnsserv/certs/ca.pem \
			 --cert-path /dnsserv/certs/client.pem \
			 --key-path /dnsserv/certs/client-key.pem \
			 --dns-server https://localhost:3242 \
			 --domain pi.joshchorlton.com

deploy: build
	scp $(PWD)/dnsserv $(PWD)/Makefile dnsserv:dnsserv/
	scp $(PWD)/certs/ca.pem $(PWD)/certs/server-key.pem $(PWD)/certs/server.pem dnsserv:dnsserv/certs/
	scp $(PWD)/scripts/dnsserv-server.service dnsserv:dnsserv/scripts/

deploy-pi: build-pi
	scp $(PWD)/dnsserv $(PWD)/Makefile nas:dnsserv/
	scp $(PWD)/certs/ca.pem $(PWD)/certs/client-key.pem $(PWD)/certs/client.pem nas:dnsserv/certs/
	scp -r $(PWD)/scripts nas:dnsserv/
