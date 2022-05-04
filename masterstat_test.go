package masterstat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueStrings(t *testing.T) {
	values := []string{"alpha", "beta", "beta", "gamma"}
	result := uniqueStrings(values)
	expect := []string{"alpha", "beta", "gamma"}

	assert.Equal(t, expect, result)
}

func TestRawServerAddressToString(t *testing.T) {

	rawAddress := rawServerAddress{
		IpParts: [4]byte{1, 2, 3, 4},
		Port:    28501,
	}

	result := rawAddress.toString()
	expect := "1.2.3.4:28501"

	assert.Equal(t, expect, result)
}
