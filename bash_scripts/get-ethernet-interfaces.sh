#!/bin/bash

# iLO 5 Redfish Ethernet Interface Lister
# Lists all network interfaces with MAC, IP, and status

# Set color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Check dependencies
command -v curl >/dev/null 2>&1 || { echo -e "${RED}Error: curl required${NC}" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo -e "${RED}Error: jq required${NC}" >&2; exit 1; }

# Initialize variables
ILO_IP=""
USERNAME=""
PASSWORD=""
INSECURE=""

usage() {
    echo "Usage: $0 -i <iLO_IP> -u <username> -p <password> [options]"
    echo "Options:"
    echo "  -k  Disable SSL verification"
    exit 1
}

# Parse arguments
while getopts ":i:u:p:k" opt; do
    case $opt in
        i) ILO_IP="$OPTARG" ;;
        u) USERNAME="$OPTARG" ;;
        p) PASSWORD="$OPTARG" ;;
        k) INSECURE="--insecure" ;;
        *) usage ;;
    esac
done

# Validate input
[[ -z "$ILO_IP" || -z "$USERNAME" || -z "$PASSWORD" ]] && usage

# Base configuration
BASE_URL="https://$ILO_IP"
CURL_OPTS=(
    --silent
    --fail
    --user "$USERNAME:$PASSWORD"
    $INSECURE
    -H "Content-Type: application/json"
)

# API request handler
send_request() {
    local endpoint="$1"
    curl "${CURL_OPTS[@]}" "$BASE_URL$endpoint"
}

# Get physical Ethernet ports
get_physical_ports() {
    echo -e "\n${YELLOW}=== Physical Ethernet Ports ===${NC}"
    
    local adapters=$(send_request "/redfish/v1/Chassis/1/NetworkAdapters/")
    echo "$adapters" | jq -r '.Members[]."@odata.id"' | while read -r adapter; do
        local ports=$(send_request "$adapter/Ports")
        echo "$ports" | jq -r '.Members[] | 
            "Port " + .Id + ":",
            "  MAC: " + (.Ethernet.AssociatedMACAddresses[0] // "N/A"),
            "  IPv4: " + "N/A",
            "  Status: " + (.Status.Health // .Status.State),
            "  Link: " + (.LinkStatus // "Unknown")'
    done
}

# Get management interfaces
get_management_interfaces() {
    echo -e "\n${YELLOW}=== Management Interfaces ===${NC}"
    
    local interfaces=$(send_request "/redfish/v1/Managers/1/EthernetInterfaces")
    echo "$interfaces" | jq -r '.Members[]."@odata.id"' | while read -r interface; do
        local data=$(send_request "$interface")
        echo "$data" | jq -r '
            "Interface " + .Id + ":",
            "  MAC: " + .MACAddress,
            "  IPv4: " + (.IPv4Addresses[0].Address // "N/A"),
            "  Status: " + (.Status.Health // .Status.State)'
    done
}

# Get system network interfaces (if available)
get_system_interfaces() {
    echo -e "\n${YELLOW}=== Host Network Interfaces ===${NC}"
    
    local interfaces=$(send_request "/redfish/v1/Systems/1/EthernetInterfaces/") 2>/dev/null
    if [[ $? -ne 0 ]]; then
        echo -e "${RED}Host interfaces unavailable (system may be powered off)${NC}"
        return
    fi
    
    echo "$interfaces" | jq -r '.Members[]."@odata.id"' | while read -r interface; do
        local data=$(send_request "$interface")
        echo "$data" | jq -r '
            "Interface " + .Id + ":",
            "  MAC: " + (.MACAddress // "N/A"),
            "  IPv4: " + (.IPv4Addresses[0].Address // "N/A"),
            "  Status: " + (.Status.Health // .Status.State)'
    done
}

# Main execution
echo -e "${GREEN}Starting iLO 5 interface discovery...${NC}"
get_physical_ports
get_management_interfaces
get_system_interfaces
echo -e "\n${GREEN}Discovery complete!${NC}"
