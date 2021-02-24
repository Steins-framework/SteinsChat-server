package lib

import "encoding/json"

type UnifiedDataFormat struct {
	Event string `json:"event"`
	Data interface{} `json:"data"`
}

func (f UnifiedDataFormat) getDataBytes() ([]byte, error) {
	return json.Marshal(f.Data)
}
