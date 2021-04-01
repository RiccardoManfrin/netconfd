/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.3.0
 * Contact: support@athonet.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// Scope scope of the object
type Scope string

// List of scope
const (
	LINK Scope = "link"
	GLOBAL Scope = "global"
	UNIVERSE Scope = "universe"
	SITE Scope = "site"
	NOWHERE Scope = "nowhere"
)

func (v *Scope) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Scope(value)
	for _, existing := range []Scope{ "link", "global", "universe", "site", "nowhere",   } {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Scope", value)
}

// Ptr returns reference to scope value
func (v Scope) Ptr() *Scope {
	return &v
}

type NullableScope struct {
	value *Scope
	isSet bool
}

func (v NullableScope) Get() *Scope {
	return v.value
}

func (v *NullableScope) Set(val *Scope) {
	v.value = val
	v.isSet = true
}

func (v NullableScope) IsSet() bool {
	return v.isSet
}

func (v *NullableScope) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScope(val *Scope) *NullableScope {
	return &NullableScope{value: val, isSet: true}
}

func (v NullableScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScope) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

