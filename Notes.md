Vehicle Routing Problem 

type VRP interface{
	Drivers()
	Loads()
}

Unbounded number of drivers. (D)

Each load has a pickup location and a dropoff location, each specified by a Cartesian point. A driver completes a load by driving to the pickup location, picking up the load, driving to the dropoff, and dropping off the load. The time required to drive from one point to another, in minutes, is the Euclidean distance between them. That is, to drive from (x1, y1) to (x2, y2) takes sqrt((x2-x1)^2 + (y2-y1)^2) minutes.

type Point struct{
	x,y int
}


func EuclideanDistance(point1, point2 Point) float64{
	return 2*sqrt(2*50^2)
}

github.com/dougwatson/Go/v3/math/geometry

Point
EuclideanDistance


p1 -> p2 -> p3
func Distance(path []Point) float64 {
    var totalDistance float64
    for i := 0; i < len(path)-1; i++ {
        totalDistance += EuclideanDistance(path[i], path[i+1])
    }
    return totalDistance
}

The total cost of a solution is given by the formula:
total_cost = 500*number_of_drivers + total_number_of_driven_minutes 


within the duration of one 12-hour shift. Your program does not have to assess problem feasibility.
●	No problem will contain more than 200 loads.
