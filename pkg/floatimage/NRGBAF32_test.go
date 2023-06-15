package floatimage

import (
	"floatimage/pkg/floatcolor"
	"math"
	"testing"
)

func TestNRGBAF32(t *testing.T) {
	nrgbaf32 := NewNRGBAF32(100, 100)

	width := nrgbaf32.Bounds().Dx()
	height := nrgbaf32.Bounds().Dy()
	diagonalMax := float32(math.Sqrt(float64(width*width + height*height)))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			diagonal := math.Sqrt(float64(x*x + y*y))

			c := floatcolor.NRGBAF32{
				R: float32(x) / float32(width),
				G: float32(y) / float32(height),
				B: float32(diagonal) / diagonalMax,
				A: float32(diagonal) / diagonalMax,
			}

			nrgbaf32.Set(x+nrgbaf32.Bounds().Min.X, y+nrgbaf32.Bounds().Min.Y, c)
		}
	}

	writeImage("../testresult/NRGBAF32.png", nrgbaf32)
	writeImage("../testresult/NRGBAF32_as_NRGBA.png", nrgbaf32.AsNRGBA())
	writeImage("../testresult/NRGBAF32_as_RGBA.png", nrgbaf32.AsRGBA())
}
