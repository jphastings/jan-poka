package mdns

import (
	"fmt"
	"net"
	"strings"

	janpoka "github.com/jphastings/jan-poka"
	"github.com/jphastings/jan-poka/pkg/shutdown"

	"github.com/grandcat/zeroconf"
)

var versionRecord = fmt.Sprintf("v=%s", janpoka.Version)

var AnnounceInterfaces []net.Interface

func SetBroadcastInterface(target string) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	if target == "" {
		AnnounceInterfaces = ifaces
		return nil
	}

	for _, ifi := range ifaces {
		if (ifi.Flags & (net.FlagUp + net.FlagMulticast)) == 0 {
			continue
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), target+"/") {
				AnnounceInterfaces = []net.Interface{ifi}
				return nil
			}
		}
	}

	return fmt.Errorf("cannot find interface")
}

func Register(name, serviceType string, port int) (func() error, error) {
	server, err := zeroconf.Register(fmt.Sprintf("Jan Poka (%s)", name), serviceType, "local.", port, []string{versionRecord}, AnnounceInterfaces)
	if err != nil {
		return nil, err
	}
	shutdown.Ensure("mDNS server", func() error { server.Shutdown(); return nil })

	return nil, nil
}
