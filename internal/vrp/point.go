package vrp

import "math"

type Point []float64

func NewPoint(f []float64) Point {
	return f
}

func (p Point) X() float64 {
	return p[0]
}

func (p Point) Y() float64 {
	return p[1]
}

func (p Point) Equal(p2 Point) bool {
	return p.X() == p2.X() && p.Y() == p2.Y()
}

func EuclideanDistance(p1 Point, p2 Point) float64 {
	n := len(p1)

	if len(p2) != n {
		return -1
	}

	var total float64 = 0

	for i, x_i := range p1 {
		// using Abs since the value could be negative but we require the magnitude
		diff := math.Abs(x_i - p2[i])
		total += diff * diff
	}

	return math.Sqrt(total)
}

func (p Point) Float64() []float64 {
	return []float64{p.X(), p.Y()}
}
