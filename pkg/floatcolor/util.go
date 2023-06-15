package floatcolor

import (
	"image/color"
	"math"
)

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func mix(c1 color.Color, c2 color.Color, mix float64) color.Color {
	cc1 := NRGBAF64Model.Convert(c1).(NRGBAF64)
	cc2 := NRGBAF64Model.Convert(c2).(NRGBAF64)

	R := cc1.R*(1.0-mix) + cc2.R*mix
	G := cc1.G*(1.0-mix) + cc2.G*mix
	B := cc1.B*(1.0-mix) + cc2.B*mix
	A := cc1.A*(1.0-mix) + cc2.A*mix

	precise := cc1.Precise || cc2.Precise

	return NRGBAF64{R: R, G: G, B: B, A: A, Precise: precise}
}

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
