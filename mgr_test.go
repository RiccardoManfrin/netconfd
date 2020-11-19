package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/nc"
)

func newConfigPutReq(reqbody string) *http.Request {
	iobody := strings.NewReader(reqbody)
	req, _ := http.NewRequest("PUT", "/api/1/config", iobody)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func TestConfigPutActiveBackupBondWithNonActiveSlave(t *testing.T) {
	req := newConfigPutReq(`{
		"global": {},
		"host_network": {
		  "links": [
			{
			  "ifname": "bond0",
			  "link_type": "ether",
			  "linkinfo": {
				"info_kind": "bond",
				"info_data": {
				  "mode": "active-backup"
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
				  "state": "BACKUP"
				}
			  },
			  "master": "bond0"
			}
		  ],
		  "routes": []
		}
	  }`)

	m := NewManager()
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HTTP Status code mismatch: got %v want %v",
			status,
			http.StatusBadRequest)
	}
	var genericError nc.GenericError
	err := json.Unmarshal(rr.Body.Bytes(), &genericError)
	if err != nil {
		t.Errorf("Err Unmarshal failure")
	}
	if genericError.Code != nc.SEMANTIC {
		t.Errorf("Err Code mismatch: got %v, want %v",
			genericError.Code,
			nc.SEMANTIC)
	}
}
