package mdns

import (
	"github.com/hashicorp/mdns"
	"github.com/jphastings/jan-poka/pkg/shutdown"
)

func Register(serviceType string, port int) (func() error, error) {
	textRecords := []string{}
	service, err := mdns.NewMDNSService("Jan Poka", serviceType, "", "", port, nil, textRecords)
	if err != nil {
		return nil, err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}
	shutdown.Ensure("mDNS server", server.Shutdown)

	return server.Shutdown, nil
}
