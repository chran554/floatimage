package floatcolor

import "image/color"

type NRGBAF32 struct {
	R, G, B, A float32
	Precise    bool
}

var (
	NRGBAF32Model = color.ModelFunc(nrgbaf32Model)
)

// NewNRGBAF32 creates a new NRGBAF32 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be using ordinary (non premultiplied) alpha.
func NewNRGBAF32(red, green, blue, alpha float32) NRGBAF32 {
	return NRGBAF32{R: red, G: green, B: blue, A: alpha, Precise: false}
}

// NewNRGBAF32PremultipliedAlpha creates a new NRGBAF32 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be already premultiplied with alpha.
func NewNRGBAF32PremultipliedAlpha(red, green, blue, alpha float32) NRGBAF32 {
	rgbaf32 := NewRGBAF32(red, green, blue, alpha)
	return NRGBAF32Model.Convert(rgbaf32).(NRGBAF32)
}

func (nrgbaf32 NRGBAF32) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf32.A * 0xffff
	return uint32(clampF32(nrgbaf32.R*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.G*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.B*conv, 0.0, 0xffff, nrgbaf32.Precise)),
		uint32(clampF32(nrgbaf32.A*0xffff, 0.0, 0xffff, nrgbaf32.Precise))
}

func (nrgbaf32 NRGBAF32) AsNRGBA() color.NRGBA {
	const conv = float32(0xff)
	return color.NRGBA{
		R: uint8(clampF32(nrgbaf32.R*conv, 0x00, 0xff, nrgbaf32.Precise)),
		G: uint8(clampF32(nrgbaf32.G*conv, 0x00, 0xff, nrgbaf32.Precise)),
		B: uint8(clampF32(nrgbaf32.B*conv, 0x00, 0xff, nrgbaf32.Precise)),
		A: uint8(clampF32(nrgbaf32.A*conv, 0x00, 0xff, nrgbaf32.Precise)),
	}
}

func (nrgbaf32 NRGBAF32) AsRGBA() color.RGBA {
	conv := nrgbaf32.A * 0xff
	return color.RGBA{
		R: uint8(clampF32(nrgbaf32.R*conv, 0x00, 0xff, nrgbaf32.Precise)),
		G: uint8(clampF32(nrgbaf32.G*conv, 0x00, 0xff, nrgbaf32.Precise)),
		B: uint8(clampF32(nrgbaf32.B*conv, 0x00, 0xff, nrgbaf32.Precise)),
		A: uint8(clampF32(nrgbaf32.A*0xff, 0x00, 0xff, nrgbaf32.Precise)),
	}
}

func (nrgbaf32 *NRGBAF32) SetPrecise(usePreciseCalculation bool) {
	nrgbaf32.Precise = usePreciseCalculation
}

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func (nrgbaf32 *NRGBAF32) Mix(c1 color.Color, mixAmount float64) color.Color {
	cc1 := NRGBAF64Model.Convert(c1).(NRGBAF64)
	return mix(nrgbaf32, cc1, mixAmount)
}

func (nrgbaf32 *NRGBAF32) SetAlpha(alpha float32) {
	nrgbaf32.A = alpha
}

func (nrgbaf32 *NRGBAF32) SetRGB(red float32, green float32, blue float32) {
	nrgbaf32.SetR(red)
	nrgbaf32.SetG(green)
	nrgbaf32.SetB(blue)
}

func (nrgbaf32 *NRGBAF32) SetR(red float32) {
	nrgbaf32.R = red
}

func (nrgbaf32 *NRGBAF32) SetG(green float32) {
	nrgbaf32.G = green
}

func (nrgbaf32 *NRGBAF32) SetB(blue float32) {
	nrgbaf32.B = blue
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
