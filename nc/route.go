package nc

import "github.com/vishvananda/netlink"

// ModelDefault This is equivalent to 0.0.0.0/0 or ::/0
type ModelDefault string

// RouteDst - struct for RouteDst
type RouteDst struct {
	Ip           CIDRAddr
	ModelDefault ModelDefault
}

// Scope scope of the object (link or global)
type Scope string

// List of scope
const (
	LINK   Scope = "link"
	GLOBAL Scope = "global"
)

// Route IP L3 Ruote entry
type Route struct {
	Dst     RouteDst `json:"dst,omitempty"`
	Gateway CIDRAddr `json:"gateway,omitempty"`
	// Interface name
	Dev      string   `json:"dev,omitempty"`
	Protocol string   `json:"protocol,omitempty"`
	Metric   int32    `json:"metric,omitempty"`
	Scope    Scope    `json:"scope,omitempty"`
	Prefsrc  CIDRAddr `json:"prefsrc,omitempty"`
	// Route flags
	Flags *[]string `json:"flags,omitempty"`
}

func routeParse(route netlink.Route) Route {
	ncroute := Route{}
	if route.Dst != nil {
		ncroute.Dst.Ip.ParseIPNet(*route.Dst)
	}
	return ncroute
}

//RouteID identifies a route
type RouteID string

func routeID(route Route) RouteID {
	return RouteID(route.Dst.Ip.String())
}

//RoutesGet returns the array of routes
func RoutesGet() ([]Route, error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, err
	}
	ncroutes := make([]Route, len(routes))
	for i, r := range routes {
		ncroutes[i] = routeParse(r)
	}
	return ncroutes, nil
}

//RouteGet Returns the list of existing link layer devices on the machine
func RouteGet(_routeID RouteID) (Route, error) {
	routes, err := RoutesGet()
	if err != nil {
		return Route{}, err
	}
	for _, r := range routes {
		if routeID(r) == _routeID {
			return r, nil
		}
	}
	return Route{}, NewRouteByIDNotFoundError(_routeID)
}
