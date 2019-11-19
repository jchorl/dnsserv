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
		--ca-path $(PWD)/certs/root.pem \
		--cert-path $(PWD)/certs/leaf.pem \
		--key-path $(PWD)/certs/leaf.key \
		--dns-port 53 \
		--https-port 443

certs-dir:
	mkdir -p certs

certs: certs-dir
	# need GOCACHE env: https://github.com/golang/go/issues/26280#issuecomment-445294378
	docker container run --rm -it \
		-u $(UID):$(GID) \
		-v $(PWD)/certs:/certs \
		-e GOCACHE=/tmp \
		-w /certs \
		golang:$(GOVERSION) \
		sh -c "rm -rf /certs/* && go get -u github.com/meterup/generate-cert && generate-cert --host dns.joshchorlton.com"

tmp-dir:
	mkdir -p tmp

dev-certs: tmp-dir
	docker container run --rm -it \
		-u $(UID):$(GID) \
		-v $(PWD)/tmp:/certs \
		-e GOCACHE=/tmp \
		-w /certs \
		golang:$(GOVERSION) \
		sh -c "rm -rf /certs/* && go get github.com/meterup/generate-cert && generate-cert --host localhost"

test-serve:
	docker container run --rm -it \
		-v $(PWD):/dnsserv \
		-w /dnsserv \
		--net=host \
		golang:$(GOVERSION) \
		./dnsserv serve \
			--ca-path /dnsserv/tmp/root.pem \
			--cert-path /dnsserv/tmp/leaf.pem \
			--key-path /dnsserv/tmp/leaf.key \
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
			 --ca-path /dnsserv/tmp/root.pem \
			 --cert-path /dnsserv/tmp/client.pem \
			 --key-path /dnsserv/tmp/client.key \
			 --dns-server https://localhost:3242 \
			 --domain pi.joshchorlton.com

deploy: build
	scp $(PWD)/dnsserv $(PWD)/Makefile dnsserv:dnsserv/
	scp $(PWD)/certs/root.pem $(PWD)/certs/leaf.key $(PWD)/certs/leaf.pem dnsserv:dnsserv/certs/
	scp $(PWD)/scripts/dnsserv-server.service dnsserv:dnsserv/scripts/

deploy-pi: build-pi
	scp $(PWD)/dnsserv $(PWD)/Makefile nas:dnsserv/
	scp $(PWD)/certs/root.pem $(PWD)/certs/client.key $(PWD)/certs/client.pem nas:dnsserv/certs/
	scp -r $(PWD)/scripts nas:dnsserv/
