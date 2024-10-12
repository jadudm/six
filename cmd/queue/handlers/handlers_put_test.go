package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutHandler(t *testing.T) {
	Init("test_data", "*/1 0 0 0 0")
	putEnqueueHandler("www.fac.gov")
	assert.Equal(t, TheMultiqueue.Length("www.fac.gov"), 1)
}
