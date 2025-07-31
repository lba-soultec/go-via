curl --header "Content-Type: application/json" \
  http://localhost:8443/api/v1/groups \
  --request POST \
  --data '{ "image_id": 1, "name": "testgrp", "password": "VMware1!", "dns": "172.16.100.4,172.16.100.5", "ntp": "172.16.100.4,172.16.100.5"}'

curl --header "Content-Type: application/json" \
  http://localhost:8443/api/v1/addresses \
  --request POST \
  --data '{ "domain": "vmlab.se", "group_id": 1, "hostname": "testhost", "ip": "172.16.100.12", "progress": 0, "progresstext": "", "reimage": false }'
