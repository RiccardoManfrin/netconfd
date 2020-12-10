package nc

// Dhcp DHCP link context to enable. When an object of this kind is specified, the DHCP protocol daemon is enabled on the  defined interface if it exists.
type Dhcp struct {
	// Interface name
	Ifname string `json:"ifname,omitempty"`
}

//DHCPsConfigure configures the DHCP for each link interface of the array.
func DHCPsConfigure([]Dhcp) error {
	return nil
}
