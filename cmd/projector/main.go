package main

import (
	"bytes"
	_ "embed"
)

import (
	"github.com/fogleman/gg"
	"github.com/jphastings/jan-poka/pkg/images"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/projections"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

const w = 1024
const h = 576

//go:embed earth_lights.jpeg
var earthLights []byte

func main() {
	im := images.New(w, h, projections.Winkel)

	im.SetAnchor(0, 512, 0)
	im.SetAnchor(1, 984, 288)
	im.SetAnchor(2, 512, 576)
	im.SetAnchor(3, 40, 288)

	night, _, err := image.Decode(bytes.NewBuffer(earthLights))
	if err != nil {
		panic(err)
	}
	if err := im.ShowImage(night, projections.Equirectangular); err != nil {
		panic(err)
	}

	im.ShowPoint(math.LLACoords{Latitude: 0, Longitude: 0})

	if err := gg.SavePNG("projector.png", im.Image()); err != nil {
		panic(err)
	}
}
