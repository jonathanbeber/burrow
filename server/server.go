package server

import (
	"log"

	"github.com/jonathanbeber/burrow/config"
	"github.com/miekg/dns"
)

var server_udp *dns.Server
var server_tcp *dns.Server

// StartServer handles the servers initialization based on the config.Config
// struct received
func StartServers(c config.Config) {
	if !c.EnableTCP && !c.EnableUDP {
		log.Fatal("Neither TCP or UDP server enabled. Exiting...")
	}
	if c.EnableTCP {
		server_tcp = &dns.Server{Addr: ":53", Net: "tcp"}
		go server_tcp.ListenAndServe()
		log.Print("Started TCP server. Listening TCP/53")
	}
	if c.EnableUDP {
		server_udp = &dns.Server{Addr: ":53", Net: "udp"}
		go server_udp.ListenAndServe()
		log.Print("Started UDP server. Listening UDP/53")
	}
}

func ShutdownServers() {
	shutdownServer(server_tcp)
	shutdownServer(server_udp)
}

// ShutdownServers is responsible for verify if the received server is initilized
// and shutdown it
func shutdownServer(s *dns.Server) {
	if s == nil {
		return
	}
	if err := s.Shutdown(); err != nil {
		log.Printf("Failed to shutdown server %s", s.Net)
	}
}
