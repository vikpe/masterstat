package masterstat_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/masterstat"
	"github.com/vikpe/udphelper"
)

func TestGetServerAddresses(t *testing.T) {
	t.Run("UDP request error", func(t *testing.T) {
		result, err := masterstat.GetServerAddresses("foo:666")
		assert.Equal(t, []string{}, result)
		assert.ErrorContains(t, err, "failure in name resolution")
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
			"245.73.111.107:28104",
			"66.69.101.148:27500",
		}
		assert.Equal(t, expect, result)
		assert.Equal(t, err, nil)
	})
}

func TestGetServerAddressesFromMany(t *testing.T) {
	t.Run("UDP request error", func(t *testing.T) {
		masterAddresses := []string{"foo:666"}
		result, errs := masterstat.GetServerAddressesFromMany(masterAddresses)

		assert.Equal(t, []string{}, result)
		assert.Len(t, errs, 1)
		assert.ErrorContains(t, errs[0], "failure in name resolution")
	})

	t.Run("Success", func(t *testing.T) {
		const master1 = ":8003"
		go func() {
			responseBody := []byte{
				0xff, 0xff, 0xff, 0xff, 0x64, 0x0a, // header
				0x42, 0x45, 0x65, 0x94, 0x6b, 0x6c, //  server 1
				0xf5, 0x49, 0x6f, 0x6b, 0x6d, 0xc8, //  server 2
			}
			udphelper.New(master1).Respond(responseBody)
		}()

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

		const masterInvalid = "foo:666"
		masterAddresses := []string{master1, masterInvalid, master2}
		result, errs := masterstat.GetServerAddressesFromMany(masterAddresses)

		expect := []string{
			"200.42.92.173:27500",
			"245.73.111.107:28104",
			"66.69.101.148:27500",
		}

		assert.Equal(t, expect, result)
		assert.Equal(t, errs, []error{errors.New("foo:666 - dial udp4: lookup foo: Temporary failure in name resolution")})
	})
}
