package http

import (
	"github.com/grandcat/zeroconf"
	"log"
)

func announce(port int) {
	announceSrv, err := zeroconf.Register("Jan-Poka", "_http._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Printf("✋ Bonjour could not start: %v\n", err.Error())
	}
	_ = announceSrv

	// TODO: Ensure the Bonjour annouce is ended when the app closes
	//sig := make(chan os.Signal, 1)
	//signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	//select {
	//case <-sig:
	//	announceSrv.Shutdown()
	//}
}
