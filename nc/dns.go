package nc

import (
	"bufio"
	"io"
	"net"
	"os"
	"regexp"
	"strings"

	"gitlab.lan.athonet.com/core/netconfd/logger"
)

// Dns Name server for DNS resolution
type Dns struct {
	// The DNS server ip address to send DNS queries to
	Nameserver net.IP `json:"nameserver,omitempty"`
	// Evaluated priority (lower value indicates higher priority)
	Priority int `json:"priority,omitempty"`
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
				Priority:   prio,
			}
			prio++
			dnss = append(dnss, dns)
		}
	}
	return dnss
}
func dumpResolv(dnss []Dns) {

}

//DNSsGet Get all DNS interfaces administrated by DNS and related config/state.
func DNSsGet() ([]Dns, error) {
	return loadResolv(), nil
}
