#!/bin/sh
/home/j/dnsserv/dnsserv update \
  --ca-path /home/j/dnsserv/certs/root.pem \
  --cert-path /home/j/dnsserv/certs/client.pem \
  --key-path /home/j/dnsserv/certs/client-key.pem \
  --dns-server https://dns.joshchorlton.com \
  --domain home.joshchorlton.com > /home/j/logs/dnsserv.log 2>&1
