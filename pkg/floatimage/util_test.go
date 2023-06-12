package floatimage

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

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
