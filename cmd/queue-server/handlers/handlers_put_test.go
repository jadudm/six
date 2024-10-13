package handlers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutHandler(t *testing.T) {
	Init("test_data", "*/1 0 0 0 0")
	job := map[string]string{
		"Host": "www.fac.gov",
		"Path": "/api",
	}
	jobjs, _ := json.Marshal(job)
	putEnqueueHandler("HEAD", jobjs)
	assert.Equal(t, Q.Length("HEAD"), 1)
	assert.Equal(t, Q.Length("DOES NOT EXIST"), -1)
}
