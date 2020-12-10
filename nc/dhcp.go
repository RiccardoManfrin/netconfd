package nc

import "os/exec"

// Dhcp DHCP link context to enable. When an object of this kind is specified, the DHCP protocol daemon is enabled on the  defined interface if it exists.
type Dhcp struct {
	// Interface name
	Ifname LinkID `json:"ifname,omitempty"`
}

//DHCPsConfigure configures the DHCP for each link interface of the array.
func DHCPsConfigure([]Dhcp) error {
	return nil
}

//DHCPsDelete stops/deletes all DHCP control managements for each interface
func DHCPsDelete() error {
	dhcps, err := DHCPsGet()
	if err != nil {
		return err
	}
	for _, d := range dhcps {
		err = DHCPDelete(d.Ifname)
		if err != nil {
			return err
		}
	}
	return nil
}

//DHCPDelete stops and delete DHCP controller for link interface
func DHCPDelete(ifname DHCPID) error {
	out, err := exec.Command("./dhcp_stop.sh", string(ifname)).Output()
	if err != nil {
		return NewCannotStopDHCPError(ifname, string(out))
	}
	return nil
}

//DHCPCreate starts and delete DHCP controller for link interface
func DHCPCreate(dhcp Dhcp) error {
	DHCPGet()
	DHCPDelete(dhcp.Ifname)
	out, err := exec.Command("./dhcp_start.sh", string(dhcp.Ifname)).Output()
	if err != nil {
		return NewCannotStartDHCPError(DHCPIDGet(dhcp), string(out))
	}
	return nil
}

func DHCPGet(ifname LinkID) (Dhcp, err) {
	d := Dhcp{Ifname: ifname}
	out, err := exec.Command("./dhcp_status.sh", string(ifname)).Output()
	if err != nil {
		return d, NewCannotStatusDHCPError(DHCPIDGet(ifname), string(out))
	}
	return nil
}

//DHCPsGet Get all DHCP interfaces administrated by DHCP and related config/state.
func DHCPsGet() ([]Dhcp, error) {

}
