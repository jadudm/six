package queues

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStruct(t *testing.T) {
	js := `{"type": "job", "domain": "www.fac.gov"}`
	j, err := ToStruct([]byte(js))
	if err != nil {
		t.Error(fmt.Sprint(err))
	}
	_, isJob := j.(Job)
	assert.True(t, isJob)

	js = `{"type": "job_response", "domain": "www.fac.gov"}`
	jr, err := ToStruct([]byte(js))
	if err != nil {
		t.Error(fmt.Sprint(err))
	}
	_, isJobResponse := jr.(JobResponse)
	assert.True(t, isJobResponse)

	js = `{"type": "job", "domain": "www.fac.gov"}`
	jr2, err := ToStruct([]byte(js))
	if err != nil {
		t.Error(fmt.Sprint(err))
	}
	_, isNotJobResponse := jr2.(JobResponse)
	assert.False(t, isNotJobResponse)
}
