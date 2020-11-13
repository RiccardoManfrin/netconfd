/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.1.0
 * Contact: support@athonet.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// Link struct for Link
type Link struct {
	// Interface name 
	Id *string `json:"__id,omitempty"`
	// Inteface index ID 
	Ifindex *int32 `json:"ifindex,omitempty"`
	// Interface name 
	Ifname string `json:"ifname"`
	// Flags of the interface Supported types:   * `BROADCAST` - Support for broadcast   * `MULTICAST` - Support for multicast   * `SLAVE` - Is slave   * `UP` - Is up   * `LOWER UP` - Is lower interface up 
	Flags *[]string `json:"flags,omitempty"`
	// Maximum Transfer Unit value 
	Mtu *int32 `json:"mtu,omitempty"`
	// Promiscuous mode flag
	Promiscuity *int32 `json:"promiscuity,omitempty"`
	// In case the interface is part of a bond or bridge, specifies the bond/bridge interface it belongs to. 
	Master *string `json:"master,omitempty"`
	Linkinfo *LinkLinkinfo `json:"linkinfo,omitempty"`
	LinkType string `json:"link_type"`
	Address *string `json:"address,omitempty"`
	AddrInfo *[]LinkAddrInfo `json:"addr_info,omitempty"`
}

// NewLink instantiates a new Link object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLink(ifname string, linkType string, ) *Link {
	this := Link{}
	this.Ifname = ifname
	this.LinkType = linkType
	return &this
}

// NewLinkWithDefaults instantiates a new Link object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinkWithDefaults() *Link {
	this := Link{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Link) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Link) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Link) SetId(v string) {
	o.Id = &v
}

// GetIfindex returns the Ifindex field value if set, zero value otherwise.
func (o *Link) GetIfindex() int32 {
	if o == nil || o.Ifindex == nil {
		var ret int32
		return ret
	}
	return *o.Ifindex
}

// GetIfindexOk returns a tuple with the Ifindex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetIfindexOk() (*int32, bool) {
	if o == nil || o.Ifindex == nil {
		return nil, false
	}
	return o.Ifindex, true
}

// HasIfindex returns a boolean if a field has been set.
func (o *Link) HasIfindex() bool {
	if o != nil && o.Ifindex != nil {
		return true
	}

	return false
}

// SetIfindex gets a reference to the given int32 and assigns it to the Ifindex field.
func (o *Link) SetIfindex(v int32) {
	o.Ifindex = &v
}

// GetIfname returns the Ifname field value
func (o *Link) GetIfname() string {
	if o == nil  {
		var ret string
		return ret
	}

	return o.Ifname
}

// GetIfnameOk returns a tuple with the Ifname field value
// and a boolean to check if the value has been set.
func (o *Link) GetIfnameOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Ifname, true
}

// SetIfname sets field value
func (o *Link) SetIfname(v string) {
	o.Ifname = v
}

// GetFlags returns the Flags field value if set, zero value otherwise.
func (o *Link) GetFlags() []string {
	if o == nil || o.Flags == nil {
		var ret []string
		return ret
	}
	return *o.Flags
}

// GetFlagsOk returns a tuple with the Flags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetFlagsOk() (*[]string, bool) {
	if o == nil || o.Flags == nil {
		return nil, false
	}
	return o.Flags, true
}

// HasFlags returns a boolean if a field has been set.
func (o *Link) HasFlags() bool {
	if o != nil && o.Flags != nil {
		return true
	}

	return false
}

// SetFlags gets a reference to the given []string and assigns it to the Flags field.
func (o *Link) SetFlags(v []string) {
	o.Flags = &v
}

// GetMtu returns the Mtu field value if set, zero value otherwise.
func (o *Link) GetMtu() int32 {
	if o == nil || o.Mtu == nil {
		var ret int32
		return ret
	}
	return *o.Mtu
}

// GetMtuOk returns a tuple with the Mtu field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetMtuOk() (*int32, bool) {
	if o == nil || o.Mtu == nil {
		return nil, false
	}
	return o.Mtu, true
}

// HasMtu returns a boolean if a field has been set.
func (o *Link) HasMtu() bool {
	if o != nil && o.Mtu != nil {
		return true
	}

	return false
}

// SetMtu gets a reference to the given int32 and assigns it to the Mtu field.
func (o *Link) SetMtu(v int32) {
	o.Mtu = &v
}

// GetPromiscuity returns the Promiscuity field value if set, zero value otherwise.
func (o *Link) GetPromiscuity() int32 {
	if o == nil || o.Promiscuity == nil {
		var ret int32
		return ret
	}
	return *o.Promiscuity
}

// GetPromiscuityOk returns a tuple with the Promiscuity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetPromiscuityOk() (*int32, bool) {
	if o == nil || o.Promiscuity == nil {
		return nil, false
	}
	return o.Promiscuity, true
}

// HasPromiscuity returns a boolean if a field has been set.
func (o *Link) HasPromiscuity() bool {
	if o != nil && o.Promiscuity != nil {
		return true
	}

	return false
}

// SetPromiscuity gets a reference to the given int32 and assigns it to the Promiscuity field.
func (o *Link) SetPromiscuity(v int32) {
	o.Promiscuity = &v
}

// GetMaster returns the Master field value if set, zero value otherwise.
func (o *Link) GetMaster() string {
	if o == nil || o.Master == nil {
		var ret string
		return ret
	}
	return *o.Master
}

// GetMasterOk returns a tuple with the Master field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetMasterOk() (*string, bool) {
	if o == nil || o.Master == nil {
		return nil, false
	}
	return o.Master, true
}

// HasMaster returns a boolean if a field has been set.
func (o *Link) HasMaster() bool {
	if o != nil && o.Master != nil {
		return true
	}

	return false
}

// SetMaster gets a reference to the given string and assigns it to the Master field.
func (o *Link) SetMaster(v string) {
	o.Master = &v
}

// GetLinkinfo returns the Linkinfo field value if set, zero value otherwise.
func (o *Link) GetLinkinfo() LinkLinkinfo {
	if o == nil || o.Linkinfo == nil {
		var ret LinkLinkinfo
		return ret
	}
	return *o.Linkinfo
}

// GetLinkinfoOk returns a tuple with the Linkinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetLinkinfoOk() (*LinkLinkinfo, bool) {
	if o == nil || o.Linkinfo == nil {
		return nil, false
	}
	return o.Linkinfo, true
}

// HasLinkinfo returns a boolean if a field has been set.
func (o *Link) HasLinkinfo() bool {
	if o != nil && o.Linkinfo != nil {
		return true
	}

	return false
}

// SetLinkinfo gets a reference to the given LinkLinkinfo and assigns it to the Linkinfo field.
func (o *Link) SetLinkinfo(v LinkLinkinfo) {
	o.Linkinfo = &v
}

// GetLinkType returns the LinkType field value
func (o *Link) GetLinkType() string {
	if o == nil  {
		var ret string
		return ret
	}

	return o.LinkType
}

// GetLinkTypeOk returns a tuple with the LinkType field value
// and a boolean to check if the value has been set.
func (o *Link) GetLinkTypeOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.LinkType, true
}

// SetLinkType sets field value
func (o *Link) SetLinkType(v string) {
	o.LinkType = v
}

// GetAddress returns the Address field value if set, zero value otherwise.
func (o *Link) GetAddress() string {
	if o == nil || o.Address == nil {
		var ret string
		return ret
	}
	return *o.Address
}

// GetAddressOk returns a tuple with the Address field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetAddressOk() (*string, bool) {
	if o == nil || o.Address == nil {
		return nil, false
	}
	return o.Address, true
}

// HasAddress returns a boolean if a field has been set.
func (o *Link) HasAddress() bool {
	if o != nil && o.Address != nil {
		return true
	}

	return false
}

// SetAddress gets a reference to the given string and assigns it to the Address field.
func (o *Link) SetAddress(v string) {
	o.Address = &v
}

// GetAddrInfo returns the AddrInfo field value if set, zero value otherwise.
func (o *Link) GetAddrInfo() []LinkAddrInfo {
	if o == nil || o.AddrInfo == nil {
		var ret []LinkAddrInfo
		return ret
	}
	return *o.AddrInfo
}

// GetAddrInfoOk returns a tuple with the AddrInfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Link) GetAddrInfoOk() (*[]LinkAddrInfo, bool) {
	if o == nil || o.AddrInfo == nil {
		return nil, false
	}
	return o.AddrInfo, true
}

// HasAddrInfo returns a boolean if a field has been set.
func (o *Link) HasAddrInfo() bool {
	if o != nil && o.AddrInfo != nil {
		return true
	}

	return false
}

// SetAddrInfo gets a reference to the given []LinkAddrInfo and assigns it to the AddrInfo field.
func (o *Link) SetAddrInfo(v []LinkAddrInfo) {
	o.AddrInfo = &v
}

func (o Link) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["__id"] = o.Id
	}
	if o.Ifindex != nil {
		toSerialize["ifindex"] = o.Ifindex
	}
	if true {
		toSerialize["ifname"] = o.Ifname
	}
	if o.Flags != nil {
		toSerialize["flags"] = o.Flags
	}
	if o.Mtu != nil {
		toSerialize["mtu"] = o.Mtu
	}
	if o.Promiscuity != nil {
		toSerialize["promiscuity"] = o.Promiscuity
	}
	if o.Master != nil {
		toSerialize["master"] = o.Master
	}
	if o.Linkinfo != nil {
		toSerialize["linkinfo"] = o.Linkinfo
	}
	if true {
		toSerialize["link_type"] = o.LinkType
	}
	if o.Address != nil {
		toSerialize["address"] = o.Address
	}
	if o.AddrInfo != nil {
		toSerialize["addr_info"] = o.AddrInfo
	}
	return json.Marshal(toSerialize)
}

type NullableLink struct {
	value *Link
	isSet bool
}

func (v NullableLink) Get() *Link {
	return v.value
}

func (v *NullableLink) Set(val *Link) {
	v.value = val
	v.isSet = true
}

func (v NullableLink) IsSet() bool {
	return v.isSet
}

func (v *NullableLink) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLink(val *Link) *NullableLink {
	return &NullableLink{value: val, isSet: true}
}

func (v NullableLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLink) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


