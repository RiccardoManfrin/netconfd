package nc

import (
	"os/exec"
)

const prefixInstallPAth string = "/opt/netconfd/"

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
func DHCPDelete(ifname LinkID) error {
	out, err := exec.Command(prefixInstallPAth+"dhcp_stop.sh", string(ifname)).Output()
	if err != nil {
		return NewCannotStopDHCPError(ifname, err)
	}
	if string(out) == "Service not running" {
		return NewDHCPRunningNotFoundError(ifname)
	}
	return nil
}

//DHCPCreate starts and delete DHCP controller for link interface
func DHCPCreate(dhcp Dhcp) error {
	_, err := DHCPGet(dhcp.Ifname)
	if err != nil {
		/* The only acceptable error is that you didn't find it. For any other error, abort */
		if _, ok := err.(*NotFoundError); !ok {
			return err
		}
	}

	out, err := exec.Command(prefixInstallPAth+"dhcp_start.sh", string(dhcp.Ifname)).Output()
	if err != nil {
		return NewCannotStartDHCPError(dhcp.Ifname, err)
	}

	if string(out) == "Service running already" {
		return NewDHCPAlreadyRunningConflictError(dhcp.Ifname)
	}
	return nil
}

//DHCPGet gets DHCP controller info for link interface
func DHCPGet(ifname LinkID) (Dhcp, error) {
	d := Dhcp{}
	out, err := exec.Command(prefixInstallPAth+"dhcp_status.sh", string(ifname)).Output()
	if err != nil {
		return d, NewCannotStatusDHCPError(ifname, err)
	}
	if string(out) != "active" {
		return d, NewDHCPRunningNotFoundError(ifname)
	}
	d.Ifname = ifname
	return d, nil
}

//DHCPsGet Get all DHCP interfaces administrated by DHCP and related config/state.
func DHCPsGet() ([]Dhcp, error) {
	var dhcps []Dhcp
	links, err := LinksGet()
	if err != nil {
		return dhcps, err
	}
	for _, l := range links {
		d, err := DHCPGet(l.Ifname)
		if err != nil {
			return dhcps, err
		}

		dhcps = append(dhcps, d)

	}
	return dhcps, nil
}