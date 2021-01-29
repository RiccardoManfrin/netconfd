package nc

import (
	"net"
)

type DnsID string

const (
	DnsPrimary   DnsID = "primary"
	DnsSecondary       = "secondary"
)

var dnsIDToString = map[int]DnsID{
	1: DnsPrimary,
	2: DnsSecondary,
}
var dnsStringToID = map[DnsID]int{
	DnsPrimary:   1,
	DnsSecondary: 2,
}

// Dns Name server for DNS resolution
type Dns struct {
	// The DNS server ip address to send DNS queries to
	Nameserver net.IP `json:"nameserver,omitempty"`
	// Evaluated priority (lower value indicates higher priority)
	Id DnsID `json:"__id,omitempty"`
}

//ResolvConf path prefix
const ResolvConf string = "/etc/resolv.conf"

func loadResolv() []Dns {
	dnss := make([]Dns, 0)

	return dnss
}
func dumpResolv(dnss []Dns) error {
	/*
	* Example:
	* resolvectl dns eth0 192.168.178.1 8.8.8.8
	* root@ngcore:~# resolvectl dns eth0
	 */
	return nil
}

func DNSCreate(dns Dns) error {
	_, err := DnsGet(dns.Id)
	if err == nil {
		return NewDNSServerExistsConflictError(dns.Id)
	}
	if _, ok := err.(*NotFoundError); ok != true {
		return err
	}
	dnss, err := DNSsGet()
	if err != nil {
		return err
	}
	dnss = append(dnss, dns)
	return dumpResolv(dnss)
}

//DNSsConfigure configures/overwrites the whole set of dnss
func DNSsConfigure(dnss []Dns) error {
	return dumpResolv(dnss)
}

//DNSsDelete deletes all DNS context
func DNSsDelete() error {
	dnss := make([]Dns, 0)
	return dumpResolv(dnss)
}

//DNSDelete delete a DNS entry in resolv.conf
func DNSDelete(dnsid DnsID) error {
	dnss, err := DNSsGet()
	if err != nil {
		return err
	}
	var resultDNSs []Dns
	/* we use knowledge on positional order */
	switch dnsid {
	case DnsPrimary:
		{
			resultDNSs = dnss[1:]
		}
	case DnsSecondary:
		{
			resultDNSs = dnss[:1]
		}
	default:
		return NewUnknownUnsupportedDNSServersIDsError(dnsid)
	}
	return dumpResolv(resultDNSs)
}

func DnsGet(dnsid DnsID) (Dns, error) {
	dnss, err := DNSsGet()
	if err != nil {
		return Dns{}, err
	}
	for _, dns := range dnss {
		if dns.Id == dnsid {
			return dns, nil
		}
	}
	return Dns{}, NewDNSServerNotFoundError(dnsid)
}

//DNSsGet Get all DNS interfaces administrated by DNS and related config/state.
func DNSsGet() ([]Dns, error) {
	return loadResolv(), nil
}
