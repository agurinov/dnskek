package docker

import (
	"fmt"
	"net"
	"net/url"
)

var (
	subdomainRegexExpression  = "[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?"
	subdomainsRegexExpression = fmt.Sprintf("(?:%s.)*", subdomainRegexExpression)
)

func getIPByTCPURL(URL string) net.IP {
	if addr, err := url.Parse(URL); err == nil {
		// no errors -> ip is the hostname
		return net.ParseIP(addr.Hostname())
	}
	// some errors occured -> cannot get IP
	return nil
}
