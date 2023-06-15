package floatcolor

import "image/color"

type Float64Color interface {
	SetPrecise(bool)

	// Mix smoothly mixes the RGB values of two color into one resulting color.
	// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
	// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
	// and a mix of 1.0 gives same color as c2.
	// Alpha value is affected.
	Mix(color.Color, float64) color.Color

	// SetAlpha sets the alpha value
	// and updates the other color components if necessary (when color is premultiplied alpha)
	SetAlpha(float64)

	// SetRGB sets the (non premultiplied alpha) red, green, and blue value
	// and updates the components if necessary (when color is premultiplied alpha type)
	SetRGB(float64, float64, float64)

	// SetR sets the (non premultiplied alpha) red value
	// and updates the red component if necessary (when color is premultiplied alpha type)
	SetR(float64)

	// SetG sets the  (non premultiplied alpha)green value
	// and updates the green component if necessary (when color is premultiplied alpha type)
	SetG(float64)

	// SetB sets the (non premultiplied alpha) blue value
	// and updates the blue component if necessary (when color is premultiplied alpha type)
	SetB(float64)

	color.Color
}

type Float32Color interface {
	// Mix smoothly mixes the RGB values of two color into one resulting color.
	// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
	// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
	// and a mix of 1.0 gives same color as c2.
	// Alpha value is affected.
	Mix(color.Color, float64) color.Color

	// SetAlpha sets the alpha value
	// and updates the other color components if necessary (when color is premultiplied alpha)
	SetAlpha(float32)

	// SetRGB sets the (non premultiplied alpha) red, green, and blue value
	// and updates the components if necessary (when color is premultiplied alpha type)
	SetRGB(float32, float32, float32)

	// SetR sets the (non premultiplied alpha) red value
	// and updates the red component if necessary (when color is premultiplied alpha type)
	SetR(float32)

	// SetG sets the  (non premultiplied alpha)green value
	// and updates the green component if necessary (when color is premultiplied alpha type)
	SetG(float32)

	// SetB sets the (non premultiplied alpha) blue value
	// and updates the blue component if necessary (when color is premultiplied alpha type)
	SetB(float32)

	color.Color
}

// Mix smoothly mixes the RGB values of two color into one resulting color.
// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
// and a mix of 1.0 gives same color as c2.
// Alpha value is affected.
func mix(c1 color.Color, c2 color.Color, mix float64) color.Color {
	cc1 := NRGBAF64Model.Convert(c1).(NRGBAF64)
	cc2 := NRGBAF64Model.Convert(c2).(NRGBAF64)

	R := cc1.R*(1.0-mix) + cc2.R*mix
	G := cc1.G*(1.0-mix) + cc2.G*mix
	B := cc1.B*(1.0-mix) + cc2.B*mix
	A := cc1.A*(1.0-mix) + cc2.A*mix

	precise := cc1.Precise || cc2.Precise

	return NRGBAF64{R: R, G: G, B: B, A: A, Precise: precise}
}
