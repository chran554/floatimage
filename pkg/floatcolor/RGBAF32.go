package floatcolor

import "image/color"

type RGBAF32 struct {
	R, G, B, A float32
	Precise    bool
}

var (
	RGBAF32Model = color.ModelFunc(rgbaf32Model)
)

// NewRGBAF32 creates a new RGBAF32 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be already premultiplied with alpha.
func NewRGBAF32(red, green, blue, alpha float32) RGBAF32 {
	return RGBAF32{R: red, G: green, B: blue, A: alpha, Precise: false}
}

// NewRGBAF32NonPremultipliedAlpha creates a new RGBAF32 color.
// It is not set to "precise".
// Red, green, and blue are supposed to be using ordinary (non premultiplied) alpha.
func NewRGBAF32NonPremultipliedAlpha(red, green, blue, alpha float32) RGBAF32 {
	nrgbaf32 := NewNRGBAF32(red, green, blue, alpha)
	return RGBAF32Model.Convert(nrgbaf32).(RGBAF32)
}

func (rgbaf32 RGBAF32) RGBA() (r, g, b, a uint32) {
	return uint32(clampF32(rgbaf32.R*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.G*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.B*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.A*0xffff, 0.0, 0xffff, rgbaf32.Precise))
}

func (rgbaf32 RGBAF32) AsNRGBA() color.NRGBA {
	conv := 0xff / rgbaf32.A
	return color.NRGBA{
		R: uint8(clampF32(rgbaf32.R*conv, 0x00, 0xff, rgbaf32.Precise)),
		G: uint8(clampF32(rgbaf32.G*conv, 0x00, 0xff, rgbaf32.Precise)),
		B: uint8(clampF32(rgbaf32.B*conv, 0x00, 0xff, rgbaf32.Precise)),
		A: uint8(clampF32(rgbaf32.A*0xff, 0x00, 0xff, rgbaf32.Precise)),
	}
}

func (rgbaf32 RGBAF32) AsRGBA() color.RGBA {
	const conv = float32(0xff)
	return color.RGBA{
		R: uint8(clampF32(rgbaf32.R*conv, 0x00, 0xff, rgbaf32.Precise)),
		G: uint8(clampF32(rgbaf32.G*conv, 0x00, 0xff, rgbaf32.Precise)),
		B: uint8(clampF32(rgbaf32.B*conv, 0x00, 0xff, rgbaf32.Precise)),
		A: uint8(clampF32(rgbaf32.A*conv, 0x00, 0xff, rgbaf32.Precise)),
	}
}

func (rgbaf32 *RGBAF32) SetPrecise(usePreciseCalculation bool) {
	rgbaf32.Precise = usePreciseCalculation
}

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func (rgbaf32 *RGBAF32) Mix(c1 color.Color, mixAmount float64) color.Color {
	cc1 := RGBAF64Model.Convert(c1).(RGBAF64)
	return mix(rgbaf32, cc1, mixAmount)
}

func (rgbaf32 *RGBAF32) SetAlpha(alpha float32) {
	nrgbaf := nrgbaf32Model(rgbaf32).(NRGBAF32)
	nrgbaf.A = alpha
	rgbaf := rgbaf32Model(nrgbaf).(RGBAF32)
	rgbaf32.R, rgbaf32.G, rgbaf32.B, rgbaf32.A = rgbaf.R, rgbaf.G, rgbaf.B, rgbaf.A
}

func (rgbaf32 *RGBAF32) SetRGB(red float32, green float32, blue float32) {
	rgbaf32.SetR(red)
	rgbaf32.SetG(green)
	rgbaf32.SetB(blue)
}

func (rgbaf32 *RGBAF32) SetR(red float32) {
	rgbaf32.R = red * rgbaf32.A
}

func (rgbaf32 *RGBAF32) SetG(green float32) {
	rgbaf32.G = green * rgbaf32.A
}

func (rgbaf32 *RGBAF32) SetB(blue float32) {
	rgbaf32.B = blue * rgbaf32.A
}

func rgbaf32Model(c color.Color) color.Color {
	if _, ok := c.(RGBAF32); ok {
		return c
	}

	if rgbaf, ok := c.(RGBAF64); ok {
		return RGBAF32{R: float32(rgbaf.R), G: float32(rgbaf.G), B: float32(rgbaf.B), A: float32(rgbaf.A)}
	}

	if nrgbaf, ok := c.(NRGBAF64); ok {
		conv := nrgbaf.A
		return RGBAF32{R: float32(nrgbaf.R * conv), G: float32(nrgbaf.G * conv), B: float32(nrgbaf.B * conv), A: float32(nrgbaf.A)}
	}

	if nrgbaf, ok := c.(NRGBAF32); ok {
		conv := nrgbaf.A
		return RGBAF32{R: nrgbaf.R * conv, G: nrgbaf.G * conv, B: nrgbaf.B * conv, A: nrgbaf.A}
	}

	r, g, b, a := c.RGBA()
	if a == 0xffff {
		conv := float32(1.0 / 0xffff)
		return RGBAF32{R: float32(r) * conv, G: float32(g) * conv, B: float32(b) * conv, A: 1.0}
	}
	if a == 0 {
		return RGBAF32{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	}

	conv := float32(1.0 / 0xffff)
	return RGBAF32{R: float32(r) * conv, G: float32(g) * conv, B: float32(b) * conv, A: float32(a) * conv}
}
