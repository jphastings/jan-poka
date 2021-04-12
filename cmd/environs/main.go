package main

import (
	"github.com/jphastings/jan-poka/pkg/environs"
	"github.com/jphastings/jan-poka/pkg/math"
)

func main() {
	e, err := environs.New("/Users/jp/src/personal/jan-poka.buntdb")
	check(err)

	//file := "/Users/jp/Downloads/allCountries/geonames-allCountries-2021-04-12T03-17.gz"
	//f, err := os.Open(file)
	//check(err)
	//check(e.BuildDB(f))

	check(e.At(math.LLACoords{
		Latitude:  51,
		Longitude: 0,
	}))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
