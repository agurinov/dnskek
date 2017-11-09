package docker

import (
	"net"
	"regexp"
	"testing"
)

func TestGetIPByTCPURL(t *testing.T) {
	tableTests := []struct {
		url        string // url for parsing
		expectedIP net.IP // expected value of IP
	}{
		{"", nil},
		{"tcp://golang.org", nil},
		{"http://golang.org", nil},
		{"tcp://192.168.99.101:2376", net.ParseIP("192.168.99.101")},
	}

	for _, tt := range tableTests {
		if ip := getIPByTCPURL(tt.url); !ip.Equal(tt.expectedIP) {
			t.Errorf("Expected %q, got %q", tt.expectedIP, ip)
		}
	}
}

func TestSubdomainRegexExpression(t *testing.T) {
	p, _ := regexp.Compile("^" + subdomainRegexExpression + "$")

	tableTests := []struct {
		subDomain string // subdomain for test
		valid     bool   // is subdomain valid?
	}{
		{"foo", true},
		{"foo-bar", true},
		{"f-------------------------------------------------------------r", true},
		{"", false},
		{"-foo", false},
		{"foo-bar-", false},
		{"foo-bar,baz", false},
		{"f--------------------------------------------------------------r", false},
	}

	for _, tt := range tableTests {
		if actual := p.MatchString(tt.subDomain); actual != tt.valid {
			t.Errorf("Expected \"%t\", got \"%t\"", tt.valid, actual)
		}
	}
}

func TestSubdomainsRegexExpression(t *testing.T) {
	p, _ := regexp.Compile("^" + subdomainsRegexExpression + "$")

	tableTests := []struct {
		subDomains string // subdomain group for test
		valid      bool   // is subdomain valid?
	}{
		{"foo.bar.baz.", true},
		{"foo..bar.", false},
	}

	for _, tt := range tableTests {
		if actual := p.MatchString(tt.subDomains); actual != tt.valid {
			t.Errorf("Expected \"%t\", got \"%t\"", tt.valid, actual)
		}
	}
}

// func TestGetMachinesByRaw(t *testing.T) {
//     // TODO look at tableTests
//     raw := `intranet.lo|virtualbox|Running|tcp://192.168.99.101:2376|v17.06.0-ce
// jetsmarter4.lo|virtualbox|Running|tcp://192.168.99.102:2376|v17.05.0-ce
// jetsmarter.lo|virtualbox|Running|tcp://192.168.99.100:2376|v17.05.0-ce
// ttt|virtualbox|Running|tcp://192.168.99.103:2376|v17.10.0-ce
// yomods.lo|virtualbox|Stopped||Unknown`
//
//
//     t.Log("========================================")
//     ms := getMachinesByRaw(raw)
//     t.Log(ms)
//     t.Log("========================================")
//
//     for _, dm := range ms {
//         t.Log(dm.Name, dm.DriverName, )
//     }
// }
