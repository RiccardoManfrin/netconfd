package openapi

import "encoding/json"

type DiscriminatedInfoData struct {
	LinkLinkinfoInfoData
}

func (v DiscriminatedInfoData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.LinkLinkinfoInfoData)
}

func (v *DiscriminatedInfoData) UnmarshalJSON(src []byte) error {
	return json.Unmarshal(src, &v.LinkLinkinfoInfoData)
}
