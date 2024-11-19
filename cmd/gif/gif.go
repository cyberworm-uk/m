package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"

	"github.com/cyberworm-uk/m"
)

func myPaletted(r image.Rectangle) *image.Paletted {
	img := image.NewPaletted(r, myPalette)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			img.Set(x, y, myPalette[0])
		}
	}
	return img
}

func generate(width, height, limit int, xstart, xend, ystart, yend float64) *image.Paletted {
	var raw = myPaletted(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	var scale = float64(width) / (xend - xstart)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// map co-ordinates to complex plane
			var real = xstart + float64(x)/scale
			var imag = ystart + float64(y)/scale
			// convert to a complex number and determine if it's bounded within our limits
			var m = m.M(complex(real, imag), limit) * 255 / limit
			// set the co-ordinate value in our RGBA
			if m > 128 {
				raw.Set(x, y, color.RGBA{255 - uint8(m), 255 - uint8(m), 255 - uint8(m), uint8(m) / 10})
			}
		}
	}
	return raw
}

func resolve(width, height, frames, limit int, xstart, xend, ystart, yend float64) *gif.GIF {
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	for i := 1; i <= frames; i++ {
		log.Printf("%v of %v, l=%v\n", i, frames, limit)
		g.Image = append(g.Image, generate(width, height, i*limit/frames, xstart, xend, ystart, yend))
		g.Delay = append(g.Delay, 20)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
	}
	return g
}

func zoom(width, height, frames, limit int, xstart, xend, ystart, yend float64) *gif.GIF {
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	for i := 1; i <= frames; i++ {
		log.Printf("%v of %v, x=%v->%v, y=%v->%v\n", i, frames, xstart, xend, ystart, yend)
		g.Image = append(g.Image, generate(width, height, limit+(i*limit/frames), xstart, xend, ystart, yend))
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
		var xdelta = math.Abs(xstart-xend) * 0.05
		var ydelta = math.Abs(ystart-yend) * 0.05
		ystart += ydelta
		yend -= ydelta
		xstart += xdelta
		xend -= xdelta
	}
	return g
}

func fname(anim string, width, frames, limit int, xstart, xend, ystart, yend float64) string {
	return fmt.Sprintf("m-%s-r%f-r%f-i%f-i%f-%d-%d-%d.gif", anim, xstart, xend, ystart, yend, width, frames, limit)
}

func main() {
	var width = flag.Int("width", 1024, "image output width")
	var xstart = flag.Float64("real-start", -2, "start of the real range")
	var xend = flag.Float64("real-end", 1, "end of the real range")
	var ystart = flag.Float64("imag-start", -1.2, "start of the imaginary range")
	var yend = flag.Float64("imag-end", 1.2, "end of imaginary range")
	var anim = flag.String("anim", "resolve", "type of animation: \"resolve\" or \"zoom\"")
	var frames = flag.Int("frames", 100, "number of frames of animation")
	var limit = flag.Int("limit", 200, "iteration limit for bounding check (precision)")
	flag.Parse()
	if *xend < *xstart {
		xend, xstart = xstart, xend
	}
	if *yend < *ystart {
		yend, ystart = ystart, yend
	}
	// calculate the aspect ratio based on the provided range...
	var aspect = math.Abs(*xstart-*xend) / math.Abs(*ystart-*yend)
	// calculate the height from the width and the aspect ratio...
	var height = int(math.Round(float64(*width) / aspect))
	// generate our GIF....
	var g *gif.GIF
	switch *anim {
	case "resolve":
		g = resolve(*width, height, *frames, *limit, *xstart, *xend, *ystart, *yend)
	case "zoom":
		g = zoom(*width, height, *frames, *limit, *xstart, *xend, *ystart, *yend)
	}
	// open file
	f, e := os.Create(fname(*anim, *width, *frames, *limit, *xstart, *xend, *ystart, *yend))
	if e != nil {
		log.Println(e)
		return
	}
	// write file
	e = gif.EncodeAll(f, g)
	if e != nil {
		log.Println(e)
	}
	// close file
	f.Close()
}
