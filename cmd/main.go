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
	"vorto/internal/greedy"

	"github.com/dougwatson/Go/v3/math/geometry"
	"github.com/spf13/viper"
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

type config struct {
	Name    string
	Ip      string
	Port    int
	LogFile string
}

func setLogging() {
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug Level Set")
	} else {
	}
}

func configParse() {
	// C := config{}
	viper.SetDefault("loglevel", "debug")
}

func init() {
	flag.Parse()
	configParse()
	setLogging()
}

func parsePoint(pointStr string) []float64 {
	pointStr = strings.Trim(pointStr, "()")
	coords := strings.Split(pointStr, ",")
	x, _ := strconv.ParseFloat(coords[0], 64)
	y, _ := strconv.ParseFloat(coords[1], 64)
	return []float64{x, y}
}

func parseLine(line string) *driver.Load {
	if len(line) == 0 { // ketch "eof or \n"
		return nil
	}
	parts := strings.Split(line, " ")
	number, _ := strconv.Atoi(parts[0])
	pickup := parsePoint(parts[1])
	dropoff := parsePoint(parts[2])
	return &driver.Load{LoadNumber: number, Pickup: pickup, Dropoff: dropoff}
}

func parseFile(file string) []driver.Load {
	parts := strings.Split(file, "\n")[1:]
	var loads []driver.Load
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
		slog.Error(err.Error())
	}
	return b, err
}

func getfiles(dir string) ([]fs.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return files, err
}

type VSP struct {
	Loads         []*driver.Load
	LoadsDistence map[int]float64
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
	vsp := &VSP{
		LoadsDistence: make(map[int]float64),
	}
	loads := parseFile(file)
	// for i := range loads {
	// 	fmt.Printf("[%d]\n", i)
	// // }
	for _, l := range loads {
		d, _ := geometry.EuclideanDistance(l.Pickup, l.Dropoff)
		vsp.LoadsDistence[l.LoadNumber] = d
	}

	// // fmt.Println(len(loads))
	c := greedy.SimpleCluster(loads, 50)
	for _, i := range c {
		loadNumbers := []int{}
		for _, p := range i.Points() {
			loadNumbers = append(loadNumbers, p.LoadNumber)
		}
		PrintAnswerStOut(loadNumbers)
		// l := []driver.Load{}
		// for _, r := range i.Points() {
		// 	l = append(l, r.Load)
		// }

		// fmt.Println(driver.CalcLoadsPrint(l))

		// fmt.Println(s + "]")
		// greedy.GreedyTSP(i.Points())
	}
}

func PrintAnswerStOut(l []int) string {
	string_integers := make([]string, len(l))
	for i, v := range l {
		string_integers[i] = fmt.Sprintf("%d", v)
	}

	// Join strings with commas
	result_string := "[" + strings.Join(string_integers, ",") + "]"
	fmt.Println(result_string)
	return result_string
}
