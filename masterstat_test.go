package masterstat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawServerAddressToString(t *testing.T) {

	rawAddress := rawServerAddress{
		IpParts: [4]byte{1, 2, 3, 4},
		Port:    28501,
	}

	result := rawAddress.toString()
	expect := "1.2.3.4:28501"

	assert.Equal(t, expect, result)
}
