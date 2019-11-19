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
