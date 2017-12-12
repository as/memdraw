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

func Line(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	if dx*dx > dy*dy {
		hLine(dst, p1, p0, thick, src, sp)
		return
	}
	vLine(dst, p1, p0, thick, src, sp)
}

// Line draws a line from q0 to q1 on dst
func vLine(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	if p0.Y > p1.Y {
		vLine(dst, p1, p0, thick, src, sp)
		return
	}
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	miny := max(p0.Y, dst.Bounds().Min.Y)
	maxy := min(p1.Y, dst.Bounds().Max.Y-1)
	x := 0
	for y := miny; y <= maxy; y++ {
		if dy == 0 {
			x = p0.X
		} else {
			x = p0.X + dx*(y-p0.Y)/dy
		}
		col := src.At(x, y)
		dstcol := dst.At(x, y)
		cr0, cg0, cb0, ca0 := col.RGBA()
		cr1, cg1, cb1, ca1 := dstcol.RGBA()
		cr := cr0 + cr1*(255-ca0)
		cg := cg0 + cg1*(255-ca0)
		cb := cb0 + cb1*(255-ca0)
		ca := ca0 + ca1*(255-ca0)
		dst.Set(x, y, color.RGBA{byte(cr), byte(cg), byte(cb), byte(ca)})
	}
}

// Line draws a line from q0 to q1 on dst
func hLine(dst draw.Image, p0, p1 image.Point, thick int, src image.Image, sp image.Point) {
	// Note: There exists a 4-step and 8-step optimization that takes advantage of repeating patterns in the
	// line. Is it worth it?
	if p1.X < p0.X {
		hLine(dst, p1, p0, thick, src, sp)
		return
	}
	dd := 1
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	if dy < 0 {
		dd = -1
		dy = -dy
	}
	e := 2*dy - dx
	dy *= 2
	dx = dy - 2*dx
	y := p0.Y
	maxx := min(p1.X, dst.Bounds().Max.X-1)
	for x := p0.X; x <= maxx; x++ {
		col := src.At(x, y)
		dstcol := dst.At(x, y)
		cr0, cg0, cb0, ca0 := col.RGBA()
		cr1, cg1, cb1, ca1 := dstcol.RGBA()
		cr := cr0 + cr1*(255-ca0)
		cg := cg0 + cg1*(255-ca0)
		cb := cb0 + cb1*(255-ca0)
		ca := ca0 + ca1*(255-ca0)
		dst.Set(x, y, color.RGBA{byte(cr), byte(cg), byte(cb), byte(ca)})
		if e > 0 {
			y += dd
			e += dx
		} else {
			e += dy
		}
	}
}
