package types

import (
	"bytes"
	"encoding/json"
	"errors"
)

type JSON = []byte

type Job struct {
	JobId  string
	Domain string
	Path   string
}

type JobResponse struct {
	Result string
	Domain string
	Pages  []string
}

// Given a JSON string, figure out which
// struct it should be, and return it
func ToStruct(js []byte) (any, error) {
	var data interface{}
	err := json.Unmarshal(js, &data)
	if err != nil {
		return nil, err
	}

	the_type := data.(map[string]interface{})["type"]
	switch the_type {
	case "job":
		var j Job
		err = json.NewDecoder(bytes.NewReader(js)).Decode(&j)
		if err == nil {
			return j, nil
		}
	case "jobresponse", "job_response":
		var jr JobResponse
		err = json.NewDecoder(bytes.NewReader(js)).Decode(&jr)
		if err == nil {
			return jr, nil
		}
	}
	return nil, errors.New("WAT")
}
