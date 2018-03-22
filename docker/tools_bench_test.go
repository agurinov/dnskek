package docker

import (
	"fmt"
	"net"
	"regexp"
	"testing"
)

func BenchmarkGetIPByTCPURL(b *testing.B) {
	tableTests := []struct {
		in  string // url for parsing
		out net.IP // expected value of IP
	}{
		{"", nil},
		{"tcp://golang.org", nil},
		{"http://golang.org", nil},
		{"tcp://192.168.99.101:2376", net.ParseIP("192.168.99.101")},
	}

	for i, tt := range tableTests {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				getIPByTCPURL(tt.in)
			}
		})
	}
}

func BenchmarkSubdomainRegexExpression(b *testing.B) {
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

	for i, tt := range tableTests {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.MatchString(tt.subDomain)
			}
		})
	}
}

func BenchmarkSubdomainsRegexExpression(b *testing.B) {
	p, _ := regexp.Compile("^" + subdomainsRegexExpression + "$")

	tableTests := []struct {
		subDomains string // subdomain group for test
		valid      bool   // is subdomain valid?
	}{
		{"foo.bar.baz.", true},
		{"foo..bar.", false},
	}

	for i, tt := range tableTests {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.MatchString(tt.subDomains)
			}
		})
	}
}
