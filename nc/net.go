package nc

// Network struct for Network
type Network struct {
	// Series of links layer interfaces to configure within the namespace
	Links []Link `json:"links,omitempty"`
	// Namespace routes
	Routes []Route `json:"routes,omitempty"`
	// DHCP context
	Dhcp []Dhcp `json:"dhcp,omitempty"`
	// DNS context
	Dnss []Dns
	//Unmanaged context
	Unmanaged []Unmanaged
}

//Patch network config
func Patch(n Network) error {
	err := UnamanagedListConfigure(n.Unmanaged)
	if err != nil {
		return err
	}
	err = LinksConfigure(n.Links)
	if err != nil {
		return err
	}
	err = RoutesConfigure(n.Routes)
	if err != nil {
		return err
	}
	err = DHCPsConfigure(n.Dhcp)
	if err != nil {
		return err
	}
	err = DNSsConfigure(n.Dnss)
	if err != nil {
		return err
	}
	return nil
}

//Put network config (wipe out and redeploy)
func Put(n Network) error {
	/* FIRST WE update the unmanaged resources, than we
	 * apply the config so we can always add unmanaged
	 * resources and use them in the same PUT/PATCH api
	 */
	err := UnmanagedListDelete()
	if err != nil {
		return err
	}
	err = UnamanagedListConfigure(n.Unmanaged)
	if err != nil {
		return err
	}

	err = Del()
	if err != nil {
		return err
	}

	err = LinksConfigure(n.Links)
	if err != nil {
		return err
	}
	err = RoutesConfigure(n.Routes)
	if err != nil {
		return err
	}
	err = DHCPsConfigure(n.Dhcp)
	if err != nil {
		return err
	}
	err = DNSsConfigure(n.Dnss)
	if err != nil {
		return err
	}
	return nil
}

//Del delete whole network config
func Del() error {
	err := LinksDelete()
	if err != nil {
		return err
	}
	err = RoutesDelete()
	if err != nil {
		return err
	}
	err = DHCPsDelete()
	if err != nil {
		return err
	}
	err = DNSsDelete()
	if err != nil {
		return err
	}
	return nil
}

//Get network config
func Get() (Network, error) {
	n := Network{}
	links, err := LinksGet()
	if err != nil {
		return n, err
	}
	routes, err := RoutesGet()
	if err != nil {
		return n, err
	}
	dhcps, err := DHCPsGet()
	if err != nil {
		return n, err
	}
	dnss, err := DNSsGet()
	if err != nil {
		return n, err
	}
	umgmts, err := UnmanagedListGet()
	if err != nil {
		return n, err
	}
	n.Links = links
	n.Routes = routes
	n.Dhcp = dhcps
	n.Dnss = dnss
	n.Unmanaged = umgmts
	return n, nil
}
