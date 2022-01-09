package http

import (
	"github.com/grandcat/zeroconf"
	"log"
)

func announce(port int) func() {
	announceSrv, err := zeroconf.Register("Jan-Poka", "_http._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Printf("✋ Bonjour could not start: %v\n", err.Error())
		return func() {}
	}

	return func() { announceSrv.Shutdown() }
}
