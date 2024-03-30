package driver

import (
	"fmt"

	"github.com/dougwatson/Go/v3/math/geometry"
)

// Each driver starts and ends his shift at a depot located at (0,0).

// but may not exceed 12 hours of total drive time

const maxtime = 12 * 60

type Driver struct {
	Loads        []Load
	Position     []float64
	ToalDistence float64
	TotalTime    int // Seconds
}

func NewDriver() *Driver {
	return &Driver{
		Position: []float64{0, 0},
	}
}

func (d *Driver) CurrentPosition() []float64 {
	return d.Position
}

func (d *Driver) ProccessLoad(load Load) {
	disToDropoff, _ := geometry.EuclideanDistance(d.Position, load.Dropoff)
	disFromDropoffToDepo, _ := geometry.EuclideanDistance(load.Dropoff, []float64{0, 0})

	if d.ToalDistence+disFromDropoffToDepo > maxtime {
		fmt.Println("ToLong")
		return
	}
	d.Loads = append(d.Loads, load)
	d.Position = load.Dropoff
	d.ToalDistence += disToDropoff
}

func (d *Driver) DistenceFromDepo() (float64, error) {
	return geometry.EuclideanDistance(d.Position, []float64{0, 0})
}
