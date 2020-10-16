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

type Link struct {

	// Inteface index ID 
	Ifindex int32 `json:"ifindex,omitempty"`

	// Interface name 
	Ifname string `json:"ifname"`

	// Maximum Transfer Unit value 
	Mtu int32 `json:"mtu,omitempty"`

	Linkinfo LinkLinkinfo `json:"linkinfo,omitempty"`

	LinkType string `json:"link_type"`

	Address string `json:"address,omitempty"`

	AddrInfo []LinkAddrInfo `json:"addr_info,omitempty"`
}


