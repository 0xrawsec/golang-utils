package main

import (
	"stats"
	"testing"
)

var (
	serie = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
)

func TestAverage(t *testing.T) {
	t.Log(stats.Average(serie))
}
