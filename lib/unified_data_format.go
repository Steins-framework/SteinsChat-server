package lib

import "encoding/json"

type UnifiedDataFormat struct {
	Event string
	Data interface{}
}

func (f UnifiedDataFormat) getDataBytes() ([]byte, error) {
	return json.Marshal(f.Data)
}
