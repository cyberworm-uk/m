package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/cyberworm-uk/m"
)

func generate(width, height, limit int, xstart, xend, ystart, yend float64) *image.RGBA {
	var raw = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// map co-ordinates to complex plane
			var real = (xstart + float64(x)/float64(width)*(xend-xstart))
			var imag = (ystart + float64(y)/float64(height)*(yend-ystart))
			// convert to a complex number and determine if it's bounded within our limits
			var m = m.M(complex(real, imag), limit) * 255 / limit
			// set the co-ordinate value in our RGBA
			raw.Set(x, y, color.RGBA{255 - uint8(m), 255 - uint8(m), 255 - uint8(m), uint8(m)})
		}
	}
	return raw
}

func fname(width, limit int, xstart, xend, ystart, yend float64) string {
	return fmt.Sprintf("m-r%f-r%f-i%f-i%f-%d-%d.png", xstart, xend, ystart, yend, width, limit)
}

func main() {
	var width = flag.Int("width", 1024, "image output width")
	var xstart = flag.Float64("real-start", -2, "start of the real range")
	var xend = flag.Float64("real-end", 1, "end of the real range")
	var ystart = flag.Float64("imag-start", -1.2, "start of the imaginary range")
	var yend = flag.Float64("imag-end", 1.2, "end of imaginary range")
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
	// generate our RGBA matrix of the mandlebrot set for our range.
	var rgba = generate(*width, height, *limit, *xstart, *xend, *ystart, *yend)
	// open file
	f, e := os.Create(fname(*width, *limit, *xstart, *xend, *ystart, *yend))
	if e != nil {
		log.Println(e)
		return
	}
	// write file
	if e = png.Encode(f, rgba); e != nil {
		log.Println(e)
		return
	}
	// close file
	f.Close()
}
