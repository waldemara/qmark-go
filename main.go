// Copyright (c) 2012-2014 Waldemar Augustyn. All rights reserved.

package main

import "fmt"
import "time"
import "math"
import "github.com/waldemara/qmark-go/qmark"

var p = fmt.Printf

func main() {
	var sum time.Duration
    p("Calculating qmark, it could take a minute... ")
	//results := qmark.RunQmark(37, 6, 3)
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
	//p("qmark.Qmark:  %d\n", qmark.Qmark())	// verify qmark() reports the same number
}
