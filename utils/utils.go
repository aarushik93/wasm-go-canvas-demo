package utils


import "math"

func Overlapping(x1, y1, x2, y2 float64) bool {
	if x1 > x2 || x2 > x1 {
		return false
	}

	// If one Shape is above other
	if y1 < y2 || y2 < y1 {
		return false
	}

	return true
}

func Distance(x1, y1, x2,y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2,2)+math.Pow(y1-y2,2))
}