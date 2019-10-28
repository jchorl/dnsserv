# Running as a systemd service
## Installing
1) Copy the service file to the systemd service dir
`sudo cp ~/dnsserv/scripts/dnsserv-server.service /etc/systemd/system/dnsserv.service`
2) Enable the service
`sudo systemctl enable /etc/systemd/system/dnsserv.service`
3) Start the service
`sudo systemctl start dnsserv.service`

## Updating
If the systemd service file is updated:
1) Copy the updated file
`sudo cp ~/dnsserv/scripts/dnsserv-server.service /etc/systemd/system/dnsserv.service`
2) Reload the daemon
`sudo systemctl daemon-reload`
