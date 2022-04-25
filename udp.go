package masterstat

import (
	"bytes"
	"errors"
	"net"
	"time"
)

func udpRequest(address string, statusPacket []byte, expectedHeader []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	const (
		BufferSize  = 8192
		Retries     = 3
		TimeoutInMs = 500
	)

	getTimeout := func() time.Time {
		return time.Now().Add(time.Duration(TimeoutInMs) * time.Millisecond)
	}
	response := make([]byte, BufferSize)
	responseLength := 0

	for i := 0; i < Retries; i++ {
		conn.SetDeadline(getTimeout())

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(getTimeout())
		responseLength, err = conn.Read(response)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	isValidHeader := bytes.Equal(response[:len(expectedHeader)], expectedHeader)
	if !isValidHeader {
		err = errors.New(address + ": Response error, invalid header.")
		return nil, err
	}

	return response[:responseLength], nil
}
