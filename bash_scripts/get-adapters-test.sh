#!/bin/bash

# iLO 5 Redfish Network Interface Lister
# Supports both physical adapters and management interfaces

# Set color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check for required commands
command -v curl >/dev/null 2>&1 || { echo -e "${RED}Error: curl is required but not installed.${NC}" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo -e "${RED}Error: jq is required for JSON parsing.${NC}" >&2; exit 1; }

# Initialize variables
ILO_IP=""
USERNAME=""
PASSWORD=""
INSECURE=""

usage() {
    echo "Usage: $0 -i <iLO_IP> -u <username> -p <password> [options]"
    echo "Options:"
    echo "  -k  Disable SSL certificate verification"
    echo "Example: $0 -i 192.168.1.100 -u admin -p password"
    exit 1
}

# Parse command line arguments
while getopts ":i:u:p:k" opt; do
    case $opt in
        i) ILO_IP="$OPTARG" ;;
        u) USERNAME="$OPTARG" ;;
        p) PASSWORD="$OPTARG" ;;
        k) INSECURE="--insecure" ;;
        *) usage ;;
    esac
done

# Validate required parameters
if [[ -z "$ILO_IP" || -z "$USERNAME" || -z "$PASSWORD" ]]; then
    usage
fi

# Base URL
BASE_URL="https://$ILO_IP"

# Curl common options
CURL_OPTS=(
    --silent
    --fail
    --user "$USERNAME:$PASSWORD"
    $INSECURE
)

# Function to handle API requests
send_request() {
    local endpoint="$1"
    curl "${CURL_OPTS[@]}" \
        -H "Content-Type: application/json" \
        -X GET \
        "$BASE_URL$endpoint"
}

# Get network adapters and ports
get_physical_interfaces() {
    echo -e "\n${YELLOW}=== Physical Network Adapters ===${NC}"
    
    # Get network adapters list
    adapters=$(send_request "/redfish/v1/Chassis/1/NetworkAdapters/")
    
    if ! echo "$adapters" | jq -e '.Members' > /dev/null; then
        echo -e "${RED}Error retrieving network adapters${NC}"
        return 1
    fi

    echo "$adapters" | jq -r '.Members[]."@odata.id"' | while read -r adapter; do
        adapter_id=$(basename "$adapter")
        echo -e "\n${GREEN}Adapter: $adapter_id${NC}"
        
        # Get adapter details
        adapter_data=$(send_request "$adapter")
        echo "$adapter_data" | jq -r '
            "Name: " + .Name,
            "Model: " + .Model,
            "Firmware: " + .Firmware.Current.VersionString'
        
        # Get ports for this adapter
        ports=$(send_request "$adapter/Ports")
        echo "$ports" | jq -r '.Members[] | 
            "Port " + (.Id | tostring) + ":",
            "  MAC: " + .AssociatedNetworkAddresses[0],
            "  Speed: " + (if .CurrentSpeedMbps == null then "N/A" else (.CurrentSpeedMbps | tostring + "Mbps") end),
            "  Status: " + (.Status.Health // .Status.State),
            "  Link: " + (.LinkStatus // .LinkStatus)'
    done
}

# Get management interfaces
get_management_interfaces() {
    echo -e "\n${YELLOW}=== Management Interfaces ===${NC}"
    
    interfaces=$(send_request "/redfish/v1/Managers/1/EthernetInterfaces")
    
    if ! echo "$interfaces" | jq -e '.Members' > /dev/null; then
        echo -e "${RED}Error retrieving management interfaces${NC}"
        return 1
    fi

    echo "$interfaces" | jq -r '.Members[]."@odata.id"' | while read -r interface; do
        interface_data=$(send_request "$interface")
        echo "$interface_data" | jq -r '
            "Interface: " + .Id,
            "  MAC: " + .MACAddress,
            "  IPv4: " + (.IPv4Addresses[0].Address // "N/A"),
            "  Speed: " + (if .SpeedMbps == null then "N/A" else (.SpeedMbps | tostring + "Mbps") end),
            "  Status: " + (.Status.Health // .Status.State)'
    done
}

# Main execution
echo -e "${GREEN}Starting iLO 5 interface discovery...${NC}"
get_physical_interfaces
get_management_interfaces
echo -e "\n${GREEN}Discovery complete!${NC}"
