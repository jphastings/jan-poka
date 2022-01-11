package mdns

import (
	"fmt"
	janpoka "github.com/jphastings/jan-poka"
	"github.com/jphastings/jan-poka/pkg/shutdown"
	"github.com/oleksandr/bonjour"
)

var versionRecord = fmt.Sprintf("v=%s", janpoka.Version)

func Register(name, serviceType string, port int) (func() error, error) {
	srv, err := bonjour.Register(fmt.Sprintf("Jan Poka (%s)", name), serviceType, "local.", port, []string{versionRecord}, nil)
	if err != nil {
		return nil, err
	}
	srv.Shutdown()

	sd := func() error { srv.Shutdown(); return nil }
	shutdown.Ensure("mDNS server", sd)

	return sd, nil
}
