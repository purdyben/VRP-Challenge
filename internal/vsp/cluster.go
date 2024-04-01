package vsp

import (
	"fmt"
	"math"

	"github.com/muesli/clusters"
)

type ClusterPoint struct {
	Load
}

func NewClusterPoint(l Load) ClusterPoint {
	return ClusterPoint{l}
}

type Cluster struct {
	P []ClusterPoint
}

func (c *Cluster) Center() []float64 {
	return clusterCenter(c)
}

func (c *Cluster) Points() []ClusterPoint {
	return c.P
}

func (c *Cluster) Loads() []Load {
	l := []Load{}
	for _, p := range c.Points() {
		l = append(l, p.Load)
	}
	return l
}

func pointDistence(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.X()-p1.X(), 2) + math.Pow(p2.Y()-p1.Y(), 2))
}

// merge cluster mergers points based on there center coordinates
func MergeCluster(loads []Load, threshold float64) []Cluster {
	clusters := make([]Cluster, len(loads))

	for i, l := range loads { // start with len(loads) number ofo clusters
		clusters[i] = Cluster{
			P: []ClusterPoint{NewClusterPoint(l)},
		}
	}
	for {
		newClusters := mergeClusters(clusters, threshold)
		if len(newClusters) == len(clusters) {
			break
		}
		clusters = newClusters
	}
	return clusters
}

func mergeClusters(clusters []Cluster, threshold float64) []Cluster {
	var mergedClusters []Cluster
	for _, cluster := range clusters {
		merged := false

		for i, mergedCluster := range mergedClusters {
			// distence between centers
			if pointDistence(cluster.Center(), mergedCluster.Center()) <= threshold {
				mergedClusters[i].P = append(mergedClusters[i].P, cluster.P...) // merge

				merged = true
				break
			}
		}

		if !merged {
			mergedClusters = append(mergedClusters, cluster)
		}
	}
	return mergedClusters
}

func calculateDistance(point1, point2 []float64) float64 {
	return math.Sqrt(math.Pow(point1[0]-point2[0], 2) + math.Pow(point1[1]-point2[1], 2))
}

// get cluster center
func clusterCenter(cluster *Cluster) []float64 {
	var sumX, sumY float64
	for _, point := range cluster.Points() {
		m := Middle(point.Pickup, point.Dropoff)
		sumX += m[0]
		sumY += m[1]
	}
	return []float64{sumX / float64(len(cluster.Points())), sumY / float64(len(cluster.Points()))}
}

// testing function
func PrintCuster(cl []Cluster) {
	for i, c := range cl {
		for _, p := range c.Points() {
			fmt.Println(i, p.LoadNumber)
		}
		fmt.Println("")
	}
}

// Another idea is to use Kmeans clustering, this has a defined cluster number which could be an issue

// kmean github.com/muesli/clusters package wrapper
type KmeansClusterObservable struct {
	Load
}

func NewClusterObservable(l Load) KmeansClusterObservable {
	return KmeansClusterObservable{l}
}

func (c KmeansClusterObservable) Coordinates() clusters.Coordinates {
	return Middle(c.Pickup, c.Dropoff)
}

func (c KmeansClusterObservable) Distance(point clusters.Coordinates) float64 {
	return c.GetDistance()
}

func (c KmeansClusterObservable) Data() Load {
	return c.Load
}
