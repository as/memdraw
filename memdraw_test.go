package memdraw

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

var (
	green  = image.NewUniform(color.RGBA{0, 255, 0, 255})
	blue   = image.NewUniform(color.RGBA{0, 0, 255, 255})
	yellow = image.NewUniform(color.RGBA{255, 255, 0, 255})
	white  = image.White
	orange = image.NewUniform(color.RGBA{255, 128, 0, 255})
	octant = image.NewUniform(color.RGBA{255, 128, 128, 255})
)

func oct(v byte) image.Image {
	return image.NewUniform(color.RGBA{v*20 + 25, v*20 + 25, v*20 + 50, 255})
}

var pt = func(x, y int) image.Point {
	return image.Pt(x+50, y+50)
}
var line = Line
var polyline = Poly
var zp = pt(0, 0)

func TestScratch(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	testClear(img)

	// testRegions(img, 100)

	draw.Draw(img, img.Bounds(), image.Black, zp, draw.Src)
	testBezier(img, 100)

	png.Encode(os.Stdout, img)

}

func testClear(img draw.Image) {
	draw.Draw(img, img.Bounds(), image.Black, zp, draw.Src)
}

func testBezspline(img draw.Image) {
	Bezspline(img, zp, image.Pt(0, 64), image.Pt(64, 0), image.Pt(512, 512), 1, 1, 1, white, zp)
}

func testPoly(img draw.Image) {
	polyline(img, []image.Point{pt(0, 0), pt(10, 10), pt(50, 50), pt(100, 100)}, 1, 1, 1, red, zp)
}

func testChromaDither(img draw.Image, u int) {
	pt := func(x, y int) image.Point { return image.Pt(x+110, y+110) }
	polyline(img, []image.Point{pt(0, 10), pt(50, 500)}, 2, 1, 5, green, zp)

	pt = func(x, y int) image.Point { return image.Pt(x+111, y+110) }
	line(img, pt(0, 10), pt(50, 500), 5, red, zp)

	pt = func(x, y int) image.Point { return image.Pt(x+112, y+110) }
	line(img, pt(0, 10), pt(50, 500), 5, blue, zp)

	pt = func(x, y int) image.Point { return image.Pt(x+113, y+110) }
	line(img, pt(0, 10), pt(50, 500), 5, blue, zp)
}

func testBezier(img draw.Image, u int) {
	Bezier(img, pt(0, 0), pt(0, 100), pt(100, 100), pt(100, 0), 1, 1, 1, green, zp)
	BezierN(img, 1, 1, 1, red, zp,
		pt(0, 0), pt(100, 100), pt(500, 65), pt(1096, 254), pt(540, 565), pt(0, 0),
	)
}

func testRegions(img draw.Image, u int) {
	// draw the straight lines, the quadrants, the octants

	line(img, zp, pt(u, 0), 1, red, zp)
	line(img, zp, pt(0, -u), 1, green, zp)
	line(img, zp, pt(-u, 0), 1, blue, zp)
	line(img, zp, pt(0, u), 1, white, zp)

	line(img, zp, pt(u, -u), 1, orange, zp)
	line(img, zp, pt(-u, -u), 1, orange, zp)
	line(img, zp, pt(-u, u), 1, orange, zp)
	line(img, zp, pt(u, u), 1, orange, zp)

	/*
		o1	(+1, -1/2)
		o2	(+1/2, -1)
		o3	(-1/2, -1)
		o4	(-1, -1/2)
		o5	(-1, +1/2)
		o6	(-1/2, +1)
		o7	(+1, +1/2)
		o8	(+1/2, +1)
	*/

	line(img, zp, pt(u, -u/2), 1, oct(1), zp)
	line(img, zp, pt(u/2, -u), 1, oct(2), zp)
	line(img, zp, pt(-u/2, -u), 1, oct(3), zp)
	line(img, zp, pt(-u, -u/2), 1, oct(4), zp)
	line(img, zp, pt(-u, u/2), 1, oct(5), zp)
	line(img, zp, pt(-u/2, u), 1, oct(6), zp)
	line(img, zp, pt(u, u/2), 1, oct(7), zp)
	line(img, zp, pt(u/2, u), 1, oct(8), zp)
}
