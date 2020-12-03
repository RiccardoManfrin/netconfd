package nc

import (
	"net"

	"github.com/riccardomanfrin/netlink"
	"gitlab.lan.athonet.com/core/netconfd/logger"
)

// ModelDefault This is equivalent to 0.0.0.0/0 or ::/0
type ModelDefault string

// RouteDst - struct for RouteDst
type RouteDst struct {
	Ip           CIDRAddr
	ModelDefault ModelDefault
}

func (r *RouteDst) String() string {
	if r.ModelDefault == "default" {
		return string(r.ModelDefault)
	}
	return r.Ip.String()
}

func (r *RouteDst) parse(dst *net.IPNet) {
	if dst == nil {
		r.ModelDefault = "default"
	} else {
		r.Ip.ParseIPNet(*dst)
	}
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
	Dev      LinkID   `json:"dev,omitempty"`
	Protocol string   `json:"protocol,omitempty"`
	Metric   int32    `json:"metric,omitempty"`
	Scope    Scope    `json:"scope,omitempty"`
	Prefsrc  CIDRAddr `json:"prefsrc,omitempty"`
	// Route flags
	Flags *[]string `json:"flags,omitempty"`
}

func routeParse(route netlink.Route) (Route, error) {
	ncroute := Route{}
	ncroute.Dst.parse(route.Dst)
	l, err := netlink.LinkByIndex(route.LinkIndex)
	if err != nil {
		return ncroute, err
	}
	ncroute.Dev = LinkID(l.Attrs().Name)
	ncroute.Gateway.SetIP(route.Gw)
	logger.Log.Warning("Convert Protocol number to string (\"dhcp\"/\"static\")")
	ncroute.Protocol = route.Protocol.String()
	ncroute.Prefsrc.SetIP(route.Src)
	ncroute.Metric = int32(route.Priority)
	ncroute.Scope = Scope(route.Scope.String())
	return ncroute, nil
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
		ncroutes[i], err = routeParse(r)
		if err != nil {
			return ncroutes, err
		}
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
