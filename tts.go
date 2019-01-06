// +build tts

package main

import (
	"fmt"
	"log"

	"github.com/jphastings/corviator/pkg/tts"
	"github.com/jphastings/corviator/pkg/tts/googletts"
)

func init() {
	if environment.UseTTS {
		ttsEngine, err := googletts.New()
		if err != nil {
			log.Fatal(err)
		}

		callbacks = append(callbacks, tts.TrackedCallback(ttsEngine))
		fmt.Println("Tracking with text-to-speech engine")
	}
}
