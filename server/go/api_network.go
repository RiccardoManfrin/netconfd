/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.1.0
 * Contact: support@athonet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A NetworkApiController binds http requests to an api service and writes the service results to the http response
type NetworkApiController struct {
	service NetworkApiServicer
}

// NewNetworkApiController creates a default api controller
func NewNetworkApiController(s NetworkApiServicer) Router {
	return &NetworkApiController{ service: s }
}

// Routes returns all of the api route for the NetworkApiController
func (c *NetworkApiController) Routes() Routes {
	return Routes{ 
		{
			"ConfigGet",
			strings.ToUpper("Get"),
			"/api/1/config",
			c.ConfigGet,
		},
		{
			"ConfigLinkDel",
			strings.ToUpper("Delete"),
			"/api/1/config/links/{ifname}",
			c.ConfigLinkDel,
		},
		{
			"ConfigLinkGet",
			strings.ToUpper("Get"),
			"/api/1/config/links/{ifname}",
			c.ConfigLinkGet,
		},
		{
			"ConfigLinkSet",
			strings.ToUpper("Post"),
			"/api/1/config/links",
			c.ConfigLinkSet,
		},
		{
			"ConfigNetNSDel",
			strings.ToUpper("Delete"),
			"/api/1/config/netns/{netnsid}",
			c.ConfigNetNSDel,
		},
		{
			"ConfigNetNSGet",
			strings.ToUpper("Get"),
			"/api/1/config/netns/{netnsid}",
			c.ConfigNetNSGet,
		},
		{
			"ConfigNetNSSet",
			strings.ToUpper("Post"),
			"/api/1/config/netns",
			c.ConfigNetNSSet,
		},
		{
			"ConfigRouteDel",
			strings.ToUpper("Delete"),
			"/api/1/config/routes/{routeid}",
			c.ConfigRouteDel,
		},
		{
			"ConfigRouteGet",
			strings.ToUpper("Get"),
			"/api/1/config/routes/{routeid}",
			c.ConfigRouteGet,
		},
		{
			"ConfigRouteSet",
			strings.ToUpper("Post"),
			"/api/1/config/routes",
			c.ConfigRouteSet,
		},
		{
			"ConfigRuleDel",
			strings.ToUpper("Delete"),
			"/api/1/config/rules/{ruleid}",
			c.ConfigRuleDel,
		},
		{
			"ConfigRuleGet",
			strings.ToUpper("Get"),
			"/api/1/config/rules/{ruleid}",
			c.ConfigRuleGet,
		},
		{
			"ConfigRuleSet",
			strings.ToUpper("Post"),
			"/api/1/config/rules",
			c.ConfigRuleSet,
		},
		{
			"ConfigSet",
			strings.ToUpper("Put"),
			"/api/1/config",
			c.ConfigSet,
		},
		{
			"ConfigVRFDel",
			strings.ToUpper("Delete"),
			"/api/1/config/vrfs/{vrfid}",
			c.ConfigVRFDel,
		},
		{
			"ConfigVRFGet",
			strings.ToUpper("Get"),
			"/api/1/config/vrfs/{vrfid}",
			c.ConfigVRFGet,
		},
		{
			"ConfigVRFSet",
			strings.ToUpper("Post"),
			"/api/1/config/vrfs",
			c.ConfigVRFSet,
		},
	}
}

// ConfigGet - Configures and enforces a new live network configuration 
func (c *NetworkApiController) ConfigGet(w http.ResponseWriter, r *http.Request) { 
	result, err := c.service.ConfigGet(r.Context())
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigLinkDel - Brings down and delete a link layer interface 
func (c *NetworkApiController) ConfigLinkDel(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	ifname := params["ifname"]
	result, err := c.service.ConfigLinkDel(r.Context(), ifname)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigLinkGet - Retrieve link layer interface information 
func (c *NetworkApiController) ConfigLinkGet(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	ifname := params["ifname"]
	result, err := c.service.ConfigLinkGet(r.Context(), ifname)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigLinkSet - Configures and brings up a link layer interface 
func (c *NetworkApiController) ConfigLinkSet(w http.ResponseWriter, r *http.Request) { 
	link := &Link{}
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigLinkSet(r.Context(), *link)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigNetNSDel - Removes an IP Rule 
func (c *NetworkApiController) ConfigNetNSDel(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	netnsid := params["netnsid"]
	result, err := c.service.ConfigNetNSDel(r.Context(), netnsid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigNetNSGet - Get a network namespace 
func (c *NetworkApiController) ConfigNetNSGet(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	netnsid := params["netnsid"]
	result, err := c.service.ConfigNetNSGet(r.Context(), netnsid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigNetNSSet - Configures an new Network Namespace 
func (c *NetworkApiController) ConfigNetNSSet(w http.ResponseWriter, r *http.Request) { 
	netns := &Netns{}
	if err := json.NewDecoder(r.Body).Decode(&netns); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigNetNSSet(r.Context(), *netns)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRouteDel - Brings down and delete an L3 IP route 
func (c *NetworkApiController) ConfigRouteDel(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	routeid := params["routeid"]
	result, err := c.service.ConfigRouteDel(r.Context(), routeid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRouteGet - Get a L3 route details 
func (c *NetworkApiController) ConfigRouteGet(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	routeid := params["routeid"]
	result, err := c.service.ConfigRouteGet(r.Context(), routeid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRouteSet - Configures a route 
func (c *NetworkApiController) ConfigRouteSet(w http.ResponseWriter, r *http.Request) { 
	route := &Route{}
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigRouteSet(r.Context(), *route)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRuleDel - Removes an IP Rule 
func (c *NetworkApiController) ConfigRuleDel(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	ruleid := params["ruleid"]
	result, err := c.service.ConfigRuleDel(r.Context(), ruleid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRuleGet - Get an IP rule details 
func (c *NetworkApiController) ConfigRuleGet(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	ruleid := params["ruleid"]
	result, err := c.service.ConfigRuleGet(r.Context(), ruleid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigRuleSet - Configures an IP rule 
func (c *NetworkApiController) ConfigRuleSet(w http.ResponseWriter, r *http.Request) { 
	body := &map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigRuleSet(r.Context(), *body)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigSet - Configures and enforces a new live network configuration 
func (c *NetworkApiController) ConfigSet(w http.ResponseWriter, r *http.Request) { 
	config := &Config{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigSet(r.Context(), *config)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigVRFDel - Removes a VRF 
func (c *NetworkApiController) ConfigVRFDel(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	vrfid := params["vrfid"]
	result, err := c.service.ConfigVRFDel(r.Context(), vrfid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigVRFGet - Get a VRF 
func (c *NetworkApiController) ConfigVRFGet(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	vrfid := params["vrfid"]
	result, err := c.service.ConfigVRFGet(r.Context(), vrfid)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// ConfigVRFSet - Configures an new VRF 
func (c *NetworkApiController) ConfigVRFSet(w http.ResponseWriter, r *http.Request) { 
	body := &map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.ConfigVRFSet(r.Context(), *body)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}
