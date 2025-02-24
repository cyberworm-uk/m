package gif

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"log"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/mazznoer/colorgrad"
)

func rawToPaletted(raw *m.Raw, grad colorgrad.Gradient, pal color.Palette) *image.Paletted {
	var img = image.NewPaletted(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{raw.Width, raw.Height},
		},
		pal,
	)
	for x := range raw.Data {
		for y := range raw.Data[x] {
			img.Set(
				x, y,
				grad.At(float64(raw.Data[x][y])/float64(raw.Limit)),
			)
		}
	}
	return img
}

func RawsToGif(rs chan *m.Raw, grad colorgrad.Gradient) []byte {
	var logger = log.New(log.Writer(), "mp4: ", log.Ldate|log.Ltime|log.Lshortfile)
	var pal = color.Palette{}
	for s := 0; s < 256; s++ {
		pal = append(pal, grad.At(float64(s)*0.00390625))
	}
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	for r := range rs {
		g.Image = append(g.Image, rawToPaletted(r, grad, pal))
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
	}
	var buf = new(bytes.Buffer)
	if e := gif.EncodeAll(buf, g); e != nil {
		logger.Printf("ERR: %s\n", e)
		return []byte{}
	}
	return buf.Bytes()
}
