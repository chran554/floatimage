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

	ConvertableColor
	color.Color
}

type Float32Color interface {
	// Mix smoothly mixes the RGB values of two color into one resulting color.
	// Parameter mix determine how much percent of color c2 is in the resulting mixed color.
	// Mix value range is [0.0, 1.0] where the resulting mix of 0.0 gives same color as c1
	// and a mix of 1.0 gives same color as c2.
	// Alpha value is affected.
	Mix(c color.Color, mixAmount float64) color.Color

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

	ConvertableColor
	color.Color
}

type ConvertableColor interface {
	AsNRGBA() color.NRGBA
	AsRGBA() color.RGBA
}
