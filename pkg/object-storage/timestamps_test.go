package object_storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Total: 8 + 4 + 5 + 5 = 24 bits
// 24-0-42-14-11-5-14
//    24    0    14    11
// 11000 0000 01110 01011

func TestBitstampString(t *testing.T) {
	bs := CreateBitstamp(24, 10, 11, 2)
	assert.Equal(t, bs.Int(), 403810)
	assert.Equal(t, "24-10-11-02", bs.String())
}

func TestBitstampBitwise(t *testing.T) {
	bs := CreateBitstamp(24, 10, 11, 2)
	//               8-------4---5----5----
	assert.Equal(t, "0001100010100101100010", bs.BitString())
}

func TestBitstampFromInt(t *testing.T) {
	bs := CreateBitstamp(24, 10, 11, 2)
	assert.Equal(t, bs.Int(), 403810)
	assert.Equal(t, bs.Int(), BitstampFromInt(403810).Int())
}
