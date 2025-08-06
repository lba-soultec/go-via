curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/pools \
  --insecure \
  --request POST \
  --data '{ "name": "sfo-m01", "netmask": 24, "net_address": "172.16.60.0", "start_address": "172.16.60.101", "end_address": "172.16.60.110", "gateway": "172.16.60.1", "only_serve_reimage": true, "lease_time": 3600 }'

  curl --header "Content-Type: application/json" \
  --user admin:VMware1! https://localhost:8443/v1/pools \
  --insecure \
  --request POST \
  --data '{ "name": "sfo-w01", "netmask": 24, "net_address": "172.16.61.0", "start_address": "172.16.61.101", "end_address": "172.16.61.110", "gateway": "172.16.61.1", "only_serve_reimage": true, "lease_time": 3600  }'