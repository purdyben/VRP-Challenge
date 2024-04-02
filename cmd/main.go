/*
Benjamin Purdy Vorto Algorithmic Challenge Submission
// 48622.74570144862
// 48643.491499678574
// 47759.118072847494
*/
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"vorto/internal/vrp"

	"github.com/parallelo-ai/kmeans"
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

	// Test 1 Using Merge Clustering
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

	// Note this line Print the output to be evaluated
	fmt.Println(strings.Trim(answer.PrintOut, "\n"))
	// fmt.Println(answer.Cost)
}

/**
 * Given a cluster threshold create driver paths and return results to the results channel
 * loads: Entire file input
 */
func EvalMergeClustering(ch chan Result, loads []vrp.Load) {
	n := len(loads)
	var wg sync.WaitGroup
	for i := range 350 {
		wg.Add(1)
		go func(threshold int) {
			defer wg.Done()

			// for safty make a copy :)
			clone := make([]vrp.Load, len(loads))
			copy(clone, loads)

			// each cluster is a subset of points near each other
			clusters := vrp.MergeCluster(clone, float64(threshold))

			driverRoutes := [][]vrp.Load{}
			for _, c := range clusters {

				// create driver paths from subset
				driverPaths := RecursivelyComputePath(c.Loads())

				// unpack
				for _, l := range driverPaths {
					driverRoutes = append(driverRoutes, l)
				}
			}

			// evaluate results
			ch <- EvalResult(driverRoutes, n)

			for i := range 3 {
				c := vrp.Copy(driverRoutes)

				ch <- EvalResult(CombineJobs(c, i), n)
			}
		}(10 + i)
	}
	wg.Wait()
}

/**
 * Given a cluster threshold create driver paths and return results to results channel
 * loads: Entire file input
 */
func EvalClusteringKmeans(ch chan Result, loads []vrp.Load) {
	n := len(loads)
	for i := 1; (i) < int(math.Max(float64((len(loads)%3)), 4)); i++ { // test multiple number of clusters
		var d kmeans.Observations
		for _, l := range loads {
			d = append(d, vrp.NewClusterObservable(l)) // wrapper for Observations interface
		}

		km := kmeans.New()
		clusters, err := km.Partition(d, i, 256)
		if err != nil {
			return
		}
		driverRoutes := [][]vrp.Load{}

		for _, c := range clusters { // Get Nodes from Cluster

			clusterLoads := []vrp.Load{}
			for _, o := range c.Observations {
				loadData := o.(vrp.KmeansClusterObservable).Data() // unwrap load data
				clusterLoads = append(clusterLoads, loadData)
			}
			// create driver paths from subset
			for _, l := range RecursivelyComputePath(clusterLoads) {
				driverRoutes = append(driverRoutes, l)
			}
		}
		// evaluate results
		ch <- EvalResult(driverRoutes, n)
		for i := range 3 {
			c := vrp.Copy(driverRoutes)
			ch <- EvalResult(CombineJobs(c, i), n)
		}
	}
}

// Stored the printout and total cost, Used for Eval
func EvalResult(drivers [][]vrp.Load, totalLoadNumber int) Result {
	r := Result{
		Cost:     0,
		PrintOut: "",
	}
	if !CheckValidSolution(drivers, totalLoadNumber) {
		return Result{Cost: math.MaxFloat64}
	}
	// Each List of Loads is the route the driver needs to drive
	for index, d := range drivers {

		// Get Print Out
		loadNumbers := []int{}
		for _, l := range d {
			loadNumbers = append(loadNumbers, l.LoadNumber)
		}

		// Unwrap Loads to a list of points
		path := vrp.ToPath(d)
		// add start 0,0 and end 0,0 to get cost
		path = vrp.AddStartAndEndPoints(path)

		r.Cost += vrp.TotalCost(1, vrp.CalcTotalDistance(path))

		// Note drivers does not include start 0,0 and end 0,0
		if len(loadNumbers) > 0 {
			r.PrintOut += CreateEvalPrintout(loadNumbers) //+ " " + fmt.Sprint(vrp.CalcTotalDistance(path))
			if index < len(drivers) {
				r.PrintOut += "\n"
			}
		}
	}
	return r
}

// Recersivly Split up loads until you have everything covered.
func RecursivelyComputePath(loads []vrp.Load) [][]vrp.Load {
	var res [][]vrp.Load
	// Creates paths based on closest nodes
	p1, p2 := vrp.OptimizeClosetPath(loads)
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
func CombineJobs(drivers [][]vrp.Load, length int) [][]vrp.Load {
	if length == 0 {
		return drivers
	}

	tempLoads := []vrp.Load{}
	indexList := map[int]int{}
	totalDis := map[int]float64{}
	for i, d := range drivers {
		path := vrp.ToPath(d)
		path = vrp.AddStartAndEndPoints(path)
		totalDis[i] = vrp.CalcTotalDistance(path)

		if len(d) <= length {
			for _, c := range d {
				tempLoads = append(tempLoads, c)
			}
			indexList[i] = i
		}
	}

	newDrivers := [][]vrp.Load{}
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

func parseLine(line string) *vrp.Load {
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
	return &vrp.Load{LoadNumber: number, Pickup: pickup, Dropoff: dropoff}
}

func parseFile(file string) []vrp.Load {
	parts := strings.Split(file, "\n")[1:]
	var loads []vrp.Load
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

// Make sure all loads are used and no duplicates
func CheckValidSolution(solution [][]vrp.Load, n int) bool {
	// Check for duplicates and to make sure ever load is present
	m := make(map[int]int)
	count := 0

	for _, d := range solution {
		count += len(d)
		for _, l := range d {
			if _, ok := m[l.LoadNumber]; ok {
				return false // duplicates
			}
			m[l.LoadNumber] = 1
		}
	}

	return count == n
}

// =====  Failed tests disregard

func TestWithOutClustering(ch chan Result, loads []vrp.Load) {
	driverPaths := [][]vrp.Load{}
	driverRoutes := RecursivelyComputePath(loads)

	for _, l := range driverRoutes {
		driverPaths = append(driverPaths, l)
	}

	ch <- EvalResult(driverPaths, len(loads))

	for i := range 3 {
		c := vrp.Copy(driverRoutes)
		ch <- EvalResult(CombineJobs(c, i), len(loads)) // optimization test
	}
}
