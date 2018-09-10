package memdraw

import (
	"image"
	"image/draw"
	"math"
)

// Bezier draws the N'th degree Bezier curve defined by the
// input points.The end styles are determined by end0
// and end1; the thickness of the curve is 1+2*thick.
// The source is aligned so sp in src corresponds to a in dst.
func BezierN(dst draw.Image, end0, end1, thick int, src image.Image, sp image.Point, pt ...image.Point) {
	// Start with a simple, but extemely slow implementation. Optimize it later when time permits
	flatcurve(dst, pt, thick, src, sp)
}

// Bezier draws the cubic Bezier curve defined by Points
// a, b, c, and d. The end styles are determined by end0
// and end1; the thickness of the curve is 1+2*thick.
// The source is aligned so sp in src corresponds to a in dst.
func Bezier(dst draw.Image, a, b, c, d image.Point, end0, end1, thick int, src image.Image, sp image.Point) {
	// Start with a simple, but extemely slow implementation. Optimize it later when time permits
	flatcurve(dst, []image.Point{a, b, c, d}, thick, src, sp)
}

type tik struct {
	t, i, k int
}

var cache = make(map[tik]int)

type point struct {
	X, Y float64
}

func (p point) Round() image.Point {
	return image.Point{int(math.Round(p.X)), int(math.Round(p.Y))}
}
func fpoint64(p ...image.Point) []point {
	pts := make([]point, 0, len(p))
	for _, v := range p {
		pts = append(pts, point{float64(v.X), float64(v.Y)})
	}
	return pts
}
func flatcurve(dst draw.Image, p []image.Point, thick int, src image.Image, sp image.Point) {
	pts := fpoint64(p...)
	pts2 := []image.Point{p[0]}
	for t := 0.0; t <= 1.00; t += 0.01 {
		q := append([]point{}, pts...)
		for len(q) > 1 {
			for i := 0; i < len(q)-1; i++ {
				q[i].X += (q[i+1].X - q[i].X) * t
				q[i].Y += (q[i+1].Y - q[i].Y) * t
			}
			q = q[:len(q)-1]
		}
		pts2 = append(pts2, q[0].Round())
	}

	Poly(dst, pts2, 1, 1, 1, src, sp)
}

func N(t, i, k int) (n int) {
	if cN, ok := cache[tik{t, i, k}]; ok {
		return cN
	}
	defer func() {
		cache[tik{t, i, k}] = n
	}()
	gt := func(i int) int {
		return N(t, i, k)
	}
	if k == 1 {
		if gt(i) <= t && t <= gt(i+1) {
			return 1
		}
		return 0
	}
	return ((t-gt(i))*N(t, i, k-1))/
		(gt(i+k-1)-gt(i)) + ((gt(i+k)-t)*N(t, i+1, k-1))/
		(gt(i+k)-gt(i+1))
}
