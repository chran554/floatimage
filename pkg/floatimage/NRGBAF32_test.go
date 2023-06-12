package floatimage

import (
	"floatimage/pkg/floatcolor"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestNRGBAF64(t *testing.T) {
	nrgbaf64 := NewNRGBAF64(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}})

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

	writeImage("test.png", nrgbaf64)
}

func writeImage(filename string, image image.Image) {
	parentPath := filepath.Dir(filename)
	os.MkdirAll(parentPath, os.ModePerm)

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Oups, no files for you today.")
		os.Exit(1)
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level then save to file
	err = png.Encode(f, image)
	if err != nil {
		fmt.Println("Oups, no image encode for you today.")
		os.Exit(1)
	}
}
