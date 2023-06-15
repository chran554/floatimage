package floatimage

import (
	"floatimage/pkg/floatcolor"
	"image"
	"image/color"
)

// RGBAF32 is an in-memory image whose At method returns floatcolor.RGBAF32 values.
type RGBAF32 struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format.
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []float32
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect    image.Rectangle
	Precise bool
}

// NewRGBAF32 returns a new RGBAF32 image with the given dimensions.
// An RGBAF32 image is an RGB image with premultiplied alpha
// where the red, green, blue, and alpha values are float32 values in the typical range [0.0, 1.0].
func NewRGBAF32(width, height int) *RGBAF32 {
	return NewRGBAF32WithBounds(0, 0, width, height)
}

// NewRGBAF32WithBounds returns a new RGBAF32 image with the given bounds.
// An RGBAF32 image is an RGB image with premultiplied alpha
// where the red, green, blue, and alpha values are float32 values in the typical range [0.0, 1.0].
func NewRGBAF32WithBounds(x0, y0, x1, y1 int) *RGBAF32 {
	r := image.Rect(x0, y0, x1, y1)
	const channels = 4
	return &RGBAF32{
		Pix:     make([]float32, pixelBufferLength(channels, r, "RGBAF32")),
		Stride:  channels * r.Dx(),
		Rect:    r,
		Precise: false,
	}
}

func (p *RGBAF32) ColorModel() color.Model { return floatcolor.RGBAF32Model }

func (p *RGBAF32) Bounds() image.Rectangle { return p.Rect }

func (p *RGBAF32) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	return floatcolor.RGBAF32{R: s[0], G: s[1], B: s[2], A: s[3], Precise: p.Precise}
}

func (p *RGBAF32) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	r := uint16(clampF32(s[0]*0xffff, 0x0000, 0xffff, p.Precise))
	g := uint16(clampF32(s[1]*0xffff, 0x0000, 0xffff, p.Precise))
	b := uint16(clampF32(s[2]*0xffff, 0x0000, 0xffff, p.Precise))
	a := uint16(clampF32(s[3]*0xffff, 0x0000, 0xffff, p.Precise))

	return color.RGBA64{R: r, G: g, B: b, A: a}
}

// PixOffset returns the index of the first element of Pix that corresponds to the pixel at (x, y).
func (p *RGBAF32) PixOffset(x, y int) int {
	const channels = 4
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*channels
}

func (p *RGBAF32) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	c1 := floatcolor.RGBAF32Model.Convert(c).(floatcolor.RGBAF32)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

func (p *RGBAF32) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	c1 := floatcolor.RGBAF32Model.Convert(c).(floatcolor.RGBAF32)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

// SubImage returns an image representing the portion of the image p visible through r.
// The returned value shares pixels with the original image.
func (p *RGBAF32) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty.
	// Without explicitly checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGBAF32{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBAF32{
		Pix:     p.Pix[i:],
		Stride:  p.Stride,
		Rect:    r,
		Precise: p.Precise,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBAF32) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}

	const channels = 4

	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for x := p.Rect.Min.X; x < p.Rect.Max.X; x++ {
			if p.Pix[y*p.Stride+x*channels+(channels-1)] != 1.0 {
				return false
			}
		}
	}

	return true
}
