package display

import (
	"github.com/jphastings/jan-poka/pkg/projections"
	"image"
)

type Drawable interface {
	ShowImage(image.Image, projections.Pos)
	DrawPoint(projections.Pos)
	Bounds() image.Rectangle
	Loop()
	Close()
}
