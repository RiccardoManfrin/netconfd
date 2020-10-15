package tenus

import (
	"bytes"
	"fmt"
	"net"

	"github.com/docker/libcontainer/netlink"
)

// Bondr embeds Linker interface and adds one extra function.
type Bondr interface {
	// Linker interface
	Linker
	// AddSlaveIfc adds network interface to the network Bond
	AddSlaveIfc(*net.Interface) error
	//RemoveSlaveIfc removes network interface from the network Bond
	RemoveSlaveIfc(*net.Interface) error
}

// Bond is Link which has zero or more slave network interfaces.
// Bond implements Bondr interface.
type Bond struct {
	Link
	slaveIfcs []net.Interface
}

// NewBond creates new network Bond on Linux host.
//
// It is equivalent of running: ip link add name br${RANDOM STRING} type Bond
// NewBond returns Bondr which is initialized to a pointer of type Bond if the
// Bond was created successfully on the Linux host. Newly created Bond is assigned
// a random name starting with "br".
// It returns error if the Bond could not be created.
func NewBond() (Bondr, error) {
	brDev := makeNetInterfaceName("bond")

	if ok, err := NetInterfaceNameValid(brDev); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(brDev); err == nil {
		return nil, fmt.Errorf("Interface name %s already assigned on the host", brDev)
	}

	if err := netlink.NetworkLinkAdd(brDev, "Bond"); err != nil {
		return nil, err
	}

	newIfc, err := net.InterfaceByName(brDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &Bond{
		Link: Link{
			ifc: newIfc,
		},
	}, nil
}

// NewBondWithName creates new network Bond on Linux host with the name passed as a parameter.
// It is equivalent of running: ip link add name ${ifcName} type Bond
// It returns error if the Bond can not be created.
func NewBondWithName(ifcName string) (Bondr, error) {
	if ok, err := NetInterfaceNameValid(ifcName); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(ifcName); err == nil {
		return nil, fmt.Errorf("Interface name %s already assigned on the host", ifcName)
	}

	if err := netlink.NetworkLinkAdd(ifcName, "bond"); err != nil {
		return nil, err
	}

	newIfc, err := net.InterfaceByName(ifcName)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &Bond{
		Link: Link{
			ifc: newIfc,
		},
	}, nil
}

// BondFromName returns a tenus network Bond from an existing Bond of given name on the Linux host.
// It returns error if the Bond of the given name cannot be found.
func BondFromName(ifcName string) (Bondr, error) {
	if ok, err := NetInterfaceNameValid(ifcName); !ok {
		return nil, err
	}

	newIfc, err := net.InterfaceByName(ifcName)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &Bond{
		Link: Link{
			ifc: newIfc,
		},
	}, nil
}

// AddToBond adds network interfaces to network Bond.
// It is equivalent of running: ip link set ${netIfc name} master ${netBond name}
// It returns error when it fails to add the network interface to Bond.
func AddToBond(netIfc, netBond *net.Interface) error {
	return netlink.NetworkSetMaster(netIfc, netBond)
}

// AddToBond adds network interfaces to network Bond.
// It is equivalent of running: ip link set dev ${netIfc name} nomaster
// It returns error when it fails to remove the network interface from the Bond.
func RemoveFromBond(netIfc *net.Interface) error {
	return netlink.NetworkSetNoMaster(netIfc)
}

// AddSlaveIfc adds network interface to network Bond.
// It is equivalent of running: ip link set ${ifc name} master ${Bond name}
// It returns error if the network interface could not be added to the Bond.
func (br *Bond) AddSlaveIfc(ifc *net.Interface) error {
	if err := netlink.NetworkSetMaster(ifc, br.ifc); err != nil {
		return err
	}

	br.slaveIfcs = append(br.slaveIfcs, *ifc)

	return nil
}

// RemoveSlaveIfc removes network interface from the network Bond.
// It is equivalent of running: ip link set dev ${netIfc name} nomaster
// It returns error if the network interface is not in the Bond or
// it could not be removed from the Bond.
func (br *Bond) RemoveSlaveIfc(ifc *net.Interface) error {
	if err := netlink.NetworkSetNoMaster(ifc); err != nil {
		return err
	}

	for index, i := range br.slaveIfcs {
		// I could reflect.DeepEqual(), but there is not point to import reflect for one operation
		if i.Name == ifc.Name && bytes.Equal(i.HardwareAddr, ifc.HardwareAddr) {
			br.slaveIfcs = append(br.slaveIfcs[:index], br.slaveIfcs[index+1:]...)
		}
	}

	return nil
}
