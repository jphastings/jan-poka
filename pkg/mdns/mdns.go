package mdns

import (
	"fmt"

	janpoka "github.com/jphastings/jan-poka"
	"github.com/jphastings/jan-poka/pkg/shutdown"

	"github.com/grandcat/zeroconf"
)

var versionRecord = fmt.Sprintf("v=%s", janpoka.Version)

func Register(name, serviceType string, port int) (func() error, error) {
	server, err := zeroconf.Register(fmt.Sprintf("Jan Poka (%s)", name), serviceType, "local.", port, []string{versionRecord}, nil)
	if err != nil {
		return nil, err
	}
	shutdown.Ensure("mDNS server", func() error { server.Shutdown(); return nil })

	return nil, nil
}
