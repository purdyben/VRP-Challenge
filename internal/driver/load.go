package driver

import (
	"fmt"
	"math"
	"sort"

	"github.com/dougwatson/Go/v3/math/geometry"
)

var (
	startnode = []float64{0, 0}
	endnode   = []float64{0, 0}
)

type Load struct {
	LoadNumber int
	Pickup     []float64
	Dropoff    []float64
}

// Note including (0,0)

type Point []float64

func (p Point) X() float64 {
	return p[0]
}

func (p Point) Y() float64 {
	return p[1]
}

func (p Point) Equal(p2 Point) bool {
	return p.X() == p2.X() && p.Y() == p2.Y()
}

func CalcLoadsDistence(path []Load) float64 {
	points := [][]float64{startnode}
	for _, l := range path {
		points = append(points, l.Pickup, l.Dropoff)
	}
	points = append(points, startnode)
	return CalcDistence(points)
}

func CalcDistence(points [][]float64) float64 {
	var dis float64
	for i := 0; i < len(points)-1; i += 1 {
		d, err := geometry.EuclideanDistance(points[i], points[i+1])
		if err != nil {
			panic(err)
		}
		dis += d
	}
	return dis
}

func CalcLoadsPrint(path []Load) float64 {
	points := [][]float64{startnode}
	for _, l := range path {
		points = append(points, l.Pickup, l.Dropoff)
	}
	points = append(points, startnode)
	s := ""
	for i := 0; i < len(points); i += 1 {
		s += fmt.Sprintf("(%.2f,%.2f)", points[i][0], points[i][1]) + " -> "
	}
	fmt.Println(s)
	return CalcDistence(points)
}

// 0,0 -> (0.3,8.9) (40.9,55.0) -> (-24.5,-19.2) (98.5,1,8) -> (5.3,-61.1) (77.8,-5.4) -> 0,0
// 00 -> (-50.1,80.0) -> (90.1,12.2) -> 0,0

func Middle(p1, p2 []float64) []float64 {
	return []float64{(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2}
}

type PointWithDistance struct {
	point       []float64
	minDistance float64
}

func distance(p1, p2 []float64) float64 {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	return math.Sqrt(dx*dx + dy*dy)
}

func calculateMinDistance(point Point, points []Point) float64 {
	minDist := math.Inf(1)
	for _, p := range points {
		if p[0] != point[0] && p[1] != point[1] {
			dist := distance(point, p)
			if dist < minDist {
				minDist = dist
			}
		}
	}
	return minDist
}

func Sort(startpoint []float64, loads []Load) []Load {
	sort.Slice(loads, func(i, j int) bool {
		d1, _ := geometry.EuclideanDistance(startpoint, loads[i].Pickup)
		d2, _ := geometry.EuclideanDistance(startpoint, loads[j].Pickup)
		return d1 < d2
	})
	return loads
}
