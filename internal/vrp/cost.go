package vrp

// Given a list of Loads, calculate the total cost of this trip
//
// - note: it's implied that this list of loads is
//   - 1 under the hr threshold
//   - 2 in order for the driver to do
func LoadsTotalCost(loads []Load) float64 {
	points := []Point{startnode}
	for _, l := range loads {
		points = append(points, l.Pickup, l.Dropoff)
	}

	points = append(points, startnode)

	return TotalCost(len(loads), CalcTotalDistance(points))
}

func TotalCost(drivers int, distance float64) float64 {
	return float64(500*drivers) + distance
}

func CalcTotalDistance(points []Point) float64 {
	var dis float64
	for i := 0; i < len(points)-1; i += 1 {
		d := EuclideanDistance(points[i], points[i+1])
		dis += d
	}
	return dis
}

func AddStartAndEndPoints(points []Point) []Point {
	points = append([]Point{startnode}, points...)
	points = append(points, endnode)
	return points
}
