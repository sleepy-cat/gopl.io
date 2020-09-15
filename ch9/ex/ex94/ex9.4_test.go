package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// Construct a pipeline that connects an arbitrary number of goroutines with channels.
// What is the maximum number of pip eline stages you can create without running out of
// memory? How lon g do es a value take to transit the ent ire pip eline?

// Estimate maximum number of goroutines:
//
// $ go test -bench=. -cpu 1 -args -n 0 2>&1 | tail
// 4327300
// signal: killed
// FAIL

// $ go test -bench=. -cpu 1 -args -n 4300000
// goos: linux
// goarch: amd64
// BenchmarkGiantPipeline 	       1	6270140106 ns/op
// PASS

var nFlag = flag.Int("n", -1, "number of goroutines or 0 to estimate maximum number of goroutines")

var last = make(chan struct{})
var out = last

func TestMain(m *testing.M) {
	flag.Parse()
	if *nFlag < 0 {
		flag.Usage()
		os.Exit(1)
	}

	n := *nFlag
	for i := 0; i < n || n == 0; i++ {
		in := make(chan struct{})
		go func(in, out chan struct{}) {
			for x := range in {
				out <- x
			}
		}(in, out)
		out = in
		if n == 0 {
			fmt.Println(i)
		}
	}

	os.Exit(m.Run())
}

func BenchmarkGiantPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out <- struct{}{}
		<-last
	}
}
