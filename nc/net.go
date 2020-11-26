package nc

// Network struct for Network
type Network struct {
	// Series of links layer interfaces to configure within the namespace
	Links []Link `json:"links,omitempty"`
	// Namespace routes
	Routes []Route `json:"routes,omitempty"`
}

//Patch network config
func Patch(n Network) error {
	err := LinksConfigure(n.Links)
	if err != nil {
		return err
	}
	return nil
}

//Put network config (wipe out and redeploy)
func Put(n Network) error {
	err := LinksConfigure(n.Links)
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
	n.Links = links
	return n, nil
}
