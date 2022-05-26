package serverlist_test

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/masterstat/commands/serverlist"
)

func TestParseResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		responseBody := []byte("")
		serverAddresses, err := serverlist.ParseResponse(responseBody, errors.New("some error"))
		assert.Equal(t, []string{}, serverAddresses)
		assert.EqualError(t, err, "some error")
	})

	t.Run("empty response body", func(t *testing.T) {
		responseBody := []byte("")
		serverAddresses, err := serverlist.ParseResponse(responseBody, nil)
		assert.Equal(t, []string{}, serverAddresses)
		assert.Nil(t, err)
	})

	t.Run("non-empty response body", func(t *testing.T) {
		responseBody, _ := hex.DecodeString(strings.Join([]string{"424565946b6c", "f5496f6b6dc8", "c82a5cad6b6c"}, ""))
		serverAddresses, err := serverlist.ParseResponse(responseBody, nil)
		expect := []string{
			"200.42.92.173:27500",
			"245.73.111.107:28104",
			"66.69.101.148:27500",
		}
		assert.Equal(t, expect, serverAddresses)
		assert.Nil(t, err)
	})
}
