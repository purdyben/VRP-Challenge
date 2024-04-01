/*
Benjamin Purdy Vorto Algorithmic Challenge Submission
// 48622.74570144862
// 48643.491499678574
// 47759.118072847494
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	debug *bool = flag.Bool("d", false, "sets log level to debug")
)

// Value Container total cost of all paths and the point out for evaluation:
// Example:
// Cost:47759.118072847494
// PrintOut: "[1]
// [4,2]
// [3]"
type Result struct {
	Cost     float64
	PrintOut string
}

func init() {
	flag.Parse()
	setLogging()
}

func main() {
	dirpath := os.Args[1] // Get Dir From Args

	b, err := readFile(dirpath)
	if err != nil {
		panic(err)
	}

	loads := parseFile(string(b))

	results := make(chan Result, 100)
	// Exit closes results chan
	exit := make(chan bool)

	answer := Result{
		Cost:     math.MaxFloat64,
		PrintOut: "",
	}

	// routine which compares all results
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

	// Test 1 Using heuristic clustering
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		EvalMergeClustering(results, loads)
	}()

	// Test 2 Kmeans Clustering
	wg.Add(1)
	go func(ch chan Result) {
		defer wg.Done()
		EvalClusteringKmeans(results, loads)
	}(results)

	wg.Wait()
	exit <- true
	fmt.Println(strings.Trim(answer.PrintOut, "\n"))
	// fmt.Println(answer.Cost)
}

/**
 * Given a cluster threshold create driver paths and return results to the results channel
 * loads: Entire file input
 */
func EvalMergeClustering(ch chan Result, loads []vsp.Load) {
	var wg sync.WaitGroup
	for i := range 350 {
		wg.Add(1)
		go func(threshold int) {
			defer wg.Done()

			// for safty make a copy :)
			clone := make([]vsp.Load, len(loads))
			copy(clone, loads)

			// each cluster is a subset of points near each other
			clusters := vsp.MergeCluster(clone, float64(threshold))

			driverRoutes := [][]vsp.Load{}
			for _, c := range clusters {

				// create driver paths from subset
				driverPaths := RecursivelyComputePath(c.Loads())

				// unpack
				for _, l := range driverPaths {
					driverRoutes = append(driverRoutes, l)
				}
			}

			// evaluate results
			ch <- EvalResult(driverRoutes)

			for i := range 3 {
				c := vsp.Copy(driverRoutes)
				ch <- EvalResult(CombineJobs(c, i)) // optimization test
			}
		}(10 + i)
	}
	wg.Wait()
}

/**
 * Given a cluster threshold create driver paths and return results to results channel
 * loads: Entire file input
 */
func EvalClusteringKmeans(ch chan Result, loads []vsp.Load) {
	for i := 1; (i) < int(math.Max(float64((len(loads)%3)), 4)); i++ { // test multiple number of clusters
		var d clusters.Observations
		for _, l := range loads {
			d = append(d, vsp.NewClusterObservable(l)) // wrapper for Observations interface
		}

		km, err := kmeans.NewWithOptions(0.05, nil)
		clusters, err := km.Partition(d, i)
		if err != nil {
			return
		}
		driverRoutes := [][]vsp.Load{}

		// n := make(map[int]int)

		for _, c := range clusters { // Get Nodes from Cluster

			l := []vsp.Load{}
			for _, o := range c.Observations {
				loadData := o.(vsp.KmeansClusterObservable).Data() // unwrap load data

				l = append(l, loadData)
				// TODO remove duplicate checker
				// if _, ok := n[loadData.LoadNumber]; !ok {
				// 	n[loadData.LoadNumber] = 1
				// l = append(l, loadData)
				// }
			}
			// create driver paths from subset
			for _, l := range RecursivelyComputePath(l) {
				driverRoutes = append(driverRoutes, l)
			}
		}

		// evaluate results
		ch <- EvalResult(driverRoutes)

		for i := range 3 {
			c := vsp.Copy(driverRoutes)
			ch <- EvalResult(CombineJobs(c, i))
		}
	}
}

// Stored the printout and total cost, Used for Eval
func EvalResult(drivers [][]vsp.Load) Result {
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

		// Unwrap Loads to a list of points
		path := vsp.ToPath(d)
		// add start 0,0 and end 0,0 to get cost
		path = vsp.AddStartAndEndPoints(path)

		r.Cost += vsp.TotalCost(1, vsp.CalcTotalDistance(path))

		// Note drivers does not include start 0,0 and end 0,0
		if len(loadNumbers) > 0 {
			r.PrintOut += CreateEvalPrintout(loadNumbers) //+ " " + fmt.Sprint(vsp.CalcTotalDistance(path))
			if index < len(drivers) {
				r.PrintOut += "\n"
			}
		}
	}

	return r
}

// Recersivly Split up loads until you have everything covered.
func RecursivelyComputePath(loads []vsp.Load) [][]vsp.Load {
	var res [][]vsp.Load
	// Creates paths based on closest nodes
	p1, p2 := vsp.OptimizeClosetPath(loads)
	if len(p2) > 0 {
		paths := RecursivelyComputePath(p2)
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

// Minior improvment by recombining short paths
// 48713.98992491605 -> 47759.118072847494 nice jump :D
// Takes a list of driver paths, and combines short paths to remove drivers and cost cost
func CombineJobs(drivers [][]vsp.Load, length int) [][]vsp.Load {
	if length == 0 {
		return drivers
	}

	tempLoads := []vsp.Load{}
	indexList := map[int]int{}
	totalDis := map[int]float64{}
	for i, d := range drivers {
		path := vsp.ToPath(d)
		path = vsp.AddStartAndEndPoints(path)
		totalDis[i] = vsp.CalcTotalDistance(path)

		if len(d) <= length {
			for _, c := range d {
				tempLoads = append(tempLoads, c)
			}
			indexList[i] = i
		}
	}

	newDrivers := [][]vsp.Load{}
	// remove short paths
	for i, l := range drivers {
		if _, ok := indexList[i]; !ok {
			newDrivers = append(newDrivers, l)
		}
	}

	// create new paths
	newLoads := RecursivelyComputePath(tempLoads)
	for _, l := range newLoads {
		newDrivers = append(newDrivers, l)
	}
	return newDrivers
}

// === File Parsing ===
func setLogging() {
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug Level Set")
	}
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

func readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		panic(err)
	}
	return b, err
}

// =====  Failed tests disregards

func TestClusteringGreedyThreshhold(ch chan Result, loads []vsp.Load) {
	var wg sync.WaitGroup
	for i := range 250 {
		wg.Add(1)
		go func(threshold int) {
			defer wg.Done()
			origJSON, err := json.Marshal(loads)
			if err != nil {
				panic(err)
			}

			clone := []vsp.Load{}
			if err = json.Unmarshal(origJSON, &clone); err != nil {
				panic(err)
			}
			clusters := vsp.MergeCluster(clone, float64(threshold))

			allPaths := [][]vsp.Load{}

			for _, c := range clusters {
				buckets := 4
				var driverPaths [][]vsp.Load
				for {
					driverPaths, err = vsp.BucketsTest(buckets, c.Loads())
					if err != nil {
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
			// 	driverRoutes := RecursivelyComputePath(c.Loads())

			// 	for _, l := range driverRoutes {
			// 		driverPaths = append(driverPaths, l)
			// 	}
			// }
			// fmt.Println(EvalResult(allPaths))
			ch <- EvalResult(allPaths)
		}(10 + i)
	}
	wg.Wait()
}

func TestWithOutClustering(ch chan Result, loads []vsp.Load) {
	driverPaths := [][]vsp.Load{}
	driverRoutes := RecursivelyComputePath(loads)

	for _, l := range driverRoutes {
		driverPaths = append(driverPaths, l)
	}
	ch <- EvalResult(driverPaths)

	for i := range 3 {
		c := vsp.Copy(driverRoutes)
		ch <- EvalResult(CombineJobs(c, i)) // optimization test
	}
}
