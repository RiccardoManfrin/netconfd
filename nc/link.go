package nc

import (
	"github.com/milosgajdos/tenus"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
)

// LinkCreate creates a link layer interface
// Link types (or kind):
// $> ip link help type
// ...
// TYPE := { vlan | veth | vcan | vxcan | dummy | ifb | macvlan | macvtap |
//	bridge | bond | team | ipoib | ip6tnl | ipip | sit | vxlan |
//	gre | gretap | erspan | ip6gre | ip6gretap | ip6erspan |
//	vti | nlmon | team_slave | bond_slave | ipvlan | geneve |
//	bridge_slave | vrf | macsec }
func LinkCreate(ifname string, kind string) error {
	switch kind {
	case "dummy":
		{
			return LinkDummyCreate(ifname)
		}
	case "bond":
		{
			return LinkBondCreate(ifname)
		}
	case "bridge":
		{
			return LinkBridgeCreate(ifname)
		}
	default:
		logger.Log.Fatal("Unknown Link Type " + kind)
	}
	return nil
}

//LinkDelete deletes a link layer interface
func LinkDelete(ifname string) error {
	return tenus.DeleteLink(ifname)
}

//LinkDummyCreate Creates a new dummy link
func LinkDummyCreate(ifname string) error {
	dl, err := tenus.NewLink(ifname)
	if err != nil {
		logger.Log.Fatal(err)
	}
	if err = dl.SetLinkUp(); err != nil {
		logger.Log.Fatal(err)
	}
	return nil
}

//LinkBondCreate Creates a new bond link
func LinkBondCreate(ifname string) error {
	dl, err := tenus.NewBondWithName(ifname)
	if err != nil {
		logger.Log.Fatal(err)
	}
	if err = dl.SetLinkUp(); err != nil {
		logger.Log.Fatal(err)
	}
	return nil
}

//LinkBridgeCreate Creates a new dummy link
func LinkBridgeCreate(ifname string) error {
	dl, err := tenus.NewBridgeWithName(ifname)
	if err != nil {
		logger.Log.Fatal(err)
	}
	if err = dl.SetLinkUp(); err != nil {
		logger.Log.Fatal(err)
	}
	return nil
}
