package memdraw

import (
	"image"
	"image/color"
	"image/draw"
)

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

func Line(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	defer func() {
		recover()
	}()
	dx := p1.X - p0.X
	dy := p0.Y - p1.Y
	if dx*dx > dy*dy {
		if p0.X > p1.X {
			hLine1(dst, p1, p0, thick, red, sp)
			//hLine(dst, p0, p1, thick, image.Black, sp)
		} else {
			hLine1(dst, p0, p1, thick, bl, sp)
			//hLine(dst, p1, p0, thick, gr, sp)
		}
	} else {
		//
		//		vLine(dst, p1, p0, thick, red, sp)
		if p1.Y > p0.Y {
			vLine1(dst, p0, p1, thick, gr, sp)
		} else {
			vLine1(dst, p1, p0, thick, bl, sp)
			//vLine(dst, p0, p1, thick, gra, sp)
		}
	}
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
		//dstcol := dst.At(x, y)
		//cr0, cg0, cb0, ca0 := col.RGBA()
		//cr1, cg1, cb1, ca1 := dstcol.RGBA()
		//cr := cr0 + cr1*(255-ca0)
		//cg := cg0 + cg1*(255-ca0)
		//cb := cb0 + cb1*(255-ca0)
		//ca := ca0 + ca1*(255-ca0)
		//dst.Set(x, y, color.RGBA{byte(cr), byte(cg), byte(cb), byte(ca)})
	}
}

// Line draws a line from q0 to q1 on dst
func hLine1(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	// Note: There exists a 4-step and 8-step optimization that takes advantage of repeating patterns in the
	// line. Is it worth it?
	if p0.X < p1.X {
		//hLine(dst, p1, p0, thick, src, sp)
		//return
	}
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
	// Note: There exists a 4-step and 8-step optimization that takes advantage of repeating patterns in the
	// line. Is it worth it?
	if p0.X < p1.X {
		//hLine(dst, p1, p0, thick, src, sp)
		//return
	}
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
