package vsp

import (
	"errors"
)

// including the return to (0,0), must be less than 12*60
const MaxThreshold int = 720

var ThreadholdErr error = errors.New("threshold reached")

// returns the full path
func OptimizePath(loads []Load) ([]Load, []Load) {
	// finalLoads := [][]Load{}
	// 1. get the current toal distence
	// 2. add the closest node
	// 3 check if we can
	// 	idea look ahead if we need to return back to depo
	// 	if all nodes are not used create a new cluster and start again

	// clusters := optimizePath(remaining)

	// fmt.Println("OptimizePath", len(loads), loads)
	path := []Load{}
	currDis := float64(0)
	currPoint := startnode
	// fmt.Println(loads)
	for len(loads) > 0 {
		loads = Sort(currPoint, loads)
		// fmt.Println("sorted len", len(loads))
		if len(loads) == 0 {
			break
		}
		nextload := loads[0]

		loads = loads[1:]

		// fmt.Println(loads)

		// fmt.Println("next sorted", len(loads), currDis)
		// fmt.Println(loads)
		err := TestNextLoad(currPoint, currDis, nextload)
		if err != nil {
			loads = append(loads, nextload)
			// fmt.Println("TestNextLoad failed", loads)
			break
		}

		path = append(path, nextload)
		currDis += LookAheadDist(currPoint, nextload)
		currPoint = nextload.Dropoff
	}
	return path, loads
}

func ClosestLoad() {
}

func ClosestLoadToDepo() {
}

func LookAheadDist(curr Point, l Load) float64 {
	dispickup := EuclideanDistance(curr, l.Pickup)
	disLoad := l.GetDistance()
	return dispickup + disLoad
}

func TestNextLoad(curr Point, currDistance float64, l Load) error {
	dispickup := EuclideanDistance(curr, l.Pickup)
	disLoad := l.GetDistance()

	if currDistance+dispickup+disLoad+DistanceFromDepo(l.Dropoff) > float64(MaxThreshold) {
		return ThreadholdErr
	}
	return nil
}
