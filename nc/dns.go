package nc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

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
	// Evaluated priority (lower value indicates higher priority)
	Id DnsID `json:"__id,omitempty"`
}

//ResolvConf path prefix
const ResolvConf string = "/etc/resolv.conf"

func loadResolv() []Dns {
	dnss := make([]Dns, 0)
	file, err := os.Open(ResolvConf)
	if err != nil {
		logger.Fatal(err)
		return dnss
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	prio := 1
	/*
	* https://man7.org/linux/man-pages/man5/resolv.conf.5.html
	*
	* The keyword and value must appear on a single line, and the
	* keyword (e.g., nameserver) must start the line.  The value
	* follows the keyword, separated by white space.
	 */
	re := regexp.MustCompile(`^nameserver {1}[^\s].*`)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if re.MatchString(line) {
			ns := strings.Replace(line, "nameserver", "", 1)
			ns = strings.ReplaceAll(ns, " ", "")
			ns = strings.ReplaceAll(ns, "\n", "")
			dns := Dns{
				Nameserver: net.ParseIP(ns),
				Id:         dnsIDToString[prio],
			}
			prio++
			dnss = append(dnss, dns)
			if prio > 2 {
				break
			}
		}
	}
	return dnss
}
func dumpResolv(dnss []Dns) error {
	if len(dnss) > 2 {
		return NewTooManyDNSServersError()
	}
	resolvConf := ""
	var pri *Dns = nil
	var sec *Dns = nil

	if len(dnss) > 0 {
		if dnss[0].Id == DnsPrimary {
			pri = &dnss[0]
		} else if dnss[0].Id == DnsSecondary {
			sec = &dnss[0]
		} else {
			return NewUnknownUnsupportedDNSServersIDsError(dnss[0].Id)
		}

		if len(dnss) == 2 {
			if dnss[1].Id == DnsPrimary {
				if pri != nil {
					return NewDuplicateDNSServersIDsError(dnss[0].Id, dnss[1].Id)
				}
				pri = &dnss[1]
			} else if dnss[1].Id == DnsSecondary {
				if sec != nil {
					return NewDuplicateDNSServersIDsError(dnss[0].Id, dnss[1].Id)
				}
				sec = &dnss[1]
			} else {
				return NewUnknownUnsupportedDNSServersIDsError(dnss[1].Id)
			}
		}
		if pri != nil {
			resolvConf = fmt.Sprintf(
				"%vnameserver %v\n",
				resolvConf,
				pri.Nameserver.String())
		}
		if sec != nil {
			resolvConf = fmt.Sprintf(
				"%vnameserver %v\n",
				resolvConf,
				sec.Nameserver.String())
		}
	}
	return ioutil.WriteFile(
		ResolvConf,
		[]byte(resolvConf), 777)
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
