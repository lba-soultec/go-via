curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/groups \
  --insecure \
  --request POST \
  --data '{ "image_id": 1, "pool_id": 1, "name": "sfo-m01", "dns": "172.16.100.4,172.16.100.5", "ntp": "172.16.100.4,172.16.100.5", "netmask": "255.255.255.0", "password": "VMw@re1!VMw@re1!", "gateway": "172.16.60.1", "syslog": "logs.sfo.rainpole.io" }'

curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/groups \
  --insecure \
  --request POST \
  --data '{ "image_id": 1, "pool_id": 2, "name": "sfo-w01", "dns": "172.16.100.4,172.16.100.5", "ntp": "172.16.100.4,172.16.100.5", "netmask": "255.255.255.0", "password": "VMw@re1!VMw@re1!", "gateway": "172.16.61.1", "syslog": "logs.sfo.rainpole.io" }'