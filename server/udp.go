package server

import (
	"fmt"
	"net"
	"runtime/debug"

	"golang.org/x/net/dns/dnsmessage"

	"bitbucket.org/agurinov/dnskek/docker"
	"bitbucket.org/agurinov/dnskek/log"
)

type Packet struct {
	addr    *net.UDPAddr
	message []byte
}

type Server struct {
	// server configuration
	conn *net.UDPConn
	ch   chan Packet
	// DNS message options
	ttl uint32
	// Docker machines options
	registry *docker.Registry
}

func New(ip net.IP, port int, reg *docker.Registry, ttl int) *Server {
	// Phase 1. Resolve udp address and create udp server listening on provided port
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	} // cannot resolve address (invalid options (ip or port))
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	} // cannot establish connection on this addr

	log.Infof("UDP server up and running on %s", log.Wrap(fmt.Sprintf("%s", udpAddr), log.Bold, log.Blink))

	// Phase 2. Get actual registry of docker machines for resolving
	if reg == nil {
		reg = docker.NewRegistry()
	}

	// Phase 3. Create and return proxy struct server
	return &Server{
		conn:     udpConn,
		ch:       make(chan Packet),
		ttl:      uint32(ttl),
		registry: reg,
	}
}

func (s *Server) getAnswers(questions []dnsmessage.Question) (answers []dnsmessage.Resource) {
	for _, q := range questions {
		// AnswerPhase 1. Resolve docker-machine by DNS requested name
		dnsRequestedName := q.Name.String()
		dm, err := s.registry.ResolveMachineByName(dnsRequestedName)
		if err != nil {
			panic(err)
		} // machine not resolved

		// AnswerPhase 2. Based on exact hit or subdomain there may be different answers
		// for example, exact hit -> [A] but subdomain -> [CNAME, A]
		// TODO some dm method that returns hit type(CNAME, A)
		// TODO in fact move switch under machine method
		rootDNSName, _ := dnsmessage.NewName(dm.DnsName())

		switch dm.DnsName() != q.Name.String() {
		case true:
			// case when it's subdomain -> prepend CNAME row
			answers = append(
				answers,
				dnsmessage.Resource{
					dnsmessage.ResourceHeader{
						Name:  q.Name,
						Type:  dnsmessage.TypeCNAME,
						Class: dnsmessage.ClassINET,
						TTL:   s.ttl,
					},
					&dnsmessage.CNAMEResource{rootDNSName},
				},
			)
			fallthrough
		default:
			// default behavior - answer root domain ip (A type)
			answers = append(
				answers,
				dnsmessage.Resource{
					dnsmessage.ResourceHeader{
						Name:  rootDNSName,
						Type:  dnsmessage.TypeA,
						Class: dnsmessage.ClassINET,
						TTL:   s.ttl,
					},
					&dnsmessage.AResource{dm.DnsIP4()},
				},
			)
		}

	}

	// AnswerPhase 3. Just return array of `type Resource interface`
	// Later in .getResponse we will have to convert to concrete type
	return
}

func (s *Server) resolve(rawRequest []byte) (rawResponse []byte) {
	var response dnsmessage.Message
	var err error

	defer func() {
		if r := recover(); r != nil {
			switch r {
			case docker.ErrMachineNotExist, docker.ErrMachineNotLocalDriver,
				docker.ErrMachineNotRunning, docker.ErrMachineNoIP:
				// no machine -> NXDOMAIN (RCode = 3) https://tools.ietf.org/html/rfc8020
				// TODO ttl = 0 !!!!!
				response.Header.RCode = dnsmessage.RCodeNameError
			default:
				// some unexpected error, return ServFail (RCode = 2) https://tools.ietf.org/html/rfc2929#page-4
				response.Header.RCode = dnsmessage.RCodeServerFailure
				// log error with traceback
				// TODO ttl = 0 !!!!!
				log.Errorf("%s\n%s", r, debug.Stack())
			}
		} else {
			// no errors -> success status code (RCode = 0)
			response.Header.RCode = dnsmessage.RCodeSuccess
		}
		response.Header.Response = true
		rawResponse, _ = response.Pack()
		// access log regardless of the answer
		log.Infof("Q: %v\t\tA: %v\t\tSTATUS: %d", response.Questions, response.Answers, response.Header.RCode)
	}()

	// Phase 1. Parse incoming rawRequest to dnsmessage
	err = response.Unpack(rawRequest)
	if err != nil {
		panic(err)
	} // cannot parse incoming message
	// Phase 2. Try to answer
	response.Answers = s.getAnswers(response.Questions)
	// Phase 3. Convert back to raw
	rawResponse, err = response.Pack()
	if err != nil {
		panic(err)
	} // cannot pack response

	// Return invoke defer trigger, which checks whether to return DNS error message
	return rawResponse
}

// create and serve
// main entrypoit, catch all unexpected errors and return somethin like `CANNOT RESOLVE`
func (s *Server) Serve() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe

	// Phase 1. Listening channel with DNS Packets and translate them to UDP socket
	go func() {
		for {
			select {
			case p := <-s.ch:
				// Phase 4. Write response back to udp socket
				s.conn.WriteToUDP(p.message, p.addr)
			}
		}
	}()

	// Phase 2. Listen infinitely UDP connection for incoming requests
	for {
		rawRequest := make([]byte, 512)                          // buffer for incoming data
		bytesRead, fromAddr, _ := s.conn.ReadFromUDP(rawRequest) // got raw UDP packet from client(fromAddr) // TODO catch and handle error
		// Phase 3. Message received, resolve this shit concurrently!
		go func() {
			s.ch <- Packet{
				fromAddr,
				s.resolve(rawRequest[:bytesRead]),
			}
		}()
	}
}
