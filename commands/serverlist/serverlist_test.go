package serverlist_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/masterstat/commands/serverlist"
)

func TestParseResponse(t *testing.T) {
	// empty response
	expect := make([]string, 0)
	actual := serverlist.ParseResponseBody([]byte(""))
	assert.Equal(t, expect, actual)

	// non-empty response
	response, _ := hex.DecodeString(strings.Join([]string{
		"424565946b6c",
		"f5496f6b6dc8",
		"c82a5cad6b6c",
	}, ""))
	actual = serverlist.ParseResponseBody(response)
	expect = []string{
		"66.69.101.148:27500",
		"245.73.111.107:28104",
		"200.42.92.173:27500",
	}
	assert.Equal(t, expect, actual)
}
