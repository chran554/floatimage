package floatcolor

import "math"

func clampF32(v float32, min float32, max float32, roundToInteger bool) float32 {
	if v > max {
		v = max
	} else if v < min {
		v = min
	}

	if roundToInteger {
		return float32(math.Round(float64(v)))
	} else {
		return v
	}
}

func clampF64(v float64, min float64, max float64, roundToInteger bool) float64 {
	if v > max {
		v = max
	} else if v < min {
		v = min
	}

	if roundToInteger {
		return math.Round(v)
	} else {
		return v
	}
}
