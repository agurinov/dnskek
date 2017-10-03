package docker

// https://github.com/dgageot/docker-machine-dns/blob/master/zone/zone.go
// d.vbmOutErr("showvminfo", d.MachineName, "--machinereadable")

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrMachineNotExist       = errors.New("Machine matching query does not exist.")
	ErrMachineNotLocalDriver = errors.New("Machine driver not local.")
	ErrMachineNotRunning     = errors.New("Machine state is not `Running`.")
	ErrMachineNoIP           = errors.New("Machine has no IP.")
)

type Machine struct {
	Name string
	// active  bool
	DriverName string
	State      string
	IP         net.IP
	URL        string
	// swarm
	DockerVersion string
	// TODO ERRORS errors
}

func newMachine(Name, DriverName, State, URL, DockerVersion string) *Machine {
	return &Machine{
		IP:            getIPByTCPURL(URL),
		Name:          Name,
		DriverName:    DriverName,
		State:         State,
		DockerVersion: DockerVersion,
		URL:           URL,
	}
}

func (dm *Machine) DnsName() string {
	// TODO zone string as param
	zone := "lo"

	if strings.HasSuffix(dm.Name, "."+zone) {
		return fmt.Sprintf("%s.", dm.Name)
	} else {
		return fmt.Sprintf("%s.%s.", dm.Name, zone)
	}
}

func (dm *Machine) DnsIP4() [4]byte {
	ip := dm.IP.To4()
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
