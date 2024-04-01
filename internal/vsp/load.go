package vsp

import (
	"fmt"
	"math"
	"sort"
)

var (
	startnode = NewPoint([]float64{0, 0})
	endnode   = NewPoint([]float64{0, 0})
)

// File Input Load
type Load struct {
	LoadNumber int
	Pickup     Point
	Dropoff    Point
}

// Internal Distance of Load
func (l *Load) GetDistance() float64 {
	return EuclideanDistance(l.Pickup, l.Dropoff)
}

func DistanceFromDepo(point []float64) float64 {
	return EuclideanDistance(point, []float64{0, 0})
}

// this function adds a start and end node
// - note: a list of Loads won't include start, end 0,0
func CalcLoadsDistance(path []Load) float64 {
	points := [][]float64{startnode}
	for _, l := range path {
		points = append(points, l.Pickup, l.Dropoff)
	}
	points = append(points, startnode)
	return CalcDistance(points)
}

func CalcDistance(points [][]float64) float64 {
	var dis float64
	for i := 0; i < len(points)-1; i += 1 {
		dis += distance(points[i], points[i+1])
	}
	return dis
}

// sorts load about point p
func Sort(p Point, loads []Load) []Load {
	sort.Slice(loads, func(i, j int) bool {
		d1 := EuclideanDistance(p, loads[i].Pickup)
		d2 := EuclideanDistance(p, loads[j].Pickup)
		return d1 < d2
	})
	return loads
}

// turns a list of Loads into a list of points in order
func ToPath(loads []Load) []Point {
	var res []Point
	for _, l := range loads {
		res = append(res, l.Pickup, l.Dropoff)
	}
	return res
}

// Middle of two points
func Middle(p1, p2 []float64) []float64 {
	return []float64{(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2}
}

// distance between points
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

func Copy(l [][]Load) (c [][]Load) {
	for i, p := range l {
		c = append(c, make([]Load, len(p)))
		copy(c[i], p)
	}
	return
}

// Testing
// output:  0,0 -> (0.3,8.9) (40.9,55.0) -> (-24.5,-19.2) (98.5,1,8) -> (5.3,-61.1) (77.8,-5.4) -> 0,0
// 00 -> (-50.1,80.0) -> (90.1,12.2) -> 0,0
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
	return CalcDistance(points)
}

// One issue is the fact that I'm passing []Loads everywhere and intern copying the list
// TODO change every over to an inter
type Loads struct {
	L map[int]Load
}

func (l *Loads) GetLoadsList() []Load {
	return nil
}

func (l *Loads) GetLoad(i int) Load {
	return l.L[i]
}
