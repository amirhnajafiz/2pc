package utils

// Average accepts an array and returns the average value of it.
func Average(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}

	sum := 0
	for _, item := range xs {
		sum += int(item)
	}

	return float64(sum / len(xs))
}
