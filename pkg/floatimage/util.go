package floatimage

import (
	"image"
	"math"
	"math/bits"
)

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(channels int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(channels, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
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
