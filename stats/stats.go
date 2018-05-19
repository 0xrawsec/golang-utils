package stats

import "math"

// Trucate truncates a float with p figures of precision
func Truncate(f float64, p int) float64 {
	return float64(int(f*math.Pow10(p))) / math.Pow10(p)
}

// Average computes the average of a table of floats
func Average(floats []float64) float64 {
	var sum float64
	var cnt float64
	for _, f := range floats {
		sum += f
		cnt++
	}
	return sum / cnt
}

// StdDev returns the standard deviation of a random variable
func StdDev(floats []float64) float64 {
	var sum float64
	a := Average(floats)
	for _, f := range floats {
		sum += math.Pow(f-a, 2)
	}
	return math.Sqrt(sum / float64(len(floats)))
}
