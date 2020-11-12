package nc

import (
	"errors"
	"net"
	"strconv"
)

// CIDRAddr is an address and a network mask
// Additionally to a net.IPNet, it allows for specifying
// further than the netmask bits. Those are intended to define
// an addresses within the IP network being defined along with.
// E.g. : 10.1.2.3/24 -> 10.1.2.3 in network 10.1.2.0/24
type CIDRAddr struct {
	ip  net.IP
	net net.IPNet
}

//NewCIDRAddr creates new CIDR address. If network is unspecified it is assumed to be /32 for ipv4 and /128 for ipv6
func NewCIDRAddr(addr string) CIDRAddr {
	a := CIDRAddr{}
	a.Parse(addr)
	return a
}

//ParseIPNet translates an IP network into a CIDRAddr
func (a *CIDRAddr) ParseIPNet(ip net.IPNet) {
	a.Parse(ip.String())
}

//Parse loads a CIDR address from a string. If network is unspecified it is assumed to be /32 for ipv4 and /128 for ipv6
func (a *CIDRAddr) Parse(straddr string) error {
	var e error
	var ipnet *net.IPNet
	a.ip, ipnet, e = net.ParseCIDR(straddr)
	if e != nil {
		a.ip = net.ParseIP(straddr)
		if a.ip == nil {
			return errors.New("Invalid Address")
		}
		if a.ip.To4() != nil {
			a.net.Mask = net.CIDRMask(32, 32)
			a.net.IP = a.ip.Mask(a.net.Mask)
		} else {
			a.net.Mask = net.CIDRMask(128, 128)
			a.net.IP = a.ip.Mask(a.net.Mask)
		}

	} else {
		a.net = *ipnet
	}
	return nil
}

func (a *CIDRAddr) String() string {
	ones, bits := a.net.Mask.Size()
	if ones == bits {
		return a.ip.String()
	}
	return a.ip.String() + "/" + strconv.Itoa(ones)
}

//Netmask returns the netmask (e.g. 255.255.255.0) of a CIDR address/network
func (a *CIDRAddr) Netmask() string {
	return a.net.IP.String()
}

//Address returns the address (e.g. 255.255.255.0) of a CIDR address/network
func (a *CIDRAddr) Address() string {
	return a.ip.String()
}
