package main

import (
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	. "github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/output/mapper"
	"github.com/wroge/wgs84"
	"math"

	"gocv.io/x/gocv"
)

const (
	width =  Meters(2.663)
	wheelRadius = Meters(0.1)
)

var mappings = []mapper.Correlation{
	// Manually grabbed from the wall
	{WallPos: WallPos{Left: 15.4, Right: 8.3}, LLACoords: LLACoords{Latitude: 51.443785958312084, Longitude: -0.3030027986397818}},
	{WallPos: WallPos{Left: 13.1, Right: 19.7}, LLACoords: LLACoords{Latitude: 51.59670164318928, Longitude: -0.29752176802369784}},
	{WallPos: WallPos{Left: 80.1, Right: 29.6}, LLACoords: LLACoords{Latitude: 51.43873804142641, Longitude: 0.005657978357681665}},
	{WallPos: WallPos{Left: 62.6, Right: 87.6}, LLACoords: LLACoords{Latitude: 51.58700664688655, Longitude: 0.012619259924454437}},
	{WallPos: WallPos{Left: 43.4, Right: 20.9}, LLACoords: LLACoords{Latitude: 51.45308010006999, Longitude: -0.14747478782252177}},
}

var testPoint = mapper.Correlation{WallPos: WallPos{Left: 43.3, Right: 56.4}, LLACoords: LLACoords{Latitude: 51.53651640508379, Longitude: -0.08640148743278318}}

var projection = 27700

func calcXY(wall WallPos, width, wheelRadius Meters) (float64, float64) {
	// TODO: This calculation doesn't work because zero isn't at the centre of the wheels and because the lengths aren't in Meters
	left2 := math.Pow(float64(wall.Left), 2)
	right2 := math.Pow(float64(wall.Right), 2)
	widthSquared := math.Pow(float64(width), 2)
	wheelRadiusSquared := math.Pow(float64(wheelRadius), 2)

	Y := math.Sqrt(wheelRadiusSquared +
		right2 +
		(left2-widthSquared-right2)/2*float64(width))
	X := math.Sqrt(wheelRadiusSquared +
		left2 -
		math.Pow(Y, 2))

	return X, Y
}

func main() {
	availableProjections := wgs84.EPSG()

	proj := availableProjections.Code(projection)

	transform := wgs84.LonLat().To(proj)

	srcPoints := make([]gocv.Point2f, len(mappings))
	dstPoints := make([]gocv.Point2f, len(mappings))

	for i, coords := range mappings {
		x, y, _ := transform(float64(coords.Latitude), float64(coords.Longitude), 0)
		X, Y := calcXY(coords.WallPos, width, wheelRadius)

		srcPoints[i] = gocv.Point2f{X: float32(x), Y: float32(y)}
		dstPoints[i] = gocv.Point2f{X: float32(X), Y: float32(Y)}
	}

	fmt.Println(dstPoints)

	src := gocv.NewPoint2fVectorFromPoints(srcPoints)
	dst := gocv.NewPoint2fVectorFromPoints(dstPoints)

	m := gocv.GetAffineTransform2f(src, dst)

	fmt.Println(m)
}
