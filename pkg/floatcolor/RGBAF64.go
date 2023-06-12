package floatcolor

import "image/color"

type RGBAF64 struct {
	R, G, B, A float64
	Precise    bool
}

var (
	RGBAF64Model = color.ModelFunc(rgbaf64Model)
)

func (rgbaf64 RGBAF64) RGBA() (r, g, b, a uint32) {
	return uint32(clampF64(rgbaf64.R*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.G*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.B*0xffff, 0.0, 0xffff, rgbaf64.Precise)),
		uint32(clampF64(rgbaf64.A*0xffff, 0.0, 0xffff, rgbaf64.Precise))
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
