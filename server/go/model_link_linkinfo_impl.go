package openapi

import "encoding/json"

type DiscriminatedLinkInfo struct {
	LinkLinkinfo
}

func (v *DiscriminatedLinkInfo) UnmarshalJSON(src []byte) error {
	return json.Unmarshal(src, &v.LinkLinkinfo)
}
