// Command icongen draws the app icon: a chip outlined in white on a dark
// oil background, labelled with silicon's symbol and atomic number. The
// artwork stays inside the middle of the canvas, so launchers that crop
// icons to circles or rounded squares don't cut it off.
//
//	go run ./cmd/icongen            # writes assets/icon.png
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/mehran9675/silicondb/assets"
)

const (
	size = 1024
	ss   = 4 // supersampling factor
	big  = size * ss
)

func main() {
	out := flag.String("o", "assets/icon.png", "output file")
	flag.Parse()

	mask := image.NewAlpha(image.Rect(0, 0, big, big))
	drawChip(mask)
	drawLabel(mask)

	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r, g, b := background(float64(x)+0.5, float64(y)+0.5)
			a := coverage(mask, x, y)
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(r + (255-r)*a),
				G: uint8(g + (255-g)*a),
				B: uint8(b + (255-b)*a),
				A: 255,
			})
		}
	}
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

// background is dark oil: near-black with a faint green-blue sheen
// toward the top left and a warmer depth toward the bottom right.
func background(x, y float64) (r, g, b float64) {
	t := (x + y) / (2 * size)
	r, g, b = lerp(22, 12, t), lerp(30, 14, t), lerp(30, 17, t)
	sheen := math.Max(0, 1-math.Hypot(x-size*0.3, y-size*0.25)/(size*0.75))
	sheen *= sheen
	r, g, b = r+10*sheen, g+24*sheen, b+26*sheen
	glow := math.Max(0, 1-math.Hypot(x-size*0.5, y-size*0.5)/(size*0.45))
	r, g, b = r+6*glow, g+8*glow, b+10*glow
	return r, g, b
}

// drawChip outlines the package and its pins.
func drawChip(m *image.Alpha) {
	const (
		lo, hi = 340.0, 684.0 // package body
		radius = 46.0
		stroke = 22.0
		pinLen = 64.0
		pinW   = 22.0
	)
	fillRoundRect(m, lo, lo, hi, hi, radius, 255)
	fillRoundRect(m, lo+stroke, lo+stroke, hi-stroke, hi-stroke, radius-stroke, 0)
	for i := 0; i < 4; i++ {
		c := lo + (hi-lo)*(float64(i)+1)/5
		fillRoundRect(m, c-pinW/2, lo-pinLen, c+pinW/2, lo+1, pinW/2, 255)
		fillRoundRect(m, c-pinW/2, hi-1, c+pinW/2, hi+pinLen, pinW/2, 255)
		fillRoundRect(m, lo-pinLen, c-pinW/2, lo+1, c+pinW/2, pinW/2, 255)
		fillRoundRect(m, hi-1, c-pinW/2, hi+pinLen, c+pinW/2, pinW/2, 255)
	}
}

// drawLabel writes "Si" in the middle and "14" in the top left corner,
// as on a periodic table tile.
func drawLabel(m *image.Alpha) {
	ttf, err := assets.Fonts.ReadFile("fonts/Inter-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := opentype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}
	text(m, f, "Si", 168, 512, 590, true)
	text(m, f, "14", 54, 392, 438, false)
}

// text draws s at size px with its baseline at y; centered centers it on
// x, otherwise x is its left edge. Coordinates are in output pixels.
func text(m *image.Alpha, f *opentype.Font, s string, px, x, y float64, centered bool) {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: px * ss, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		log.Fatal(err)
	}
	defer face.Close()
	d := &font.Drawer{Dst: m, Src: image.NewUniform(color.Alpha{A: 255}), Face: face}
	w := d.MeasureString(s)
	startX := fixed.I(int(x * ss))
	if centered {
		startX -= w / 2
	}
	d.Dot = fixed.Point26_6{X: startX, Y: fixed.I(int(y * ss))}
	d.DrawString(s)
}

// fillRoundRect sets the pixels of a rounded rectangle, given in output
// pixels, to a.
func fillRoundRect(m *image.Alpha, x0, y0, x1, y1, r float64, a uint8) {
	X0, Y0, X1, Y1, R := x0*ss, y0*ss, x1*ss, y1*ss, r*ss
	for y := int(Y0); y < int(math.Ceil(Y1)); y++ {
		for x := int(X0); x < int(math.Ceil(X1)); x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			cx := math.Max(X0+R, math.Min(fx, X1-R))
			cy := math.Max(Y0+R, math.Min(fy, Y1-R))
			if math.Hypot(fx-cx, fy-cy) <= R {
				m.SetAlpha(x, y, color.Alpha{A: a})
			}
		}
	}
}

// coverage averages the supersampled mask over one output pixel.
func coverage(m *image.Alpha, x, y int) float64 {
	sum := 0
	for j := 0; j < ss; j++ {
		off := m.PixOffset(x*ss, y*ss+j)
		for i := 0; i < ss; i++ {
			sum += int(m.Pix[off+i])
		}
	}
	return float64(sum) / float64(255*ss*ss)
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }
