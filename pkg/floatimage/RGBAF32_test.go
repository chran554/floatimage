package floatimage

import (
	"floatimage/pkg/floatcolor"
	"math"
	"testing"
)

func TestRGBAF32(t *testing.T) {
	rgbaf32 := NewRGBAF32(100, 100)

	width := rgbaf32.Bounds().Dx()
	height := rgbaf32.Bounds().Dy()
	diagonalMax := float32(math.Sqrt(float64(width*width + height*height)))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			diagonal := float32(math.Sqrt(float64(x*x + y*y)))

			alpha := diagonal / diagonalMax
			c := floatcolor.RGBAF32{
				R: alpha * float32(x) / float32(width),
				G: alpha * float32(y) / float32(height),
				B: alpha * diagonal / diagonalMax,
				A: alpha,
			}

			rgbaf32.Set(x+rgbaf32.Bounds().Min.X, y+rgbaf32.Bounds().Min.Y, c)
		}
	}

	writeImage("../testresult/RGBAF32.png", rgbaf32)
}
