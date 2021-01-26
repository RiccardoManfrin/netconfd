package openapi

import (
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

func ncDhcpFormat(dhcp Dhcp) nc.Dhcp {
	d := nc.Dhcp{
		Ifname: nc.LinkID(dhcp.Ifname),
	}
	return d
}

func ncDhcpParse(ncdhcp nc.Dhcp) Dhcp {
	d := Dhcp{
		Ifname: string(ncdhcp.Ifname),
	}
	return d
}

func ncRouteFormat(route Route) nc.Route {
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
	return ncroute
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
		Nameserver: &ns,
		Id:         &prio,
	}
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
	isd := LinkLinkinfoInfoSlaveData{}
	lli.InfoSlaveData = &isd
	if nclink.Master != "" {
		master := string(nclink.Master)
		link.Master = &master
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
		icisd := &nclink.Linkinfo.InfoSlaveData
		link.SetMaster(string(nclink.Master))
		isd.SetState(icisd.State)
		isd.SetLinkFailureCount(int32(icisd.LinkFailureCount))
		isd.SetMiiStatus(icisd.MiiStatus)
		isd.SetPermHwaddr(icisd.PermHwaddr)
	}

	return link
}

func ncLinkFormat(link Link) nc.Link {

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
			cidrnet.SetPrefixLen(int(addr.Prefixlen))
			lai := nc.LinkAddrInfo{
				Local:   cidrnet,
				Address: addr.Address,
			}
			nclink.AddrInfo[i] = lai
		}
	}
	return nclink
}

func ncNetFormat(config Config) nc.Network {
	network := nc.Network{}
	if config.Network != nil {
		if config.Network.Links != nil {
			network.Links = make([]nc.Link, len(*config.Network.Links))
			for i, l := range *config.Network.Links {
				network.Links[i] = ncLinkFormat(l)
			}

		}
		if config.Network.Routes != nil {
			network.Routes = make([]nc.Route, len(*config.Network.Routes))
			for i, l := range *config.Network.Routes {
				network.Routes[i] = ncRouteFormat(l)
			}
		}
		if config.Network.Dhcp != nil {
			network.Dhcp = make([]nc.Dhcp, len(*config.Network.Dhcp))
			for i, d := range *config.Network.Dhcp {
				network.Dhcp[i] = ncDhcpFormat(d)
			}
		}
	}
	return network
}

func ncNetParse(net nc.Network) Network {
	links := make([]Link, len(net.Links))
	routes := make([]Route, len(net.Routes))
	dhcps := make([]Dhcp, len(net.Dhcp))
	for i, l := range net.Links {
		links[i] = ncLinkParse(l)
	}
	for i, r := range net.Routes {
		routes[i] = ncRouteParse(r)
	}
	for i, d := range net.Dhcp {
		dhcps[i] = ncDhcpParse(d)
	}
	return Network{
		Links:  &links,
		Routes: &routes,
		Dhcp:   &dhcps,
	}
}
