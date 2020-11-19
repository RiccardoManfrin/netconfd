package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/nc"
	oas "gitlab.lan.athonet.com/riccardo.manfrin/netconfd/server/go"
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
				  "miimon" : 200,
				  "xmit_hash_policy" : "layer2+3",
				  "ad_lacp_rate" : "fast"
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
		t.Errorf("HTTP Status code mismatch: got %v want %v",
			status,
			http.StatusBadRequest)
	}
	var genericError nc.GenericError
	err := json.Unmarshal(rr.Body.Bytes(), &genericError)
	if err != nil {
		t.Errorf("Err Unmarshal failure")
	}
	if ncErrorCode != nc.RESERVED {
		if genericError.Code != ncErrorCode {
			t.Errorf("Err Code mismatch: got %v, want %v",
				genericError.Code,
				nc.SEMANTIC)
		}
	}
	if ncreason != "" {
		if genericError.Reason != ncreason {
			t.Errorf("Err Reason mismatch: got %v, want %v",
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
	if lid.GetXmitHashPolicy() != rid.GetXmitHashPolicy() {
		return "link_info->info_data->xmit_hash_policy"
	}
	if lid.GetAdLacpRate() != rid.GetAdLacpRate() {
		return "link_info->info_data->ad_lacp_rate"
	}
	return ""

}

//Test004 - OK-005 Bond params check
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
