package memdraw

import (
	"image"
	"image/draw"

//	"github.com/as/frame/font"
)

// Border draws an outline of a rectangle on dst
func Border(dst draw.Image, r image.Rectangle, thick int, sp image.Point, src image.Image) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}

// Line draws a line from q0 to q1 on dst
func Line(dst draw.Image, q0, q1 image.Point, thick int, src image.Image, sp image.Point) {
	if q1.X < q0.X {
		q0, q1 = q1, q0
	}
	p0 := q0.Sub(q0)
	p1 := q1.Sub(q0)
	dy := p1.Y
	dx := p1.X
	s := 1
	if dy < 0 {
		dy = -dy
		s = -1
	}
	if thick == 0 {
		thick++
	}
	if thick > 0 {
		thick--
	}
	Ellipse(dst, q0, thick/2, thick/2, 1, src, q0, 0, 0)
	if dx > dy {
		d := 2*dy - dx
		for p0.X < dx {
			r := image.Rect(p0.X+q0.X, s*p0.Y+q0.Y-thick/2, p0.X+q0.X+1, s*p0.Y+q0.Y+1+thick/2)
			draw.Draw(dst, r, src, r.Min, draw.Src)
			if d > 0 {
				p0.Y++
				d -= dx
			}
			p0.X++
			d += dy
		}

	} else {
		d := 2*dx - dy
		for p0.Y < dy {
			r := image.Rect(p0.X+q0.X-thick/2, s*p0.Y+q0.Y, p0.X+q0.X+1+thick/2, s*p0.Y+q0.Y+1)
			draw.Draw(dst, r, src, r.Min, draw.Src)
			if d > 0 {
				p0.X++
				d -= dy
			}
			p0.Y++
			d += dx
		}
	}
	q1 = image.Pt(p0.X+q0.X, s*p0.Y+q0.Y)
	Ellipse(dst, q1, thick/2, thick/2, 1, src, q1, 0, 0)
}

// Poly draws a polygon
func Poly(dst draw.Image, p []image.Point, end0, end1, thick int, src image.Image, sp image.Point) {
	if len(p) < 2 {
		return
	}
	for i := 1; i < len(p); i++ {
		Line(dst, p[i-1], p[i], thick, src, sp)
	}
}

// Bezier draws the cubic Bezier curve defined by Points
// a, b, c, and d. The end styles are determined by end0
// and end1; the thickness of the curve is 1+2*thick.
// The source is aligned so sp in src corresponds to a in dst.
func Bezier(dst draw.Image, a, b, c, d image.Point, end0, end1, thick int, src image.Image, sp image.Point) {

	// start with a slow implementation and 
	// optimize it later when time permits
	for t := float64(0); t < 1.0; t+= 0.5{
		curve(dst, []image.Point{a,b,c,d}, t, thick, src, sp)
	}
}

func flatcurve(dst draw.Image, p []image.Point, thick int, src image.Image, sp image.Point){
	
}

func curve(dst draw.Image, p []image.Point, t float64, thick int, src image.Image, sp image.Point){
	if len(p) == 1{
		r := image.Rect(-1,-1,1,1).Inset(-thick+1).Add(p[0])
		draw.Draw(dst, r, src, sp, draw.Src)
		return
	}
	p2 := make([]image.Point, 0, len(p)-1)
	for i := 0; i < len(p)-1; i++{
		x := int((1-t) * float64(p[i].X) + t * float64(p[i+1].X))
		y := int((1-t) * float64(p[i].Y) + t * float64(p[i+1].Y))
		p2 = append(p2, image.Pt(x,y))
	}
	curve(dst, p2, t, thick, src, sp)
}

// Bezspline takes the same arguments as poly but draws a
// quadratic B-spline (despite its name) rather than a
// polygon. If the first and last points in p are equal,
// the spline has periodic end conditions.
func Bezspline(dst draw.Image, p []image.Point, end0, end1, thick int, src image.Image, sp image.Point) {

}

// Ellipse draws a filled ellipse at center point c
//
// The method uses an efficient integer-based rasterization
// technique originally described in:
//
// McIlroy, M.D.: There is no royal road to programs: a trilogy
// on raster ellipses and programming methodology,
// Computer Science TR155, AT&T Bell Laboratories, 1990
//
func Ellipse(dst draw.Image, c image.Point, a, b, thick int, src image.Image, sp image.Point, alpha, phi int) {
	xc, yc := c.X, c.Y
	var (
		x, y       = 0, b
		a2, b2     = a * a, b * b
		crit1      = -(a2/4 + a%2 + b2)
		crit2      = -(b2/4 + b%2 + a2)
		crit3      = -(b2/4 + b%2)
		t          = -a2 * y
		dxt, dyt   = 2 * b2 * x, -2 * a2 * y
		d2xt, d2yt = 2 * b2, 2 * a2
		incx       = func() { x++; dxt += d2xt; t += dxt }
		incy       = func() { y--; dyt += d2yt; t += dyt }
	)
	point := func(x, y int) {
		draw.Draw(dst, image.Rect(x, y, x+1, yc), src, sp, draw.Over)
		//draw.Draw(dst, image.Rect(x, y, x+1, y-1), src, sp, draw.Over)
		// Perspective-retaining lines
		//draw.Draw(dst, image.Rect(x, y, x+1, yc/2), src, sp, draw.Over)
	}

	for y >= 0 && x <= a {
		point(xc+x, yc+y)
		if x != 0 || y != 0 {
			point(xc-x, yc-y)
		}
		if x != 0 && y != 0 {
			point(xc+x, yc-y)
			point(xc-x, yc+y)
		}
		if t+b2*x <= crit1 || t+a2*y <= crit3 {
			incx()
		} else if t-a2*y > crit2 {
			incy()
		} else {
			incx()
			incy()
		}
	}
}

/*
func String(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *font.Font, s []byte) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}
*/
