package floatcolor

import "image/color"

type NRGBAF64 struct {
	R, G, B, A float64
	Precise    bool
}

var (
	NRGBAF64Model = color.ModelFunc(nrgbaf64Model)
)

func (nrgbaf64 NRGBAF64) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf64.A * 0xffff
	return uint32(clampF64(nrgbaf64.R*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.G*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.B*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.A*0xffff, 0x0000, 0xffff, nrgbaf64.Precise))
}

// Add adds the RGB values to the color making it brighter.
// No clamping of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf64 *NRGBAF64) Add(c color.Color) {
	c2 := NRGBAF64Model.Convert(c).(NRGBAF64)

	nrgbaf64.R += c2.R
	nrgbaf64.G += c2.G
	nrgbaf64.B += c2.B
}

// Sub subtracts the RGB values to the color making it darker.
// No clamping of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf64 *NRGBAF64) Sub(c color.Color) {
	c2 := NRGBAF64Model.Convert(c).(NRGBAF64)

	nrgbaf64.R -= c2.R
	nrgbaf64.G -= c2.G
	nrgbaf64.B -= c2.B
}

// Mul multiplies the RGB values to the color.
// No clamping of the result of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf64 *NRGBAF64) Mul(c color.Color) {
	c2 := NRGBAF64Model.Convert(c).(NRGBAF64)

	nrgbaf64.R *= c2.R
	nrgbaf64.G *= c2.G
	nrgbaf64.B *= c2.B
}

func nrgbaf64Model(c color.Color) color.Color {
	if _, ok := c.(NRGBAF64); ok {
		return c
	}

	if nrgbaf, ok := c.(NRGBAF32); ok {
		return NRGBAF64{R: float64(nrgbaf.R), G: float64(nrgbaf.G), B: float64(nrgbaf.B), A: float64(nrgbaf.A)}
	}

	if rgbaf, ok := c.(RGBAF64); ok {
		conv := 1.0 / rgbaf.A
		return NRGBAF64{R: rgbaf.R * conv, G: rgbaf.G * conv, B: rgbaf.B * conv, A: rgbaf.A}
	}

	if rgbaf, ok := c.(RGBAF32); ok {
		conv := 1.0 / rgbaf.A
		return NRGBAF64{R: float64(rgbaf.R * conv), G: float64(rgbaf.G * conv), B: float64(rgbaf.B * conv), A: float64(rgbaf.A)}
	}

	if nrgba, ok := c.(color.NRGBA); ok {
		conv := 1.0 / 0xff
		return NRGBAF64{R: float64(nrgba.R) * conv, G: float64(nrgba.G) * conv, B: float64(nrgba.B) * conv, A: float64(nrgba.A) * conv}
	}

	if rgba64, ok := c.(color.RGBA64); ok {
		alpha := float64(rgba64.A) / 0xffff
		conv := (1.0 / 0xffff) / alpha // Premultiply alpha channel
		return NRGBAF64{R: float64(rgba64.R) * conv, G: float64(rgba64.G) * conv, B: float64(rgba64.B) * conv, A: alpha}
	}

	r, g, b, a := c.RGBA()
	if a == 0xffff {
		conv := 1.0 / 0xffff
		return NRGBAF64{R: float64(r) * conv, G: float64(g) * conv, B: float64(b) * conv, A: 1.0}
	}
	if a == 0 {
		return NRGBAF64{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	}

	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	conv := float64(a) / (0xffff * 0xffff) // (1.0 / 0xffff) * (float64(a) / 0xffff)
	return NRGBAF64{R: float64(r) * conv, G: float64(g) * conv, B: float64(b) * conv, A: float64(a) / 0xffff}
}
