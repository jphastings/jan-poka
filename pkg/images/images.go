package images

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/projections"
	"image"
	"image/color"
	"image/draw"
)

type ImageMap struct {
	bounds image.Rectangle

	prj       projections.Projection
	anchors   []float64
	anchorSet [4]bool
	trans     projections.Transformer

	dc     *gg.Context
	target draw.Image

	PointColor color.Color
	PointSize  float64
}

func New(bounds image.Rectangle, prj projections.Projection) *ImageMap {
	dc := gg.NewContext(bounds.Max.X, bounds.Max.Y)
	dc.SetRGB255(0, 0, 0)
	dc.Clear()

	return &ImageMap{
		bounds:  bounds,
		prj:     prj,
		anchors: make([]float64, 8),
		dc:      dc,

		PointColor: color.RGBA{R: 255, A: 255},
		PointSize:  2,
	}
}

func (im *ImageMap) Target(img draw.Image) {
	im.target = img
}

func (im *ImageMap) SetAnchor(n, x, y int) bool {
	im.anchorSet[n] = true
	im.anchors[n*2] = float64(x)
	im.anchors[n*2+1] = float64(y)

	if !im.anchorSet[0] || !im.anchorSet[1] || !im.anchorSet[2] || !im.anchorSet[3] {
		return false
	}

	trans, err := projections.CreateMatrix(im.prj, im.anchors)
	if err != nil {
		return false
	}

	im.trans = trans
	return true
}

func (im *ImageMap) ShowImage(img image.Image, prj projections.Projection) error {
	dstImg, err := projections.Reproject(img, prj, im.prj, im.trans, im.bounds)
	if err != nil {
		return err
	}

	im.dc.DrawImage(dstImg, 0, 0)
	return nil
}

func (im *ImageMap) ShowPoint(target math.LLACoords) {
	im.dc.SetColor(im.PointColor)
	screenPos := im.trans(im.prj.Normalize(target))
	im.dc.DrawPoint(screenPos.X, screenPos.Y, im.PointSize)
	im.dc.Fill()
}

func (im *ImageMap) Image() image.Image {
	return im.dc.Image()
}

func (im *ImageMap) Draw() error {
	if im.target == nil {
		return fmt.Errorf("no target defined")
	}

	draw.Draw(im.target, im.bounds, im.dc.Image(), image.Point{}, draw.Src)
	return nil
}
