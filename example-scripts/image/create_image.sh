curl -X POST -H "Content-Type: multipart/form-data" -F "file[]=@/Users/kimjohansson/Downloads/VMware-VMvisor-Installer-8.0U3-24022510.x86_64.iso" http://localhost:8443/api/v1/images

curl -H "Content-Type: multipart/form-data" \
  --user admin:VMware1! https://localhost:8443/api/v1/images \
  --insecure \
  --request POST \
  -F "file[]=@/home/kim/VMware-VMvisor-Installer-9.0.0.0.24755229.x86_64.iso"
