package nc

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/riccardomanfrin/netlink"
	"gitlab.lan.athonet.com/core/netconfd/logger"
	"golang.org/x/sys/unix"
)

// LinkAddrInfo struct for LinkAddrInfo
type LinkAddrInfo struct {
	Local CIDRAddr `json:"local,omitempty"`
	//Prefixlen int32    `json:"prefixlen,omitempty"`
	//Broadcast CIDRAddr `json:"broadcast,omitempty"`
	Address *net.IP `json:"local,omitempty"`
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
	// Specify the delay, in milliseconds, between each peer notification (gratuitous ARP and unsolicited IPv6 Neighbor Advertisement) when they are issued after a failover event. This delay should be a multiple of the link monitor interval (arp_interval or miimon, whichever is active). The default value is 0 which means to match the value of the link monitor interval.
	PeerNotifyDelay int32 `json:"peer_notify_delay,omitempty"`
	// Specifies whether or not miimon should use MII or ETHTOOL ioctls vs. netif_carrier_ok() to determine the link status. The MII or ETHTOOL ioctls are less efficient and utilize a deprecated calling sequence within the kernel.  The netif_carrier_ok() relies on the device driver to maintain its state with netif_carrier_on/off; at this writing, most, but not all, device drivers support this facility. If bonding insists that the link is up when it should not be, it may be that your network device driver does not support netif_carrier_on/off.  The default state for netif_carrier is \"carrier on,\" so if a driver does not support netif_carrier, it will appear as if the link is always up.  In this case, setting use_carrier to 0 will cause bonding to revert to the MII / ETHTOOL ioctl method to determine the link state. A value of 1 enables the use of netif_carrier_ok(), a value of 0 will use the deprecated MII / ETHTOOL ioctls.  The default value is 1.
	UseCarrier int32 `json:"use_carrier,omitempty"`
	// Specifies the ARP link monitoring frequency in milliseconds. The ARP monitor works by periodically checking the slave devices to determine whether they have sent or received traffic recently (the precise criteria depends upon the bonding mode, and the state of the slave).  Regular traffic is generated via ARP probes issued for the addresses specified by the arp_ip_target option. This behavior can be modified by the arp_validate option, below. If ARP monitoring is used in an etherchannel compatible mode (modes 0 and 2), the switch should be configured in a mode that evenly distributes packets across all links. If the switch is configured to distribute the packets in an XOR fashion, all replies from the ARP targets will be received on the same link which could cause the other team members to fail.  ARP monitoring should not be used in conjunction with miimon.  A value of 0 disables ARP monitoring.  The default value is 0.
	ArpInterval int32 `json:"arp_interval,omitempty"`
	// Specifies whether or not ARP probes and replies should be validated in any mode that supports arp monitoring, or whether non-ARP traffic should be filtered (disregarded) for link monitoring purposes. Possible values are: * `none` - or 0 No validation or filtering is performed. * `active` - or 1 Validation is performed only for the active slave. * `backup` - or 2 Validation is performed only for backup slaves. * `all` - or 3 Validation is performed for all slaves. * `filter` - or 4 Filtering is applied to all slaves. No validation is performed. * `filter_active` - or 5 Filtering is applied to all slaves, validation is performed only for the active slave. * `filter_backup` - or 6 Filtering is applied to all slaves, validation is performed only for backup slaves.
	ArpValidate string `json:"arp_validate,omitempty"`
	// Specifies the quantity of arp_ip_targets that must be reachable in order for the ARP monitor to consider a slave as being up. This option affects only active-backup mode for slaves with arp_validation enabled. Possible values are: * `any` - or 0   consider the slave up only when any of the arp_ip_targets   is reachable  * `all` - or 1   consider the slave up only when all of the arp_ip_targets   are reachable
	ArpAllTargets string `json:"arp_all_targets,omitempty"`
	// Specifies the reselection policy for the primary slave.  This affects how the primary slave is chosen to become the active slave when failure of the active slave or recovery of the primary slave occurs.  This option is designed to prevent flip-flopping between the primary slave and other slaves.  Possible values are:    * `always` - or 0 (default)     The primary slave becomes the active slave whenever it     comes back up.   * `better` - or 1     The primary slave becomes the active slave when it comes     back up, if the speed and duplex of the primary slave is     better than the speed and duplex of the current active     slave.   * `failure` - or 2     The primary slave becomes the active slave only if the     current active slave fails and the primary slave is up.  The primary_reselect setting is ignored in two cases:    * If no slaves are active, the first slave to recover is     made the active slave.    * When initially enslaved, the primary slave is always made     the active slave.  Changing the primary_reselect policy via sysfs will cause an immediate selection of the best active slave according to the new policy.  This may or may not result in a change of the active slave, depending upon the circumstances. This option was added for bonding version 3.6.0.
	PrimaryReselect string `json:"primary_reselect,omitempty"`
	// Specifies whether active-backup mode should set all slaves to the same MAC address at enslavement (the traditional behavior), or, when enabled, perform special handling of the bond's MAC address in accordance with the selected policy. The default policy is none, unless the first slave cannot change its MAC address, in which case the active policy is selected by default. This option may be modified via sysfs only when no slaves are present in the bond. This option was added in bonding version 3.2.0.  The \"follow\" policy was added in bonding version 3.3.0. Possible values are:   * `none` - or 0   This setting disables fail_over_mac, and causes   bonding to set all slaves of an active-backup bond to   the same MAC address at enslavement time.  This is the   default.   * `active` - or 1   The \"active\" fail_over_mac policy indicates that the   MAC address of the bond should always be the MAC   address of the currently active slave.  The MAC   address of the slaves is not changed; instead, the MAC   address of the bond changes during a failover.   This policy is useful for devices that cannot ever   alter their MAC address, or for devices that refuse   incoming broadcasts with their own source MAC (which   interferes with the ARP monitor).   The down side of this policy is that every device on   the network must be updated via gratuitous ARP,   vs. just updating a switch or set of switches (which   often takes place for any traffic, not just ARP   traffic, if the switch snoops incoming traffic to   update its tables) for the traditional method.  If the   gratuitous ARP is lost, communication may be   disrupted.   When this policy is used in conjunction with the mii   monitor, devices which assert link up prior to being   able to actually transmit and receive are particularly   susceptible to loss of the gratuitous ARP, and an   appropriate updelay setting may be required.   * `follow` - or 2   The \"follow\" fail_over_mac policy causes the MAC   address of the bond to be selected normally (normally   the MAC address of the first slave added to the bond).   However, the second and subsequent slaves are not set   to this MAC address while they are in a backup role; a   slave is programmed with the bond's MAC address at   failover time (and the formerly active slave receives   the newly active slave's MAC address).   This policy is useful for multiport devices that   either become confused or incur a performance penalty   when multiple ports are programmed with the same MAC   address.
	FailOverMac string `json:"fail_over_mac,omitempty"`
	// Hash policy to route packets on different bond interfaces.  Supported Modes:   * `layer2` - Hash is made on L2 metadata   * `layer2+3` - Hash is made on L2 and L3 metadata   * `layer3+4` - Hash is made on L3 and L4 metadata
	XmitHashPolicy string `json:"xmit_hash_policy,omitempty"`
	// Specifies the number of IGMP membership reports to be issued after a failover event. One membership report is issued immediately after the failover, subsequent packets are sent in each 200ms interval.  The valid range is 0 - 255; the default value is 1. A value of 0 prevents the IGMP membership report from being issued in response to the failover event.  This option is useful for bonding modes balance-rr (0), active-backup (1), balance-tlb (5) and balance-alb (6), in which a failover can switch the IGMP traffic from one slave to another.  Therefore a fresh IGMP report must be issued to cause the switch to forward the incoming IGMP traffic over the newly selected slave.  This option was added for bonding version 3.7.0.
	ResendIgmp int32 `json:"resend_igmp,omitempty"`
	// Specifies that duplicate frames (received on inactive ports) should be dropped (0) or delivered (1).  Normally, bonding will drop duplicate frames (received on inactive ports), which is desirable for most users. But there are some times it is nice to allow duplicate frames to be delivered.  The default value is 0 (drop duplicate frames received on inactive ports).
	AllSlavesActive int32 `json:"all_slaves_active,omitempty"`
	// Specifies the minimum number of links that must be active before asserting carrier. It is similar to the Cisco EtherChannel min-links feature. This allows setting the minimum number of member ports that must be up (link-up state) before marking the bond device as up (carrier on). This is useful for situations where higher level services such as clustering want to ensure a minimum number of low bandwidth links are active before switchover. This option only affect 802.3ad mode.  The default value is 0. This will cause carrier to be asserted (for 802.3ad mode) whenever there is an active aggregator, regardless of the number of available links in that aggregator. Note that, because an aggregator cannot be active without at least one available link, setting this option to 0 or to 1 has the exact same effect.
	MinLinks int32 `json:"min_links,omitempty"`
	// Specifies the number of seconds between instances where the bonding driver sends learning packets to each slaves peer switch.  The valid range is 1 - 0x7fffffff; the default value is 1. This Option has effect only in balance-tlb and balance-alb modes.
	LpInterval int32 `json:"lp_interval,omitempty"`
	// Specify the number of packets to transmit through a slave before moving to the next one. When set to 0 then a slave is chosen at random.  The valid range is 0 - 65535; the default value is 1. This option has effect only in balance-rr mode.
	PacketsPerSlave int32 `json:"packets_per_slave,omitempty"`
	// Rate at which LACP control packets are sent to an LACP-supported interface Supported Modes:   * `slow` - LACP Slow Rate (less bandwidth)   * `fast` - LACP Fast Rate (faster fault detection)
	AdLacpRate string `json:"ad_lacp_rate,omitempty"`
	// Specifies the 802.3ad aggregation selection logic to use.  The possible values and their effects are:   * `stable` - or 0     The active aggregator is chosen by largest aggregate     bandwidth.     Reselection of the active aggregator occurs only when all     slaves of the active aggregator are down or the active     aggregator has no slaves.     This is the default value.   * `bandwidth` or 1     The active aggregator is chosen by largest aggregate     bandwidth.  Reselection occurs if:     - A slave is added to or removed from the bond     - Any slave's link state changes     - Any slave's 802.3ad association state changes     - The bond's administrative state changes to up   * `count` - or 2     The active aggregator is chosen by the largest number of     ports (slaves).  Reselection occurs as described under the     \"bandwidth\" setting, above.      The bandwidth and count selection policies permit failover of 802.3ad aggregations when partial failure of the active aggregator occurs.  This keeps the aggregator with the highest availability (either in bandwidth or in number of ports) active at all times. This option was added in bonding version 3.4.0.
	AdSelect string `json:"ad_select,omitempty"`
	// Specifies if dynamic shuffling of flows is enabled in tlb mode. The value has no effect on any other modes.  The default behavior of tlb mode is to shuffle active flows across slaves based on the load in that interval. This gives nice lb characteristics but can cause packet reordering. If re-ordering is a concern use this variable to disable flow shuffling and rely on load balancing provided solely by the hash distribution. xmit-hash-policy can be used to select the appropriate hashing for the setup.  The sysfs entry can be used to change the setting per bond device and the initial value is derived from the module parameter. The sysfs entry is allowed to be changed only if the bond device is down.  The default value is \"1\" that enables flow shuffling while value \"0\" disables it. This option was added in bonding driver 3.7.1
	TlbDynamicLb int32 `json:"tlb_dynamic_lb,omitempty"`
	// VLAN protocols. Supported protocols:   * `802.1Q` - 802.1Q protocol
	Protocol string `json:"protocol,omitempty"`
	// VLAN TAG ID
	Id int32 `json:"id,omitempty"`
	// Flags of the virtual device
	Flags []string `json:"flags,omitempty"`
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

//IsValid checks whether a link is valid
func (l LinkID) IsValid() bool {
	return string(l) != "" /* && len(string(l)) <= 15 */
}

// LinkFlag the model 'LinkFlag'
type LinkFlag string

//LinkFlags is a slice of flags
type LinkFlags []LinkFlag

// List of link_flag
const (
	BROADCAST LinkFlag = "broadcast"
	MULTICAST LinkFlag = "multicast"
	LOOPBACK  LinkFlag = "loopback"
	UP        LinkFlag = "up"
)

// Link definition
// For Bond parameters information please refer to
// https://www.kernel.org/doc/Documentation/networking/bonding.txt
// https://www.kernel.org/doc/Documentation/networking/operstates.txt
type Link struct {
	// Inteface index ID
	Ifindex int32 `json:"ifindex,omitempty"`
	// Interface name identifier
	Ifname LinkID `json:"ifname"`
	// Specify what is the physical device the virtual device is linked to. Applies to vlan type virtual devices
	Link LinkID `json:"link,omitempty"`
	// Maximum Transfer Unit value
	Mtu int32 `json:"mtu,omitempty"`
	// In case the interface is part of a bond or bridge, specifies the bond/bridge interface it belongs to.
	Master   LinkID         `json:"master,omitempty"`
	Linkinfo LinkLinkinfo   `json:"linkinfo,omitempty"`
	LinkType string         `json:"link_type"`
	Address  string         `json:"address,omitempty"`
	AddrInfo []LinkAddrInfo `json:"addr_info,omitempty"`
	Flags    LinkFlags      `json:"flags,omitempty"`
	// Readonly state of the interface.  Provides information on the state being for example UP of an interface.  It is ignored when applying the config
	Operstate string `json:"operstate,omitempty"`
}

//Print implements route print
func (r *Link) Print() string {
	return fmt.Sprintf("%v", r)
}

func linkParse(link netlink.Link) Link {
	nclink := Link{}
	la := link.Attrs()
	nclink.Ifname = LinkID(la.Name)
	nclink.Ifindex = int32(la.Index)
	nclink.Mtu = int32(la.MTU)
	nclink.Linkinfo.InfoKind = link.Type()
	nclink.LinkType = la.EncapType
	nclink.Operstate = la.OperState.String()
	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	nclink.Address = la.HardwareAddr.String()
	if err == nil {
		nclink.AddrInfo = make([]LinkAddrInfo, len(addrs))
		for i, a := range addrs {
			//nclink.AddrInfo[i].Local.Parse(a.IPNet.String())
			//ones, bits := a.IPNet.Mask.Size()
			nclink.AddrInfo[i].Local.SetIP(a.IPNet.IP)

			var ones int
			if a.Peer != nil {
				nclink.AddrInfo[i].Address = &a.Peer.IP
				ones, _ = a.Peer.Mask.Size()
			} else {
				ones, _ = a.IPNet.Mask.Size()
			}

			nclink.AddrInfo[i].Local.SetPrefixLen(ones)

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
			id.UseCarrier = int32(bond.UseCarrier)
			id.ArpInterval = int32(bond.ArpInterval)
			id.ArpValidate = bond.ArpValidate.String()
			id.LpInterval = int32(bond.LpInterval)
			id.ArpAllTargets = bond.ArpAllTargets.String()
			id.PacketsPerSlave = int32(bond.PacketsPerSlave)
			id.FailOverMac = bond.FailOverMac.String()
			id.XmitHashPolicy = bond.XmitHashPolicy.String()
			id.ResendIgmp = int32(bond.ResendIgmp)
			id.MinLinks = int32(bond.MinLinks)
			id.ArpInterval = int32(bond.ArpInterval)
			id.PrimaryReselect = bond.PrimaryReselect.String()
			id.TlbDynamicLb = int32(bond.TlbDynamicLb)
			id.AdSelect = bond.AdSelect.String()
			id.AdLacpRate = bond.LacpRate.String()
			id.AllSlavesActive = int32(bond.AllSlavesActive)
			id.UseCarrier = int32(bond.UseCarrier)
		}
	case "device":
	case "bridge":
	case "dummy":
	case "vlan":
		{
			id := &nclink.Linkinfo.InfoData
			vlan := link.(*netlink.Vlan)
			mkink, err := netlink.LinkByIndex(vlan.Attrs().ParentIndex)
			if err == nil {
				logger.Log.Warning("Virtual vlan link device parent link not found by index", vlan.Attrs().ParentIndex)
			}

			nclink.Link = LinkID(mkink.Attrs().Name)
			id.Protocol = vlan.VlanProtocol.String()
			id.Id = int32(vlan.VlanId)
		}
	case "tuntap":
		{
			nctuntap, ok := link.(*netlink.Tuntap)
			if !ok {
				logger.Log.Warning("Unmatched tuntap info kind vs link typecast")
			}
			nclink.Linkinfo.InfoKind = nctuntap.Mode.String()
		}
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

	nclink.Flags = linkFlagsParse(link.Attrs().Flags)

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

//SetFlag return true if the searched flag is found
func (flags LinkFlags) SetFlag(flag LinkFlag) LinkFlags {
	if !flags.HaveFlag(flag) {
		return append(flags, flag)
	}
	return flags
}

//ClearFlag return true if the searched flag is found
func (flags LinkFlags) ClearFlag(flag LinkFlag) LinkFlags {
	var outFlags []LinkFlag
	for _, f := range flags {
		if f != flag {
			outFlags = append(outFlags, f)
		}
	}
	return outFlags
}

//HaveFlag return true if the searched flag is found
func (flags LinkFlags) HaveFlag(flag LinkFlag) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}

func isLinkRemovable(link netlink.Link) (bool, string) {
	if link.Attrs().Index == 1 {
		return false, fmt.Sprintf("Skipping loopback iface %v removal/creation as per %v", link.Attrs().Name, loopbackUniquenessRef)
	}
	if _, ok := link.(*netlink.Device); ok {
		return false, fmt.Sprintf("Skipping physical iface %v removal/creation as per %v", link.Attrs().Name, ethernetNoRemovalRef)
	}
	return true, ""
}

//LinksDelete remove all non physical and non loopback links
// Refs:
// Loopback uniqueness:
// https://elixir.bootlin.com/linux/latest/source/drivers/net/loopback.c#L195
// Phy interfaces can't be removed if not for modprobe -r or Hot-Plug events
// https://github.com/ryoon/e1000e-linux/blob/master/src/netdev.c#L7968
func LinksDelete() error {
	links, err := netlink.LinkList()
	if err != nil {
		return mapNetlinkError(err, nil)
	}
	for _, link := range links {
		if isUnmanaged(UnmanagedID(link.Attrs().Name), LINKTYPE) {
			logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v removal", link.Attrs().Name))
			continue
		}
		removable, why := isLinkRemovable(link)
		if !removable {
			logger.Log.Debug(why)
			continue
		}
		logger.Log.Warning("Removing link " + link.Attrs().Name)
		err := netlink.LinkDel(link)
		if err != nil {
			l := linkParse(link)
			return mapNetlinkError(err, &l)
		}
	}
	return nil
}

//LinksConfigure configures the whole set of links to manage in the correct sequential order
//for example some of the link properties require other links to be established already or
//to have the link down/up etc..
//This function tries to wipe out every type of conflicting in place configuration such as
//existing links whose ifname LinkID collides with the ones being created.
func LinksConfigure(links []Link) error {
	//Recreate all links
	for _, link := range links {
		if isUnmanaged(UnmanagedID(link.Ifname), LINKTYPE) {
			logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v configuration", link.Ifname))
			continue
		}
		l, _ := netlink.LinkByName(string(link.Ifname))
		if l != nil {
			logger.Log.Debug("Setting link %v down", link.Ifname)
			LinkSetDown(link.Ifname)
			removable, why := isLinkRemovable(l)
			if !removable {
				//Just set addresses
				logger.Log.Debug("Setting link %v addresses", link.Ifname)
				if err := LinkSetAddresses(link); err != nil {
					return err
				}

				logger.Log.Debug(why)
				continue
			}
			logger.Log.Debug("Deleting link %v", link.Ifname)
			if err := LinkDelete(link.Ifname); err != nil {
				logger.Log.Warning("Link Delete Error:", err)
			}
		}
		/* You cannot enslave a link if it is UP */
		logger.Log.Debug("Deleting link %v", link.Ifname)
		if err := LinkCreateDown(link); err != nil {
			return err
		}
	}

	//Set active-backup bond links active slaves (apparently you need to do this before setting the backups)
	for _, link := range links {
		logger.Log.Debug("Setting slave active/backup properties for link %v", link.Ifname)
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
					if err = LinkSetMaster(activeSlave.Ifname, link.Master); err != nil {
						logger.Log.Warning("Link Set Master Error:", err)
					}
				}
			}
		}
	}

	//Set all links cross properties (e.g. being slave of some master link interface)
	for _, link := range links {
		logger.Log.Debug("Setting master/slave cross properties for link %v", link.Ifname)
		if link.Master != "" {
			l, err := LinkGet(link.Master)
			if err != nil {
				return err
			}
			if l.Linkinfo.InfoKind == "bond" {
				if l.Linkinfo.InfoData.Mode == netlink.BOND_MODE_ACTIVE_BACKUP.String() {
					if link.Linkinfo.InfoSlaveData.State == netlink.BondStateBackup.String() {
						if err = LinkSetBondSlave(link.Ifname, link.Master); err != nil {
							logger.Log.Warning("Link Set Bond Slave Error:", err)
						}
					}
				} else {
					if link.Linkinfo.InfoSlaveData.State == netlink.BondStateBackup.String() {
						return NewBackupSlaveIfaceFoundForNonActiveBackupBondError(link.Ifname, link.Master)
					}
					if err = LinkSetBondSlave(link.Ifname, link.Master); err != nil {
						logger.Log.Warning("Link Set Bond Slave Error:", err)
					}
				}
			}
		}
	}
	//Because slave interface cannot set up BEFORE enslaving them
	//we just set all links with LinkCreateDown and then set them up at the end with no
	//distinction
	for _, link := range links {
		if link.Flags.HaveFlag(LinkFlag(net.FlagUp.String())) {
			logger.Log.Debug("Setting link %v up", link.Ifname)
			if err := LinkSetUp(link.Ifname); err != nil {
				logger.Log.Warning("Link Set Up Error:", err)
			}
		}
	}

	for _, link := range links {
		if link.Mtu > 0 {
			err := LinkSetMTU(link.Ifname, int(link.Mtu))
			if err != nil {
				logger.Log.Warning("Link Set MTU Error:", err)
			}
		}
	}

	return nil
}

//LinkGet Returns the list of existing link layer devices on the machine
func LinkGet(ifname LinkID) (Link, error) {
	nclink := Link{}
	link, err := netlink.LinkByName(string(ifname))
	if err == nil {
		nclink = linkParse(link)
	} else {
		//Translate to our err codes
		err = NewLinkNotFoundError(ifname)
	}
	return nclink, err
}

//LinkCreateDown Creates a link interface but does not bring it up
func LinkCreateDown(link Link) error {
	if isUnmanaged(UnmanagedID(link.Ifname), LINKTYPE) {
		logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v configuration", link.Ifname))
		return nil
	}
	isFlagUp := link.Flags.HaveFlag(LinkFlag(net.FlagUp.String()))
	link.Flags = link.Flags.ClearFlag(LinkFlag(net.FlagUp.String()))
	err := LinkCreate(link)
	if isFlagUp {
		link.Flags = link.Flags.SetFlag(LinkFlag(net.FlagUp.String()))
	}
	return err
}

//LinkSetAddresses assignes all addresses of a link (erase and recreate them)
func LinkSetAddresses(link Link) error {
	ifname := link.Ifname

	l, _ := netlink.LinkByName(string(ifname))
	if l != nil {
		addrlist, err := netlink.AddrList(l, netlink.FAMILY_ALL)
		if err != nil {
			return NewGenericErrorWithReason(fmt.Sprintf("Failed to get address list: error %v", err))
		}

		for _, addr := range addrlist {
			if err = netlink.AddrDel(l, &addr); err != nil {
				return mapNetlinkError(err, nil)
			}
		}
	}

	for _, a := range link.AddrInfo {
		logger.Log.Debug("Adding addr", a.Local.String(), "to iface", ifname)
		if err := linkAddrAdd(ifname, a.Local, a.Address); err != nil {
			return err
		}
	}
	return nil
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

	if isUnmanaged(UnmanagedID(link.Ifname), LINKTYPE) {
		logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v configuration", link.Ifname))
		return NewUnmanagedLinkCannotBeModifiedError(ifname)
	}

	logger.Log.Debug("Creating link", ifname)
	l, _ := netlink.LinkByName(string(ifname))

	removable := true
	if l != nil {
		removable, _ = isLinkRemovable(l)
	}
	if removable && (l != nil) {
		return NewLinkExistsConflictError(ifname)
	}
	if (l == nil) && link.Linkinfo.InfoKind == "device" {
		return NewLinkDeviceDoesNotExistError(ifname)
	}

	if removable {
		nllink, err := linkFormat(link)
		if err != nil {
			return err
		}
		if err = netlink.LinkAdd(nllink); err != nil {
			return mapNetlinkError(err, &link)
		}
	}

	err = LinkSetAddresses(link)

	return mapNetlinkError(err, nil)
}

func linkAddrAdd(ifname LinkID, addr CIDRAddr, peer *net.IP) error {
	l, _ := netlink.LinkByName(string(ifname))
	if l == nil {
		return NewLinkNotFoundError(ifname)
	}
	/* resolves ambiguity on prefixLen. it comes from network address.
	 * peer is only /32 remote peer endpoint */
	addrNet := addr.ToIPNet()
	var peerNet *net.IPNet
	if peer != nil {
		_peerNet := net.IPNet{
			IP:   *peer,
			Mask: addr.ToIPNet().Mask,
		}
		peerNet = &_peerNet
	}

	nladdr := netlink.Addr{
		IPNet: &addrNet,
		Label: string(ifname),
		Scope: unix.RT_SCOPE_UNIVERSE,
		Peer:  peerNet,
		Flags: unix.IFA_F_PERMANENT}
	if err := netlink.AddrAdd(l, &nladdr); err != nil {
		link := linkParse(l)
		return mapNetlinkError(err, &link)
	}
	return nil
}

//LinkSetUp set a link up
func LinkSetUp(ifname LinkID) error {
	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	return netlink.LinkSetUp(link)
}

//LinkSetMTU set a link MTU
func LinkSetMTU(ifname LinkID, mtu int) error {
	link, _ := netlink.LinkByName(string(ifname))
	if link == nil {
		return NewLinkNotFoundError(ifname)
	}
	return netlink.LinkSetMTU(link, mtu)
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
		return netlink.LinkSetBondSlave(link, masterLink.(*netlink.Bond))
	}
	return NewNonBondMasterLinkTypeError(masterIfname)
}

//LinkDelete deletes a link layer interface
func LinkDelete(ifname LinkID) error {
	l, err := LinkGet(ifname)
	if err != nil {
		return err
	}
	if isUnmanaged(UnmanagedID(ifname), LINKTYPE) {
		return NewUnmanagedLinkCannotBeModifiedError(ifname)
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = string(ifname)
	// Fool it with a dummy.. it should use ifname and ignore the rest
	nllink := &netlink.Dummy{
		LinkAttrs: attrs,
	}
	err = netlink.LinkDel(nllink)
	if err != nil {
		return mapNetlinkError(err, &l)
	}
	return nil
}

func linkFlagsParse(flags net.Flags) []LinkFlag {
	var linkFlags []LinkFlag

	if (flags & net.FlagUp) > 0 {
		linkFlags = append(linkFlags, LinkFlag(net.FlagUp.String()))
	}
	if (flags & net.FlagBroadcast) > 0 {
		linkFlags = append(linkFlags, LinkFlag(net.FlagBroadcast.String()))
	}
	if (flags & net.FlagMulticast) > 0 {
		linkFlags = append(linkFlags, LinkFlag(net.FlagMulticast.String()))
	}
	if (flags & net.FlagLoopback) > 0 {
		linkFlags = append(linkFlags, LinkFlag(net.FlagLoopback.String()))
	}
	if (flags & net.FlagPointToPoint) > 0 {
		linkFlags = append(linkFlags, LinkFlag(net.FlagPointToPoint.String()))
	}
	return linkFlags

}

func linkFlagsFormat(link Link) (net.Flags, error) {
	var flags net.Flags
	for _, f := range link.Flags {
		switch string(f) {
		case net.FlagUp.String():
			{
				flags += net.FlagUp
			}
		case net.FlagBroadcast.String():
			{
				flags += net.FlagBroadcast
			}
		case net.FlagMulticast.String():
			{
				flags += net.FlagMulticast
			}
		case net.FlagLoopback.String():
			{
				flags += net.FlagLoopback
			}
		case net.FlagPointToPoint.String():
			{
				flags += net.FlagPointToPoint
			}
		default:
			return flags, NewLinkUnknownFlagTypeError(f)
		}

	}
	return flags, nil
}

func linkFormat(link Link) (netlink.Link, error) {
	ifname := link.Ifname
	//linkType := link.LinkType
	kind := link.Linkinfo.InfoKind
	attrs := netlink.NewLinkAttrs()
	attrs.Name = string(ifname)
	attrs.Index = int(link.Ifindex)
	var err error
	var nllink netlink.Link = nil
	if kind != "" {
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
				nlbondlink.Miimon = int(link.Linkinfo.InfoData.Miimon)
				nlbondlink.DownDelay = int(link.Linkinfo.InfoData.Downdelay)
				nlbondlink.UpDelay = int(link.Linkinfo.InfoData.Updelay)

				id := &link.Linkinfo.InfoData

				if id.UseCarrier != -1 {
					nlbondlink.UseCarrier = int(id.UseCarrier)
				}
				if id.ArpInterval != -1 {
					nlbondlink.ArpInterval = int(id.ArpInterval)
				}
				if id.LpInterval != -1 {
					nlbondlink.LpInterval = int(id.LpInterval)
				}
				if id.PacketsPerSlave != -1 {
					nlbondlink.PacketsPerSlave = int(id.PacketsPerSlave)
				}
				if id.ResendIgmp != -1 {
					nlbondlink.ResendIgmp = int(id.ResendIgmp)
				}
				if id.MinLinks != -1 {
					nlbondlink.MinLinks = int(id.MinLinks)
				}
				if id.ArpInterval != -1 {
					nlbondlink.ArpInterval = int(id.ArpInterval)
				}
				if id.TlbDynamicLb != -1 {
					nlbondlink.TlbDynamicLb = int(id.TlbDynamicLb)
				}
				if id.AllSlavesActive != -1 {
					nlbondlink.AllSlavesActive = int(id.AllSlavesActive)
				}
				if id.UseCarrier != -1 {
					nlbondlink.UseCarrier = int(id.UseCarrier)
				}
				if id.XmitHashPolicy != "" {
					nlbondlink.XmitHashPolicy = netlink.StringToBondXmitHashPolicy(id.XmitHashPolicy)
				}
				if id.AdLacpRate != "" {
					nlbondlink.LacpRate = netlink.StringToBondLacpRate(id.AdLacpRate)
				}
				if id.ArpValidate != "" {
					nlbondlink.ArpValidate = netlink.StringToBondArpValidateMap[id.ArpValidate]
				}
				if id.ArpAllTargets != "" {
					nlbondlink.ArpAllTargets = netlink.StringToBondArpAllTargetsMap[id.ArpAllTargets]
				}
				if id.FailOverMac != "" {
					nlbondlink.FailOverMac = netlink.StringToBondFailOverMacMap[id.FailOverMac]
				}
				if id.XmitHashPolicy != "" {
					nlbondlink.XmitHashPolicy = netlink.StringToBondXmitHashPolicyMap[id.XmitHashPolicy]
				}
				if id.PrimaryReselect != "" {
					nlbondlink.PrimaryReselect = netlink.StringToBondPrimaryReselectMap[id.PrimaryReselect]
				}
				if id.AdSelect != "" {
					nlbondlink.AdSelect = netlink.StringToBondAdSelectMap[id.AdSelect]
				}
				if id.AdLacpRate != "" {
					nlbondlink.LacpRate = netlink.StringToBondLacpRateMap[id.AdLacpRate]
				}

			}
		case "bridge":
			{
				nllink = &netlink.Bridge{
					LinkAttrs: attrs,
				}
			}
		case "vlan":
			{
				id := &link.Linkinfo.InfoData
				parentLink, err := netlink.LinkByName(string(link.Link))
				if err != nil {
					return nllink, NewParentLinkNotFoundForVlan(ifname, link.Link)
				}
				attrs.ParentIndex = parentLink.Attrs().Index
				nllink = &netlink.Vlan{
					LinkAttrs:    attrs,
					VlanId:       int(id.Id),
					VlanProtocol: netlink.StringToVlanProtocol(id.Protocol),
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
		case "tun", "tap":
			{
				nllink = &netlink.Tuntap{
					LinkAttrs: attrs,
					Mode:      netlink.StringToTuntapModeMap[kind],
				}
			}
		case "device":
			{
				nllink = &netlink.Device{
					LinkAttrs: attrs,
				}
			}
		default:
			return nil, NewUnknownLinkKindError(kind)
		}
	} else {
		logger.Log.Warning("Unspecified link kind: assuming device link")
		nllink = &netlink.Device{
			LinkAttrs: attrs,
		}
	}
	netFlags, err := linkFlagsFormat(link)
	if err != nil {
		return nil, err
	}
	if link.Mtu > 0 {
		nllink.Attrs().MTU = int(link.Mtu)
	}
	nllink.Attrs().Flags = netFlags
	return nllink, err
}

//LinksVMReorder renames link devices to reflect hypervisor order on vmware or
//at least to be consistent over hypervisor changes.
//Beware that after renaming the interfaces are turned off.
func LinksVMReorder() error {
	err := vmwareLinksReorder()
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	logger.Log.Info("Non vmware/KVM deployed Virtual Appliance: reordering skipped")
	// Consider as a fallback plan to at least preserve order over changes,
	// although this would initially be unpredictable
	// Some literature:
	//
	// https://en.wikipedia.org/wiki/Consistent_Network_Device_Naming
	// https://wiki.debian.org/NetworkInterfaceNames#predictable
	// https://www.debian.org/releases/buster/amd64/release-notes/ch-information.en.html#migrate-interface-names
	// https://libvirt.org/pci-hotplug.html
	// https://www.freedesktop.org/wiki/Software/systemd/PredictableNetworkInterfaceNames/
	// https://github.com/systemd/systemd/blob/main/src/udev/udev-builtin-net_id.c#L20
	//
	// Aparently for a device named enpXsY
	// X is the PCI "path" (bus?)
	// Y is the PCI "slot"
	// X and Y assigned to a NIC do not change over time/reboots/network modification,
	// but at the same time they are in general not predictable (cannot be guessed a priori).
	//
	// As fallback plan, in case no matching (label based) can be used (e.g. Virtualbox),
	// a more generalized approach could be to start all NICs in DHCP (one will get the address),
	// so that reachability is granted, than delegate the ordering to the user via first
	// configuration:
	//
	//"network": {
	//	"devmap" : [
	//		{
	//			"mac" : "00:11:22:33:44:55",
	//			"ifname": "eth0"
	//		}
	//	]
	//}

	return nil
}

//LinkRename Rename a NIC Link Ifname
func LinkRename(currNICIface LinkID, remappedNICIface LinkID) error {
	//ip link set $link name $name

	l, err := LinkGet(currNICIface)
	if err != nil {
		logger.Log.Warning("Could not find iface %v by Ifname", currNICIface)
		return err
	}
	nclink, err := linkFormat(l)
	if err != nil {
		return err
	}
	err = netlink.LinkSetName(nclink, string(remappedNICIface))
	if err != nil {
		logger.Log.Warning(fmt.Sprintf("Failed to rename Link %+v to Link %v", nclink, remappedNICIface))
	}
	return err
}

type nICInfo struct {
	MACAddr string
	Ifname  LinkID
}

func linksRename(mis map[LinkID]nICInfo) error {
	for _, mi := range mis {
		if err := LinkSetDown(mi.Ifname); err != nil {
			return err
		}
		if err := LinkRename(mi.Ifname, LinkID("old"+mi.Ifname)); err != nil {
			return err
		}
	}
	for remappedEthName, mi := range mis {
		if err := LinkRename(LinkID("old"+mi.Ifname), remappedEthName); err != nil {
			return err
		}
		logger.Log.Info(fmt.Sprintf("Remapped NIC %v MAC %v to %v", mi.Ifname, mi.MACAddr, remappedEthName))
	}

	return nil
}

func vmwareLinksReorder() error {
	VNICMap := make(map[int]nICInfo)
	DNICMap := make(map[int]nICInfo)

	links, err := LinksGet()
	if err != nil {
		return err
	}

	reVNIC := regexp.MustCompile(`^Ethernet([0-9]+)$`)
	reDNIC := regexp.MustCompile(`^pciPassthru([0-9]+)$`)

	for _, l := range links {
		//Skip loopback
		if l.Ifindex == 1 {
			logger.Log.Info("Skipping loopback reordering")
			continue
		}
		//Remapping is only for devices
		if l.Linkinfo.InfoKind != "device" {
			logger.Log.Info("Skipping virtual network device %v reordering", string(l.Ifname))
			continue
		}
		path := "/sys/class/net/" + string(l.Ifname) + "/device/label"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		virtual := true
		matches := reVNIC.FindStringSubmatch(string(data))
		if len(matches) != 2 {
			matches := reDNIC.FindStringSubmatch(string(data))
			if len(matches) != 2 {
				return NewUnknownLinkDeviceLabel(string(data))
			}
			virtual = false
		}
		ethIndex, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		if virtual {
			VNICMap[ethIndex] = nICInfo{
				MACAddr: l.Address,
				Ifname:  l.Ifname,
			}
		} else {
			DNICMap[ethIndex] = nICInfo{
				MACAddr: l.Address,
				Ifname:  l.Ifname,
			}
		}
	}
	NICReMap := make(map[LinkID]nICInfo)

	for i, nicinfo := range VNICMap {
		remappedEthName := "eth" + strconv.Itoa(i)
		NICReMap[LinkID(remappedEthName)] = nICInfo{
			MACAddr: nicinfo.MACAddr,
			Ifname:  nicinfo.Ifname,
		}
	}
	for i, nicinfo := range DNICMap {
		remappedEthName := "eth" + strconv.Itoa(i+len(VNICMap))
		NICReMap[LinkID(remappedEthName)] = nICInfo{
			MACAddr: nicinfo.MACAddr,
			Ifname:  nicinfo.Ifname,
		}
	}

	return linksRename(NICReMap)
}
