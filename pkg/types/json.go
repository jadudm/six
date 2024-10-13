package types

import (
	"encoding/json"
	"log"

	"github.com/tidwall/gjson"
)

type JSON = []byte

var EmptyObject JSON = []byte("{}")

func BytesToMap(msg []byte) map[string]interface{} {
	return gjson.ParseBytes(msg).Value().(map[string]interface{})
}

func StringToMap(msg string) map[string]interface{} {
	return gjson.Parse(msg).Value().(map[string]interface{})
}

func MapToBytes(msg map[string]interface{}) []byte {
	v, _ := json.Marshal(msg)
	return v
}

func StructToJson(s interface{}) JSON {
	js, err := json.Marshal(s)
	if err != nil {
		log.Printf("Could not marshal %s to JSON", s)
	}
	return js
}
