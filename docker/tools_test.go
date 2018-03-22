package docker

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
)

func TestGetIPByTCPURL(t *testing.T) {
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
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if ip := getIPByTCPURL(tt.in); !ip.Equal(tt.out) {
				t.Errorf("Expected %q, got %q", tt.out, ip)
			}
		})
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

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if actual := p.MatchString(tt.subDomain); actual != tt.valid {
				t.Errorf("Expected \"%t\", got \"%t\"", tt.valid, actual)
			}
		})
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

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if actual := p.MatchString(tt.subDomains); actual != tt.valid {
				t.Errorf("Expected \"%t\", got \"%t\"", tt.valid, actual)
			}
		})
	}
}

func TestGetMachinesByRaw(t *testing.T) {
	tableTests := []struct {
		in  string     // raw output from exec
		out []*Machine // list of machines
	}{
		{"", []*Machine{}},
		{"intranet.lo|virtualbox|Running|tcp://192.168.99.101:2376|v17.06.0-ce", []*Machine{
			&Machine{"intranet.lo", "virtualbox", "Running", net.ParseIP("192.168.99.101"), "tcp://192.168.99.101:2376", "v17.06.0-ce"},
		}},
		{`ttt|virtualbox|Running|tcp://192.168.99.103:2376|v17.10.0-ce
yomods.lo|virtualbox|Stopped||Unknown`, []*Machine{
			&Machine{"ttt", "virtualbox", "Running", net.ParseIP("192.168.99.103"), "tcp://192.168.99.103:2376", "v17.10.0-ce"},
			&Machine{
				Name:          "yomods.lo",
				DriverName:    "virtualbox",
				State:         "Stopped",
				DockerVersion: "Unknown",
			},
		}},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := getMachinesByRaw(tt.in)
			for j := range out {
				t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
					if !reflect.DeepEqual(out[j], tt.out[j]) {
						t.Errorf("Expected %q, got %q", tt.out[j], out[j])
					}
				})
			}
		})
	}
}
