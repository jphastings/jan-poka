package main

import (
	_ "embed"
	"github.com/jphastings/jan-poka/pkg/display"
	"github.com/jphastings/jan-poka/pkg/math"
	"log"
)

import (
	"github.com/jphastings/jan-poka/pkg/images"
	"github.com/jphastings/jan-poka/pkg/projections"
	_ "image/jpeg"
	_ "image/png"
)

//go:embed earth_lights.jpeg
var earthLights []byte

func main() {
	scr, err := display.Get()
	if err != nil {
		log.Fatalf("Could not get a display: %v", err)
	}

	im := images.New(scr, projections.Winkel)

	im.SetAnchor(0, 640, 0)
	im.SetAnchor(1, 1280, 0)
	im.SetAnchor(2, 1280, 1080)
	im.SetAnchor(3, 640, 1080)

	//night, _, err := image.Decode(bytes.NewBuffer(earthLights))
	//if err != nil {
	//	log.Fatalf("Could not decode the earth at night image: %v", err)
	//}
	//if err := im.ShowImage(night, projections.Equirectangular); err != nil {
	//	log.Fatalf("Could not decode the earth at night image: %v", err)
	//}

	im.ShowPoint(math.LLACoords{Latitude: 0, Longitude: 0})

	defer scr.Close()
	scr.Loop()
}
