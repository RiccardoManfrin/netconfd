package openapi

import (
	"gitlab.lan.athonet.com/core/netconfd/logger"
	"gitlab.lan.athonet.com/core/netconfd/nc"
)

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

func ncLinkParse(nclink nc.Link) Link {
	link := Link{
		Ifname:    string(nclink.Ifname),
		Ifindex:   &nclink.Ifindex,
		LinkType:  nclink.LinkType,
		Operstate: &nclink.Operstate,
	}

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
			ip := a.Local.Address()
			lai[i].Local = Ip{string: &ip}
			lai[i].Prefixlen = int32(a.Local.PrefixLen())
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
	default:
		{
			logger.Log.Warning("Unknown Link Kind")
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
			cidrnet.ParseIP(*(addr.Local.string))
			cidrnet.ParsePrefixLen(int(addr.Prefixlen))
			lai := nc.LinkAddrInfo{
				Local: cidrnet,
			}
			nclink.AddrInfo[i] = lai
		}
	}
	return nclink
}

func ncNetFormat(config Config) nc.Network {
	network := nc.Network{}
	if config.HostNetwork != nil {
		if config.HostNetwork.Links != nil {
			network.Links = make([]nc.Link, len(*config.HostNetwork.Links))
			for i, l := range *config.HostNetwork.Links {
				network.Links[i] = ncLinkFormat(l)
			}

		}
		if config.HostNetwork.Routes != nil {
			//TODO
		}
	}
	return network
}

func ncNetParse(net nc.Network) Config {
	links := make([]Link, len(net.Links))
	for i, l := range net.Links {
		links[i] = ncLinkParse(l)
	}
	return Config{
		HostNetwork: &Network{
			Links: &links,
		},
	}
}
