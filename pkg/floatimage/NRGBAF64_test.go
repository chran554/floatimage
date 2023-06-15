package floatimage

import (
	"floatimage/pkg/floatcolor"
	"math"
	"testing"
)

func TestNRGBAF64(t *testing.T) {
	nrgbaf64 := NewNRGBAF64(100, 100)

	width := nrgbaf64.Bounds().Dx()
	height := nrgbaf64.Bounds().Dy()
	diagonalMax := math.Sqrt(float64(width*width + height*height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			diagonal := math.Sqrt(float64(x*x + y*y))

			c := floatcolor.NRGBAF64{
				R: float64(x) / float64(width),
				G: float64(y) / float64(height),
				B: diagonal / diagonalMax,
				A: diagonal / diagonalMax,
			}

			nrgbaf64.Set(x+nrgbaf64.Bounds().Min.X, y+nrgbaf64.Bounds().Min.Y, c)
		}
	}

	writeImage("../testresult/NRGBAF64.png", nrgbaf64)
	writeImage("../testresult/NRGBAF64_as_NRGBA.png", nrgbaf64.AsNRGBA())
	writeImage("../testresult/NRGBAF64_as_RGBA.png", nrgbaf64.AsRGBA())
}
