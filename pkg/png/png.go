package png

import (
	"bytes"
	"image"
	"image/png"
	"log"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/mazznoer/colorgrad"
)

func RawToPng(raw *m.Raw, grad colorgrad.Gradient) []byte {
	var logger = log.New(log.Writer(), "png: ", log.Ldate|log.Ltime|log.Lshortfile)
	var rgba = image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{raw.Width, raw.Height},
		},
	)
	for x := range raw.Data {
		for y := range raw.Data[x] {
			rgba.Set(x, y,
				grad.At(float64(raw.Data[x][y])/float64(raw.Limit)),
			)
		}
	}
	var buf = new(bytes.Buffer)
	if e := png.Encode(buf, rgba); e != nil {
		logger.Printf("ERR: %s\n", e)
		return []byte{}
	}
	return buf.Bytes()
}
