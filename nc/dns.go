package nc

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"

	"gitlab.lan.athonet.com/core/netconfd/logger"
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
	// Evaluated priority
	Id DnsID `json:"__id,omitempty"`
}

//ResolvConf path prefix
const ResolvConf string = "/etc/resolv.conf"

func loadResolv() ([]Dns, error) {
	dnss := make([]Dns, 0)
	links, err := LinksGet()
	if err != nil {
		return dnss, err
	}
	//Credits to https://regex101.com/
	re := regexp.MustCompile(`^.*:\ ([^ \n]*)\ ?([^ \n]*).*`)

	for _, l := range links {
		if l.Ifindex == 1 {
			logger.Log.Debug("Skipping lo interface DNS nameserver config")
			continue
		}
		out, err := exec.Command(prefixInstallPAth + "dns_status.sh").Output()
		if err != nil {
			return dnss, err
		}
		matches := re.FindStringSubmatch(string(out))
		if len(matches) == 3 {
			ip := net.ParseIP(matches[1])
			if ip != nil {
				dnss = append(dnss, Dns{Nameserver: ip, Id: DnsPrimary})
			}
			ip = net.ParseIP(matches[2])
			if ip != nil {
				dnss = append(dnss, Dns{Nameserver: ip, Id: DnsSecondary})
			}
		}
		//do it once only
		break
	}
	return dnss, nil
}

func dumpResolv(dnss []Dns) error {
	/*
	* Example:
	* resolvectl dns eth0 192.168.178.1 8.8.8.8
	* root@ngcore:~# resolvectl dns eth0
	 */
	links, err := LinksGet()
	if err != nil {
		return err
	}
	if dnss == nil {
		return nil
	}

	primary := ""
	secondary := ""
	for _, d := range dnss {
		if d.Id == DnsPrimary {
			if isUnmanaged(UnmanagedID(DnsPrimary), DNSTYPE) {
				logger.Log.Info(fmt.Sprintf("Skipping Unmanaged DNS %v configuration", DnsPrimary))
				continue
			}
			primary = d.Nameserver.String()
		}
		if d.Id == DnsSecondary {
			if isUnmanaged(UnmanagedID(DnsSecondary), DNSTYPE) {
				logger.Log.Info(fmt.Sprintf("Skipping Unmanaged DNS %v configuration", DnsSecondary))
				continue
			}
			secondary = d.Nameserver.String()
		}
	}
	if primary == "" {
		primary = secondary
		secondary = ""
	}

	for _, l := range links {
		if l.Ifindex == 1 {
			logger.Log.Debug("Skipping lo interface DNS nameserver config")
			continue
		}
		if err := dnsConfigure(l.Ifname, primary, secondary); err != nil {
			return err
		}

	}
	return nil
}

func dnsConfigure(ifname LinkID, primary string, secondary string) error {
	out, err := exec.Command(prefixInstallPAth+"dns_configure.sh", string(ifname), primary, secondary).Output()
	if err != nil {
		return err
	}
	logger.Log.Debug(string(out))
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
	return loadResolv()
}
