package floatcolor

import "image/color"

type RGBAF64 struct {
	R, G, B, A float64
	Precise    bool
}

var (
	RGBAF64Model = color.ModelFunc(rgbaf64Model)
)

// NewRGBAF64 creates a new RGBAF64 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be already premultiplied with alpha.
func NewRGBAF64(red, green, blue, alpha float64) RGBAF64 {
	return RGBAF64{R: red, G: green, B: blue, A: alpha, Precise: false}
}

// NewRGBAF64NonPremultipliedAlpha creates a new RGBAF64 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be using ordinary (non premultiplied) alpha.
func NewRGBAF64NonPremultipliedAlpha(red, green, blue, alpha float64) RGBAF64 {
	nrgbaf64 := NewNRGBAF64(red, green, blue, alpha)
	return RGBAF64Model.Convert(nrgbaf64).(RGBAF64)
}

func (rgbaf64 RGBAF64) RGBA() (r, g, b, a uint32) {
	return uint32(clampF64(rgbaf64.R*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.G*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.B*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.A*0xffff, 0.0, 0xffff, rgbaf64.Precise))
}

func (rgbaf64 RGBAF64) AsNRGBA() color.NRGBA {
	conv := 0xff / rgbaf64.A
	return color.NRGBA{
		R: uint8(clampF64(rgbaf64.R*conv, 0x00, 0xff, rgbaf64.Precise)),
		G: uint8(clampF64(rgbaf64.G*conv, 0x00, 0xff, rgbaf64.Precise)),
		B: uint8(clampF64(rgbaf64.B*conv, 0x00, 0xff, rgbaf64.Precise)),
		A: uint8(clampF64(rgbaf64.A*0xff, 0x00, 0xff, rgbaf64.Precise)),
	}
}

func (rgbaf64 RGBAF64) AsRGBA() color.RGBA {
	const conv = float64(0xff)
	return color.RGBA{
		R: uint8(clampF64(rgbaf64.R*conv, 0x00, 0xff, rgbaf64.Precise)),
		G: uint8(clampF64(rgbaf64.G*conv, 0x00, 0xff, rgbaf64.Precise)),
		B: uint8(clampF64(rgbaf64.B*conv, 0x00, 0xff, rgbaf64.Precise)),
		A: uint8(clampF64(rgbaf64.A*conv, 0x00, 0xff, rgbaf64.Precise)),
	}
}

func (rgbaf64 *RGBAF64) SetPrecise(usePreciseCalculation bool) {
	rgbaf64.Precise = usePreciseCalculation
}

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func (rgbaf64 *RGBAF64) Mix(c1 color.Color, mixAmount float64) color.Color {
	cc1 := RGBAF64Model.Convert(c1).(RGBAF64)
	return mix(rgbaf64, cc1, mixAmount)
}

func (rgbaf64 *RGBAF64) SetAlpha(alpha float64) {
	nrgbaf := nrgbaf64Model(rgbaf64).(NRGBAF64)
	nrgbaf.A = alpha
	rgbaf := rgbaf64Model(nrgbaf).(RGBAF64)
	rgbaf64.R, rgbaf64.G, rgbaf64.B, rgbaf64.A = rgbaf.R, rgbaf.G, rgbaf.B, rgbaf.A
}

func (rgbaf64 *RGBAF64) SetRGB(red float64, green float64, blue float64) {
	rgbaf64.SetR(red)
	rgbaf64.SetG(green)
	rgbaf64.SetB(blue)
}

func (rgbaf64 *RGBAF64) SetR(red float64) {
	rgbaf64.R = red * rgbaf64.A
}

func (rgbaf64 *RGBAF64) SetG(green float64) {
	rgbaf64.G = green * rgbaf64.A
}

func (rgbaf64 *RGBAF64) SetB(blue float64) {
	rgbaf64.B = blue * rgbaf64.A
}

func rgbaf64Model(c color.Color) color.Color {
	if _, ok := c.(RGBAF64); ok {
		return c
	}

	if rgbaf, ok := c.(RGBAF32); ok {
		return RGBAF64{R: float64(rgbaf.R), G: float64(rgbaf.G), B: float64(rgbaf.B), A: float64(rgbaf.A)}
	}

	if nrgbaf, ok := c.(NRGBAF64); ok {
		conv := nrgbaf.A
		return RGBAF64{R: nrgbaf.R * conv, G: nrgbaf.G * conv, B: nrgbaf.B * conv, A: nrgbaf.A}
	}

	if nrgbaf, ok := c.(NRGBAF32); ok {
		conv := nrgbaf.A
		return RGBAF64{R: float64(nrgbaf.R * conv), G: float64(nrgbaf.G * conv), B: float64(nrgbaf.B * conv), A: float64(nrgbaf.A)}
	}

	r, g, b, a := c.RGBA()
	if a == 0xffff {
		conv := 1.0 / 0xffff
		return RGBAF64{R: float64(r) * conv, G: float64(g) * conv, B: float64(b) * conv, A: 1.0}
	}
	if a == 0 {
		return RGBAF64{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	}

	conv := 1.0 / 0xffff
	return RGBAF64{R: float64(r) * conv, G: float64(g) * conv, B: float64(b) * conv, A: float64(a) * conv}
}
