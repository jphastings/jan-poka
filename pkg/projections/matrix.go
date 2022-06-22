package projections

import (
	"gonum.org/v1/gonum/mat"
)

type Pos struct {
	X float64
	Y float64
}

type Transformer func(Pos) Pos

// CreateMatrix see: https://stackoverflow.com/a/66101610/201130
func CreateMatrix(prj Projection, points []float64) (Transformer, error) {
	p0 := prj.Normalize(prj.Anchors[0].Coords)
	p1 := prj.Normalize(prj.Anchors[1].Coords)
	p2 := prj.Normalize(prj.Anchors[2].Coords)
	p3 := prj.Normalize(prj.Anchors[3].Coords)

	u0, v0, u1, v1, u2, v2, u3, v3 := points[0], points[1], points[2], points[3], points[4], points[5], points[6], points[7]

	// The data must be arranged in row-major order, i.e. the (i*c + j)-th
	// element in the data slice is the {i, j}-th element in the matrix.
	Adata := []float64{
		p0.X, p0.Y, 1, 0, 0, 0, -p0.X * u0, -p0.Y * u0,
		p1.X, p1.Y, 1, 0, 0, 0, -p1.X * u1, -p1.Y * u1,
		p2.X, p2.Y, 1, 0, 0, 0, -p2.X * u2, -p2.Y * u2,
		p3.X, p3.Y, 1, 0, 0, 0, -p3.X * u3, -p3.Y * u3,
		0, 0, 0, p0.X, p0.Y, 1, -p0.X * v0, -p0.Y * v0,
		0, 0, 0, p1.X, p1.Y, 1, -p1.X * v1, -p1.Y * v1,
		0, 0, 0, p2.X, p2.Y, 1, -p2.X * v2, -p2.Y * v2,
		0, 0, 0, p3.X, p3.Y, 1, -p3.X * v3, -p3.Y * v3,
	}

	A := mat.NewDense(8, 8, Adata)
	bb := mat.NewVecDense(8, []float64{u0, u1, u2, u3, v0, v1, v2, v3})
	result := mat.VecDense{}
	var err error
	err = result.SolveVec(A, bb)
	if err != nil {
		return nil, err
	}

	return func(pos Pos) Pos {
		den := result.At(6, 0)*pos.X + result.At(7, 0) + 1

		return Pos{
			X: (result.At(0, 0)*pos.X + result.At(1, 0)*pos.Y + result.At(2, 0)) / den,
			Y: (result.At(3, 0)*pos.X + result.At(4, 0)*pos.Y + result.At(5, 0)) / den,
		}
	}, nil
}
