package nc

// Network struct for Network
type Network struct {
	// Series of links layer interfaces to configure within the namespace
	Links []Link `json:"links,omitempty"`
	// Namespace routes
	Routes []Route `json:"routes,omitempty"`
	// DHCP context
	Dhcp []Dhcp `json:"dhcp,omitempty"`
}

//Patch network config
func Patch(n Network) error {
	err := LinksConfigure(n.Links)
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
	return nil
}

//Put network config (wipe out and redeploy)
func Put(n Network) error {
	err := Del()
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
	n.Links = links
	n.Routes = routes
	n.Dhcp = dhcps

	return n, nil
}
