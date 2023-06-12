package floatcolor

import "image/color"

type NRGBAF64 struct {
	R, G, B, A float64
}

type NRGBAF32 struct {
	R, G, B, A float32
}

type RGBAF64 struct {
	R, G, B, A float64
}

type RGBAF32 struct {
	R, G, B, A float32
}

var (
	NRGBAF64Model = color.ModelFunc(nrgbaf64Model)
	NRGBAF32Model = color.ModelFunc(nrgbaf32Model)
	RGBAF64Model  = color.ModelFunc(rgbaf64Model)
	RGBAF32Model  = color.ModelFunc(rgbaf32Model)
)

func (nrgbaf64 NRGBAF64) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf64.A * 0xffff
	return uint32(clampF64(nrgbaf64.R*conv, 0.0, 0xffff)),
		uint32(clampF64(nrgbaf64.G*conv, 0.0, 0xffff)),
		uint32(clampF64(nrgbaf64.B*conv, 0.0, 0xffff)),
		uint32(clampF64(nrgbaf64.A*0xffff, 0.0, 0xffff))
}

func (nrgbaf32 NRGBAF32) RGBA() (r, g, b, a uint32) {
	conv := nrgbaf32.A * 0xffff
	return uint32(clampF32(nrgbaf32.R*conv, 0.0, 0xffff)),
		uint32(clampF32(nrgbaf32.G*conv, 0.0, 0xffff)),
		uint32(clampF32(nrgbaf32.B*conv, 0.0, 0xffff)),
		uint32(clampF32(nrgbaf32.A*0xffff, 0.0, 0xffff))
}

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

func clampF32(v float32, min float32, max float32) float32 {
	if v > max {
		return max
	} else if v < min {
		return min
	} else {
		return v
	}
}

func clampF64(v float64, min float64, max float64) float64 {
	if v > max {
		return max
	} else if v < min {
		return min
	} else {
		return v
	}
}

func nrgbaf64Model(c color.Color) color.Color {
	if _, ok := c.(NRGBAF64); ok {
		return c
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

func nrgbaf32Model(c color.Color) color.Color {
	if _, ok := c.(NRGBAF32); ok {
		return c
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
