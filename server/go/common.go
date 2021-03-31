package openapi

import (
	"net"

	"gitlab.lan.athonet.com/core/netconfd/logger"
	"gitlab.lan.athonet.com/core/netconfd/nc"
)

func dnssGet() ([]Dns, error) {
	var dnss []Dns
	ncdnss, err := nc.DNSsGet()
	if err == nil {
		dnss = make([]Dns, len(ncdnss))
		for i, l := range ncdnss {
			dnss[i] = ncDnsParse(l)
		}
	}
	return dnss, err
}

func linksGet() ([]Link, error) {
	var links []Link
	nclinks, err := nc.LinksGet()
	if err == nil {
		links = make([]Link, len(nclinks))
		for i, l := range nclinks {
			links[i] = ncLinkParse(l)
		}
	}
	return links, err
}

func dhcpsGet() ([]Dhcp, error) {
	var dhcps []Dhcp
	ncdhcps, err := nc.DHCPsGet()
	if err == nil {
		dhcps = make([]Dhcp, len(ncdhcps))
		for i, d := range ncdhcps {
			dhcps[i] = ncDhcpParse(d)
		}
	}
	return dhcps, err
}

func unmanagedListGet() ([]Unmanaged, error) {
	var umgmts []Unmanaged
	ncumgmts, err := nc.UnmanagedListGet()
	if err == nil {
		umgmts = make([]Unmanaged, len(ncumgmts))
		for i, u := range ncumgmts {
			umgmts[i] = ncUnmanagedParse(u)
		}
	}
	return umgmts, err
}

func ncUnmanagedFormat(u Unmanaged) (nc.Unmanaged, error) {
	d := nc.Unmanaged{
		Type: nc.Type(u.GetType()),
		ID:   nc.UnmanagedID(u.GetId()),
	}
	return d, nil
}

func ncUnmanagedParse(ncunmanaged nc.Unmanaged) Unmanaged {
	d := Unmanaged{}
	d.SetId(string(ncunmanaged.ID))
	d.SetType(string(ncunmanaged.Type))
	return d
}

func ncDhcpFormat(dhcp Dhcp) (nc.Dhcp, error) {
	d := nc.Dhcp{
		Ifname: nc.LinkID(dhcp.Ifname),
	}
	return d, nil
}

func ncDhcpParse(ncdhcp nc.Dhcp) Dhcp {
	d := Dhcp{
		Ifname: string(ncdhcp.Ifname),
	}
	return d
}

func ncRouteFormat(route Route) (nc.Route, error) {
	ncroute := nc.Route{}
	ncroute.Dst = route.Dst
	if route.Gateway != nil {
		ncroute.Gateway = *route.Gateway
	}
	if route.Prefsrc != nil {
		ncroute.Prefsrc = *route.Prefsrc
	}
	if route.Dev != nil {
		ncroute.Dev = nc.LinkID(*route.Dev)
	}
	if route.Metric != nil {
		ncroute.Metric = *route.Metric
	}
	return ncroute, nil
}

func ncRouteParse(ncroute nc.Route) Route {
	var route Route
	id := string(ncroute.ID)
	route.Id = &id
	prefsrc := ncroute.Prefsrc.String()
	if prefsrc != "" {
		route.SetPrefsrc(prefsrc)
	}
	dst := ncroute.Dst.String()
	if dst != "" {
		route.SetDst(dst)
	}
	gw := ncroute.Gateway.String()
	if gw != "" {
		route.SetGateway(gw)
	}
	route.SetDev(string(ncroute.Dev))
	route.SetProtocol(ncroute.Protocol)
	route.SetMetric(ncroute.Metric)
	route.SetScope(Scope(ncroute.Scope))
	return route
}

func routesGet() ([]Route, error) {
	var routes []Route
	ncroutes, err := nc.RoutesGet()
	if err == nil {
		routes = make([]Route, len(ncroutes))
		for i, r := range ncroutes {
			routes[i] = ncRouteParse(r)
		}
	}
	return routes, err
}

func ncDnsParse(ncdns nc.Dns) Dns {
	ns := ncdns.Nameserver.String()
	prio := Dnsid(ncdns.Id)
	return Dns{
		Nameserver: ns,
		Id:         prio,
	}
}

func ncDnsFormat(dns Dns) (nc.Dns, error) {
	d := nc.Dns{
		Nameserver: net.ParseIP(dns.GetNameserver()),
		Id:         nc.DnsID(dns.GetId()),
	}
	return d, nil
}

func ncLinkParse(nclink nc.Link) Link {
	link := Link{
		Ifname:    string(nclink.Ifname),
		Ifindex:   &nclink.Ifindex,
		LinkType:  nclink.LinkType,
		Operstate: &nclink.Operstate,
	}

	link.Mtu = &nclink.Mtu

	flagsLen := len(nclink.Flags)
	if flagsLen > 0 {
		lfs := make([]LinkFlag, flagsLen)
		link.Flags = &lfs
		for i, lf := range nclink.Flags {
			(*link.Flags)[i] = LinkFlag(lf)
		}
	}

	if len(nclink.AddrInfo) > 0 {
		lai := make([]LinkAddrInfo, len(nclink.AddrInfo))
		for i, a := range nclink.AddrInfo {
			lai[i].Local = a.Local.ToIPNet().IP
			lai[i].Prefixlen = int32(a.Local.PrefixLen())
			if a.Address != nil {
				lai[i].Address = a.Address
			}
		}
		link.AddrInfo = &lai
	}

	lli := LinkLinkinfo{}
	if nclink.Linkinfo.InfoKind != "" {
		lli.InfoKind = &nclink.Linkinfo.InfoKind
		link.Linkinfo = &lli
	}

	switch nclink.Linkinfo.InfoKind {
	case "vlan":
		{
			id := LinkLinkinfoInfoData{}
			id.Protocol = &nclink.Linkinfo.InfoData.Protocol
			id.Id = &nclink.Linkinfo.InfoData.Id
			lli.InfoData = &id
			parentLink := string(nclink.Link)
			link.Link = &parentLink
		}
	case "gre":
		{
			id := LinkLinkinfoInfoData{}
			lli.InfoData = &id
			lli.InfoData.SetLocal(nclink.Linkinfo.InfoData.Local.String())
			lli.InfoData.SetRemote(nclink.Linkinfo.InfoData.Remote.String())
		}
	case "bond":
		{
			id := LinkLinkinfoInfoData{}
			id.Mode = &nclink.Linkinfo.InfoData.Mode
			id.Miimon = &nclink.Linkinfo.InfoData.Miimon
			id.Downdelay = &nclink.Linkinfo.InfoData.Downdelay
			id.Updelay = &nclink.Linkinfo.InfoData.Updelay
			id.PeerNotifyDelay = &nclink.Linkinfo.InfoData.PeerNotifyDelay
			id.UseCarrier = &nclink.Linkinfo.InfoData.UseCarrier
			id.ArpInterval = &nclink.Linkinfo.InfoData.ArpInterval
			id.ArpValidate = &nclink.Linkinfo.InfoData.ArpValidate
			id.LpInterval = &nclink.Linkinfo.InfoData.LpInterval
			id.ArpAllTargets = &nclink.Linkinfo.InfoData.ArpAllTargets
			id.PacketsPerSlave = &nclink.Linkinfo.InfoData.PacketsPerSlave
			id.FailOverMac = &nclink.Linkinfo.InfoData.FailOverMac
			id.XmitHashPolicy = &nclink.Linkinfo.InfoData.XmitHashPolicy
			id.ResendIgmp = &nclink.Linkinfo.InfoData.ResendIgmp
			id.MinLinks = &nclink.Linkinfo.InfoData.MinLinks
			id.ArpInterval = &nclink.Linkinfo.InfoData.ArpInterval
			id.PrimaryReselect = &nclink.Linkinfo.InfoData.PrimaryReselect
			id.TlbDynamicLb = &nclink.Linkinfo.InfoData.TlbDynamicLb
			id.AdSelect = &nclink.Linkinfo.InfoData.AdSelect
			id.AdLacpRate = &nclink.Linkinfo.InfoData.AdLacpRate
			id.Mode = &nclink.Linkinfo.InfoData.Mode
			id.AllSlavesActive = &nclink.Linkinfo.InfoData.AllSlavesActive
			id.UseCarrier = &nclink.Linkinfo.InfoData.UseCarrier
			lli.InfoData = &id
		}
	case "device":
	case "bridge":
	case "dummy":
	case "ppp":
	case "tun":
	case "tap":
		{
		}
	default:
		{
			logger.Log.Warning("Unknown Link Kind : %v", nclink.Linkinfo.InfoKind)
		}
	}

	if nclink.Master != "" {
		master := string(nclink.Master)
		link.Master = &master
		isd := LinkLinkinfoInfoSlaveData{}
		lli.InfoSlaveData = &isd
		icisd := &nclink.Linkinfo.InfoSlaveData
		link.SetMaster(string(nclink.Master))
		isd.SetState(icisd.State)
		isd.SetLinkFailureCount(int32(icisd.LinkFailureCount))
		isd.SetMiiStatus(icisd.MiiStatus)
		isd.SetPermHwaddr(icisd.PermHwaddr)
	}

	return link
}

func ncLinkFormat(link Link) (nc.Link, error) {

	nclink := nc.Link{
		Ifname:   nc.LinkID(link.GetIfname()),
		Linkinfo: nc.LinkLinkinfo{},
		Mtu:      link.GetMtu(),
		LinkType: link.GetLinkType(),
		Master:   nc.LinkID(link.GetMaster()),
	}

	if link.Link != nil {
		nclink.Link = nc.LinkID(*link.Link)
	}

	if link.Flags != nil {
		flagsLen := len(*link.Flags)
		if flagsLen > 0 {
			nclink.Flags = make([]nc.LinkFlag, flagsLen)
			for i, lf := range *link.Flags {
				nclink.Flags[i] = nc.LinkFlag(lf)
			}
		}
	}

	li := link.GetLinkinfo()

	if li.InfoData != nil {
		nclink.Linkinfo.InfoData = nc.LinkLinkinfoInfoData{
			Mode:            li.InfoData.GetMode(),
			Miimon:          li.InfoData.GetMiimon(),
			Downdelay:       li.InfoData.GetDowndelay(),
			Updelay:         li.InfoData.GetUpdelay(),
			PeerNotifyDelay: li.InfoData.GetPeerNotifyDelay(),
			UseCarrier:      li.InfoData.GetUseCarrier(),
			ArpInterval:     li.InfoData.GetArpInterval(),
			ArpValidate:     li.InfoData.GetArpValidate(),
			LpInterval:      li.InfoData.GetLpInterval(),
			ArpAllTargets:   li.InfoData.GetArpAllTargets(),
			PacketsPerSlave: li.InfoData.GetPacketsPerSlave(),
			FailOverMac:     li.InfoData.GetFailOverMac(),
			XmitHashPolicy:  li.InfoData.GetXmitHashPolicy(),
			ResendIgmp:      li.InfoData.GetResendIgmp(),
			MinLinks:        li.InfoData.GetMinLinks(),
			PrimaryReselect: li.InfoData.GetPrimaryReselect(),
			TlbDynamicLb:    li.InfoData.GetTlbDynamicLb(),
			AdSelect:        li.InfoData.GetAdSelect(),
			AdLacpRate:      li.InfoData.GetAdLacpRate(),
			AllSlavesActive: li.InfoData.GetAllSlavesActive(),
			Protocol:        li.InfoData.GetProtocol(),
			Id:              li.InfoData.GetId(),
			Local:           net.ParseIP(li.InfoData.GetLocal()),
			Remote:          net.ParseIP(li.InfoData.GetRemote()),
		}

	}
	if li.InfoKind != nil {
		nclink.Linkinfo.InfoKind = *li.InfoKind
	}

	if li.InfoSlaveData != nil {
		nclink.Linkinfo.InfoSlaveData.State = li.InfoSlaveData.GetState()
	}
	if link.AddrInfo != nil {
		nclink.AddrInfo = make([]nc.LinkAddrInfo, len(*link.AddrInfo))
		for i, addr := range *link.AddrInfo {
			var cidrnet nc.CIDRAddr
			cidrnet.SetIP(addr.Local)
			err := cidrnet.SetPrefixLen(int(addr.Prefixlen))
			if err != nil {
				return nclink, err
			}
			lai := nc.LinkAddrInfo{
				Local:   cidrnet,
				Address: addr.Address,
			}
			nclink.AddrInfo[i] = lai
		}
	}
	return nclink, nil
}

func ncNetFormat(config Config) (nc.Network, error) {
	network := nc.Network{}
	if config.Network != nil {
		if config.Network.Links != nil {
			network.Links = make([]nc.Link, len(*config.Network.Links))
			for i, l := range *config.Network.Links {
				var err error
				network.Links[i], err = ncLinkFormat(l)
				if err != nil {
					return network, err
				}
			}

		}
		if config.Network.Routes != nil {
			network.Routes = make([]nc.Route, len(*config.Network.Routes))
			for i, l := range *config.Network.Routes {
				network.Routes[i], _ = ncRouteFormat(l)
			}
		}
		if config.Network.Dhcp != nil {
			network.Dhcp = make([]nc.Dhcp, len(*config.Network.Dhcp))
			for i, d := range *config.Network.Dhcp {
				network.Dhcp[i], _ = ncDhcpFormat(d)
			}
		}
		if config.Network.Dns != nil {
			network.Dnss = make([]nc.Dns, len(*config.Network.Dns))
			for i, s := range *config.Network.Dns {
				network.Dnss[i], _ = ncDnsFormat(s)
			}
		}
		if config.Network.Unmanaged != nil {
			network.Unmanaged = make([]nc.Unmanaged, len(*config.Network.Unmanaged))
			for i, s := range *config.Network.Unmanaged {
				network.Unmanaged[i], _ = ncUnmanagedFormat(s)
			}
		}
	}
	return network, nil
}

func ncNetParse(net nc.Network) Network {
	links := make([]Link, len(net.Links))
	routes := make([]Route, len(net.Routes))
	dhcps := make([]Dhcp, len(net.Dhcp))
	dnss := make([]Dns, len(net.Dnss))
	unmanaged := make([]Unmanaged, len(net.Unmanaged))
	for i, l := range net.Links {
		links[i] = ncLinkParse(l)
	}
	for i, r := range net.Routes {
		routes[i] = ncRouteParse(r)
	}
	for i, d := range net.Dhcp {
		dhcps[i] = ncDhcpParse(d)
	}
	for i, s := range net.Dnss {
		dnss[i] = ncDnsParse(s)
	}
	for i, s := range net.Unmanaged {
		unmanaged[i] = ncUnmanagedParse(s)
	}
	return Network{
		Links:     &links,
		Routes:    &routes,
		Dhcp:      &dhcps,
		Dns:       &dnss,
		Unmanaged: &unmanaged,
	}
}
