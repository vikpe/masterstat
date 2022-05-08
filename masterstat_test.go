package masterstat_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/masterstat"
	"github.com/vikpe/udphelper"
)

func TestGetServerAddresses(t *testing.T) {
	t.Run("UDP request error", func(t *testing.T) {
		result, err := masterstat.GetServerAddresses("foo:666")
		expect := []string{}
		assert.Equal(t, expect, result)
		assert.ErrorContains(t, err, "dial udp4: lookup foo:")
	})

	t.Run("Success", func(t *testing.T) {
		const addr = ":8001"

		go func() {
			responseBody := []byte{
				0xff, 0xff, 0xff, 0xff, 0x64, 0x0a, // header
				0x42, 0x45, 0x65, 0x94, 0x6b, 0x6c, //  server 1
				0xf5, 0x49, 0x6f, 0x6b, 0x6d, 0xc8, //  server 2
			}
			udphelper.New(addr).Respond(responseBody)
		}()

		time.Sleep(10 * time.Millisecond)

		result, err := masterstat.GetServerAddresses(addr)
		expect := []string{
			"66.69.101.148:27500",
			"245.73.111.107:28104",
		}
		assert.Equal(t, expect, result)
		assert.Equal(t, err, nil)
	})
}

func TestGetServerAddressesFromMany(t *testing.T) {
	t.Run("UDP request error", func(t *testing.T) {
		const master = ":8002"
		go func() {
			net.ListenPacket("udp", master)
		}()

		time.Sleep(10 * time.Millisecond)

		masterAddresses := []string{master}
		result, err := masterstat.GetServerAddressesFromMany(masterAddresses)

		assert.Equal(t, []string{}, result)
		assert.ErrorContains(t, err, ":8002: i/o timeout")
	})

	t.Run("Success", func(t *testing.T) {
		// master 1
		const master1 = ":8003"
		go func() {
			responseBody := []byte{
				0xff, 0xff, 0xff, 0xff, 0x64, 0x0a, // header
				0x42, 0x45, 0x65, 0x94, 0x6b, 0x6c, //  server 1
				0xf5, 0x49, 0x6f, 0x6b, 0x6d, 0xc8, //  server 2
			}
			udphelper.New(master1).Respond(responseBody)
		}()

		// master 2
		const master2 = ":8004"
		go func() {
			responseBody := []byte{
				0xff, 0xff, 0xff, 0xff, 0x64, 0x0a, // header
				0x42, 0x45, 0x65, 0x94, 0x6b, 0x6c, //  server 1
				0xc8, 0x2a, 0x5c, 0xad, 0x6b, 0x6c, //  server 3
			}
			udphelper.New(master2).Respond(responseBody)
		}()

		time.Sleep(10 * time.Millisecond)

		masterAddresses := []string{master1, master2}
		result, err := masterstat.GetServerAddressesFromMany(masterAddresses)
		assert.Contains(t, result, "66.69.101.148:27500")
		assert.Contains(t, result, "245.73.111.107:28104")
		assert.Contains(t, result, "200.42.92.173:27500")
		assert.Equal(t, err, nil)
	})
}
