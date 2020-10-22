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

// Route IP L3 Ruote entry
type Route struct {
	Dst *RouteDst `json:"dst,omitempty"`
	Gateway *Ip `json:"gateway,omitempty"`
	// Interface name 
	Dev *string `json:"dev,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	Metric *int32 `json:"metric,omitempty"`
	Scope *Scope `json:"scope,omitempty"`
	Prefsrc *Ip `json:"prefsrc,omitempty"`
	// Route flags
	Flags *[]string `json:"flags,omitempty"`
}

// NewRoute instantiates a new Route object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRoute() *Route {
	this := Route{}
	return &this
}

// NewRouteWithDefaults instantiates a new Route object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRouteWithDefaults() *Route {
	this := Route{}
	return &this
}

// GetDst returns the Dst field value if set, zero value otherwise.
func (o *Route) GetDst() RouteDst {
	if o == nil || o.Dst == nil {
		var ret RouteDst
		return ret
	}
	return *o.Dst
}

// GetDstOk returns a tuple with the Dst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetDstOk() (*RouteDst, bool) {
	if o == nil || o.Dst == nil {
		return nil, false
	}
	return o.Dst, true
}

// HasDst returns a boolean if a field has been set.
func (o *Route) HasDst() bool {
	if o != nil && o.Dst != nil {
		return true
	}

	return false
}

// SetDst gets a reference to the given RouteDst and assigns it to the Dst field.
func (o *Route) SetDst(v RouteDst) {
	o.Dst = &v
}

// GetGateway returns the Gateway field value if set, zero value otherwise.
func (o *Route) GetGateway() Ip {
	if o == nil || o.Gateway == nil {
		var ret Ip
		return ret
	}
	return *o.Gateway
}

// GetGatewayOk returns a tuple with the Gateway field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetGatewayOk() (*Ip, bool) {
	if o == nil || o.Gateway == nil {
		return nil, false
	}
	return o.Gateway, true
}

// HasGateway returns a boolean if a field has been set.
func (o *Route) HasGateway() bool {
	if o != nil && o.Gateway != nil {
		return true
	}

	return false
}

// SetGateway gets a reference to the given Ip and assigns it to the Gateway field.
func (o *Route) SetGateway(v Ip) {
	o.Gateway = &v
}

// GetDev returns the Dev field value if set, zero value otherwise.
func (o *Route) GetDev() string {
	if o == nil || o.Dev == nil {
		var ret string
		return ret
	}
	return *o.Dev
}

// GetDevOk returns a tuple with the Dev field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetDevOk() (*string, bool) {
	if o == nil || o.Dev == nil {
		return nil, false
	}
	return o.Dev, true
}

// HasDev returns a boolean if a field has been set.
func (o *Route) HasDev() bool {
	if o != nil && o.Dev != nil {
		return true
	}

	return false
}

// SetDev gets a reference to the given string and assigns it to the Dev field.
func (o *Route) SetDev(v string) {
	o.Dev = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *Route) GetProtocol() string {
	if o == nil || o.Protocol == nil {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetProtocolOk() (*string, bool) {
	if o == nil || o.Protocol == nil {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *Route) HasProtocol() bool {
	if o != nil && o.Protocol != nil {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *Route) SetProtocol(v string) {
	o.Protocol = &v
}

// GetMetric returns the Metric field value if set, zero value otherwise.
func (o *Route) GetMetric() int32 {
	if o == nil || o.Metric == nil {
		var ret int32
		return ret
	}
	return *o.Metric
}

// GetMetricOk returns a tuple with the Metric field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetMetricOk() (*int32, bool) {
	if o == nil || o.Metric == nil {
		return nil, false
	}
	return o.Metric, true
}

// HasMetric returns a boolean if a field has been set.
func (o *Route) HasMetric() bool {
	if o != nil && o.Metric != nil {
		return true
	}

	return false
}

// SetMetric gets a reference to the given int32 and assigns it to the Metric field.
func (o *Route) SetMetric(v int32) {
	o.Metric = &v
}

// GetScope returns the Scope field value if set, zero value otherwise.
func (o *Route) GetScope() Scope {
	if o == nil || o.Scope == nil {
		var ret Scope
		return ret
	}
	return *o.Scope
}

// GetScopeOk returns a tuple with the Scope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetScopeOk() (*Scope, bool) {
	if o == nil || o.Scope == nil {
		return nil, false
	}
	return o.Scope, true
}

// HasScope returns a boolean if a field has been set.
func (o *Route) HasScope() bool {
	if o != nil && o.Scope != nil {
		return true
	}

	return false
}

// SetScope gets a reference to the given Scope and assigns it to the Scope field.
func (o *Route) SetScope(v Scope) {
	o.Scope = &v
}

// GetPrefsrc returns the Prefsrc field value if set, zero value otherwise.
func (o *Route) GetPrefsrc() Ip {
	if o == nil || o.Prefsrc == nil {
		var ret Ip
		return ret
	}
	return *o.Prefsrc
}

// GetPrefsrcOk returns a tuple with the Prefsrc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetPrefsrcOk() (*Ip, bool) {
	if o == nil || o.Prefsrc == nil {
		return nil, false
	}
	return o.Prefsrc, true
}

// HasPrefsrc returns a boolean if a field has been set.
func (o *Route) HasPrefsrc() bool {
	if o != nil && o.Prefsrc != nil {
		return true
	}

	return false
}

// SetPrefsrc gets a reference to the given Ip and assigns it to the Prefsrc field.
func (o *Route) SetPrefsrc(v Ip) {
	o.Prefsrc = &v
}

// GetFlags returns the Flags field value if set, zero value otherwise.
func (o *Route) GetFlags() []string {
	if o == nil || o.Flags == nil {
		var ret []string
		return ret
	}
	return *o.Flags
}

// GetFlagsOk returns a tuple with the Flags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Route) GetFlagsOk() (*[]string, bool) {
	if o == nil || o.Flags == nil {
		return nil, false
	}
	return o.Flags, true
}

// HasFlags returns a boolean if a field has been set.
func (o *Route) HasFlags() bool {
	if o != nil && o.Flags != nil {
		return true
	}

	return false
}

// SetFlags gets a reference to the given []string and assigns it to the Flags field.
func (o *Route) SetFlags(v []string) {
	o.Flags = &v
}

func (o Route) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Dst != nil {
		toSerialize["dst"] = o.Dst
	}
	if o.Gateway != nil {
		toSerialize["gateway"] = o.Gateway
	}
	if o.Dev != nil {
		toSerialize["dev"] = o.Dev
	}
	if o.Protocol != nil {
		toSerialize["protocol"] = o.Protocol
	}
	if o.Metric != nil {
		toSerialize["metric"] = o.Metric
	}
	if o.Scope != nil {
		toSerialize["scope"] = o.Scope
	}
	if o.Prefsrc != nil {
		toSerialize["prefsrc"] = o.Prefsrc
	}
	if o.Flags != nil {
		toSerialize["flags"] = o.Flags
	}
	return json.Marshal(toSerialize)
}

type NullableRoute struct {
	value *Route
	isSet bool
}

func (v NullableRoute) Get() *Route {
	return v.value
}

func (v *NullableRoute) Set(val *Route) {
	v.value = val
	v.isSet = true
}

func (v NullableRoute) IsSet() bool {
	return v.isSet
}

func (v *NullableRoute) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRoute(val *Route) *NullableRoute {
	return &NullableRoute{value: val, isSet: true}
}

func (v NullableRoute) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRoute) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


