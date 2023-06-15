package floatcolor

import "image/color"

type NRGBAF64 struct {
	R, G, B, A float64
	Precise    bool
}

var (
	NRGBAF64Model = color.ModelFunc(nrgbaf64Model)
)

// NewNRGBAF64 creates a new NRGBAF64 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be using ordinary (non premultiplied) alpha.
func NewNRGBAF64(red, green, blue, alpha float64) NRGBAF64 {
	return NRGBAF64{R: red, G: green, B: blue, A: alpha, Precise: false}
}

// NewNRGBAF64PremultipliedAlpha creates a new NRGBAF64 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be already premultiplied with alpha.
func NewNRGBAF64PremultipliedAlpha(red, green, blue, alpha float64) NRGBAF64 {
	rgbaf64 := NewRGBAF64(red, green, blue, alpha)
	return NRGBAF64Model.Convert(rgbaf64).(NRGBAF64)
}

func (nrgbaf64 NRGBAF64) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf64.A * 0xffff
	return uint32(clampF64(nrgbaf64.R*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.G*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.B*conv, 0x0000, 0xffff, nrgbaf64.Precise)),
		uint32(clampF64(nrgbaf64.A*0xffff, 0x0000, 0xffff, nrgbaf64.Precise))
}

func (nrgbaf64 *NRGBAF64) SetPrecise(usePreciseCalculation bool) {
	nrgbaf64.Precise = usePreciseCalculation
}

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func (nrgbaf64 *NRGBAF64) Mix(c1 color.Color, mixAmount float64) color.Color {
	cc1 := NRGBAF64Model.Convert(c1).(NRGBAF64)
	return mix(nrgbaf64, cc1, mixAmount)
}

func (nrgbaf64 *NRGBAF64) SetAlpha(alpha float64) {
	nrgbaf64.A = alpha
}

func (nrgbaf64 *NRGBAF64) SetRGB(red float64, green float64, blue float64) {
	nrgbaf64.SetR(red)
	nrgbaf64.SetG(green)
	nrgbaf64.SetB(blue)
}

func (nrgbaf64 *NRGBAF64) SetR(red float64) {
	nrgbaf64.R = red
}

func (nrgbaf64 *NRGBAF64) SetG(green float64) {
	nrgbaf64.G = green
}

func (nrgbaf64 *NRGBAF64) SetB(blue float64) {
	nrgbaf64.B = blue
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
