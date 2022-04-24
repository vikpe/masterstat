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
	expect := "1.2.3.4:2850z1"

	assert.Equal(t, expect, result)
}
