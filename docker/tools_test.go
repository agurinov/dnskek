package docker

import (
	"net"
	"regexp"
	"testing"
)

func TestGetIPByTCPURL(t *testing.T) {
	if ip := getIPByTCPURL(""); ip != nil {
		t.Errorf("Expecting %d, got: %d", nil, ip)
	}

	if ip := getIPByTCPURL("tcp://golang.org"); ip != nil {
		t.Errorf("Expecting %d, got: %d", nil, ip)
	}

	if ip := getIPByTCPURL("tcp://192.168.99.101:2376"); ip == nil || !ip.Equal(net.ParseIP("192.168.99.101")) {
		t.Errorf("Expecting %d, got: %d", net.ParseIP("192.168.99.101"), ip)
	}
}

func TestSubdomainRegexExpression(t *testing.T) {
	validDomains := []string{
		"foo", "foo-bar",
		"f-------------------------------------------------------------r", // max length check
	}
	for _, d := range validDomains {
		if matched, _ := regexp.MatchString("^"+subdomainRegexExpression+"$", d); matched == false {
			t.Errorf("%d is valid domain", d)
		}
	}

	invalidDomains := []string{
		"", "-foo", "foo-bar-", "foo-bar,baz",
		"f--------------------------------------------------------------r", // max length check
	}
	for _, d := range invalidDomains {
		if matched, _ := regexp.MatchString("^"+subdomainRegexExpression+"$", d); matched == true {
			t.Errorf("%d is invalid domain", d)
		}
	}
}
