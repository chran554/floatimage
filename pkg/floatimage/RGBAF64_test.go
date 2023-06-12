package floatimage

import (
	"floatimage/pkg/floatcolor"
	"image"
	"math"
	"testing"
)

func TestRGBAF64(t *testing.T) {
	rgbaf64 := NewRGBAF64(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}})

	width := rgbaf64.Bounds().Dx()
	height := rgbaf64.Bounds().Dy()
	diagonalMax := math.Sqrt(float64(width*width + height*height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			diagonal := math.Sqrt(float64(x*x + y*y))

			alpha := diagonal / diagonalMax
			c := floatcolor.RGBAF64{
				R: alpha * float64(x) / float64(width),
				G: alpha * float64(y) / float64(height),
				B: alpha * diagonal / diagonalMax,
				A: alpha,
			}

			rgbaf64.Set(x+rgbaf64.Bounds().Min.X, y+rgbaf64.Bounds().Min.Y, c)
		}
	}

	writeImage("../testresult/RGBAF64.png", rgbaf64)
}
