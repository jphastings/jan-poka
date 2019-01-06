// +build libnova

package locator

import "github.com/jphastings/corviator/pkg/locator/celestial"

func init() {
	celestial.Load()
}
