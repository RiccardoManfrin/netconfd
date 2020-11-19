package nc

import (
	"github.com/vishvananda/netlink"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
)

// LinkAddrInfo struct for LinkAddrInfo
type LinkAddrInfo struct {
	Local CIDRAddr `json:"local,omitempty"`
	//Prefixlen int32    `json:"prefixlen,omitempty"`
	//Broadcast CIDRAddr `json:"broadcast,omitempty"`
}

// LinkLinkinfoInfoSlaveData Info about slave state/config
type LinkLinkinfoInfoSlaveData struct {
	// State of the link:   * `ACTIVE` - Link is actively used   * `BACKUP` - Link is used for failover
	State string `json:"state,omitempty"`
	// MII Status:   * `UP`    * `DOWN`
	MiiStatus string `json:"mii_status,omitempty"`
	// Number of link failures
	LinkFailureCount uint32 `json:"link_failure_count,omitempty"`
	// Hardware address
	PermHwaddr string `json:"perm_hwaddr,omitempty"`
	// Queue Identifier
	QueueId uint16 `json:"queue_id,omitempty"`
}

// LinkLinkinfoInfoData Additional information on the link
type LinkLinkinfoInfoData struct {
	// Bonding modes. Supported Modes:   * `balance-rr` - Round-robin: Transmit network packets in sequential order from the first available network interface (NIC) slave through the last. This mode provides load balancing and fault tolerance.   * `active-backup` - Active-backup: Only one NIC slave in the bond is active. A different slave becomes active if, and only if, the active slave fails. The single logical bonded interface's MAC address is externally visible on only one NIC (port) to avoid distortion in the network switch. This mode provides fault tolerance.   * `balance-xor` - XOR: Transmit network packets based on a hash of the packet's source and destination. The default algorithm only considers MAC addresses (layer2). Newer versions allow selection of additional policies based on IP addresses (layer2+3) and TCP/UDP port numbers (layer3+4). This selects the same NIC slave for each destination MAC address, IP address, or IP address and port combination, respectively. This mode provides load balancing and fault tolerance.   * `broadcast` - Broadcast: Transmit network packets on all slave network interfaces. This mode provides fault tolerance.   * `802.3ad` - IEEE 802.3ad Dynamic link aggregation: Creates aggregation groups that share the same speed and duplex settings. Utilizes all slave network interfaces in the active aggregator group according to the 802.3ad specification. This mode is similar to the XOR mode above and supports the same balancing policies. The link is set up dynamically between two LACP-supporting peers.   * `balance-tlb` - Adaptive transmit load balancing: Linux bonding driver mode that does not require any special network-switch support. The outgoing network packet traffic is distributed according to the current load (computed relative to the speed) on each network interface slave. Incoming traffic is received by one currently designated slave network interface. If this receiving slave fails, another slave takes over the MAC address of the failed receiving slave.   * `balance-alb` - Adaptive load balancing: includes balance-tlb plus receive load balancing (rlb) for IPV4 traffic, and does not require any special network switch support. The receive load balancing is achieved by ARP negotiation. The bonding driver intercepts the ARP Replies sent by the local system on their way out and overwrites the source hardware address with the unique hardware address of one of the NIC slaves in the single logical bonded interface such that different network-peers use different MAC addresses for their network packet traffic.
	Mode string `json:"mode,omitempty"`
	// Specifies the MII link monitoring frequency in milliseconds.  The default value is 0, and this will disable the MII monitor
	Miimon int32 `json:"miimon,omitempty"`
	// Specifies the time, in milliseconds, to wait before enabling a slave after a  link recovery has been detected. The updelay value should be a multiple of the miimon value
	Updelay int32 `json:"updelay,omitempty"`
	// Specifies the time, in milliseconds, to wait before disabling a slave after a  link failure has been detected. The downdelay value should be a multiple of the miimon value.
	Downdelay int32 `json:"downdelay,omitempty"`
}

//LinkLinkinfo definition
type LinkLinkinfo struct {
	// Type of link layer interface. Supported Types:   * `dummy` - Dummy link type interface for binding intenal services   * `bridge` - Link layer virtual switch type interface   * `bond` - Bond type interface letting two interfaces be seen as one   * `vlan` - Virtual LAN (TAG ID based) interface   * `veth` - Virtual ethernet (with virtual MAC and IP address)   * `macvlan` - Direct virtual eth interface connected to the physical interface,      with owned mac address   * `ipvlan` - Direct virtual eth interface connected to the physical interface.     Physical interface MAC address is reused. L2 type directly connects the lan to      the host phyisical device. L3 type adds a routing layer in between.
	InfoKind string `json:"info_kind,omitempty"`
	// FILL ME
	InfoSlaveKind string                    `json:"info_slave_kind,omitempty"`
	InfoSlaveData LinkLinkinfoInfoSlaveData `json:"info_slave_data,omitempty"`
	InfoData      LinkLinkinfoInfoData      `json:"info_data,omitempty"`
}

//LinkID type
type LinkID string

//Link definition
type Link struct {
	// Inteface index ID
	Ifindex int32 `json:"ifindex,omitempty"`
	// Interface name identifier
	Ifname LinkID `json:"ifname"`
	// Maximum Transfer Unit value
	Mtu int32 `json:"mtu,omitempty"`
	// In case the interface is part of a bond or bridge, specifies the bond/bridge interface it belongs to.
	Master   LinkID         `json:"master,omitempty"`
	Linkinfo LinkLinkinfo   `json:"linkinfo,omitempty"`
	LinkType string         `json:"link_type"`
	Address  string         `json:"address,omitempty"`
	AddrInfo []LinkAddrInfo `json:"addr_info,omitempty"`
}

func linkParse(link netlink.Link) Link {
	nclink := Link{}
	la := link.Attrs()
	nclink.Ifname = LinkID(la.Name)
	nclink.Mtu = int32(la.MTU)
	nclink.Linkinfo.InfoKind = link.Type()
	nclink.LinkType = la.EncapType
	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err == nil {
		nclink.AddrInfo = make([]LinkAddrInfo, len(addrs))
		for i, a := range addrs {
			nclink.AddrInfo[i].Local.Parse(a.IPNet.String())
		}
	} else {
		logger.Log.Warning(err)
	}
	switch nclink.Linkinfo.InfoKind {
	case "bond":
		{
			id := &nclink.Linkinfo.InfoData
			bond := link.(*netlink.Bond)
			id.Mode = bond.Mode.String()
			id.Miimon = int32(bond.Miimon)
			id.Updelay = int32(bond.UpDelay)
			id.Downdelay = int32(bond.DownDelay)
		}
	case "device":
	case "bridge":
	case "dummy":
	case "ppp":
	default:
		{
			logger.Log.Warning("Unknown Link Kind: " + nclink.Linkinfo.InfoKind)
		}
	}
	if la.Slave != nil {
		mkink, err := netlink.LinkByIndex(la.MasterIndex)
		if err == nil {
			nclink.Master = LinkID(mkink.Attrs().Name)
			nclink.Linkinfo.InfoSlaveKind = la.Slave.SlaveType()
			switch la.Slave.(type) {
			case *netlink.BondSlave:
				{
					bondslave := la.Slave.(*netlink.BondSlave)
					ids := &nclink.Linkinfo.InfoSlaveData
					ids.State = bondslave.State.String()
					ids.MiiStatus = bondslave.MiiStatus.String()
					ids.PermHwaddr = bondslave.PermHardwareAddr.String()
					ids.QueueId = bondslave.QueueId
					ids.LinkFailureCount = bondslave.LinkFailureCount
				}
			default:
				{
					logger.Log.Warning("Unsupported type of slave/master type interface")
				}
			}
		} else {
			logger.Log.Warning(err)
		}
	}
	return nclink
}

//LinksGet Returns the list of existing link layer devices on the machine
func LinksGet() ([]Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}
	nclinks := make([]Link, len(links))
	for i, l := range links {
		nclinks[i] = linkParse(l)
	}
	return nclinks, nil
}

func findActiveBackupBondActiveSlave(links []Link, bondIfname LinkID) (*Link, error) {
	var foundLink *Link = nil
	var secondFoundLink *Link = nil
	for _, link := range links {
		if link.Master == bondIfname && link.Linkinfo.InfoSlaveData.State == netlink.BondStateActive.String() {
			if foundLink == nil {
				foundLink = &link
			} else {
				secondFoundLink = &link
			}

		}
	}
	if foundLink == nil {
		return nil, NewActiveSlaveIfaceNotFoundForActiveBackupBondError(bondIfname)
	}
	if secondFoundLink != nil {
		return nil, NewMultipleActiveSlaveIfacesFoundForActiveBackupBondError(bondIfname)
	}
	return foundLink, nil
}

//LinksConfigure configures the whole set of links to manage in the correct sequential order
//for example some of the link properties require other links to be established already or
//to have the link down/up etc..
//This function tries to wipe out every type of conflicting in place configuration such as
//existing links whose ifname LinkID collides with the ones being created.
func LinksConfigure(links []Link) error {
	//Recreate all links
	for _, link := range links {
		LinkSetDown(link.Ifname)
		LinkDelete(link.Ifname)
		if err := LinkCreate(link); err != nil {
			return err
		}
	}

	//Set active-backup bond links active slaves (apparently you need to do this before setting the backups)
	for _, link := range links {
		if link.Master != "" {
			l, err := LinkGet(link.Master)
			if err != nil {
				return err
			}
			if l.Linkinfo.InfoKind == "bond" {
				if l.Linkinfo.InfoData.Mode == netlink.BOND_MODE_ACTIVE_BACKUP.String() {
					activeSlave, err := findActiveBackupBondActiveSlave(links, link.Master)
					if err != nil {
						return err
					}
					// Apparently first Link Set becomes master so do it first.
					LinkSetMaster(activeSlave.Ifname, link.Master)
				}
			}
		}
	}

	//Set all links cross properties (e.g. being slave of some master link interface)
	for _, link := range links {
		if link.Master != "" {
			l, err := LinkGet(link.Master)
			if err != nil {
				return err
			}
			if l.Linkinfo.InfoKind == "bond" {
				if l.Linkinfo.InfoData.Mode == netlink.BOND_MODE_ACTIVE_BACKUP.String() {
					if link.Linkinfo.InfoSlaveData.State == netlink.BondStateBackup.String() {
						LinkSetBondSlave(link.Ifname, link.Master)
					}
				} else {
					if link.Linkinfo.InfoSlaveData.State == netlink.BondStateBackup.String() {
						return NewBackupSlaveIfaceFoundForNonActiveBackupBondError(link.Ifname, link.Master)
					}
					LinkSetBondSlave(link.Ifname, link.Master)
				}
			}
		}
	}
	//Set all links up
	for _, link := range links {
		LinkSetUp(link.Ifname)
	}

	return nil
}

//LinkGet Returns the list of existing link layer devices on the machine
func LinkGet(ifname LinkID) (Link, error) {
	nclink := Link{}
	link, err := netlink.LinkByName(string(ifname))
	if err == nil {
		nclink = linkParse(link)
	}
	return nclink, err
}

// LinkCreate creates a link layer interface
// Link types (or kind):
// $> ip link help type
// ...
// TYPE := { vlan | veth | vcan | vxcan | dummy | ifb | macvlan | macvtap |
//	bridge | bond | team | ipoib | ip6tnl | ipip | sit | vxlan |
//	gre | gretap | erspan | ip6gre | ip6gretap | ip6erspan |
//	vti | nlmon | team_slave | bond_slave | ipvlan | geneve |
//	bridge_slave | vrf | macsec }
func LinkCreate(link Link) error {
	var err error = nil
	ifname := link.Ifname
	l, _ := netlink.LinkByName(string(ifname))
	if l != nil {
		return NewLinkExistsConflictError(ifname)
	}

	nllink, err := linkFormat(link)
	if err != nil {
		return err
	}
	return netlink.LinkAdd(nllink)
}

//LinkSetUp set a link up
func LinkSetUp(ifname LinkID) error {
	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	return netlink.LinkSetUp(link)
}

//LinkSetDown set a link up
func LinkSetDown(ifname LinkID) error {
	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	return netlink.LinkSetDown(link)
}

//LinkSetMaster specifies for a given interface (by ifname) the master to federate with (by masterIfname)
func LinkSetMaster(ifname LinkID, masterIfname LinkID) error {
	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	masterLink, _ := netlink.LinkByName(string(masterIfname))
	if masterLink == nil {
		return NewLinkNotFoundError(masterIfname)
	}
	return netlink.LinkSetMaster(link, masterLink)
}

//LinkSetBondSlave enslaves an interface to a master one
func LinkSetBondSlave(ifname LinkID, masterIfname LinkID) error {

	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	masterLink, _ := netlink.LinkByName(string(masterIfname))
	if masterLink == nil {
		return NewLinkNotFoundError(masterIfname)
	}
	if masterLink.Type() == "bond" {
		logger.Log.Debug("Setting bond slave")
		return netlink.LinkSetBondSlave(link, masterLink.(*netlink.Bond))
	}
	return NonBondMasterLinkTypeError(masterIfname)
}

//LinkDelete deletes a link layer interface
func LinkDelete(ifname LinkID) error {
	l, _ := netlink.LinkByName(string(ifname))
	if l == nil {
		return NewLinkNotFoundError(ifname)
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = string(ifname)
	// Fool it with a dummy.. it should use ifname and ignore the rest
	nllink := &netlink.Dummy{
		LinkAttrs: attrs,
	}
	return netlink.LinkDel(nllink)
}

func linkFormat(link Link) (netlink.Link, error) {
	ifname := link.Ifname
	kind := link.Linkinfo.InfoKind
	attrs := netlink.NewLinkAttrs()
	attrs.Name = string(ifname)
	var err error
	var nllink netlink.Link = nil
	switch kind {
	case "dummy":
		{
			nllink = &netlink.Dummy{
				LinkAttrs: attrs,
			}
		}
	case "bond":
		{

			nllink = netlink.NewLinkBond(attrs)
			nlbondlink := nllink.(*netlink.Bond)
			nlbondlink.Mode = netlink.StringToBondMode(link.Linkinfo.InfoData.Mode)
		}
	case "bridge":
		{
			nllink = &netlink.Bridge{
				LinkAttrs: attrs,
			}
		}
	case "vlan":
		{
			nllink = &netlink.Vlan{
				LinkAttrs: attrs,
			}
		}
	case "veth":
		{
			nllink = &netlink.Veth{
				LinkAttrs: attrs,
			}
		}
	case "ipvlan":
		{
			nllink = &netlink.IPVlan{
				LinkAttrs: attrs,
			}
		}
	case "macvlan":
		{
			nllink = &netlink.Macvlan{
				LinkAttrs: attrs,
			}
		}
	case "tuntap":
		{
			nllink = &netlink.Tuntap{
				LinkAttrs: attrs,
			}
		}
	default:
		err = NewUnknownLinkKindError(kind)
	}
	return nllink, err
}
