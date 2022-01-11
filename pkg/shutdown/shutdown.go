package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type shutdownable struct {
	name string
	fn   func() error
}

var shutdownables []*shutdownable

func Ensure(name string, fn func() error) {
	shutdownables = append(shutdownables, &shutdownable{name: name, fn: fn})
}

// TODO: Add timeout
func Await() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	for _, s := range shutdownables {
		if err := s.fn(); err != nil {
			log.Printf("⚠️ Could not shutdown %s\n", s.name)
		}
	}
}
