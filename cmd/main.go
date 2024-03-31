package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"vorto/internal/vsp"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// Cli
var (
	// Default api bindings
	// IPFlag         *string = flag.String("i", "0.0.0.0", "Server Address Default: 0.0.0.0")
	// PortFlag       *uint   = flag.Uint("p", 8001, "Server Port Default: 8001")
	// LogPathFlag    *string = flag.String("l", "", "Service Log File")
	// ProductionFlag *bool   = flag.Bool("prod", false, "Production Mode")
	debug *bool = flag.Bool("d", false, "sets log level to debug")
	// --cmd //{command to run your program}

	// problemDir *string = flag.String("problemDir", "trainingProblems", "Specify the directory for the problem")
)

func setLogging() {
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug Level Set")
	}
}

func init() {
	flag.Parse()
	setLogging()
}

func parsePoint(pointStr string) []float64 {
	pointStr = strings.Trim(pointStr, "()")
	coords := strings.Split(pointStr, ",")
	x, _ := strconv.ParseFloat(coords[0], 64)
	y, _ := strconv.ParseFloat(coords[1], 64)
	return []float64{x, y}
}

func parseLine(line string) *vsp.Load {
	if len(line) == 0 { // ketch "eof or \n"
		return nil
	}
	parts := strings.Split(line, " ")
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	pickup := parsePoint(parts[1])
	dropoff := parsePoint(parts[2])
	return &vsp.Load{LoadNumber: number, Pickup: pickup, Dropoff: dropoff}
}

func parseFile(file string) []vsp.Load {
	parts := strings.Split(file, "\n")[1:]
	var loads []vsp.Load
	for _, line := range parts {
		load := parseLine(line)
		if load == nil {
			continue
		}
		loads = append(loads, *load)

	}
	return loads
}

func ReadFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		panic(err)
	}
	return b, err
}

func getfiles(dir string) ([]fs.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}

	return files, err
}

type Result struct {
	Cost     float64
	PrintOut string
}

func main() {
	path := os.Args[1] // Get Dir From Args
	b, err := ReadFile(path)
	if err != nil {
		panic(err)
	}
	loads := parseFile(string(b))
	results := make(chan Result, 20)
	exit := make(chan bool)
	answer := Result{
		Cost:     math.MaxFloat64,
		PrintOut: "",
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		TestClusteringThreshhold(results, loads)
	}()
	go func() {
		for {
			select {
			case a := <-results:
				if a.Cost < answer.Cost {
					answer = a
				}
			case <-exit:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		TestClusteringGreedyThreshhold(results, loads)
		// buckets := 4
		// var driverPaths [][]vsp.Load
		// for {
		// 	driverPaths, err = vsp.Greedy(buckets, loads)
		// 	if err != nil {
		// 		// fmt.Println(err)
		// 		buckets += 1
		// 		continue
		// 	}
		// 	break
		// }
		// for _, l := range driverRoutes {
		// 	driverPaths = append(driverPaths, l)
		// }

		// results <- CalcResult(driverPaths)
	}()

	wg.Wait()
	exit <- true
	// fmt.Println(strings.Trim(answer.PrintOut, "\n"))
	fmt.Println(answer.Cost)

	var d clusters.Observations
	for _, l := range loads {
		d = append(d, vsp.NewClusterObservable(l))
	}

	km := kmeans.New()
	clusters, err := km.Partition(d, 4)

	for _, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	}
	for _, c := range clusters {
		for _, o := range c.Observations {
			fmt.Println(o.(vsp.KmeansClusterObservable).Data())
		}
	}
}

func TestClusteringGreedyThreshhold(ch chan Result, loads []vsp.Load) {
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(threshhold int) {
			defer wg.Done()
			origJSON, err := json.Marshal(loads)
			if err != nil {
				panic(err)
			}

			clone := []vsp.Load{}
			if err = json.Unmarshal(origJSON, &clone); err != nil {
				panic(err)
			}
			clusters := vsp.MergeCluster(clone, float64(10+i*5))

			allPaths := [][]vsp.Load{}

			for _, c := range clusters {
				buckets := 4
				var driverPaths [][]vsp.Load
				for {
					driverPaths, err = vsp.Greedy(buckets, c.Loads())
					if err != nil {
						// fmt.Println(err)
						buckets += 1
						continue
					}
					break
				}
				for _, d := range driverPaths {
					allPaths = append(allPaths, d)
				}
			}

			// // pathsOfPoints := [][]vsp.Point{}
			// driverPaths := [][]vsp.Load{}
			// for _, c := range clusters {
			// 	driverRoutes := TestRec(c.Loads())

			// 	for _, l := range driverRoutes {
			// 		driverPaths = append(driverPaths, l)
			// 	}
			// }
			// fmt.Println(CalcResult(allPaths))
			ch <- CalcResult(allPaths)
		}(10 + i)
	}
	wg.Wait()
}

func TestClusteringThreshhold(ch chan Result, loads []vsp.Load) {
	var wg sync.WaitGroup
	for i := range 350 {
		wg.Add(1)
		go func(threshhold int) {
			defer wg.Done()
			origJSON, err := json.Marshal(loads)
			if err != nil {
				panic(err)
			}

			clone := []vsp.Load{}
			if err = json.Unmarshal(origJSON, &clone); err != nil {
				panic(err)
			}
			clusters := vsp.MergeCluster(clone, float64(10+i))

			// pathsOfPoints := [][]vsp.Point{}
			driverPaths := [][]vsp.Load{}
			for _, c := range clusters {
				driverRoutes := TestRec(c.Loads())

				for _, l := range driverRoutes {
					driverPaths = append(driverPaths, l)
				}
			}

			ch <- CalcResult(driverPaths)
		}(10 + i)
	}
	wg.Wait()
}

// Stored the printout and total cost, Used for Eval
func CalcResult(drivers [][]vsp.Load) Result {
	r := Result{
		Cost:     0,
		PrintOut: "",
	}
	// pathsOfPoints := [][]vsp.Point{}
	// Each List of Loads is the route the driver needs to drive
	for index, d := range drivers {

		// Get Print Out
		loadNumbers := []int{}
		for _, l := range d {
			loadNumbers = append(loadNumbers, l.LoadNumber)
		}

		// Note drivers does not include start 0,0 and end 0,0
		if len(loadNumbers) > 0 {
			r.PrintOut += CreateEvalPrintout(loadNumbers)
			if index < len(drivers) {
				r.PrintOut += "\n"
			}
		}

		// Unwrap Loads to a list of points
		path := vsp.ToPath(d)
		// add start 0,0 and end 0,0 to get cost
		path = vsp.AddStartAndEndPoints(path)

		r.Cost += vsp.TotalCost(1, vsp.CalcTotalDistance(path))
		// pathsOfPoints = append(pathsOfPoints, path)
	}

	return r
}

// mean cost: 48713.98992491605
// mean run time: 293.98015567234586ms

// Recersivly Split up loads until you have everything covered.
func TestRec(loads []vsp.Load) [][]vsp.Load {
	var res [][]vsp.Load
	p1, p2 := vsp.OptimizeClosetPath(loads)
	if len(p2) > 0 {
		paths := TestRec(p2)
		for _, p := range paths {
			res = append(res, p)
		}
	}
	res = append(res, p1)
	return res
}

// Create Printout for Eval
func CreateEvalPrintout(l []int) string {
	if len(l) < 1 {
		return ""
	}
	string_integers := make([]string, len(l))
	for i, v := range l {
		string_integers[i] = fmt.Sprintf("%d", v)
	}

	// Join strings with commas
	result_string := "[" + strings.Join(string_integers, ",") + "]"
	return result_string
}
