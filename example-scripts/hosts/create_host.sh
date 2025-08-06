curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/hosts \
  --insecure \
  --request POST \
  --data '{ "domain": "sfo.rainpole.io", "group_id": 1, "pool_id": 1, "hostname": "sfo01-m01-esx01", "ip": "172.16.60.101", "mac": "00:50:56:8a:7c:09", "reimage": true }'

curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/hosts \
  --insecure \
  --request POST \
  --data '{ "domain": "sfo.rainpole.io", "group_id": 2, "pool_id": 2, "hostname": "sfo01-w01-esx01", "ip": "172.16.61.101", "mac": "00:50:56:8a:73:65", "reimage": true }'
