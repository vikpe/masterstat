package serverlist

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0x63, 0x0a, 0x00},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 0x64, 0x0a},
}

func ParseResponse(responseBody []byte, err error) ([]string, error) {
	if err != nil {
		return []string{}, err
	}

	return ParseResponseBody(responseBody), nil
}

func ParseResponseBody(responseBody []byte) []string {
	reader := bytes.NewReader(responseBody)
	serverAddresses := make([]string, 0)

	for {
		var rawAddress rawServerAddress

		err := binary.Read(reader, binary.BigEndian, &rawAddress)
		if err != nil {
			break
		}

		serverAddresses = append(serverAddresses, rawAddress.ToString())
	}

	return serverAddresses
}

type rawServerAddress struct {
	IpParts [4]byte
	Port    uint16
}

func (addr rawServerAddress) ToString() string {
	ip := fmt.Sprintf("%d.%d.%d.%d", addr.IpParts[0], addr.IpParts[1], addr.IpParts[2], addr.IpParts[3])
	return fmt.Sprintf("%s:%d", ip, addr.Port)
}
