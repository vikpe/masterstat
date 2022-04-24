package masterstat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func timeInFuture(delta int) time.Time {
	return time.Now().Add(time.Duration(delta) * time.Millisecond)
}

type rawServerAddress struct {
	IpParts [4]byte
	Port    uint16
}

func (addr rawServerAddress) toString() string {
	ip := net.IPv4(addr.IpParts[0], addr.IpParts[1], addr.IpParts[2], addr.IpParts[3]).String()
	return fmt.Sprintf("%s:%d", ip, addr.Port)
}

func Stat(masterAddress string, retries int, timeout int) ([]string, error) {
	serverAddresses := make([]string, 0)

	conn, err := net.Dial("udp4", masterAddress)
	if err != nil {
		return serverAddresses, err
	}

	defer conn.Close()

	statusPacket := []byte{0x63, 0x0a, 0x00}
	buffer := make([]byte, 8192)
	bufferLength := 0

	for i := 0; i < retries; i++ {
		conn.SetDeadline(timeInFuture(timeout))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return serverAddresses, err
		}

		conn.SetDeadline(timeInFuture(timeout))
		bufferLength, err = conn.Read(buffer)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return serverAddresses, err
	}

	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 0x64, 0x0a}
	responseHeader := buffer[:len(expectedHeader)]
	isValidHeader := bytes.Equal(responseHeader, expectedHeader)

	if !isValidHeader {
		err = errors.New(masterAddress + ": Response error")
		return serverAddresses, err
	}

	reader := bytes.NewReader(buffer[6:bufferLength])

	for {
		var rawAddress rawServerAddress

		err = binary.Read(reader, binary.BigEndian, &rawAddress)
		if err != nil {
			break
		}

		serverAddresses = append(serverAddresses, rawAddress.toString())
	}

	return serverAddresses, nil
}

func StatMany(masterAddresses []string, retries int, timeout int) []string {
	var (
		wg              sync.WaitGroup
		mutex           sync.Mutex
		serverAddresses = make([]string, 0)
	)

	for _, masterAddress := range masterAddresses {
		wg.Add(1)

		go func(masterAddress string) {
			defer wg.Done()

			addresses, err := Stat(masterAddress, retries, timeout)

			if err != nil {
				log.Println(fmt.Sprintf("ERROR: unable to stat %s", masterAddresses), err)
				return
			}

			mutex.Lock()
			serverAddresses = append(serverAddresses, addresses...)
			mutex.Unlock()
		}(masterAddress)
	}

	wg.Wait()

	addressMap := make(map[string]bool, 0)
	uniqueAddresses := make([]string, 0)

	for _, address := range serverAddresses {
		if !addressMap[address] {
			uniqueAddresses = append(uniqueAddresses, address)
			addressMap[address] = true
		}
	}

	return uniqueAddresses
}
