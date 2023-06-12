package floatimage

import (
	"floatimage/pkg/floatcolor"
	"image"
	"image/color"
	"math/bits"
)

/*

// Image is an image.Image with a Set method to change a single pixel.
type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

// RGBA64Image extends both the Image and image.RGBA64Image interfaces with a
// SetRGBA64 method to change a single pixel. SetRGBA64 is equivalent to
// calling Set, but it can avoid allocations from converting concrete color
// types to the color.Color interface type.
type RGBA64Image interface {
	image.RGBA64Image
	Set(x, y int, c color.Color)
	SetRGBA64(x, y int, c color.RGBA64)
}

*/

// NRGBAF64 is an in-memory image whose At method returns floatcolor.NRGBAF64 values.
type NRGBAF64 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format.
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []float64
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *NRGBAF64) ColorModel() color.Model { return floatcolor.NRGBAF64Model }

func (p *NRGBAF64) Bounds() image.Rectangle { return p.Rect }

func (p *NRGBAF64) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	return floatcolor.NRGBAF64{R: s[0], G: s[1], B: s[2], A: s[3]}
}

func (p *NRGBAF64) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	r := uint16(clampF64(s[0]*0xffff*s[3], 0.0, 0xffff))
	g := uint16(clampF64(s[1]*0xffff*s[3], 0.0, 0xffff))
	b := uint16(clampF64(s[2]*0xffff*s[3], 0.0, 0xffff))
	a := uint16(clampF64(s[3]*0xffff, 0.0, 0xffff))

	return color.RGBA64{R: r, G: g, B: b, A: a}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *NRGBAF64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
}

func (p *NRGBAF64) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	c1 := floatcolor.NRGBAF64Model.Convert(c).(floatcolor.NRGBAF64)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

func (p *NRGBAF64) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	c1 := floatcolor.NRGBAF64Model.Convert(c).(floatcolor.NRGBAF64)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *NRGBAF64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty.
	// Without explicitly checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &NRGBAF64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &NRGBAF64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *NRGBAF64) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}

	const bytesPerBpp = 64 / 8
	const bytesPerPixel = bytesPerBpp * 4

	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for x := p.Rect.Min.X; x < p.Rect.Max.X; x++ {
			if p.Pix[y*p.Stride+x*bytesPerPixel+3] != 1.0 {
				return false
			}
		}
	}

	return true
}

// NewNRGBAF64 returns a new NRGBAF64 image with the given bounds.
func NewNRGBAF64(r image.Rectangle) *NRGBAF64 {
	const bytesPerBpp = 64 / 8
	const bytesPerPixel = bytesPerBpp * 4
	return &NRGBAF64{
		Pix:    make([]float64, pixelBufferLength(bytesPerPixel, r, "NRGBAF64")),
		Stride: bytesPerPixel * r.Dx(),
		Rect:   r,
	}
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
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

func clampF64(v float64, min float64, max float64) float64 {
	if v > max {
		return max
	} else if v < min {
		return min
	} else {
		return v
	}
}
