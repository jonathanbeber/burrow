package server

import (
	"testing"

	"github.com/jonathanbeber/burrow/config"
)

func resetScenario() {
	server_udp, server_tcp = nil, nil
}

func TestServerInitializationBothProtocols(t *testing.T) {
	cfg := config.Config{
		EnableUDP: true,
		EnableTCP: true,
	}
	StartServers(cfg)
	if server_udp == nil || server_tcp == nil {
		t.Errorf("Expected both servers to be not nil")
	}
	resetScenario()
}

func TestServerInitializationOnlyTCP(t *testing.T) {
	cfg := config.Config{
		EnableUDP: false,
		EnableTCP: true,
	}
	StartServers(cfg)
	if server_udp != nil {
		t.Errorf("Serverr UDP should not be initialized")
	}
	if server_tcp == nil {
		t.Errorf("Serverr TCP should be initialized")
	}
	resetScenario()
}

func TestServerInitializationOnlyUDP(t *testing.T) {
	cfg := config.Config{
		EnableUDP: true,
		EnableTCP: false,
	}
	StartServers(cfg)
	if server_tcp != nil {
		t.Errorf("Serverr TCP should not be initialized")
	}
	if server_udp == nil {
		t.Errorf("Serverr UDP should be initialized")
	}
	resetScenario()
}
