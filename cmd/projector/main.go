package main

import (
	"bytes"
	_ "embed"
	"github.com/fogleman/gg"
	"github.com/gonutz/framebuffer"
	"github.com/jphastings/jan-poka/pkg/rpi"
	"log"
)

import (
	"github.com/jphastings/jan-poka/pkg/images"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/projections"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

//go:embed earth_lights.jpeg
var earthLights []byte

func main() {
	fb, err := framebuffer.Open(rpi.HDMI)
	if err != nil {
		log.Fatalf("Could not connect to the framebuffer: %v", err)
	}

	im := images.New(fb.Bounds(), projections.Winkel)
	im.Target(fb)

	im.SetAnchor(0, 512, 0)
	im.SetAnchor(1, 984, 288)
	im.SetAnchor(2, 512, 576)
	im.SetAnchor(3, 40, 288)

	night, _, err := image.Decode(bytes.NewBuffer(earthLights))
	if err != nil {
		log.Fatalf("Could not decode the earth at night image: %v", err)
	}
	if err := im.ShowImage(night, projections.Equirectangular); err != nil {
		log.Fatalf("Could not decode the earth at night image: %v", err)
	}

	im.ShowPoint(math.LLACoords{Latitude: 0, Longitude: 0})
	if err := im.Draw(); err != nil {
		log.Fatalf("Couldn't display map")
	}

	gg.SavePNG("projection.png", im.Image())
}
