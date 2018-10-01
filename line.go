package memdraw

import (
	"image"
	"image/color"
	"image/draw"
)

func Line(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	if !lineInRect(dst.Bounds(), p0, p1) {
		return
	}
	dx := p1.X - p0.X
	dy := p0.Y - p1.Y

	if dx*dx > dy*dy {
		if p0.X > p1.X {
			hLine1(dst, p1, p0, thick, src, sp)
		} else {
			hLine1(dst, p0, p1, thick, src, sp)
		}
	} else {
		if p1.Y > p0.Y {
			vLine1(dst, p0, p1, thick, src, sp)
		} else {
			vLine1(dst, p1, p0, thick, src, sp)
		}
	}
}

var lineOverRect = LineInRect

func LineInRect(r image.Rectangle, p0, p1 image.Point) bool {
	return lineInRect(r, p0, p1)
}

func lineInRect(r image.Rectangle, p0, p1 image.Point) bool {
	if r.Empty() {
		return false // can't cross an empty rectangle
	}
	c0 := oc(r, p0)
	c1 := oc(r, p1)
	if c0&c1 != 0 {
		return false
	}
	if (c0|c1) == 3 || (c0|c1) == 12 {
//	if c0 == 0 || c1 == 0 || (c0|c1) == 3 || (c0|c1) == 12 {
		return true
	}
	c := [4]image.Point{
		r.Min,
		image.Point{r.Max.X, r.Min.Y},
		image.Point{r.Min.X, r.Max.Y},
		r.Max,
	}
	p := [2]image.Point{p0, p1}
	q := lut[byte(c0<<4 | c1)]
	
	if q[2] >= 0x80 {
		return b3(p[q[0]], p[q[1]], c[(q[2]&^0x80)>>4], c[q[2]&3])
	}
	return b2(p[q[0]], p[q[1]], c[q[2]])
}

func slope(p1, p2 image.Point) (m float64) {
	return float64(p2.Y-p1.Y) / float64(p2.X-p1.X)
}
func b3(v, p, a, b image.Point) bool {
	dp := slope(v, p)
	return !(dp > slope(v, b) || dp < slope(v, a))
}
func b2(v, p, a image.Point) bool {
	return !(slope(v, p) < slope(v, a))
}

func Oc(r image.Rectangle, p image.Point) int {
	c := 0
	if p.X < r.Min.X {
		c = 1
	} else if p.X >= r.Max.X{
		c = 2
	}
	if p.Y < r.Min.Y {
		c |= 4
	} else if p.Y >= r.Max.Y {
		c |= 8
	}
	return c
}
var lut = [256][3]byte{
	0x14: {0, 1, 0},
	0x16: {0, 1, 0},
	0x18: {1, 0, 2},
	0x1a: {1, 0, 2},

	0x24: {1, 0, 1},
	0x25: {1, 0, 1},
	0x28: {0, 1, 3},
	0x29: {0, 1, 3},

	0x49: {1, 0, 0},
	0x4a: {0, 1, 1},

	0x58: {1, 0, 2},
	0x68: {0, 1, 3},
	0x5a: {0, 1, 0x80 + 0x12},
	0x69: {1, 0, 0x80 + 0x03},
}

func init() {
	for i := byte(0); i < 127; i++ {
		if i > 127 {
			break
		}
		v := lut[i]
		lut[i<<4|i>>4] = [3]byte{v[1], v[0], v[2]}
	}
	//	lut[0x86] = [3]byte{1, 0, 3}
	//	lut[0x68] = [3]byte{0, 1, 3}
}

func oc(r image.Rectangle, p image.Point) int {
	return Oc(r, p)
}


// Line draws a line from q0 to q1 on dst
func vLine1(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	clipr := dst.Bounds()
	dd := 1
	if dx < 0 {
		dd = -1
		dx = -dx
	}
	maxy := min(p1.Y, clipr.Max.Y-1)
	e := 2*dx - dy
	x := p0.X
	dx *= 2
	dy = dx - 2*dy
	for y := p0.Y; y <= maxy; y++ {
		if clipr.Min.Y <= y && clipr.Min.X <= x && x < clipr.Max.X {
			col := src.At(x, y)
			dst.Set(x, y, col)
		}
		if e > 0 {
			x += dd
			e += dy
		} else {
			e += dx
		}
	}
}

// Line draws a line from q0 to q1 on dst
func vLine(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	clipr := dst.Bounds()
	miny := max(p0.Y, clipr.Min.Y)
	maxy := min(p1.Y, clipr.Max.Y-1)
	x := 0
	for y := miny; y <= maxy; y++ {
		if dy == 0 {
			x = p0.X
		} else {
			x = p0.X + dx*(y-p0.Y)/dy
		}
		if clipr.Min.X <= x && x < clipr.Max.X {
			col := src.At(x, y)
			dst.Set(x, y, col)
		}
	}
}

// Line draws a line from q0 to q1 on dst
func hLine1(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	// Note: There exists a 4-step and 8-step optimization that takes advantage of repeating patterns in the
	// line. Is it worth it?
	dd := 1
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	if dy < 0 {
		dd = -1
		dy = -dy
	}
	clipr := dst.Bounds()
	maxx := min(p1.X, clipr.Max.X-1)
	e := 2*dy - dx
	y := p0.Y
	dy *= 2
	dx = dy - 2*dx
	for x := p0.X; x <= maxx; x++ {
		if clipr.Min.X <= x && clipr.Min.Y <= y && y < clipr.Max.Y {
			col := src.At(x, y)
			dst.Set(x, y, col)
		}
		if e > 0 {
			y += dd
			e += dx
		} else {
			e += dy
		}
	}
}

// Line draws a line from q0 to q1 on dst
func hLine(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	clipr := dst.Bounds()
	minx := min(p0.X, clipr.Min.X)
	maxx := min(p1.X, clipr.Max.X-1)

	for x := minx; x <= maxx; x++ {
		y := p0.Y + (dy*(x-p0.X)+dy/2)/dx
		if clipr.Min.Y <= y && y < clipr.Max.Y {
			col := src.At(x, y)
			dst.Set(x, y, col)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	red = image.NewUniform(color.RGBA{255, 0, 0, 255})
	gr  = image.NewUniform(color.RGBA{0, 255, 0, 255})
	bl  = image.NewUniform(color.RGBA{0, 0, 255, 255})
	gra = image.NewUniform(color.RGBA{222, 222, 222, 255})
)
