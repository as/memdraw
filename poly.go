package memdraw

import (
	"image"
	"image/draw"
)

// Poly draws a polygon
func Poly(dst draw.Image, p []image.Point, end0, end1, thick int, src image.Image, sp image.Point) {
	if len(p) < 2 {
		return
	}
	for i := 1; i < len(p); i++ {
		Line(dst, p[i-1], p[i], thick, src, sp)
	}
}
