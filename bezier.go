package memdraw

import (
	"fmt"
	"image"
	"image/draw"
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

func flatcurve(dst draw.Image, p []image.Point, thick int, src image.Image, sp image.Point) {
	const (
		seg  = 29.0
		step = 1 / seg
	)
	t := 0.0
	q := make([]image.Point, seg)
	{
		for i := 1; i <= int(seg); i++ {
			t = float64(i-1) * step
			q[len(q)-i] = curve(dst, p, t, thick, src, sp)
		}
	}
	q[0] = p[len(p)-1]
	q[len(q)-1] = p[0]
	Poly(dst, q[:], 1, 1, thick, src, sp)
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
