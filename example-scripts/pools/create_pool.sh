curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/api/v1/pools \
  --insecure \
  --request POST \
  --data '{ "name": "sfo-m01", "netmask": 24, "net_address": "172.16.60.0", "gateway": "172.16.60.1", "only_serve_reimage": true }'