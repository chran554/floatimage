package floatcolor

import "image/color"

type RGBAF64 struct {
	R, G, B, A float64
}

type RGBAF32 struct {
	R, G, B, A float32
}

var (
	RGBAF64Model  = color.ModelFunc(rgbaf64Model)
	RGBAF32Model  = color.ModelFunc(rgbaf32Model)
)

func (rgbaf64 RGBAF64) RGBA() (r, g, b, a uint32) {
	return uint32(clampF64(rgbaf64.R*0xffff, 0.0, 0xffff)),
		uint32(clampF64(rgbaf64.G*0xffff, 0.0, 0xffff)),
		uint32(clampF64(rgbaf64.B*0xffff, 0.0, 0xffff)),
		uint32(clampF64(rgbaf64.A*0xffff, 0.0, 0xffff))
}

func (rgbaf32 RGBAF32) RGBA() (r, g, b, a uint32) {
	return uint32(clampF32(rgbaf32.R*0xffff, 0.0, 0xffff)),
		uint32(clampF32(rgbaf32.G*0xffff, 0.0, 0xffff)),
		uint32(clampF32(rgbaf32.B*0xffff, 0.0, 0xffff)),
		uint32(clampF32(rgbaf32.A*0xffff, 0.0, 0xffff))
}

func rgbaf64Model(c color.Color) color.Color {
	if _, ok := c.(RGBAF64); ok {
		return c
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

func rgbaf32Model(c color.Color) color.Color {
	if _, ok := c.(RGBAF32); ok {
		return c
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
