package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestEncodeUId(t *testing.T) {
	uid, _ := EncodeUId(1)
	assert.Equal(t, uid, "27906184352")
}

func TestDecodeUId(t *testing.T) {
	uid, _ := DecodeUId("27906184352")
	assert.Equal(t, uid, int64(1))
}
