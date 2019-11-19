# DNSServ
DNSServ is a dynamic dns client and server built using golang

## Deployment
### Server
```bash
$ make certs
$ ssh dnsserv
```
```bash
jchorlton@dnsserv:~/dnsserv/certs$ sudo systemctl stop dnsserv
```
```bash
$ make deploy
```
```bash
jchorlton@dnsserv:~/dnsserv/certs$ sudo systemctl start dnsserv
```
See [scripts/README.md](scripts/README.md) for instructions on installing the systemd service.
### Client
```bash
$ make deploy-pi
$ ssh nas
```
```bash
j@orangepipcplus:~$ /home/j/dnsserv/scripts/install.sh
```

## Development
1. Generate certs
```bash
make dev-certs
```
2. Build for linux
```bash
make build
```
2. Start the server
```bash
make test-serve
```
3. Hit it
```bash
make test-client
```

## Generating certs
```bash
go get -u github.com/cloudflare/cfssl/cmd/...
cfssl print-defaults csr > ca-csr.json
cfssl print-defaults csr > server-csr.json
cfssl print-defaults csr > client-csr.json
cfssl print-defaults csr > devserver-csr.json
cfssl print-defaults config > ca-config.json
cfssl print-defaults config > server-config.json
cfssl print-defaults config > client-config.json

# now update the csrs and configs

cfssl gencert -initca ca-csr.json | cfssljson -bare ca
cfssl gencert  \
    -ca=ca.pem \
    -ca-key=ca-key.pem \
    -config=server-config.json \
    -hostname=dns.joshchorlton.com \
    -profile=www \
    server-csr.json | cfssljson -bare server
cfssl gencert \
    -ca=ca.pem \
    -ca-key=ca-key.pem \
    -config=server-config.json \
    -profile=www \
    devserver-csr.json | cfssljson -bare devserver
cfssl gencert \
    -ca=ca.pem \
    -ca-key=ca-key.pem \
    -config=client-config.json \
    -profile=client \
    client-csr.json | cfssljson -bare client
```
