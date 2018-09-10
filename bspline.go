package memdraw

import (
	"image"
	"image/draw"
)

// Bezspline takes the same arguments as poly but draws a
// quadratic B-spline (despite its name) rather than a
// polygon. If the first and last points in p are equal,
// the spline has periodic end conditions.
func Bezspline(dst draw.Image, a, b, c, d image.Point, end0, end1, thick int, src image.Image, sp image.Point) {
	p := []image.Point{a, b, c, d}
	//	for t := float64(0); t < 1.0; t+= 0.1{
	//		curve(dst, interp(pt, t), t, thick, src, sp)
	//	}
	const (
		seg  = 10000.0
		step = 1 / seg
	)
	t := 0.0
	q := make([]image.Point, 0, seg+2)
	for i := 0; i < seg; i++ {
		t = float64(i) * step
		q = append(q, curve(dst, interp(p, t), t, thick, src, sp))
	}
	Poly(dst, q, 1, 1, thick, src, sp)
}

// TODO: Remove/cleanup below

func interp(p []image.Point, t float64) []image.Point {
	deg := 2
	n := len(p)
	ndeg := len(p) - 1 - deg

	w := make([]float64, n)
	for i := range w {
		w[i] = 1.0
	}

	k := make([]float64, len(p)+2)
	for i := range k {
		k[i] = float64(i)
	}

	lo, hi := float64(deg), float64(ndeg)
	t *= hi - lo
	t += lo

	q := make([]image.Point, len(p))
	for i := range p {
		q[i] = p[i]
		q[i].X = int(float64(q[i].X) * float64(w[i]))
		q[i].Y = int(float64(q[i].Y) * float64(w[i]))
	}

	s := deg
	for ; s < ndeg; s++ {
		if t >= k[s] && t <= k[s+1] {
			break
		}
	}

	for l := 1; l <= deg+1; l++ {
		for i := s; i > s-deg+l; i-- {
			a := (t - k[i]) / (k[i+deg+1-l] - k[i])
			q[i].X = int((1-a)*float64(q[i-1].X) + a*float64(q[i].X))
			q[i].Y = int((1-a)*float64(q[i-1].Y) + a*float64(q[i].Y))
		}
	}
	return q
}
func curve(dst draw.Image, p []image.Point, t float64, thick int, src image.Image, sp image.Point) image.Point {
	for len(p) != 1 {
		p = op(p, t, dst)
	}
	dst.Set(p[0].X, p[0].Y, src.At(p[0].X, p[0].Y))
	//r := image.ZR.Inset(-1).Add(p[0])
	//draw.Draw(dst, r, src, image.ZP, draw.Src)
	return p[0]
}
func op(p []image.Point, t float64, dst ...draw.Image) []image.Point {
	p2 := make([]image.Point, 0, len(p)-1)

	for i := 0; i < len(p)-1; i++ {
		x := p[i].X + int(float64(p[i+1].X-p[i].X)*t)
		y := p[i].Y + int(float64(p[i+1].Y-p[i].Y)*t)
		p2 = append(p2, image.Pt(x, y))
	}
	return p2
}
