package images

import (
	"github.com/fogleman/gg"
	"github.com/jphastings/jan-poka/pkg/display"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/projections"
	"image"
	"image/color"
)

type ImageMap struct {
	bounds image.Rectangle

	prj       projections.Projection
	anchors   []float64
	anchorSet [4]bool
	trans     projections.Transformer

	dc     *gg.Context
	target display.Drawable

	PointColor color.Color
	PointSize  float64
}

func New(d display.Drawable, prj projections.Projection) *ImageMap {
	return &ImageMap{
		prj:     prj,
		anchors: make([]float64, 8),
		target:  d,
	}
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
	im.target.DrawPoint(im.trans(im.prj.Normalize(target)))
}
