package memdraw

import (
	"image"
	"image/color"
)

var rainbow = color.RGBA{255, 0, 0, 255}
var uniform = image.NewUniform(rainbow)

func next() *image.Uniform {
	rainbow = nextcolor(rainbow)
	uniform = image.NewUniform(rainbow)
	return uniform
}

// nextcolor steps through a gradient
func nextcolor(c color.RGBA) color.RGBA {
	switch {
	case c.R == 255 && c.G == 0 && c.B == 0:
		c.G += 5
	case c.R == 255 && c.G != 255 && c.B == 0:
		c.G += 5
	case c.G == 255 && c.R != 0:
		c.R -= 5
	case c.R == 0 && c.B != 255:
		c.B += 5
	case c.B == 255 && c.G != 0:
		c.G -= 5
	case c.G == 0 && c.R != 255:
		c.R += 5
	default:
		c.B -= 5
	}
	return c
}
