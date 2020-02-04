package handler

import (
	"fmt"
	"log"

	"github.com/jonathanbeber/burrow/config"
	"github.com/miekg/dns"
)

// Handler is the unique and main Handler, responsible for all the received requests
type Handler struct {
	client *dns.Client
	config config.Config
}

// NewHandler returns a Handler struct responsible for manage DNS the requests
func NewHandler(client *dns.Client, config config.Config) Handler {
	return Handler{
		client: client,
		config: config,
	}
}

// ServerDNS is the function responsible for receive the DNS requests, send it
// to upstream server and proxy the response to the clients
func (h Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	rString := ""
	for _, v := range r.Question {
		rString += v.String()
	}
	log.Printf("Received request: '%s'", rString)
	a, rtt, err := h.client.Exchange(r, fmt.Sprintf("%s:%s", h.config.UpstreamServer, h.config.UpstreamPort))
	if err != nil {
		log.Printf("failed to communicate with upstream: %s", err)
		return
	}
	log.Printf("Answer for '%s' received in %s", rString, rtt.String())
	w.WriteMsg(a)
}
