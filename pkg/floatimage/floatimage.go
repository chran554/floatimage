package floatimage

import "image"

type FloatImage interface {
	AsRGBA() *image.RGBA
	AsNRGBA() *image.NRGBA

	AsRGBAForRange(min, max float64) *image.RGBA
	AsNRGBAForRange(min, max float64) *image.NRGBA

	image.Image
	image.RGBA64Image
}
