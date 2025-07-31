curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/api/v1/groups \
  --insecure \
  --request POST \
  --data '{ "image_id": 1, "pool_id": 1, "name": "sfo-m01", "dns": "172.16.100.4,172.16.100.5", "ntp": "172.16.100.4,172.16.100.5", "netmask": "255.255.255.0", "password": "VMware1!", "gateway": "172.16.60.1", "device": "vmnic0", "syslog": "logs.sfo.rainpole.io" }'

curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/api/v1/hosts \
  --insecure \
  --request POST \
  --data '{ "domain": "sfo.rainpole.io", "group_id": 1, "pool_id": 1, "hostname": "sfo01-m01-esx01", "ip": "172.16.60.101", "mac": "00:50:56:8a:7c:09", "reimage": true }'
