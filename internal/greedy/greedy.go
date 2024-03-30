package greedy

import (
	"fmt"
	"math"

	"vorto/internal/driver"

	"github.com/dougwatson/Go/v3/math/geometry"
)

type Point struct {
	driver.Load
}

func newPoint(l driver.Load) Point {
	return Point{l}
}

type Cluster struct {
	points []Point
}

func (c *Cluster) Points() []Point {
	return c.points
}

func (c *Cluster) AdPdPoint(p Point) {
	c.points = append(c.points, p)
}

// func (c *Cluster) TotalDistence() float64 {
// 	d := float64(0)
// 	// for _, p := range c.points {
// 	// 	d += EuclideanDistance
// 	// }
// 	return d
// }

// Merge c2 into c1.
// return err if not posable
// func (c1 *Cluster) merge(c2 *Cluster, threshhold float64) error {
// 	return nil
// }

// Cluster(points []Point,12*60):
// threshold 12*60
func SimpleCluster(loads []driver.Load, threshold float64) []Cluster {
	// Try start with N number of cluster and merge until unable to.

	clusters := make([]Cluster, len(loads))

	loads = driver.Sort([]float64{0, 0}, loads)
	for i, l := range loads {
		clusters[i] = Cluster{points: []Point{newPoint(l)}}
	}

	for {
		merged := false
		for i := range clusters {
			for j := i + 1; j < len(clusters); j++ {
				centroid1 := calculateCentroid(clusters[i])
				centroid2 := calculateCentroid(clusters[j])

				distance := calculateDistance(centroid1, centroid2)

				if distance < threshold {
					clusters[i].points = append(clusters[i].points, clusters[j].points...)
					clusters = append(clusters[:j], clusters[j+1:]...)
					merged = true
					break
				}
			}
			if merged {
				break
			}
		}
		if !merged {
			break
		}
	}

	return clusters
}

func calculateDistance(point1, point2 []float64) float64 {
	return math.Sqrt(math.Pow(point1[0]-point2[0], 2) + math.Pow(point1[1]-point2[1], 2))
}

func calculateCentroid(cluster Cluster) []float64 {
	var sumX, sumY float64
	for _, point := range cluster.points {
		m := driver.Middle(point.Pickup, point.Dropoff)
		sumX += m[0]
		sumY += m[1]
	}
	return []float64{sumX / float64(len(cluster.points)), sumY / float64(len(cluster.points))}
}

func PrintCuster(cl []Cluster) {
	for i, c := range cl {
		for _, p := range c.Points() {
			fmt.Println(i, p.LoadNumber)
		}
		fmt.Println("")
	}
}

func GreedyTSP(loads []Point) float64 {
	for _, l := range loads {
		d, _ := geometry.EuclideanDistance([]float64{0, 0}, l.Pickup)
		fmt.Println([]float64{0, 0}, l.Pickup, d)

		// for _, l := range loads {

		// }
	}

	return 0
}

// func calculateCentroid(cluster Cluster) Point {
// 	var sumX, sumY float64
// 	for _, point := range cluster {
// 		sumX += point.X
// 		sumY += point.Y
// 	}
// 	return NewPoint(sumX/float64(len(cluster)), sumY/float64(len(cluster)))
// }

// func calculateDistance(point1, point2 Point) float64 {
// 	return math.Sqrt(math.Pow(point1.X-point2.X, 2) + math.Pow(point1.Y-point2.Y, 2))
// }

// func (d *Driver) DistenceFromDepo() (float64, error) {
// 	return geometry.EuclideanDistance(d.Position, []float64{0, 0})
// }

// func main() {
// 	points := []Point{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}}
// 	threshold := 3.0
// 	clusters := greedyCluster(points, threshold)
// 	for i, cluster := range clusters {
// 		fmt.Printf("Cluster %d: %v\n", i+1, cluster)
// 	}
// }
