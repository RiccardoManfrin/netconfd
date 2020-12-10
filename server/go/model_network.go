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

// Network struct for Network
type Network struct {
	// Series of links layer interfaces to configure within the namespace
	Links *[]Link `json:"links,omitempty"`
	// Namespace routes
	Routes *[]Route `json:"routes,omitempty"`
	// DHCP context
	Dhcp *[]Dhcp `json:"dhcp,omitempty"`
}

// NewNetwork instantiates a new Network object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNetwork() *Network {
	this := Network{}
	return &this
}

// NewNetworkWithDefaults instantiates a new Network object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNetworkWithDefaults() *Network {
	this := Network{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise.
func (o *Network) GetLinks() []Link {
	if o == nil || o.Links == nil {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Network) GetLinksOk() (*[]Link, bool) {
	if o == nil || o.Links == nil {
		return nil, false
	}
	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *Network) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *Network) SetLinks(v []Link) {
	o.Links = &v
}

// GetRoutes returns the Routes field value if set, zero value otherwise.
func (o *Network) GetRoutes() []Route {
	if o == nil || o.Routes == nil {
		var ret []Route
		return ret
	}
	return *o.Routes
}

// GetRoutesOk returns a tuple with the Routes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Network) GetRoutesOk() (*[]Route, bool) {
	if o == nil || o.Routes == nil {
		return nil, false
	}
	return o.Routes, true
}

// HasRoutes returns a boolean if a field has been set.
func (o *Network) HasRoutes() bool {
	if o != nil && o.Routes != nil {
		return true
	}

	return false
}

// SetRoutes gets a reference to the given []Route and assigns it to the Routes field.
func (o *Network) SetRoutes(v []Route) {
	o.Routes = &v
}

// GetDhcp returns the Dhcp field value if set, zero value otherwise.
func (o *Network) GetDhcp() []Dhcp {
	if o == nil || o.Dhcp == nil {
		var ret []Dhcp
		return ret
	}
	return *o.Dhcp
}

// GetDhcpOk returns a tuple with the Dhcp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Network) GetDhcpOk() (*[]Dhcp, bool) {
	if o == nil || o.Dhcp == nil {
		return nil, false
	}
	return o.Dhcp, true
}

// HasDhcp returns a boolean if a field has been set.
func (o *Network) HasDhcp() bool {
	if o != nil && o.Dhcp != nil {
		return true
	}

	return false
}

// SetDhcp gets a reference to the given []Dhcp and assigns it to the Dhcp field.
func (o *Network) SetDhcp(v []Dhcp) {
	o.Dhcp = &v
}

func (o Network) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Links != nil {
		toSerialize["links"] = o.Links
	}
	if o.Routes != nil {
		toSerialize["routes"] = o.Routes
	}
	if o.Dhcp != nil {
		toSerialize["dhcp"] = o.Dhcp
	}
	return json.Marshal(toSerialize)
}

type NullableNetwork struct {
	value *Network
	isSet bool
}

func (v NullableNetwork) Get() *Network {
	return v.value
}

func (v *NullableNetwork) Set(val *Network) {
	v.value = val
	v.isSet = true
}

func (v NullableNetwork) IsSet() bool {
	return v.isSet
}

func (v *NullableNetwork) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNetwork(val *Network) *NullableNetwork {
	return &NullableNetwork{value: val, isSet: true}
}

func (v NullableNetwork) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNetwork) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


