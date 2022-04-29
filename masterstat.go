package masterstat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)

func GetServerAddresses(masterAddress string) ([]string, error) {
	statusPacket := []byte{0x63, 0x0a, 0x00}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 0x64, 0x0a}
	response, err := udpRequest(masterAddress, statusPacket, expectedHeader)

	if err != nil {
		return nil, err
	}

	responseBody := response[len(expectedHeader):]
	reader := bytes.NewReader(responseBody)
	serverAddresses := make([]string, 0)

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

func GetServerAddressesFromMany(masterAddresses []string) []string {
	var (
		wg              sync.WaitGroup
		mutex           sync.Mutex
		serverAddresses = make([]string, 0)
	)

	for _, masterAddress := range masterAddresses {
		wg.Add(1)

		go func(masterAddress string) {
			defer wg.Done()

			addresses, err := GetServerAddresses(masterAddress)

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

	return uniqueStrings(serverAddresses)
}

type rawServerAddress struct {
	IpParts [4]byte
	Port    uint16
}

func (addr rawServerAddress) toString() string {
	ip := net.IPv4(addr.IpParts[0], addr.IpParts[1], addr.IpParts[2], addr.IpParts[3]).String()
	return fmt.Sprintf("%s:%d", ip, addr.Port)
}

func uniqueStrings(values []string) []string {
	valueMap := make(map[string]bool, 0)
	uniqueValues := make([]string, 0)

	for _, value := range values {
		if !valueMap[value] {
			uniqueValues = append(uniqueValues, value)
			valueMap[value] = true
		}
	}

	return uniqueValues
}
