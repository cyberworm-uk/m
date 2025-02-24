package mp4

import (
	"image"
	"io"
	"log"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/gen2brain/x264-go"
	"github.com/mazznoer/colorgrad"
)

func rawToYCbCr(raw *m.Raw, grad colorgrad.Gradient) *x264.YCbCr {
	var img = x264.NewYCbCr(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{raw.Width, raw.Height + raw.Height%2},
		},
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

func RawsToMp4(rs chan *m.Raw, grad colorgrad.Gradient, out io.Writer) {
	var logger = log.New(log.Writer(), "mp4: ", log.Ldate|log.Ltime|log.Lshortfile)
	sample := <-rs
	opts := &x264.Options{
		Width:     sample.Width,
		Height:    sample.Height + sample.Height%2,
		FrameRate: 25,
		Tune:      "grain",
		Preset:    "veryfast",
		Profile:   "baseline",
		LogLevel:  x264.LogWarning,
	}
	var enc, e = x264.NewEncoder(out, opts)
	if e != nil {
		logger.Printf("ERR: %s\n", e)
		return
	}

	for r := range rs {
		e = enc.Encode(rawToYCbCr(r, grad))
		if e != nil {
			logger.Printf("ERR: %s\n", e)
			return
		}
	}

	e = enc.Flush()
	if e != nil {
		log.Printf("ERR: %s\n", e)
		return
	}

	e = enc.Close()
	if e != nil {
		log.Printf("ERR: %s\n", e)
		return
	}
}
