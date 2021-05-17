// Copyright (c) 2021, Athonet S.r.l. All rights reserved.
// riccardo.manfrin@athonet.com

package nc

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/riccardomanfrin/netconfd/logger"
)

const prefixInstallPAth string = "/opt/netconfd/"

// Dhcp DHCP link context to enable. When an object of this kind is specified, the DHCP protocol daemon is enabled on the  defined interface if it exists.
type Dhcp struct {
	// Interface name
	Ifname LinkID `json:"ifname,omitempty"`
}

//DHCPsConfigure configures the DHCP for each link interface of the array.
func DHCPsConfigure(dhcp []Dhcp) error {
	for _, d := range dhcp {
		if isUnmanaged(UnmanagedID(d.Ifname), LINKTYPE) {
			logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v DHCP configuration", d.Ifname))
			continue
		}
		err := DHCPDelete(d.Ifname)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok != true {
				return err
			}
		}
		if err := DHCPCreate(d); err != nil {
			return err
		}
	}
	return nil
}

//DHCPsDelete stops/deletes all DHCP control managements for each interface
func DHCPsDelete() error {
	dhcps, err := DHCPsGet()
	if err != nil {
		return err
	}
	for _, d := range dhcps {
		if isUnmanaged(UnmanagedID(d.Ifname), LINKTYPE) {
			logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v DHCP configuration", d.Ifname))
			continue
		}
		err = DHCPDelete(d.Ifname)
		if err != nil {
			return err
		}
	}
	return nil
}

//DHCPDelete stops and delete DHCP controller for link interface
func DHCPDelete(ifname LinkID) error {
	if isUnmanaged(UnmanagedID(ifname), LINKTYPE) {
		return NewUnmanagedLinkDHCPCannotBeModifiedError(ifname)
	}
	out, err := exec.Command(prefixInstallPAth+"dhcp_stop.sh", string(ifname)).Output()
	if err != nil {
		return NewCannotStopDHCPError(ifname, err)
	}
	if string(out) == "Service not running" {
		return NewDHCPRunningNotFoundError(ifname)
	}
	return nil
}

//DHCPStaticAddressesManage manages arrangements to also bring up static addresses
func DHCPStaticAddressesManage(ifname LinkID) error {
	file := "/opt/netconfd/version/etc/" + string(ifname)
	l, err := LinkGet(ifname)
	if err != nil {
		return err
	}
	if l.AddrInfo == nil || len(l.AddrInfo) == 0 {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	addresses := ""
	for _, a := range l.AddrInfo {
		addresses += a.Local.String() + "\n"
	}
	err = ioutil.WriteFile(file, []byte(addresses), 0644)
	return err
}

//DHCPCreate starts and delete DHCP controller for link interface
func DHCPCreate(dhcp Dhcp) error {
	if isUnmanaged(UnmanagedID(dhcp.Ifname), LINKTYPE) {
		return NewUnmanagedLinkDHCPCannotBeModifiedError(dhcp.Ifname)
	}
	_, err := DHCPGet(dhcp.Ifname)
	if err != nil {
		/* The only acceptable error is that you didn't find it. For any other error, abort */
		if _, ok := err.(*NotFoundError); !ok {
			return err
		}
	}

	err = DHCPStaticAddressesManage(dhcp.Ifname)
	if err != nil {
		return err
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
	_, err := LinkGet(ifname)
	if err != nil {
		return d, err
	}

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
			if _, ok := err.(*NotFoundError); ok == true {
				continue
			}
			return dhcps, err
		}

		dhcps = append(dhcps, d)

	}
	return dhcps, nil
}
