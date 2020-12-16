package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	comm "gitlab.lan.athonet.com/core/netconfd/common"
	"gitlab.lan.athonet.com/core/netconfd/logger"
	"gitlab.lan.athonet.com/core/netconfd/nc"
	oas "gitlab.lan.athonet.com/core/netconfd/server/go"
)

func parseSampleConfig(t *testing.T, sampleConfig string) oas.Config {
	var config oas.Config
	err := json.Unmarshal([]byte(sampleConfig), &config)
	if err != nil {
		t.Error(err)
	}
	return config
}

var sampleConfig string = `
{
"global": {},
"host_network": {
	"links": [
		{
			"ifname": "bond0",
			"link_type": "ether",
			"flags": ["up"],
			"linkinfo": {
			"info_kind": "bond",
			"info_data": {
				"mode": "active-backup",
				"downdelay": 800,
				"updelay" : 400,
				"miimon" : 200
			}
			}
		},
		{
			"ifname": "dummy0",
			"link_type": "ether",
			"flags": ["up"],
			"linkinfo": {
				"info_kind": "dummy",
				"info_slave_data": {
					"state": "BACKUP"
				}
			},
			"addr_info": [
				{
					"local": "10.6.7.8",
					"prefixlen": 24
				}
			],
			"master": "bond0"
		},
		{
			"ifname": "dummy1",
			"link_type": "ether",
			"linkinfo": {
			"info_kind": "dummy",
			"info_slave_data": {
				"state": "ACTIVE"
			}
			},
			"master": "bond0"
		}
	],
	"routes": [
		{
			"dev": "dummy0",
			"dst": "10.8.9.0/24",
			"gateway": "10.6.7.8",
			"metric": 50,
			"protocol": "boot",
			"scope": "universe"
		}
	]
}
}`

func genSampleConfig(t *testing.T) oas.Config {
	return parseSampleConfig(t, sampleConfig)
}

func newConfigPatchReq(config oas.Config) *http.Request {
	reqbody, _ := json.Marshal(config)
	iobody := bytes.NewReader(reqbody)
	req, _ := http.NewRequest("PATCH", "/api/1/mgmt/config", iobody)
	req.Header.Add("Content-Type", "application/json")
	return req
}
func newConfigGetReq() *http.Request {
	req, _ := http.NewRequest("GET", "/api/1/mgmt/config", nil)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func NewTestManager() *Manager {
	mgr := NewManager()
	logger.LoggerInit("-")
	logger.LoggerSetLevel("WRN")
	return mgr
}

var m *Manager = NewTestManager()

func checkResponse(t *testing.T, rr *httptest.ResponseRecorder, httpStatusCode int, ncErrorCode nc.ErrorCode, ncreason string) {
	if status := rr.Code; status != httpStatusCode {
		t.Errorf("HTTP Status code mismatch: got [%v] want [%v]",
			status,
			httpStatusCode)
	}
	var genericError nc.GenericError
	if ncErrorCode != nc.RESERVED {
		err := json.Unmarshal(rr.Body.Bytes(), &genericError)
		if err != nil {
			t.Errorf("Err Unmarshal failure")
		}
		if genericError.Code != ncErrorCode {
			t.Errorf("Err Code mismatch: got [%v], want [%v]",
				genericError.Code,
				nc.SEMANTIC)
		}
		if ncreason != "" {
			if genericError.Reason != ncreason {
				t.Errorf("Err Reason mismatch: got [%v], want [%v]",
					genericError.Reason,
					ncreason)
			}
		}
	} else {
		rr.Body.Bytes()
	}
}

func runConfigSet(config oas.Config) *httptest.ResponseRecorder {
	req := newConfigPatchReq(config)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	return rr
}

func runConfigGet(t *testing.T) oas.Config {
	req := newConfigGetReq()
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	return parseSampleConfig(t, string(rr.Body.Bytes()))
}

/* Tests are divided by OK and EC where
 * - OK are checks on a correct action
 * - EC are checks on faulty behavior/requests (edge cases)
 */

//Test001 - EC-001 Active-Backup Bond Without ActiveSlave
func Test001(t *testing.T) {
	c := genSampleConfig(t)
	*(*c.HostNetwork.Links)[2].Linkinfo.InfoSlaveData.State = "BACKUP"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Active Slave Iface not found for Active-Backup type bond bond0")
}

//Test002 - EC-002 Active-Backup Bond With Multiple Active Slaves
func Test002(t *testing.T) {
	c := genSampleConfig(t)
	*(*c.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.State = "ACTIVE"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Multiple Active Slave Ifaces found for Active-Backup type bond bond0")
}

//Test003 - EC-003 Non Active-Backup Bond With Backup Slave
func Test003(t *testing.T) {
	c := genSampleConfig(t)
	*(*c.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-rr"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Backup Slave Iface dummy0 found for non Active-Backup type bond bond0")
}

func linkStateMatch(setLinkData oas.Link, getLinkData oas.Link) bool {
	//Check for up request to correspond to up interface
	upIsUp := false
	lfs := setLinkData.GetFlags()
	if lfs != nil {
		for _, lf := range lfs {
			if lf == "UP" {
				rfs := getLinkData.GetFlags()
				for _, rf := range rfs {
					if rf == "UP" {
						upIsUp = true
					}
				}
				if upIsUp == false {
					return false
				}

				// Dummy interfaces report unknown operstate
				// https://serverfault.com/questions/629676/dummy-network-interface-in-linux
				// which according to  the kernel doc
				// https://www.kernel.org/doc/Documentation/networking/operstates.txt
				// just tells that the setting of the operational state was not implemented by
				// the below driver (can be a bug / lack of compliance)
				if getLinkData.Linkinfo != nil && *getLinkData.Linkinfo.InfoKind == "dummy" {
					//Let's also check operstate:
					operstate := getLinkData.GetOperstate()
					if operstate != "up" {
						return false
					}
				}
			}
		}
	}
	return true
}

func linkMatch(_setLinkData interface{}, _getLinkData interface{}) bool {
	setLinkData := _setLinkData.(oas.Link)
	getLinkData := _getLinkData.(oas.Link)
	if setLinkData.GetMaster() != getLinkData.GetMaster() {
		return false
	}
	if setLinkData.LinkType != getLinkData.LinkType {
		return false
	}
	lli := setLinkData.GetLinkinfo()
	rli := getLinkData.GetLinkinfo()
	lid := lli.GetInfoData()
	rid := rli.GetInfoData()
	lisd := lli.GetInfoSlaveData()
	risd := rli.GetInfoSlaveData()
	if lisd.GetState() != "" && lisd.GetState() != risd.GetState() {
		if lisd.GetState() == "ACTIVE" && risd.GetState() == "BACKUP" {
			//Accept this, it's just the bond not beeing up (in which case)
		} else {
			return false
		}
	}

	if lid.GetMode() != rid.GetMode() {
		return false
	}
	if lid.GetMiimon() != -1 && lid.GetMiimon() != rid.GetMiimon() {
		return false
	}
	if lid.GetDowndelay() != -1 && lid.GetDowndelay() != rid.GetDowndelay() {
		return false
	}
	if lid.GetUpdelay() != -1 && lid.GetUpdelay() != rid.GetUpdelay() {
		return false
	}
	if lid.GetXmitHashPolicy() != "" && (lid.GetXmitHashPolicy() != rid.GetXmitHashPolicy()) {
		return false
	}
	if lid.GetAdLacpRate() != "" && lid.GetAdLacpRate() != rid.GetAdLacpRate() {
		return false
	}
	if lid.GetPeerNotifyDelay() != -1 && lid.GetPeerNotifyDelay() != rid.GetPeerNotifyDelay() {
		return false
	}
	if lid.GetUseCarrier() != -1 && lid.GetUseCarrier() != rid.GetUseCarrier() {
		return false
	}
	if lid.GetLpInterval() != -1 && lid.GetLpInterval() != rid.GetLpInterval() {
		return false
	}
	if lid.GetArpAllTargets() != "" && lid.GetArpAllTargets() != rid.GetArpAllTargets() {
		return false
	}
	if lid.GetXmitHashPolicy() != "" && lid.GetXmitHashPolicy() != rid.GetXmitHashPolicy() {
		return false
	}
	if lid.GetResendIgmp() != -1 && lid.GetResendIgmp() != rid.GetResendIgmp() {
		return false
	}
	if lid.GetMinLinks() != -1 && lid.GetMinLinks() != rid.GetMinLinks() {
		return false
	}
	if lid.GetPrimaryReselect() != "" && lid.GetPrimaryReselect() != rid.GetPrimaryReselect() {
		return false
	}
	if lid.GetAdSelect() != "" && lid.GetAdSelect() != rid.GetAdSelect() {
		return false
	}
	if lid.GetAllSlavesActive() != -1 && lid.GetAllSlavesActive() != rid.GetAllSlavesActive() {
		return false
	}

	return linkStateMatch(setLinkData, getLinkData)
}

//Test004 - OK-004 Bond Active-Backup params check
func Test004(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}

}

//Test005 - OK-005 Bond Balance-RR Xmit Hash Policy params check
func Test005(t *testing.T) {
	cset := genSampleConfig(t)
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-rr"
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.SetState("ACTIVE")
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test006 - OK-006 Bond 802.3ad mix
func Test006(t *testing.T) {
	cset := genSampleConfig(t)
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "802.3ad"
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdLacpRate("fast")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPeerNotifyDelay(2000)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUseCarrier(0)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpInterval(500)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpValidate("backup")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetLpInterval(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpAllTargets("all")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPacketsPerSlave(2)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetFailOverMac()
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetResendIgmp(3)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMinLinks(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPrimaryReselect("better")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetTlbDynamicLb(1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdSelect("bandwidth")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAllSlavesActive(1)
	(*cset.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.SetState("ACTIVE")
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test007 - OK-007 Bond Balance-RR Mix
func Test007(t *testing.T) {
	cset := genSampleConfig(t)
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-rr"
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdLacpRate("fast")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPeerNotifyDelay(2000)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUseCarrier(0)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpInterval(500)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpValidate("backup")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetLpInterval(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpAllTargets("all")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPacketsPerSlave(2)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetFailOverMac()
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetResendIgmp(3)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMinLinks(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPrimaryReselect("better")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetTlbDynamicLb(1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdSelect("bandwidth")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAllSlavesActive(1)
	(*cset.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.SetState("ACTIVE")
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test008 - OK-008 Bond Balance-TLB
func Test008(t *testing.T) {
	cset := genSampleConfig(t)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMiimon(-1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUpdelay(-1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetDowndelay(-1)
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-tlb"
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdLacpRate("fast")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPeerNotifyDelay(2000)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUseCarrier(0)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpInterval(500)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpValidate("backup")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetLpInterval(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpAllTargets("all")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPacketsPerSlave(2)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetFailOverMac()
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetResendIgmp(3)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMinLinks(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPrimaryReselect("better")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetTlbDynamicLb(0)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdSelect("bandwidth")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAllSlavesActive(1)
	(*cset.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.SetState("ACTIVE")
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test009 - OK-009 Bond Active-Backup mix
func Test009(t *testing.T) {
	cset := genSampleConfig(t)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMiimon(-1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUpdelay(-1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetDowndelay(-1)
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "active-backup"
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdLacpRate("fast")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPeerNotifyDelay(2000)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetUseCarrier(0)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpInterval(500)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpValidate("backup")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetLpInterval(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpAllTargets("all")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPacketsPerSlave(2)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetFailOverMac()
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetResendIgmp(3)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetMinLinks(2)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPrimaryReselect("better")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetTlbDynamicLb(1)
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAdSelect("bandwidth")
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetAllSlavesActive(1)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test010 - OK-010 Up/Down flag and operstate
func Test010(t *testing.T) {
	cset := genSampleConfig(t)
	(*cset.HostNetwork.Links)[0].Flags = nil
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)
	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

var sampleRouteConfig string = `
{
  "__id": "498b44c3999f2edfa715123748696ad8",
  "dev": "dummy0",
  "dst": "10.1.2.0/24",
  "gateway": "10.1.2.3",
  "metric": 50,
  "protocol": "boot",
  "scope": "universe"
}
`

func runRequest(method string, uri string, body string) *httptest.ResponseRecorder {
	iobody := bytes.NewReader([]byte(body))
	req, _ := http.NewRequest(method, uri, iobody)
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	return rr
}

//Test011 - EC-011 Route network not found
func Test011(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	rr = runRequest("POST", "/api/1/routes", sampleRouteConfig)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC,
		`Got ENETUNREACH error: network is not reachable for route {498b44c3999f2edfa715123748696ad8 {[10 1 2 0] 24} 10.1.2.3 dummy0  50  <nil> <nil>}`)
}

//Test012 - EC-012 Link not found for route to create
func Test012(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	rr = runRequest("DELETE", "/api/1/links/dummy0", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	rr = runRequest("POST", "/api/1/routes", sampleRouteConfig)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC,
		`Route 498b44c3999f2edfa715123748696ad8 Link Device dummy0 not found`)
}

func checkBody(t *testing.T, rr *httptest.ResponseRecorder, body string) {
	res, err := comm.JSONBytesEqual(rr.Body.Bytes(), []byte(body))
	if err != nil {
		t.Error(err)
	}
	if res != true {
		t.Errorf("Body mismatch: got [%v], want [%v]",
			string(rr.Body.Bytes()),
			body)
	}
}

//Test013 - EC-013 Route Creation + Route Check + Route Exists
func Test013(t *testing.T) {
	cset := genSampleConfig(t)
	lai := oas.LinkAddrInfo{
		Local:     net.IPv4(10, 1, 2, 3),
		Prefixlen: 24,
	}

	*(*cset.HostNetwork.Links)[1].AddrInfo = append(*(*cset.HostNetwork.Links)[1].AddrInfo, lai)

	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	rr = runRequest("DELETE", "/api/1/routes/498b44c3999f2edfa715123748696ad8", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	rr = runRequest("POST", "/api/1/routes", sampleRouteConfig)
	checkResponse(t, rr, http.StatusCreated, nc.RESERVED, "")
	checkBody(t, rr, `"498b44c3999f2edfa715123748696ad8"`)
	rr = runRequest("POST", "/api/1/routes", sampleRouteConfig)
	checkResponse(t, rr, http.StatusConflict, nc.CONFLICT,
		`Route 498b44c3999f2edfa715123748696ad8 exists`)
}

func routeMatch(_setRoute interface{}, _getRoute interface{}) bool {
	setRoute := _setRoute.(oas.Route)
	getRoute := _getRoute.(oas.Route)

	if setRoute.Dev != nil && getRoute.Dev != nil {
		if *setRoute.Dev != *getRoute.Dev {
			return false
		}
	}

	if (setRoute.Gateway != nil && getRoute.Gateway == nil) ||
		(setRoute.Gateway == nil && getRoute.Gateway != nil) {
		return false
	}
	if setRoute.Gateway != nil && getRoute.Gateway != nil {
		if setRoute.Gateway.Equal(*getRoute.Gateway) == false {
			return false
		}
	}

	return true
}

//Test014 - OK-014 Batch Link + Route config
func Test14(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet(t)

	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
	if delta := comm.ListCompare(*cset.HostNetwork.Routes, *cget.HostNetwork.Routes, routeMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test015 - OK-015 Route without dev
func Test15(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

	newroute := `
{
	"dst": "10.102.2.0/24",
	"gateway": "10.6.7.8",
	"metric": 50
}`
	rr = runRequest("POST", "/api/1/routes", newroute)
	checkResponse(t, rr, http.StatusCreated, nc.RESERVED, "")
	cget := runConfigGet(t)

	if delta := comm.ListCompare(*cset.HostNetwork.Links, *cget.HostNetwork.Links, linkMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
	if delta := comm.ListCompare(*cset.HostNetwork.Routes, *cget.HostNetwork.Routes, routeMatch); delta != nil {
		t.Errorf("Mismatch on %v", delta)
	}
}

//Test016 - OK-016 MTU Set/Get
func Test16(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

	//Cleanup
	runRequest("DELETE", "/api/1/links/dummy3", "")

	newlink := `{
			"ifname": "dummy3",
			"link_type": "ether",
			"flags": ["up"],
			"linkinfo": {
				"info_kind": "dummy",
				"info_slave_data": {
					"state": "BACKUP"
				}
			},
			"mtu": 3000,
			"addr_info": [
				{
					"local": "10.6.7.8",
					"prefixlen": 24
				}
			],
			"master": "bond0"
		}`
	rr = runRequest("POST", "/api/1/links", newlink)
	checkResponse(t, rr, http.StatusCreated, nc.RESERVED, "")

	rr = runRequest("GET", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	var l oas.Link
	err := json.Unmarshal(rr.Body.Bytes(), &l)
	if err != nil {
		t.Errorf(err.Error())
	}
	if l.Mtu == nil {
		t.Errorf("MTU was not found in LinkGet")
	} else {
		if *l.Mtu != 3000 {
			t.Errorf("MTU was not properly configured")
		}
	}

	rr = runRequest("DELETE", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
}

//Test017 - OK-017 Checks that single POSTED interface is effectively enslave to bond
func Test17(t *testing.T) {
	cset := genSampleConfig(t)
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

	//Cleanup
	runRequest("DELETE", "/api/1/links/dummy3", "")

	newlink := `{
			"ifname": "dummy3",
			"link_type": "ether",
			"flags": ["up"],
			"linkinfo": {
				"info_kind": "dummy",
				"info_slave_data": {
					"state": "BACKUP"
				}
			},
			"mtu": 3000,
			"addr_info": [
				{
					"local": "10.6.7.8",
					"prefixlen": 24
				}
			],
			"master": "bond0"
		}`
	rr = runRequest("POST", "/api/1/links", newlink)
	checkResponse(t, rr, http.StatusCreated, nc.RESERVED, "")

	rr = runRequest("GET", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	var l oas.Link
	err := json.Unmarshal(rr.Body.Bytes(), &l)
	if err != nil {
		t.Errorf(err.Error())
	}
	if l.Linkinfo.InfoSlaveData.MiiStatus == nil {
		t.Errorf("Interface is not enslaved (missing slave data)")
	} else if *l.Linkinfo.InfoSlaveData.MiiStatus != "UP" &&
		*l.Linkinfo.InfoSlaveData.MiiStatus != "GOING_BACK" {
		t.Errorf(fmt.Sprintf("Interface Slave MiiStatus is not up/going up: %v", *l.Linkinfo.InfoSlaveData.MiiStatus))
	}

	rr = runRequest("DELETE", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

}

//Test018 - OK-018 Check overlapping route does not stop config PATCH
func Test18(t *testing.T) {
	cfg := `{
			"global": {},
			"host_network": {
				"links": [
					{
						"ifname": "dummy3",
						"link_type": "ether",
						"flags": ["up"],
						"linkinfo": {
							"info_kind": "dummy"
						},
						"mtu": 3000,
						"addr_info": [
							{
								"local": "10.16.7.8",
								"prefixlen": 24
							}
						]
					}
				],
				"routes":[
					{
						"dev": "dummy3",
						"dst": "10.16.7.0/24",
						"metric": 0,
						"prefsrc": "10.16.7.8",
						"protocol": "kernel",
						"scope": "link"
					},
					{
						"dev": "dummy3",
						"dst": "fe80::/64",
						"metric": 256,
						"protocol": "kernel",
						"scope": "global"
					}
				]
			}
		}`
	rr := runRequest("PATCH", "/api/1/mgmt/config", cfg)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

	rr = runRequest("GET", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	var l oas.Link
	err := json.Unmarshal(rr.Body.Bytes(), &l)
	if err != nil {
		t.Errorf(err.Error())
	}

	rr = runRequest("DELETE", "/api/1/links/dummy3", "")
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")

}
