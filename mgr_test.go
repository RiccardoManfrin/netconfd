package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.lan.athonet.com/core/netconfd/nc"
	oas "gitlab.lan.athonet.com/core/netconfd/server/go"
)

func parseSampleConfig(sampleConfig string) oas.Config {
	var config oas.Config
	json.Unmarshal([]byte(sampleConfig), &config)
	return config
}

func genSampleConfig() oas.Config {
	sampleConfig := `{
		"global": {},
		"host_network": {
		  "links": [
			{
			  "ifname": "bond0",
			  "link_type": "ether",
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
			  "linkinfo": {
				"info_kind": "dummy",
				"info_slave_data": {
				  "state": "BACKUP"
				}
			  },
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
		  "routes": []
		}
	  }`
	return parseSampleConfig(sampleConfig)
}

func newConfigSetReq(config oas.Config) *http.Request {
	reqbody, _ := json.Marshal(config)
	iobody := bytes.NewReader(reqbody)
	req, _ := http.NewRequest("PUT", "/api/1/config", iobody)
	req.Header.Add("Content-Type", "application/json")
	return req
}
func newConfigGetReq() *http.Request {
	req, _ := http.NewRequest("GET", "/api/1/config", nil)
	req.Header.Add("Content-Type", "application/json")
	return req
}

var m *Manager = NewManager()

func checkResponse(t *testing.T, rr *httptest.ResponseRecorder, httpStatusCode int, ncErrorCode nc.ErrorCode, ncreason string) {
	if status := rr.Code; status != httpStatusCode {
		t.Errorf("HTTP Status code mismatch: got [%v] want [%v]",
			status,
			httpStatusCode)
	}
	var genericError nc.GenericError
	err := json.Unmarshal(rr.Body.Bytes(), &genericError)
	if err != nil {
		t.Errorf("Err Unmarshal failure")
	}
	if ncErrorCode != nc.RESERVED {
		if genericError.Code != ncErrorCode {
			t.Errorf("Err Code mismatch: got [%v], want [%v]",
				genericError.Code,
				nc.SEMANTIC)
		}
	}
	if ncreason != "" {
		if genericError.Reason != ncreason {
			t.Errorf("Err Reason mismatch: got [%v], want [%v]",
				genericError.Reason,
				ncreason)
		}
	}
}

func runConfigSet(config oas.Config) *httptest.ResponseRecorder {
	req := newConfigSetReq(config)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	return rr
}

func runConfigGet() oas.Config {
	req := newConfigGetReq()
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	return parseSampleConfig(string(rr.Body.Bytes()))
}

/* Tests are divided by OK and EC where
 * - OK are checks on a correct action
 * - EC are checks on faulty behavior/requests (edge cases)
 */

func listToMap(slice interface{}, key string) map[string]interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	if s.IsNil() {
		return nil
	}
	trueSlice := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		trueSlice[i] = s.Index(i).Interface()
	}

	mappedList := make(map[string]interface{})

	for _, l := range trueSlice {
		val := reflect.ValueOf(l)
		kval := val.FieldByName(key).String()
		mappedList[kval] = l
	}
	return mappedList
}

//Test001 - EC-001 Active-Backup Bond Without ActiveSlave
func Test001(t *testing.T) {
	c := genSampleConfig()
	*(*c.HostNetwork.Links)[2].Linkinfo.InfoSlaveData.State = "BACKUP"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Active Slave Iface not found for Active-Backup type bond bond0")
}

//Test002 - EC-002 Active-Backup Bond With Multiple Active Slaves
func Test002(t *testing.T) {
	c := genSampleConfig()
	*(*c.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.State = "ACTIVE"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Multiple Active Slave Ifaces found for Active-Backup type bond bond0")
}

//Test003 - EC-003 Non Active-Backup Bond With Backup Slave
func Test003(t *testing.T) {
	c := genSampleConfig()
	*(*c.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-rr"
	rr := runConfigSet(c)
	checkResponse(t, rr, http.StatusBadRequest, nc.SEMANTIC, "Backup Slave Iface dummy0 found for non Active-Backup type bond bond0")
}

func deltaLink(l oas.Link, r oas.Link) string {
	if l.GetMaster() != r.GetMaster() {
		return "master"
	}
	if l.LinkType != r.LinkType {
		return "link_type"
	}
	lli := l.GetLinkinfo()
	rli := r.GetLinkinfo()
	lid := lli.GetInfoData()
	rid := rli.GetInfoData()
	if lid.GetMode() != rid.GetMode() {
		return "link_info->info_data->mode"
	}
	if lid.GetDowndelay() != rid.GetDowndelay() {
		return "link_info->info_data->downdelay"
	}
	if lid.GetUpdelay() != rid.GetUpdelay() {
		return "link_info->info_data->updelay"
	}
	if lid.GetXmitHashPolicy() != "" && (lid.GetXmitHashPolicy() != rid.GetXmitHashPolicy()) {
		return "link_info->info_data->xmit_hash_policy"
	}
	if lid.GetAdLacpRate() != "" && (lid.GetAdLacpRate() != rid.GetAdLacpRate()) {
		return "link_info->info_data->ad_lacp_rate"
	}

	if lid.GetAdLacpRate() != "" && lid.GetAdLacpRate() != rid.GetAdLacpRate() {
		return fmt.Sprintf("link_info->info_data->AdLacpRate: l:[%v], r:[%v]", lid.GetAdLacpRate(), rid.GetAdLacpRate())
	}
	if lid.GetPeerNotifyDelay() != -1 && lid.GetPeerNotifyDelay() != rid.GetPeerNotifyDelay() {
		return fmt.Sprintf("link_info->info_data->PeerNotifyDelay: l:[%v], r:[%v]", lid.GetPeerNotifyDelay(), rid.GetPeerNotifyDelay())
	}
	if lid.GetUseCarrier() != -1 && lid.GetUseCarrier() != rid.GetUseCarrier() {
		return fmt.Sprintf("link_info->info_data->UseCarrier: l:[%v], r:[%v]", lid.GetUseCarrier(), rid.GetUseCarrier())
	}
	if lid.GetLpInterval() != -1 && lid.GetLpInterval() != rid.GetLpInterval() {
		return fmt.Sprintf("link_info->info_data->LpInterval: l:[%v], r:[%v]", lid.GetLpInterval(), rid.GetLpInterval())
	}
	if lid.GetArpAllTargets() != "" && lid.GetArpAllTargets() != rid.GetArpAllTargets() {
		return fmt.Sprintf("link_info->info_data->ArpAllTargets: l:[%v], r:[%v]", lid.GetArpAllTargets(), rid.GetArpAllTargets())
	}
	if lid.GetXmitHashPolicy() != "" && lid.GetXmitHashPolicy() != rid.GetXmitHashPolicy() {
		return fmt.Sprintf("link_info->info_data->XmitHashPolicy: l:[%v], r:[%v]", lid.GetXmitHashPolicy(), rid.GetXmitHashPolicy())
	}
	if lid.GetResendIgmp() != -1 && lid.GetResendIgmp() != rid.GetResendIgmp() {
		return fmt.Sprintf("link_info->info_data->ResendIgmp: l:[%v], r:[%v]", lid.GetResendIgmp(), rid.GetResendIgmp())
	}
	if lid.GetMinLinks() != -1 && lid.GetMinLinks() != rid.GetMinLinks() {
		return fmt.Sprintf("link_info->info_data->MinLinks: l:[%v], r:[%v]", lid.GetMinLinks(), rid.GetMinLinks())
	}
	if lid.GetPrimaryReselect() != "" && lid.GetPrimaryReselect() != rid.GetPrimaryReselect() {
		return fmt.Sprintf("link_info->info_data->PrimaryReselect: l:[%v], r:[%v]", lid.GetPrimaryReselect(), rid.GetPrimaryReselect())
	}
	if lid.GetAdSelect() != "" && lid.GetAdSelect() != rid.GetAdSelect() {
		return fmt.Sprintf("link_info->info_data->AdSelect: l:[%v], r:[%v]", lid.GetAdSelect(), rid.GetAdSelect())
	}
	if lid.GetAllSlavesActive() != -1 && lid.GetAllSlavesActive() != rid.GetAllSlavesActive() {
		return fmt.Sprintf("link_info->info_data->AllSlavesActive: l:[%v], r:[%v]", lid.GetAllSlavesActive(), rid.GetAllSlavesActive())
	}

	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpInterval(500)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetArpValidate("backup")
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetPacketsPerSlave(2)
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetFailOverMac()
	//(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetTlbDynamicLb(1)

	return ""
}

//Test004 - OK-004 Bond Active-Backup params check
func Test004(t *testing.T) {
	cset := genSampleConfig()
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet()
	cLinksSetMap := listToMap(*cset.HostNetwork.Links, "Ifname")
	cLinksGetMap := listToMap(*cget.HostNetwork.Links, "Ifname")
	for ifname, setLink := range cLinksSetMap {
		getLink := cLinksGetMap[ifname].(oas.Link)
		if delta := deltaLink(setLink.(oas.Link), getLink); delta != "" {
			t.Errorf("Mismatch on %v", delta)
		}
	}
}

//Test005 - OK-005 Bond Balance-RR Xmit Hash Policy params check
func Test005(t *testing.T) {
	cset := genSampleConfig()
	*(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.Mode = "balance-rr"
	(*cset.HostNetwork.Links)[0].Linkinfo.InfoData.SetXmitHashPolicy("layer2+3")
	(*cset.HostNetwork.Links)[1].Linkinfo.InfoSlaveData.SetState("ACTIVE")
	rr := runConfigSet(cset)
	checkResponse(t, rr, http.StatusOK, nc.RESERVED, "")
	cget := runConfigGet()
	cLinksSetMap := listToMap(*cset.HostNetwork.Links, "Ifname")
	cLinksGetMap := listToMap(*cget.HostNetwork.Links, "Ifname")
	for ifname, setLink := range cLinksSetMap {
		getLink := cLinksGetMap[ifname].(oas.Link)
		if delta := deltaLink(setLink.(oas.Link), getLink); delta != "" {
			t.Errorf("Mismatch on %v", delta)
		}
	}
}

//Test006 - OK-006 Bond Balance-RR Xmit Hash Policy params check
func Test006(t *testing.T) {
	cset := genSampleConfig()
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
	cget := runConfigGet()
	cLinksSetMap := listToMap(*cset.HostNetwork.Links, "Ifname")
	cLinksGetMap := listToMap(*cget.HostNetwork.Links, "Ifname")
	for ifname, setLink := range cLinksSetMap {
		getLink := cLinksGetMap[ifname].(oas.Link)
		if delta := deltaLink(setLink.(oas.Link), getLink); delta != "" {
			t.Errorf("Mismatch on %v", delta)
		}
	}
}
