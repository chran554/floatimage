package floatcolor

import "image/color"

type RGBAF32 struct {
	R, G, B, A float32
	Precise    bool
}

var (
	RGBAF32Model = color.ModelFunc(rgbaf32Model)
)

func (rgbaf32 RGBAF32) RGBA() (r, g, b, a uint32) {
	return uint32(clampF32(rgbaf32.R*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.G*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.B*0xffff, 0.0, 0xffff, rgbaf32.Precise)),
		uint32(clampF32(rgbaf32.A*0xffff, 0.0, 0xffff, rgbaf32.Precise))
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
