package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/nc"
)

func TestActiveBackupBondWithNonActiveSlave(t *testing.T) {
	reqbody := `{
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
	  }`
	iobody := strings.NewReader(reqbody)
	m := NewManager()
	req, err := http.NewRequest("PUT", "/api/1/config", iobody)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	t.Errorf("%v", string(rr.Body.Bytes()))
	var genericError nc.GenericError
	err = json.Unmarshal(rr.Body.Bytes(), &genericError)
	if err != nil {
		t.Errorf(err.Error())
	}
}
