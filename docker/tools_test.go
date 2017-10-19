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
			t.Errorf("Expected %s, got %s", tt.expectedIP, ip)
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
			t.Errorf("Expected %t, got %t", tt.valid, actual)
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
			t.Errorf("Expected %t, got %t", tt.valid, actual)
		}
	}
}
