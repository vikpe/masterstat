package masterstat

import (
	"bytes"
	"errors"
	"net"
	"time"
)

func UdpRequest(address string, statusPacket []byte, expectedHeader []byte) ([]byte, error) {
	nullResponse := make([]byte, 0)

	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nullResponse, err
	}
	defer conn.Close()

	const (
		BufferSize  = 8192
		Retries     = 3
		TimeoutInMs = 500
	)

	response := make([]byte, BufferSize)

	for i := 0; i < Retries; i++ {
		conn.SetDeadline(timeInFuture(TimeoutInMs))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nullResponse, err
		}

		conn.SetDeadline(timeInFuture(TimeoutInMs))
		_, err = conn.Read(response)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nullResponse, err
	}

	isValidHeader := bytes.Equal(response[:len(expectedHeader)], expectedHeader)
	if !isValidHeader {
		err = errors.New(address + ": Response error, invalid header.")
		return nullResponse, err
	}

	return response, nil
}

func timeInFuture(delta int) time.Time {
	return time.Now().Add(time.Duration(delta) * time.Millisecond)
}
