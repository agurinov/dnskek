package docker

import (
	"fmt"
	"net"
	"net/url"
)

var (
	subdomainRegexExpression  = "[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?"              //each DNS label can contain up to 63 characters
	subdomainsRegexExpression = fmt.Sprintf("(?:%s[.])+", subdomainRegexExpression) // subdivision can go down to 127 levels deep
)

func getIPByTCPURL(URL string) net.IP {
	if addr, err := url.Parse(URL); err == nil {
		// no errors -> ip is the hostname
		return net.ParseIP(addr.Hostname())
	}
	// some errors occured -> cannot get IP
	return nil
}
