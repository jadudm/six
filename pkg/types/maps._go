package types

import "encoding/json"

func ToMap(msg JSON) map[string]interface{} {
	var data interface{}
	json.Unmarshal(msg, &data)
	return data.(map[string]interface{})
}
