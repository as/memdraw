package memdraw

import (
	"image"
	"image/color"
)

var rainbow = color.RGBA{255, 0, 0, 255}

func next() *image.Uniform {
	rainbow = nextcolor(rainbow)
	return image.NewUniform(rainbow)
}

// nextcolor steps through a gradient
func nextcolor(c color.RGBA) color.RGBA {
	switch {
	case c.R == 255 && c.G == 0 && c.B == 0:
		c.G += 17
	case c.R == 255 && c.G != 255 && c.B == 0:
		c.G += 17
	case c.G == 255 && c.R != 0:
		c.R -= 17
	case c.R == 0 && c.B != 255:
		c.B += 17
	case c.B == 255 && c.G != 0:
		c.G -= 17
	case c.G == 0 && c.R != 255:
		c.R += 17
	default:
		c.B -= 17
	}
	return c
}
