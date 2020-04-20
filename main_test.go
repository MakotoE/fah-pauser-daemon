package main

import (
	"github.com/mitchellh/go-ps"
	"testing"
)

func BenchmarkProcesses(b *testing.B) {
	// When fah is running
	// BenchmarkProcesses-8   	     188	   5893616 ns/op
	// That is 5.89 ms

	// When fah is paused
	// BenchmarkProcesses-8   	     358	   3083836 ns/op

	for i := 0; i < b.N; i++ {
		_, err := ps.Processes()
		if err != nil {
			b.Fatal(err)
		}
	}
}
