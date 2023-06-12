package floatcolor

import "image/color"

type NRGBAF32 struct {
	R, G, B, A float32
	Precise    bool
}

var (
	NRGBAF32Model = color.ModelFunc(nrgbaf32Model)
)

func (nrgbaf32 NRGBAF32) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf32.A * 0xffff
	return uint32(clampF32(nrgbaf32.R*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.G*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.B*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.A*0xffff, 0.0, 0xffff, nrgbaf32.Precise))
}

// Add adds the RGB values to the color making it brighter.
// No clamping of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf32 *NRGBAF32) Add(c color.Color) {
	c2 := NRGBAF32Model.Convert(c).(NRGBAF32)

	nrgbaf32.R += c2.R
	nrgbaf32.G += c2.G
	nrgbaf32.B += c2.B
}

// Sub subtracts the RGB values to the color making it darker.
// No clamping of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf32 *NRGBAF32) Sub(c color.Color) {
	c2 := NRGBAF32Model.Convert(c).(NRGBAF32)

	nrgbaf32.R -= c2.R
	nrgbaf32.G -= c2.G
	nrgbaf32.B -= c2.B
}

// Mul multiplies the RGB values to the color.
// No clamping of the result of the RGB components is made.
// Result for each component may be outside the [0.0, 1.0] range.
// Alpha value is not affected.
func (nrgbaf32 *NRGBAF32) Mul(c color.Color) {
	c2 := NRGBAF32Model.Convert(c).(NRGBAF32)

	nrgbaf32.R *= c2.R
	nrgbaf32.G *= c2.G
	nrgbaf32.B *= c2.B
}

func nrgbaf32Model(c color.Color) color.Color {
	if _, ok := c.(NRGBAF32); ok {
		return c
	}

	if nrgbaf, ok := c.(NRGBAF64); ok {
		return NRGBAF32{R: float32(nrgbaf.R), G: float32(nrgbaf.G), B: float32(nrgbaf.B), A: float32(nrgbaf.A)}
	}

	if rgbaf, ok := c.(RGBAF64); ok {
		conv := 1.0 / rgbaf.A
		return NRGBAF32{R: float32(rgbaf.R * conv), G: float32(rgbaf.G * conv), B: float32(rgbaf.B * conv), A: float32(rgbaf.A)}
	}

	if rgbaf, ok := c.(RGBAF32); ok {
		conv := 1.0 / rgbaf.A
		return NRGBAF32{R: rgbaf.R * conv, G: rgbaf.G * conv, B: rgbaf.B * conv, A: rgbaf.A}
	}

	if nrgba, ok := c.(color.NRGBA); ok {
		conv := float32(1.0 / 0xff)
		return NRGBAF32{R: float32(nrgba.R) * conv, G: float32(nrgba.G) * conv, B: float32(nrgba.B) * conv, A: float32(nrgba.A) * conv}
	}

	r, g, b, a := c.RGBA()
	if a == 0xffff {
		conv := float32(1.0 / 0xffff)
		return NRGBAF32{R: float32(r) * conv, G: float32(g) * conv, B: float32(b) * conv, A: 1.0}
	}
	if a == 0 {
		return NRGBAF32{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	}

	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	conv := float32(a) / (0xffff * 0xffff) // (1.0 / 0xffff) * (float64(a) / 0xffff)
	return NRGBAF32{R: float32(r) * conv, G: float32(g) * conv, B: float32(b) * conv, A: float32(a) / 0xffff}
}
