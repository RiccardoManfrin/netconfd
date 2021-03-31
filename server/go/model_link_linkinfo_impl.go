package openapi

import (
	"encoding/json"
	"fmt"
)

type DiscriminatedLinkInfo struct {
	LinkLinkinfo
}

func (v *LinkLinkinfo) Validate() error {
	return fmt.Errorf("Don't user protocol for gre")
}

func (v *DiscriminatedLinkInfo) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &v.LinkLinkinfo)
	if err != nil {
		return err
	}
	err = v.LinkLinkinfo.Validate()
	return err
}
