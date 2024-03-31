package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"vorto/internal/driver"
	"vorto/internal/vsp"
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
		// fmt.Println(line)
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

func main() {
	path := os.Args[1] // Get Dir From Args
	b, err := ReadFile(path)
	if err != nil {
		panic(err)
	}
	proccess(string(b))
}

func proccess(file string) {
	loads := parseFile(file)
	paths := [][]vsp.Load{}
	var cost float64
	for len(loads) > 0 {
		p, r := vsp.OptimizePath(loads)
		paths = append(paths, p)
		loads = r
	}
	for _, p := range paths {
		n := []int{}
		for _, i := range p {
			n = append(n, i.LoadNumber)
		}
		PrintAnswerStOut(n)
		path := vsp.ToPath(p)
		path = vsp.AddStartAndEndPoints(path)
		cost += vsp.TotalCost(1, vsp.CalcTotalDistance(path))
	}
	fmt.Println("Costs", cost)

	fmt.Println()
	// for i := range 10 {
	loads = parseFile(file)
	for i := range 10 {
		clusters := vsp.MergeCluster(loads, float64(i*10))
		cost = 0
		for _, c := range clusters {
			p1, p2 := vsp.OptimizePath(c.Loads())
			loadNumbers := []int{}
			for _, l := range p1 {
				loadNumbers = append(loadNumbers, l.LoadNumber)
			}
			path := vsp.ToPath(p1)
			path = vsp.AddStartAndEndPoints(path)
			cost += vsp.TotalCost(1, vsp.CalcTotalDistance(path))
			PrintAnswerStOut(loadNumbers)

			loadNumbers = []int{}
			for _, l := range p2 {
				loadNumbers = append(loadNumbers, l.LoadNumber)
			}
			path = vsp.ToPath(p2)
			path = vsp.AddStartAndEndPoints(path)
			cost += vsp.TotalCost(1, vsp.CalcTotalDistance(path))
			PrintAnswerStOut(loadNumbers)
		}
		fmt.Println("Costs", cost)
	}

	// }
	// fmt.Println(c.Loads())
	// path1, r := vsp.OptimizePath(c.Loads())
	// n := []int{}
	// for _, i := range path1 {
	// 	n = append(n, i.LoadNumber)
	// }
	// PrintAnswerStOut(n)
	// n = []int{}
	// for _, i := range r {
	// 	n = append(n, i.LoadNumber)
	// }
	// PrintAnswerStOut(n)
	// }
	// for _, l := range loads {
	// 	vsp.Loads[l.LoadNumber] = l
	// }

	// set up a random two-dimensional data set (float64 values between 0.0 and 1.0)
	// var d clusters.Observations
	// // // var
	// for _, l := range loads {
	// 	// vsp.
	// 	d = append(d, vsp.KmeansClusterObservable{Load: l})
	// }

	// // Partition the data points into 16 clusters
	// km := kmeans.New()
	// clusters, _ := km.Partition(d, len(loads)/2)
	// for _, c := range clusters {
	// 	fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
	// 	fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	// }

	// for _, c := range clusters {
	// 	l := []vsp.Load{}
	// 	for _, o := range c.Observations {
	// 		a := o.(vsp.KmeansClusterObservable)
	// 		l = append(l, a.Load)
	// 		// fmt.Println(a)
	// 	}
	// 	path1, _ := vsp.OptimizePath(l)

	// 	n := []int{}
	// 	for _, i := range path1 {
	// 		n = append(n, i.LoadNumber)
	// 	}
	// 	PrintAnswerStOut(n)
	// }

	// OptimizePath
	// Cluster all of the nodes

	// Within each cluster check paths
	// 	Return all paths needed

	// c := greedy.SimpleCluster(loads, 1)
	// var L [][]driver.Load

	// for index, i := range c {
	// 	loadNumbers := []int{}
	// 	L = append(L, []driver.Load{})
	// 	for _, p := range i.Points() {
	// 		loadNumbers = append(loadNumbers, p.LoadNumber)
	// 		L[index] = append(L[index], p.Load)
	// 	}
	// 	PrintAnswerStOut(loadNumbers)
	// }
	// fmt.Println(TotalCost(L))
}

func PrintAnswerStOut(l []int) string {
	if len(l) < 1 {
		return ""
	}
	string_integers := make([]string, len(l))
	for i, v := range l {
		string_integers[i] = fmt.Sprintf("%d", v)
	}

	// Join strings with commas
	result_string := "[" + strings.Join(string_integers, ",") + "]"
	fmt.Println(result_string)
	return result_string
}

func TotalCost(load [][]driver.Load) float64 {
	cost := float64(0)

	for _, l := range load {
		cost += driver.CalcLoadsDistence(l)
	}
	total_cost := float64(500*len(load)) + cost
	return total_cost
}
