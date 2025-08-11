package ilomapi

type IlomApi interface {
	// GetEndpoint
	GetEndpoint() string
	// GetFlavour of ilom (e.g. HPE, Dell, etc.)
	GetFlavour() string

	// GetHostConfig returns the host configuration for the given iloIpAddr and port
	GetHostConfig() ([]IfaceConfig, error)

	// SetVLANID sets the VLAN ID in the BIOS
	SetVLANID(vlanID int) error

	// SetOneTimeHTTPBoot sets the boot source to HTTP for one-time boot
	SetOneTimeHTTPBoot() error

	// RebootServer reboots the server
	RebootServer() error

	// StartServer starts the server
	StartServer() error

	// StopServer stops the server
	StopServer() error
}

type IfaceConfig struct {
	IfaceName  string `json:"ifaceName"`
	IpAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
	Speed      string `json:"speed"`  // Added Speed field
	Status     string `json:"status"` // Added Status field
}
