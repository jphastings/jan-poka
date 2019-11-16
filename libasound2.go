// +build libasound2

package main

import (
	"fmt"
	"log"

	"github.com/jphastings/jan-poka/pkg/tts"
	"github.com/jphastings/jan-poka/pkg/tts/googletts"
)

func init() {
	if environment.UseTTS {
		ttsEngine, err := googletts.New()
		if err != nil {
			log.Fatal(err)
		}

		callbacks = append(callbacks, tts.TrackedCallback(ttsEngine))
		fmt.Println("TTS tracking: on")
	} else {
		fmt.Println("TTS tracking: off")
	}
}
