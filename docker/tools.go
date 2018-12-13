package docker

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var (
	subdomainRegexExpression  = "[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?"              //each DNS label can contain up to 63 characters
	subdomainsRegexExpression = fmt.Sprintf("(?:%s[.])*", subdomainRegexExpression) // subdivision can go down to 127 levels deep
)

func getIPByTCPURL(URL string) net.IP {
	if addr, err := url.Parse(URL); err == nil {
		// no errors -> ip is the hostname
		return net.ParseIP(addr.Hostname())
	}
	// some errors occured -> cannot get IP
	return nil
}

func getMachinesByRaw(raw string) (machines []*Machine) {
	if raw == "" {
		// empty string will split into slice with len == 1
		// no need this case
		return
	}

	for _, row := range strings.Split(raw, "\n") {
		if row == "" {
			// empty string will split into slice with len == 1
			// no need this case
			return
		}

		splitted := strings.SplitN(row, "|", 5)
		// allocate machine info (contains 5 parts)
		var bits [5]string
		// copy from slice (may be less tahn 5 parts to limited size array)
		copy(bits[:], splitted)
		// parse data and get machine struct
		dm := newMachine(bits[0], bits[1], bits[2], bits[3], bits[4])
		// register machine
		machines = append(machines, dm)
	}
	return
}
