/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.2.0
 * Contact: support@athonet.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// Dnsid ID of the DNS servers, defining their priority 
type Dnsid string

// List of dnsid
const (
	PRIMARY Dnsid = "primary"
	SECONDARY Dnsid = "secondary"
)

func (v *Dnsid) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Dnsid(value)
	for _, existing := range []Dnsid{ "primary", "secondary",   } {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Dnsid", value)
}

// Ptr returns reference to dnsid value
func (v Dnsid) Ptr() *Dnsid {
	return &v
}

type NullableDnsid struct {
	value *Dnsid
	isSet bool
}

func (v NullableDnsid) Get() *Dnsid {
	return v.value
}

func (v *NullableDnsid) Set(val *Dnsid) {
	v.value = val
	v.isSet = true
}

func (v NullableDnsid) IsSet() bool {
	return v.isSet
}

func (v *NullableDnsid) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnsid(val *Dnsid) *NullableDnsid {
	return &NullableDnsid{value: val, isSet: true}
}

func (v NullableDnsid) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnsid) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

