package docker

import (
	"net"
	"regexp"
	"testing"
)

func TestGetIPByTCPURL(t *testing.T) {
	if ip := getIPByTCPURL(""); ip != nil {
		t.Errorf("Expected nil, got: %v", ip)
	}

	if ip := getIPByTCPURL("tcp://golang.org"); ip != nil {
		t.Errorf("Expected nil, got: %v", ip)
	}

	if ip := getIPByTCPURL("tcp://192.168.99.101:2376"); ip == nil || !ip.Equal(net.ParseIP("192.168.99.101")) {
		t.Errorf("Expected 192.168.99.101, got: %v", ip)
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
