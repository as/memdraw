package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/as/frame"
	"github.com/as/memdraw"
	"github.com/as/shiny/event/key"
	"github.com/as/shiny/event/lifecycle"
	"github.com/as/shiny/event/mouse"
	"github.com/as/shiny/screen"
	"github.com/as/ui"
)

var (
	size   = image.Pt(1024, 868)
	w      = dev.Window()
	dev, _ = ui.Init(&ui.Config{
		Width: size.X, Height: size.Y,
		Title: "harness",
	})
	img, _ = dev.NewBuffer(size)
	fb     = img.RGBA()
	down   uint

	sem = make(chan bool, 1)
)

func blank() {
	if fb == nil {
		panic("graphics: fb == nil")
	}
	draw.Draw(fb, fb.Bounds(), image.Black, image.ZP, draw.Src)
}
func redraw() {
	w.Upload(image.ZP, img, img.Bounds())
	w.Publish()
	if recording {
		capture(img.RGBA())
	}
}

var (
	black  = image.Black
	red    = image.NewUniform(color.RGBA{255, 0, 0, 255})
	green  = image.NewUniform(color.RGBA{0, 255, 0, 255})
	blue   = image.NewUniform(color.RGBA{0, 0, 255, 255})
	yellow = image.NewUniform(color.RGBA{255, 255, 0, 255})
	white  = image.White
	orange = image.NewUniform(color.RGBA{255, 128, 0, 255})
	octant = image.NewUniform(color.RGBA{255, 128, 128, 255})
)

var current int

// grab selects the element of the list nearest to p
// and sets current to the index if such a point is found
func grab(p image.Point, list []image.Point) bool {
	r := image.ZR.Inset(-6).Add(p)
	for i, q := range list {
		if q.In(r) {
			current = i
			return true
		}
	}
	return false
}

func main() {
	d := screen.Dev
	blank()
	redraw()

	var (
		lpts = make([]image.Point, 2)
		pts  = make([]image.Point, 5)
	)

	hz := time.NewTicker(time.Second / 128).C
	fr := frame.New(fb, image.Rect(0, 768, 1024, 868), nil)

	for {
		select {
		case e := <-d.Mouse:
			e = readmouse(e)
			switch {
			case HasButton(3, down):
				copy(lpts, lpts[1:])
				lpts[0] = pts[1]
				lpts[1] = pt(e)
				memdraw.Line(fb, lpts[0], lpts[1], 1, blue, image.ZP)
				redraw()
			case HasButton(2, down):
				select {
				case <-hz:
					// drawCurve(fb, black, pts...)

					copy(pts, pts[1:])
					pts[len(pts)-1] = pt(e)

					drawCurve(fb, next(), pts...)
					redraw()
				default:
				}
			case HasButton(1, down):
				grab(pt(e), pts)
				for down != 0 {
					select {
					case e = <-d.Mouse:
						e = readmouse(e)
					case <-d.Key:
						continue
					}
					drawCurve(fb, black, pts...)
					pts[current] = pt(e)
					drawCurve(fb, uniform, pts...)
					redraw()
				}
			}
		case e := <-d.Key:
			if e.Direction == key.DirRelease {
				continue
			}
			if specialkey(e.Rune) {
				continue
			}
			blank()
			switch e.Code {
			case key.CodeUpArrow:
				pts = append(pts, pts[len(pts)-1])
				fr.Delete(0, fr.Len())
				fr.Insert([]byte(fmt.Sprintf("Degree: %d\n", len(pts)-1)), 0)
			case key.CodeDownArrow:
				if len(pts)-1 != 0 {
					pts = pts[:len(pts)-1]
					if current == len(pts) {
						current--
					}
				}
				fr.Delete(0, fr.Len())
				fr.Insert([]byte(fmt.Sprintf("Degree: %d\n", len(pts)-1)), 0)
			}
			drawCurve(fb, uniform, pts...)
			redraw()
		case e := <-d.Size:
			size = image.Pt(e.WidthPx, e.HeightPx)
		case e := <-d.Paint:
			e = e
			redraw()
		case e := <-d.Lifecycle:
			if e.To == lifecycle.StageDead {
				os.Exit(0)
			}
		}
	}
}

/*
 * Animation Helpers
 */

var (
	raw       [30 * 60]image.RGBA
	pending   [30 * 60]*image.Paletted
	fp        int
	recording bool
	camera    image.Rectangle
	g         *gif.GIF
)

func init() {
	g = &gif.GIF{
		Delay: make([]int, len(pending), len(pending)),
		Image: pending[:],
	}
	for i := range pending {
		tmp := image.NewRGBA(fb.Bounds())
		raw[i] = *tmp
		pending[i] = image.NewPaletted(fb.Bounds(), palette.WebSafe)
		g.Delay[i] = 2
	}
}

func capture(img *image.RGBA) {
	rmin := img.Rect.Min
	s := camera.Bounds().Min
	sp := (s.Y-rmin.Y)*img.Stride + (s.X-rmin.X)*4
	e := camera.Bounds().Max
	ep := (e.Y-rmin.Y)*img.Stride + (e.X-rmin.X)*4
	copy(raw[fp].Pix, img.Pix[sp:ep])

	// draw.Draw(&raw[fp], r, img, camera.Bounds().Min, draw.Src)
	fp++
	if fp == len(pending) {
		fp--
		println("buffer overrun")
	}
}

func specialkey(r rune) bool {
	d := screen.Dev
	switch r {
	case 'r':
		select {
		case e := <-d.Mouse:
			e = readmouse(e)
			camera.Min = pt(e)
		case <-d.Key:
			return true
		}
		for down == 0 {
			e := readmouse(<-d.Mouse)
			memdraw.Border(fb, camera, 1, image.ZP, black)
			camera.Max = pt(e)
			memdraw.Border(fb, camera, 1, image.ZP, red)
			redraw()
		}
		memdraw.Border(fb, camera, 1, image.ZP, black)
		camera = camera.Canon() // lol
		redraw()
	case 's':
		fp = 0
		recording = true
	case 'e':
		recording = false
	case 'p':
		if fp == 0 {
			println("nothing to put")
			break
		}
		fd, err := os.Create("harness.gif")
		if err != nil {
			log.Println("save: error", err)
			break
		}

		r := image.Rectangle{image.ZP, camera.Bounds().Size()}
		for i := 0; i < fp; i++ {
			pending[i] = image.NewPaletted(r, palette.Plan9)
			draw.Draw(pending[i], r, &raw[i], r.Min, draw.Src)
		}
		g.Image = pending[:fp]
		g.Delay = g.Delay[:fp]
		gif.EncodeAll(fd, g)
		fd.Close()
		for i := range raw[:fp] {
			draw.Draw(&raw[i], r, image.Transparent, image.ZP, draw.Src)
		}
		fp = 0
	default:
		return false
	}
	return true
}

/*
 * Curve Helpers
 */

func drawCurve(dst draw.Image, src image.Image, p ...image.Point) {
	drawControls(dst, src, p...)
	// memdraw.Bezspline(fb, p[0],p[1],p[2],p[3],1,1,1, src, image.ZP)
	memdraw.BezierN(fb, 1, 1, 1, src, dst.Bounds().Min, p...)
}

func drawControls(dst draw.Image, src image.Image, p ...image.Point) {
	for _, p := range p {
		memdraw.Ellipse(dst, p, 3, 3, 1, src, dst.Bounds().Min, 1, 1)
	}
}

func drawEllipses(dst draw.Image, src image.Image, c, b, a int, p ...image.Point) {
	for _, p := range p {
		memdraw.Ellipse(dst, p, c, b, a, src, dst.Bounds().Min, 1, 1)
	}
}

/*
 * Mouse Stuff
 */

func pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

const (
	KShift = 1 << iota
	KCtrl
	KAlt
	KMeta
)

func Button(n uint) uint {
	return 1 << n
}
func HasButton(n, mask uint) bool {
	return Button(n)&mask != 0
}

func readmouse(e mouse.Event) mouse.Event {
	if e.Button == 1 {
		if km := e.Modifiers; km&KCtrl != 0 {
			e.Button = 3
		} else if km&KAlt != 0 {
			e.Button = 2
		}
	}
	if dir := e.Direction; dir == 1 {
		down |= 1 << uint(e.Button)
	} else if dir == 2 {
		down &^= 1 << uint(e.Button)
	}
	return e
}

/*
 * Debug Stuff
 */

var rainbow = color.RGBA{255, 0, 0, 255}
var prime = color.RGBA{255, 0, 0, 255}
var uniform = image.NewUniform(rainbow)
var uprime = image.NewUniform(prime)

func next() *image.Uniform {
	rainbow = nextcolor(rainbow)
	prime = inverse(rainbow)
	uniform = image.NewUniform(rainbow)
	uprime = image.NewUniform(prime)
	return uniform
}

// inverse steps through a gradient
func inverse(c color.RGBA) color.RGBA {
	c.R = ^c.R - 147
	c.G = ^c.G - 147
	c.B = ^c.B - 147
	return c
}

// nextcolor steps through a gradient
func nextcolor(c color.RGBA) color.RGBA {
	switch {
	case c.R == 255 && c.G == 0 && c.B == 0:
		c.G += 5
	case c.R == 255 && c.G != 255 && c.B == 0:
		c.G += 5
	case c.G == 255 && c.R != 0:
		c.R -= 5
	case c.R == 0 && c.B != 255:
		c.B += 5
	case c.B == 255 && c.G != 0:
		c.G -= 5
	case c.G == 0 && c.R != 255:
		c.R += 5
	default:
		c.B -= 5
	}
	return c
}
