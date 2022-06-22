package projections

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

func Reproject(img image.Image, srcPrj Projection, dstPrj Projection, dstTrans Transformer, bounds image.Rectangle) (image.Image, error) {
	if srcPrj.Reverse == nil {
		return nil, fmt.Errorf("source projection does not support reversing")
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	// TODO: Break this out
	srcTrans := func(pos Pos) Pos {
		return Pos{
			X: (pos.X/float64(w-1) - 0.5) * 360,
			Y: (pos.Y/float64(h-1) - 0.5) * -180,
		}
	}

	chain := func(pos Pos) Pos {
		s := srcTrans(pos)
		r := srcPrj.Reverse(s)
		n := dstPrj.Normalize(r)
		d := dstTrans(n)
		return d
	}

	weights := make(map[Pos]weight)

	for px := 0; px < w; px++ {
		for py := 0; py < h; py++ {
			updateWeights(weights, img.At(px, py), chain, px, py, bounds)
		}
	}

	outImg := image.NewRGBA(bounds)
	for pos, w := range weights {
		col := color.RGBA{
			R: uint8(w.r * 255 / w.total),
			G: uint8(w.g * 255 / w.total),
			B: uint8(w.b * 255 / w.total),
			A: 255,
		}
		outImg.Set(int(pos.X), int(pos.Y), col)
	}
	return outImg, nil
}

type weight struct {
	r     float64
	g     float64
	b     float64
	total float64
}

func updateWeights(weights map[Pos]weight, col color.Color, chain func(Pos) Pos, px, py int, bounds image.Rectangle) {
	cPos := chain(Pos{X: float64(px), Y: float64(py)})
	r, g, b, _ := col.RGBA()

	// Will all be overwritten by the loop
	minX := cPos.X
	minY := cPos.Y
	maxX := cPos.X
	maxY := cPos.Y

	for _, eX := range []float64{-0.5, 0.5} {
		for _, eY := range []float64{-0.5, 0.5} {
			srcPos := Pos{X: float64(px) + eX, Y: float64(py) + eY}
			dstPos := chain(srcPos)

			// TODO: This is too dumb a way to determine the affected pixels
			// If the projection is at 45º, it'll include pixels that aren't covered
			if dstPos.X < minX {
				minX = dstPos.X
			} else if dstPos.X > maxX {
				maxX = dstPos.X
			}

			if dstPos.Y < minY {
				minY = dstPos.Y
			} else if dstPos.Y > maxY {
				maxY = dstPos.Y
			}
		}
	}

	// Limit to output image bounds
	minXi := math.Floor(minX)
	if minXi < 0 {
		minXi = 0
	}
	maxXi := math.Ceil(maxX)
	if maxXi > float64(bounds.Max.X) {
		maxXi = float64(bounds.Max.X)
	}
	minYi := math.Floor(minY)
	if minYi < 0 {
		minYi = 0
	}
	maxYi := math.Ceil(maxY)
	if maxYi > float64(bounds.Max.Y) {
		maxYi = float64(bounds.Max.Y)
	}

	for x := minXi; x <= maxXi; x++ {
		for y := minYi; y <= maxYi; y++ {
			dx := x - cPos.X
			dy := y - cPos.Y
			d := math.Sqrt(dx*dx + dy*dy)

			// Ensure we're never dividing by zero
			s := math.MaxFloat64
			if d != 0 {
				s = 1 / d
			}

			key := Pos{X: x, Y: y}
			w := weights[key]

			// Scale to "full red" being 1 * scale factor
			w.r += s * float64(r) / 65535
			w.g += s * float64(g) / 65535
			w.b += s * float64(b) / 65535
			w.total += s

			weights[key] = w
		}
	}
}
