// +build !libnova

package celestial

import "log"

func init() {
	log.Println("❌ Provider: Celestial positions unavailable.")
}
