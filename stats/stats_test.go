package stats

import (
	"testing"
)

var (
	serie = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
)

func TestAverage(t *testing.T) {
	avg := Average(serie)
	t.Logf("Average of population: %f", avg)
	if avg != 5.5 {
		t.Fail()
	}
}

func TestStdDev(t *testing.T) {
	sd := StdDev(serie)
	t.Logf("Standard deviation of population: %f", sd)
	if Truncate(sd, 2) != 2.87 {
		t.Fail()
	}
}
