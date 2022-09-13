package masterstat

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/vikpe/masterstat/commands/serverlist"
	"github.com/vikpe/masterstat/qutil"
	"github.com/vikpe/udpclient"
)

func GetServerAddresses(masterAddress string) ([]string, error) {
	return serverlist.ParseResponse(
		udpclient.New().SendCommand(masterAddress, serverlist.Command),
	)
}

func GetServerAddressesFromMany(masterAddresses []string) ([]string, []error) {
	var (
		wg              sync.WaitGroup
		mutex           sync.Mutex
		serverAddresses = make([]string, 0)
		errs            = make([]error, 0)
	)

	for _, masterAddress := range masterAddresses {
		wg.Add(1)

		go func(masterAddress string) {
			defer wg.Done()

			addresses, err := GetServerAddresses(masterAddress)

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, masterAddress, err)))
				return
			}

			mutex.Lock()
			serverAddresses = append(serverAddresses, addresses...)
			mutex.Unlock()
		}(masterAddress)
	}

	wg.Wait()

	serverAddresses = qutil.UniqueStrings(serverAddresses)
	sort.Strings(serverAddresses)

	return serverAddresses, errs
}
