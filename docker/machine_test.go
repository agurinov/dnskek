package docker

import (
	"net"
	"testing"
)

var dm = Machine{
	Name:          "playground.lo",
	DriverName:    "virtualbox",
	State:         "Running",
	IP:            net.ParseIP("192.168.99.110"),
	URL:           "tcp://192.168.99.110:2376",
	DockerVersion: "17.0.0.3",
}

func TestNewMachine(t *testing.T) {
	machine := newMachine(
		"playground.lo",
		"virtualbox",
		"Running",
		"tcp://192.168.99.110:2376",
		"17.0.0.3",
	)
	if machine.Name != dm.Name {
		t.Errorf("Expected %q, got %q", dm.Name, machine.Name)
	}
	if machine.DriverName != dm.DriverName {
		t.Errorf("Expected %q, got %q", dm.DriverName, machine.DriverName)
	}
	if machine.State != dm.State {
		t.Errorf("Expected %q, got %q", dm.State, machine.State)
	}
	if machine.URL != dm.URL {
		t.Errorf("Expected %q, got %q", dm.URL, machine.URL)
	}
	if machine.DockerVersion != dm.DockerVersion {
		t.Errorf("Expected %q, got %q", dm.DockerVersion, machine.DockerVersion)
	}
	if !machine.IP.Equal(dm.IP) {
		t.Errorf("Expected %q, got %q", dm.IP, machine.IP)
	}
}

func TestDnsIP4(t *testing.T) {
	if dnsIP4 := dm.DnsIP4(); dnsIP4 != [4]byte{192, 168, 99, 110} {
		t.Errorf("Expected \"%v\", got \"%v\"", [4]byte{192, 168, 99, 110}, dnsIP4)
	}
}

func TestDnsName(t *testing.T) {
	// playground.lo
	if dnsName := dm.DnsName(); dnsName != "playground.lo." {
		t.Errorf("Expected \"playground.lo\", got `%s`", dnsName)
	}

	// playground
	dm.Name = "playground"
	if dnsName := dm.DnsName(); dnsName != "playground.lo." {
		t.Errorf("Expected \"playground.lo\", got `%s`", dnsName)
	}
}
