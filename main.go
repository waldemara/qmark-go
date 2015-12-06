package main

import (
	"fmt"
	"github.com/waldemara/qmark-go/qmark"
	"math"
	"runtime"
	"time"
)

var p = fmt.Printf

func main() {
	runtime.GOMAXPROCS(1)

	var sum time.Duration
	p("Calculating qmark, it could take a minute... ")
	results := qmark.RunQmark(qmark.CLIENTS, qmark.SERVERS, qmark.RUNS)
	for _, res := range results {
		sum += res
	}
	avg := float64(int(sum)/len(results)) / 1000000000.0 // [s]
	sumsqr := 0.0
	for _, res := range results {
		diff := float64(res)/1000000000.0 - avg
		sumsqr += diff * diff
	}
	stdev := math.Sqrt(sumsqr / float64(len(results)))
	p("completed\n")
	p("results [s]:")
	for _, res := range results {
		p("  %5.3f", float64(res)/1000000000.0)
	}
	p("\n")
	p("average [s]:  %5.3f\n", avg)
	p("stdev [s]:    %5.3f\n", stdev)
	p("qmark:        %d\n", int(1000.0/avg))
}
