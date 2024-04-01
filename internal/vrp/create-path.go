package vrp

import (
	"errors"
	"math"

	"gopkg.in/karalabe/cookiejar.v2/collections/stack"
)

// including the return to (0,0), must be less than 12*60
const MaxDistance = float64(720)

var ThreadholdErr error = errors.New("threshold reached")

// Idea 1
// Simple idea, Fron a given set of close points, pick the closest load,
// This solution maximizes the driver's capacity as if the closest node cannot be added return to depo.
//
// - note: This does not guarantee to use all loads
// - return (driver path, leftover loads )
func OptimizeClosetPath(loads []Load) ([]Load, []Load) {
	// 1. get the current total distance
	// 2. add the closest load
	// 	idea look ahead if we need to return to depo
	// 	if all nodes are not used create a new cluster and start again

	path := []Load{}
	currDis := float64(0)
	currPoint := startnode // <- current starting point is 0,0

	for len(loads) > 0 {
		loads = Sort(currPoint, loads) //<- Sort see Load.go, sorts by closes load
		if len(loads) == 0 {
			break
		}
		nextload := loads[0]

		loads = loads[1:]
		err := TestNextLoadViability(currPoint, currDis, nextload)
		if err != nil {
			loads = append(loads, nextload)
			break
		}

		path = append(path, nextload)
		currDis += DistanceNextLoad(currPoint, nextload)
		currPoint = nextload.Dropoff
	}
	return path, loads
}

func DistanceNextLoad(curr Point, l Load) float64 {
	dispickup := EuclideanDistance(curr, l.Pickup)
	disLoad := l.GetDistance()
	return dispickup + disLoad
}

// Test if you can the next load or if you need to return to depo
func TestNextLoadViability(curr Point, currDistance float64, l Load) error {
	dispickup := EuclideanDistance(curr, l.Pickup)
	disLoad := l.GetDistance()

	// distance to return to the depo
	if currDistance+dispickup+disLoad+DistanceFromDepo(l.Dropoff) > MaxDistance {
		return ThreadholdErr
	}
	return nil
}

// =============== i'm leaving in failed ideas to showcase other ideas

// Idea 2 Fail
// 1. Go to the farthest nodes first and work backward.
func OptimizeFurthestPath(loads []Load) ([]Load, []Load) {
	s := stack.New()
	for i := 0; i < len(loads); i++ {
		s.Push(loads[i])
	}
	path := []Load{}
	currDis := float64(0)
	currPoint := startnode // <- current starting point is 0,0

	// Get furthest first
	loads = Sort(currPoint, loads) //<- sort by closes node first
	nextload := loads[len(loads)-1]

	loads = loads[:len(loads)-1]
	err := TestNextLoadViability(currPoint, currDis, nextload)
	if err != nil {
		loads = append(loads, nextload)
	}
	path = append(path, nextload)
	currDis += DistanceNextLoad(currPoint, nextload)
	currPoint = nextload.Dropoff

	for !s.Empty() {
		nextload := s.Pop().(Load)

		err := TestNextLoadViability(currPoint, currDis, nextload)
		if err != nil {
			// loads = append(loads, nextload)
			break
		}

		path = append(path, nextload)
		currDis += DistanceNextLoad(currPoint, nextload)
		currPoint = nextload.Dropoff
	}
	return path, loads
}

type Bucket struct {
	Loads     []Load
	currDis   float64
	currPoint Point
}

// Idea 3 Fail
// Given a Load add it to a bucket which returns the lost cost,
func BucketsTest(driverNum int, loads []Load) ([][]Load, error) {
	s := stack.New()
	for i := 0; i < len(loads); i++ {
		s.Push(loads[i])
	}

	res := [][]Load{}
	buckets := []*Bucket{}
	for range driverNum {
		buckets = append(buckets, &Bucket{
			Loads:     []Load{},
			currPoint: startnode,
		})
	}

	for !s.Empty() {
		nextload := s.Pop().(Load)
		selectedBucket := -1
		selectedBucketDistance := math.MaxFloat64

		for i := range buckets {

			b := buckets[i]
			dis := GetDistanceWithNextNode(b.currDis, b.currPoint, nextload)
			if dis+DistanceFromDepo(nextload.Dropoff) > MaxDistance {
				continue
			}
			if dis < selectedBucketDistance {
				selectedBucketDistance = dis
				selectedBucket = i
			}
		}
		if selectedBucket == -1 {
			return nil, errors.New("unable to proceed with this number of buckets")
		}

		b := buckets[selectedBucket]
		b.Loads = append(b.Loads, nextload)
		b.currDis = GetDistanceWithNextNode(b.currDis, b.currPoint, nextload)
		b.currPoint = nextload.Dropoff
	}

	for i := range buckets {
		res = append(res, buckets[i].Loads)
	}
	return res, nil
}

func GetDistanceWithNextNode(currDis float64, curr Point, l Load) float64 {
	dispickup := EuclideanDistance(curr, l.Pickup)
	disLoad := l.GetDistance()
	return currDis + dispickup + disLoad
}

// Idea 4 Fail
// Did not improve results under 8 loads in a row.
func OptimizePath(loads []Load) []Load {
	// Optimize path via permutations
	opLoads := loads
	newpaths := permutations(loads)
	cost := math.MaxFloat64
	for _, p := range newpaths {
		path := ToPath(p)
		path = AddStartAndEndPoints(path)
		if c := TotalCost(1, CalcTotalDistance(path)); c < cost {
			opLoads = p
			cost = c
		}
	}
	return opLoads
}

func permutations(arr []Load) [][]Load {
	var helper func([]Load, int)
	res := [][]Load{}

	helper = func(arr []Load, n int) {
		if n == 1 {
			tmp := make([]Load, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
