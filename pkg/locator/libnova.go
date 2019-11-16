// +build libnova

package locator

import "github.com/jphastings/jan-poka/pkg/locator/celestial"

func init() {
	celestial.Load()
}
